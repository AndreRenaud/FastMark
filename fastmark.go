package main

import (
	"bufio"
	"bytes"
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"path/filepath"
	"slices"
	"strings"
	"sync"

	_ "embed"
	_ "image/jpeg"
	_ "image/png"

	"github.com/AndreRenaud/fastmark/storage"
	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/dialog"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"golang.design/x/clipboard"
	"golang.org/x/image/draw"
)

// clipboardErr records whether clipboard.Init succeeded; Write panics if
// called after a failed Init, so copyToClipboard checks it first.
var clipboardErr error

func copyToClipboard(text string) {
	if clipboardErr != nil {
		log.Printf("Clipboard unavailable: %s", clipboardErr)
		return
	}
	clipboard.Write(clipboard.FmtText, []byte(text))
}

//go:embed icon-128.png
var iconData []byte

// maxImageDimension caps decoded images so they stay within common GPU
// texture size limits.
const maxImageDimension = 8192

type Metadata struct {
	Total          int
	Categorised    int
	Scanned        int
	TotalRegions   int
	CategoryTotals []int
}

func (m Metadata) Summary() string {
	return fmt.Sprintf("Total: %d, Scanned %d (%d%%) Categorised: %d (%d%%)", m.Total, m.Scanned, m.ScannedPercent(), m.Categorised, m.Percent())
}

func (m Metadata) Percent() int {
	if m.Total == 0 {
		return 0
	}
	return m.Categorised * 100 / m.Total
}

func (m Metadata) ScannedPercent() int {
	if m.Total == 0 {
		return 0
	}
	return m.Scanned * 100 / m.Total
}

// decodedImage is the result of an asynchronous image decode. display is what
// should be shown (possibly contrast-stretched); source is the original
// decode, kept so auto-contrast can be re-applied without re-reading the file.
// source is nil when only the contrast setting changed.
type decodedImage struct {
	gen     uint64
	source  image.Image
	display image.Image
}

// appModel holds all application state. It is only mutated on the main
// goroutine (input handlers and Tick); background goroutines communicate
// results back over the decoded/chosenDirs channels. metadata is the
// exception: scan workers update it directly under metadataMu.
type appModel struct {
	backend storage.Storage

	files    []string
	labels   []string
	filesGen int

	selectedIndex int
	filter        string

	currentImage image.Image   // decoded source of the selected file
	displayImage *ebiten.Image // texture currently shown by the editor
	imageGen     uint64        // bumped whenever displayImage changes
	loadGen      uint64        // bumped on navigation; stale decodes are dropped

	currentRegions RegionList
	drawingIndex   int

	// autoContrast stretches each displayed image's histogram so the
	// darkest value maps to 0 and the brightest to 255.
	autoContrast bool

	metadataMu  sync.Mutex
	metadata    Metadata
	metadataGen int

	decoded    chan decodedImage
	chosenDirs chan string
}

func (m *appModel) labelName(index int) string {
	if index >= 0 && index < len(m.labels) {
		return m.labels[index]
	}
	return "unknown"
}

func (m *appModel) metadataSnapshot() Metadata {
	m.metadataMu.Lock()
	defer m.metadataMu.Unlock()
	snap := m.metadata
	snap.CategoryTotals = slices.Clone(snap.CategoryTotals)
	return snap
}

func (m *appModel) categorySummary(meta Metadata) string {
	summary := ""
	for i := range len(meta.CategoryTotals) {
		percent := 0.0
		if meta.Total > 0 {
			percent = float64(meta.CategoryTotals[i]*100) / float64(meta.Total)
		}
		summary += fmt.Sprintf("%s: %d %.1f%%\n", m.labelName(i), meta.CategoryTotals[i], percent)
	}
	return summary
}

// startMetadataScan rescans every label file in the background, updating
// m.metadata as it goes so progress can be displayed live.
func (m *appModel) startMetadataScan() {
	m.metadataMu.Lock()
	m.metadataGen++
	gen := m.metadataGen
	m.metadata = Metadata{Total: len(m.files), CategoryTotals: make([]int, len(m.labels))}
	m.metadataMu.Unlock()

	files := slices.Clone(m.files)
	backend := m.backend

	go func() {
		filesChan := make(chan string, len(files))
		var wg sync.WaitGroup
		// This is mostly blocked by file I/O, especially on network drives,
		// so run a bunch of parallel workers to compensate
		const workerCount = 50
		wg.Add(workerCount)
		for range workerCount {
			go func() {
				defer wg.Done()
				for file := range filesChan {
					ext := filepath.Ext(file)
					labelFile := filepath.Join("labels", strings.TrimSuffix(file, ext)+".txt")
					// An error just means the image has no label file yet;
					// count it as scanned but uncategorised.
					regions, _ := LoadRegionList(backend, labelFile)
					m.metadataMu.Lock()
					if m.metadataGen == gen {
						for _, region := range regions.Regions {
							if region.index >= 0 && region.index < len(m.metadata.CategoryTotals) {
								m.metadata.CategoryTotals[region.index]++
							}
							m.metadata.TotalRegions++
						}
						if len(regions.Regions) > 0 {
							m.metadata.Categorised++
						}
						m.metadata.Scanned++
					}
					m.metadataMu.Unlock()
				}
			}()
		}
		for _, file := range files {
			filesChan <- file
		}
		close(filesChan)
		wg.Wait()
	}()
}

// autoContrastImage returns a copy of src with its histogram stretched so
// that the darkest channel value present maps to 0 and the brightest to 255,
// applying the same linear mapping to every channel to preserve colour.
func autoContrastImage(src image.Image) image.Image {
	bounds := src.Bounds()

	// RGBA() reports channels in the range [0, 0xffff].
	minV := uint32(0xffff)
	maxV := uint32(0)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, _ := src.At(x, y).RGBA()
			for _, c := range [3]uint32{r, g, b} {
				if c < minV {
					minV = c
				}
				if c > maxV {
					maxV = c
				}
			}
		}
	}

	// Flat or single-colour image: nothing to stretch.
	if maxV <= minV {
		return src
	}

	scale := float64(0xffff) / float64(maxV-minV)
	stretch := func(c uint32) uint8 {
		v := float64(c-minV) * scale
		if v > 0xffff {
			v = 0xffff
		}
		return uint8(uint32(v) >> 8)
	}

	dst := image.NewRGBA(bounds)
	for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
		for x := bounds.Min.X; x < bounds.Max.X; x++ {
			r, g, b, a := src.At(x, y).RGBA()
			dst.SetRGBA(x, y, color.RGBA{stretch(r), stretch(g), stretch(b), uint8(a >> 8)})
		}
	}
	return dst
}

func loadImage(backend storage.Storage, filename string) (image.Image, error) {
	f, err := backend.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	if err != nil {
		return nil, err
	}
	return capImageSize(img, maxImageDimension), nil
}

// capImageSize downscales img so neither dimension exceeds maxDim.
func capImageSize(img image.Image, maxDim int) image.Image {
	b := img.Bounds()
	if b.Dx() <= maxDim && b.Dy() <= maxDim {
		return img
	}
	scale := min(float64(maxDim)/float64(b.Dx()), float64(maxDim)/float64(b.Dy()))
	w := max(1, int(float64(b.Dx())*scale))
	h := max(1, int(float64(b.Dy())*scale))
	log.Printf("Downscaling %dx%d image to %dx%d", b.Dx(), b.Dy(), w, h)
	dst := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.ApproxBiLinear.Scale(dst, dst.Bounds(), img, b, draw.Over, nil)
	return dst
}

// loadFile starts an asynchronous decode of filename and synchronously loads
// its regions. The decode result arrives on m.decoded and is applied in
// Root.Tick (ebiten.NewImageFromImage must run on the main goroutine).
func (m *appModel) loadFile(filename string) {
	m.currentImage = nil
	m.displayImage = nil
	m.imageGen++
	m.loadGen++
	gen := m.loadGen
	backend := m.backend
	autoContrast := m.autoContrast

	go func() {
		img, err := loadImage(backend, "images/"+filename)
		if err != nil {
			log.Printf("Error loading image %s: %s", filename, err)
			return
		}
		display := img
		if autoContrast {
			display = autoContrastImage(img)
		}
		m.decoded <- decodedImage{gen: gen, source: img, display: display}
	}()

	ext := filepath.Ext(filename)
	labelFile := filepath.Join("labels", strings.TrimSuffix(filename, ext)+".txt")

	var err error
	m.currentRegions, err = LoadRegionList(m.backend, labelFile)
	if err != nil {
		log.Printf("Error loading regions for %s: %s", filename, err)
	}
}

// regenerateDisplayImage re-applies the auto-contrast setting to the cached
// source image in the background.
func (m *appModel) regenerateDisplayImage() {
	img := m.currentImage
	if img == nil {
		return
	}
	gen := m.loadGen
	autoContrast := m.autoContrast
	go func() {
		display := img
		if autoContrast {
			display = autoContrastImage(img)
		}
		m.decoded <- decodedImage{gen: gen, display: display}
	}()
}

func (m *appModel) getClosestRegion(click image.Point, imageWidth int, imageHeight int) int {
	for i, region := range m.currentRegions.Regions {
		w := int(float32(region.width) * float32(imageWidth))
		h := int(float32(region.height) * float32(imageHeight))
		x := int(float32(region.xMid)*float32(imageWidth)) - w/2
		y := int(float32(region.yMid)*float32(imageHeight)) - h/2
		if click.X >= x && click.X <= x+w && click.Y >= y && click.Y <= y+h {
			log.Printf("Clicked on region %d", i)
			return i
		}
	}

	// If we're close to a region, and it's small, then assume we just missed and select it
	for i, region := range m.currentRegions.Regions {
		w := int(float32(region.width) * float32(imageWidth))
		h := int(float32(region.height) * float32(imageHeight))
		x := int(float32(region.xMid)*float32(imageWidth)) - w/2
		y := int(float32(region.yMid)*float32(imageHeight)) - h/2
		if w > 10 || h > 10 {
			continue
		}
		if click.X >= x-5 && click.X <= x+w+5 && click.Y >= y-5 && click.Y <= y+h+5 {
			log.Printf("Clicked near region %d", i)
			return i
		}
	}
	return -1
}

type Root struct {
	guigui.DefaultWidget

	model appModel

	background  basicwidget.Background
	statusText  basicwidget.Text
	jumpLabel   basicwidget.Text
	jumpInput   basicwidget.TextInput
	fileList    basicwidget.List[int]
	split       splitter
	editorPanel basicwidget.Panel
	pane        editorPane

	sidebarWidth   int
	dragStartWidth int
	contentWidth   int

	builtFilesGen int

	rootItems    []guigui.LinearLayoutItem
	jumpRowItems []guigui.LinearLayoutItem
	mainRowItems []guigui.LinearLayoutItem
}

// WriteStateKey exposes the state that can change outside input handlers
// (decode results and directory changes applied in Tick, metadata updated by
// scan workers) so the framework rebuilds when it changes.
func (r *Root) WriteStateKey(w *guigui.StateKeyWriter) {
	m := &r.model
	w.WriteInt(m.selectedIndex)
	w.WriteInt(m.filesGen)
	w.WriteInt(len(m.files))
	w.WriteUint64(m.imageGen)
	w.WriteInt(m.drawingIndex)
	w.WriteBool(m.autoContrast)
	w.WriteInt(len(m.currentRegions.Regions))
	if m.backend != nil {
		w.WriteString(m.backend.Describe())
	}
	meta := m.metadataSnapshot()
	w.WriteInt(meta.Total)
	w.WriteInt(meta.Scanned)
	w.WriteInt(meta.Categorised)
	w.WriteInt(meta.TotalRegions)
	for _, c := range meta.CategoryTotals {
		w.WriteInt(c)
	}
}

func (r *Root) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&r.background)
	adder.AddWidget(&r.statusText)
	adder.AddWidget(&r.jumpLabel)
	adder.AddWidget(&r.jumpInput)
	adder.AddWidget(&r.fileList)
	adder.AddWidget(&r.split)
	adder.AddWidget(&r.editorPanel)

	m := &r.model
	context.SetButtonInputReceptive(r, true)

	r.statusText.SetValue(fmt.Sprintf("Fast Mark image tagging %d/%d images", m.selectedIndex, len(m.files)))

	r.jumpLabel.SetValue("Jump to")
	r.jumpLabel.SetVerticalAlign(basicwidget.VerticalAlignMiddle)
	r.jumpInput.SetPlaceholder("Substring to jump to")
	r.jumpInput.OnValueChanged(func(context *guigui.Context, text string, committed bool) {
		m.filter = text
		r.jumpTo()
	})

	if r.builtFilesGen != m.filesGen {
		r.fileList.SetItemsByStrings(m.files)
		r.builtFilesGen = m.filesGen
	}
	r.fileList.OnItemSelected(func(context *guigui.Context, index int) {
		if index != m.selectedIndex {
			r.selectFile(index)
		}
	})

	r.split.OnDragStart(func() {
		r.dragStartWidth = r.sidebarWidth
	})
	r.split.OnDrag(func(deltaX int) {
		u := basicwidget.UnitSize(context)
		width := r.dragStartWidth + deltaX
		width = max(width, u*4)
		if r.contentWidth > 0 {
			width = min(width, r.contentWidth-u*8)
		}
		r.sidebarWidth = width
	})

	r.editorPanel.SetContent(&r.pane)
	r.editorPanel.SetContentConstraints(basicwidget.PanelContentConstraintsFixedWidth)
	r.pane.SetModel(m)

	return nil
}

func (r *Root) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	u := basicwidget.UnitSize(context)

	r.contentWidth = widgetBounds.Bounds().Dx()
	if r.sidebarWidth == 0 {
		r.sidebarWidth = u * 12
	}

	r.jumpRowItems = slices.Delete(r.jumpRowItems, 0, len(r.jumpRowItems))
	r.jumpRowItems = append(r.jumpRowItems,
		guigui.LinearLayoutItem{Widget: &r.jumpLabel},
		guigui.LinearLayoutItem{Widget: &r.jumpInput, Size: guigui.FlexibleSize(1)},
	)
	jumpRow := guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     r.jumpRowItems,
		Gap:       u / 4,
	}

	r.mainRowItems = slices.Delete(r.mainRowItems, 0, len(r.mainRowItems))
	r.mainRowItems = append(r.mainRowItems,
		guigui.LinearLayoutItem{Widget: &r.fileList, Size: guigui.FixedSize(r.sidebarWidth)},
		guigui.LinearLayoutItem{Widget: &r.split, Size: guigui.FixedSize(u / 3)},
		guigui.LinearLayoutItem{Widget: &r.editorPanel, Size: guigui.FlexibleSize(1)},
	)
	mainRow := guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     r.mainRowItems,
	}

	r.rootItems = slices.Delete(r.rootItems, 0, len(r.rootItems))
	r.rootItems = append(r.rootItems,
		guigui.LinearLayoutItem{Widget: &r.statusText},
		guigui.LinearLayoutItem{Layout: &jumpRow},
		guigui.LinearLayoutItem{Layout: &mainRow, Size: guigui.FlexibleSize(1)},
	)

	layouter.LayoutWidget(&r.background, widgetBounds.Bounds())
	(guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     r.rootItems,
		Gap:       u / 4,
		Padding:   guigui.Padding{Start: u / 2, Top: u / 2, End: u / 2, Bottom: u / 2},
	}).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

// Tick applies results from background goroutines on the main goroutine.
func (r *Root) Tick(context *guigui.Context, widgetBounds *guigui.WidgetBounds) error {
	m := &r.model
	for {
		select {
		case d := <-m.decoded:
			if d.gen != m.loadGen {
				continue // stale result from a file we've navigated away from
			}
			if d.source != nil {
				m.currentImage = d.source
			}
			m.displayImage = ebiten.NewImageFromImage(d.display)
			m.imageGen++
		case dir := <-m.chosenDirs:
			m.backend = storage.NewStorage(dir)
			r.updateFiles()
		default:
			return nil
		}
	}
}

// keyRepeating reports whether key was just pressed or is being held long
// enough to auto-repeat, matching basicwidget's repeat timing.
func keyRepeating(key ebiten.Key) bool {
	if !ebiten.IsKeyPressed(key) {
		return false
	}
	duration := inpututil.KeyPressDuration(key)
	if duration == 1 {
		return true
	}
	delay := ebiten.TPS() * 2 / 5
	if duration < delay {
		return false
	}
	return (duration-delay)%4 == 0
}

func (r *Root) HandleButtonInput(context *guigui.Context, widgetBounds *guigui.WidgetBounds) guigui.HandleInputResult {
	// Don't treat typing in the jump-to filter as navigation.
	if context.IsFocusedOrHasFocusedDescendant(&r.jumpInput) {
		return guigui.HandleInputResult{}
	}

	m := &r.model

	if keyRepeating(ebiten.KeyDown) || keyRepeating(ebiten.KeyJ) {
		r.selectFile(m.selectedIndex + 1)
		return guigui.HandleInputByWidget(r)
	}
	if keyRepeating(ebiten.KeyUp) || keyRepeating(ebiten.KeyK) {
		r.selectFile(m.selectedIndex - 1)
		return guigui.HandleInputByWidget(r)
	}
	if keyRepeating(ebiten.KeyLeft) || keyRepeating(ebiten.KeyH) {
		if m.drawingIndex > 0 {
			m.drawingIndex--
		}
		return guigui.HandleInputByWidget(r)
	}
	if keyRepeating(ebiten.KeyRight) || keyRepeating(ebiten.KeyL) {
		if m.drawingIndex < len(m.labels)-1 {
			m.drawingIndex++
		}
		return guigui.HandleInputByWidget(r)
	}
	for i := range 10 {
		if inpututil.IsKeyJustPressed(ebiten.KeyDigit0+ebiten.Key(i)) && i < len(m.labels) {
			m.drawingIndex = i
			return guigui.HandleInputByWidget(r)
		}
	}
	if keyRepeating(ebiten.KeyN) {
		direction := 1
		if ebiten.IsKeyPressed(ebiten.KeyShiftLeft) || ebiten.IsKeyPressed(ebiten.KeyShiftRight) {
			direction = -1
		}
		// Find the next image that's not labeled
		for i := m.selectedIndex + direction; i < len(m.files) && i >= 0; i += direction {
			filename := m.files[i]
			ext := filepath.Ext(filename)
			labelFile := filepath.Join("labels", strings.TrimSuffix(filename, ext)+".txt")
			regions, err := LoadRegionList(m.backend, labelFile)
			if err != nil || len(regions.Regions) == 0 {
				log.Printf("Found unlabeled image %s", filename)
				r.selectFile(i)
				break
			}
		}
		return guigui.HandleInputByWidget(r)
	}

	return guigui.HandleInputResult{}
}

func (r *Root) selectFile(i int) {
	m := &r.model
	if i < 0 || i >= len(m.files) {
		log.Printf("Invalid file index %d (max %d)", i, len(m.files))
		return
	}
	m.selectedIndex = i
	r.pane.editor.cancelDrawing()
	// Set the model index before syncing the list so the OnItemSelected
	// callback this triggers sees an up-to-date model and doesn't recurse.
	r.fileList.SelectItemByIndex(i)
	r.fileList.EnsureItemVisibleByIndex(i)
	m.loadFile(m.files[i])
}

// jumpTo selects the first file that contains the filter value as a
// substring, leaving the full list displayed so adjacent images can be
// compared with the up/down keys.
func (r *Root) jumpTo() {
	m := &r.model
	if m.filter == "" {
		return
	}
	needle := strings.ToLower(m.filter)
	for i, f := range m.files {
		if strings.Contains(strings.ToLower(f), needle) {
			r.selectFile(i)
			return
		}
	}
}

// selectDirectory shows the native directory picker on a goroutine (it
// marshals itself to the main thread) and delivers the result to Tick.
func (m *appModel) selectDirectory() {
	go func() {
		newDirectory, err := dialog.Directory().Title("Load images").Browse()
		if err != nil {
			if !errors.Is(err, dialog.ErrCancelled) {
				log.Printf("Error selecting directory: %s", err)
			}
			return
		}
		m.chosenDirs <- newDirectory
	}()
}

func (r *Root) updateFiles() {
	m := &r.model
	m.files = nil

	match, err := m.backend.Glob("images", "*")
	if err != nil {
		log.Printf("Error listing files: %s", err)
	} else {
		for _, f := range match {
			ext := strings.ToLower(filepath.Ext(f))
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				m.files = append(m.files, filepath.Base(f))
			}
		}
	}
	slices.Sort(m.files)

	file, err := m.backend.Open("labels.txt")
	if err != nil {
		log.Printf("Error opening labels file: %s", err)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		m.labels = nil
		for scanner.Scan() {
			m.labels = append(m.labels, scanner.Text())
		}
	}

	m.filesGen++
	m.startMetadataScan()
	r.selectFile(0)
}

func main() {
	directory := flag.String("directory", "", "Directory to load images from")
	flag.Parse()

	if err := RegionsInit(); err != nil {
		log.Printf("Error loading regions: %s", err)
	}

	if clipboardErr = clipboard.Init(); clipboardErr != nil {
		log.Printf("Error initialising clipboard: %s", clipboardErr)
	}

	root := &Root{}
	m := &root.model
	m.decoded = make(chan decodedImage, 8)
	m.chosenDirs = make(chan string, 1)
	if *directory != "" {
		m.backend = storage.NewStorage(*directory)
	} else {
		m.backend = &storage.DummyStorage{}
	}

	if icon, _, err := image.Decode(bytes.NewReader(iconData)); err == nil {
		ebiten.SetWindowIcon([]image.Image{icon})
	} else {
		log.Printf("Error setting icon: %s", err)
	}

	root.updateFiles()

	if err := guigui.Run(root, &guigui.RunOptions{
		Title:      "Fast Mark Image Tagging",
		WindowSize: image.Pt(1024, 768),
	}); err != nil {
		log.Fatal(err)
	}
}
