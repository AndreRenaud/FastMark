package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"path/filepath"
	"slices"
	"strings"

	_ "embed"
	_ "image/jpeg"
	_ "image/png"

	"github.com/sqweek/dialog"

	"github.com/AllenDang/giu"
)

//go:embed icon-128.png
var iconData []byte

var (
	files      []string
	fileRows   []*giu.TableRowWidget
	fileLabels []*giu.SelectableWidget

	currentImageTexture *giu.Texture
	currentImage        image.Image
	selectedIndex       int
	currentRegions      RegionList

	drawingRect  bool
	drawingStart image.Point
	drawingIndex int

	splitPos = float32(300)

	labels []string

	wnd *giu.MasterWindow

	backend Storage
)

func labelName(index int) string {
	if index >= 0 && index < len(labels) {
		return labels[index]
	}
	return "unknown"
}

func loadImage(filename string) (image.Image, error) {
	f, err := backend.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, _, err := image.Decode(f)
	return img, nil
}

func drawFile(filename string) {
	currentImageTexture = nil
	currentImage = nil

	if img, err := loadImage("images/" + filename); err != nil {
		log.Printf("Error loading image %s: %s", filename, err)
	} else {
		giu.EnqueueNewTextureFromRgba(img, func(t *giu.Texture) {
			currentImageTexture = t
			currentImage = img
		})
	}

	ext := filepath.Ext(filename)

	labelFile := filepath.Join("labels", strings.TrimSuffix(filename, ext)+".txt")

	var err error
	currentRegions, err = LoadRegionList(backend, labelFile)
	if err != nil {
		log.Printf("Error loading regions for %s: %s", filename, err)
	}
}

func selectDirectory() {
	newDirectory, err := dialog.Directory().Title("Load images").Browse()
	if err != nil {
		log.Printf("Error selecting directory: %s", err)
		return
	}
	backend = Storage{newDirectory}
	updateFiles()
}

func selectFile(i int) {
	if i < 0 || i >= len(files) {
		log.Printf("Invalid file index %d (max %d)", i, len(files))
		return
	}
	if selectedIndex > 0 && selectedIndex < len(fileLabels) {
		old := fileLabels[selectedIndex]
		if old != nil {
			old.Selected(false)
		}
	}
	file := files[i]
	selectedIndex = i
	fileLabels[i].Selected(true)
	drawingRect = false
	drawFile(file)
}

func updateFiles() {
	files = []string{}

	match, err := backend.Glob("images", "*")
	if err != nil {
		log.Printf("Error listing files: %s", err)
	} else {
		for _, m := range match {
			ext := filepath.Ext(m)
			if ext == ".jpg" || ext == ".jpeg" || ext == ".png" {
				files = append(files, filepath.Base(m))
			}
		}
	}
	slices.Sort(files)
	fileRows = make([]*giu.TableRowWidget, len(files))
	fileLabels = make([]*giu.SelectableWidget, len(files))
	for i, file := range files {
		fileLabels[i] = giu.Selectable(file)
		fileLabels[i].OnClick(func() {
			selectFile(i)
		})
		fileRows[i] = giu.TableRow(fileLabels[i])
	}

	file, err := backend.Open("labels.txt")
	if err != nil {
		log.Printf("Error opening labels file: %s", err)
	} else {
		defer file.Close()
		scanner := bufio.NewScanner(file)
		labels = nil
		for scanner.Scan() {
			labels = append(labels, scanner.Text())
		}
	}

	selectFile(0)
}

func getClosestRegion(click image.Point, imageWidth int, imageHeight int) int {
	for i, region := range currentRegions.Regions {
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
	for i, region := range currentRegions.Regions {
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

func loop() {
	window := giu.SingleWindow()
	var file string
	if selectedIndex >= 0 && selectedIndex < len(files) {
		file = files[selectedIndex]
	}

	windowWidth, windowHeight := wnd.GetSize()

	var regionSummary string
	for i, region := range currentRegions.Regions {
		if i > 0 {
			regionSummary += ", "
		}
		regionSummary += fmt.Sprintf("%d: %s mid=%.3f,%.3f size=%.3fx%.3f", i, labelName(region.index), region.xMid, region.yMid, region.width, region.height)
	}

	window.Layout(
		giu.Label(fmt.Sprintf("Fast Mark image tagging %d/%d images", selectedIndex, len(files))),
		giu.SplitLayout(giu.DirectionVertical, &splitPos,
			giu.Table().
				FastMode(true).
				Columns(giu.TableColumn("Files")).
				Rows(fileRows...).
				Size(giu.Auto, giu.Auto),
			giu.Column(
				giu.Row(
					giu.Button("Change Directory").OnClick(selectDirectory),
					giu.Label(backend.Describe()),
				),
				giu.Labelf("Current file: %s", file),
				giu.Labelf("Regions: %s", regionSummary),
				giu.Style().SetColor(giu.StyleColorText, RegionIndexColor(drawingIndex)).To(
					giu.Labelf("Drawing label: %d %s (Press 1-9 to select new type)\n", drawingIndex, labelName(drawingIndex)),
				),
				giu.Custom(func() {
					canvas := giu.GetCanvas()
					pos := giu.GetCursorScreenPos()
					wx, wy := window.CurrentPosition()

					canvasPos := image.Point{X: pos.X - int(wx), Y: pos.Y - int(wy)}

					// Draw the image as big as the remaining space with a small border
					imageWidth := windowWidth - canvasPos.X - 10
					imageHeight := windowHeight - canvasPos.Y - 10

					if currentImageTexture != nil {
						// Make sure we maintain the ratio
						size := currentImage.Bounds()
						if float32(imageWidth)/float32(size.Dx()) < float32(imageHeight)/float32(size.Dy()) {
							imageHeight = int(float32(size.Dy()) * float32(imageWidth) / float32(size.Dx()))
						} else {
							imageWidth = int(float32(size.Dx()) * float32(imageHeight) / float32(size.Dy()))
						}
						max := pos.Add(image.Point{X: imageWidth, Y: imageHeight})
						canvas.AddImage(currentImageTexture, pos, max)
					}
					// Check if the user has stopped drawing
					if drawingRect {
						end := giu.GetMousePos().Sub(pos)
						if !giu.IsMouseDown(giu.MouseButtonLeft) {
							// Create a new well formed region clamped within the image
							newRect := image.Rect(drawingStart.X, drawingStart.Y, end.X, end.Y)
							newRect = newRect.Intersect(image.Rect(0, 0, imageWidth, imageHeight)).Canon()
							log.Printf("New rect: %v", newRect)
							newRegion := Region{
								xMid:   (float64(newRect.Dx())/2 + float64(newRect.Min.X)) / float64(imageWidth),
								yMid:   (float64(newRect.Dy())/2 + float64(newRect.Min.Y)) / float64(imageHeight),
								width:  float64(newRect.Dx()) / float64(imageWidth),
								height: float64(newRect.Dy()) / float64(imageHeight),
								index:  drawingIndex, // TODO: Make this configurable
							}
							currentRegions.AddRegion(newRegion)
							drawingRect = false
						}
						canvas.AddRect(pos.Add(drawingStart), pos.Add(end), color.RGBA{255, 0, 0, 255}, 0, 0, 2)
					} else if giu.IsMouseDown(giu.MouseButtonLeft) {
						// Get the current screen position of the cursor
						scr := giu.GetMousePos().Sub(pos)
						if scr.X >= 0 && scr.X <= imageWidth && scr.Y >= 0 && scr.Y <= imageHeight {
							drawingRect = true
							drawingStart = scr
						}
					}

					if giu.IsMouseClicked(giu.MouseButtonRight) {
						changeRegion := -1
						// If we're pressing a number key, change the region type, otherwise delete it
						for key := giu.Key0; key <= giu.Key9; key++ {
							if giu.IsKeyDown(key) {
								changeRegion = int(key - giu.Key0)
							}
						}
						// Find the region that was clicked
						click := giu.GetMousePos().Sub(pos)
						index := getClosestRegion(click, imageWidth, imageHeight)
						if index >= 0 {
							if changeRegion >= 0 {
								currentRegions.Regions[index].index = changeRegion
								currentRegions.Save()
							} else {
								currentRegions.Remove(index)
							}
						}
					}

					for _, region := range currentRegions.Regions {
						w := int(float32(region.width) * float32(imageWidth))
						h := int(float32(region.height) * float32(imageHeight))
						x := int(float32(region.xMid)*float32(imageWidth)) - w/2
						y := int(float32(region.yMid)*float32(imageHeight)) - h/2
						color := region.Color()
						canvas.AddRect(pos.Add(image.Point{X: x, Y: y}), pos.Add(image.Point{X: x + w, Y: y + h}), color, 0, 0, 2)
						canvas.AddText(pos.Add(image.Point{X: x + w/2, Y: y - 20}), color, fmt.Sprintf("%s - %d", labelName(region.index), region.index))
					}
				}),
			),
		),
	)

	if giu.IsKeyPressed(giu.KeyDown) || giu.IsKeyPressed(giu.KeyJ) {
		selectFile(selectedIndex + 1)
	}
	if giu.IsKeyPressed(giu.KeyUp) || giu.IsKeyPressed(giu.KeyK) {
		selectFile(selectedIndex - 1)
	}
	if giu.IsKeyPressed(giu.KeyLeft) || giu.IsKeyPressed(giu.KeyH) {
		drawingIndex--
		if drawingIndex < 0 {
			drawingIndex = 0
		}
	}
	if giu.IsKeyPressed(giu.KeyRight) || giu.IsKeyPressed(giu.KeyL) {
		drawingIndex++
		if drawingIndex >= len(labels) {
			drawingIndex = len(labels) - 1
		}
	}
	for i := 0; i < 9; i++ {
		if giu.IsKeyPressed(giu.Key(int(giu.Key0) + i)) {
			drawingIndex = i
		}
	}
	if giu.IsKeyPressed(giu.KeyN) {
		direction := 1
		if giu.IsKeyDown(giu.KeyLeftShift) || giu.IsKeyDown(giu.KeyRightShift) {
			direction = -1
		}
		// Find the next region that's not labeled
		for i := selectedIndex + direction; i < len(files) && i >= 0; i += direction {
			filename := files[i]
			ext := filepath.Ext(filename)

			labelFile := filepath.Join("labels", strings.TrimSuffix(filename, ext)+".txt")

			regions, err := LoadRegionList(backend, labelFile)
			log.Printf("%d Loaded %d regions: %s", i, len(regions.Regions), labelFile)
			if err != nil || len(regions.Regions) == 0 {
				selectFile(i)
				break
			}
		}
	}
}

func main() {
	directory := flag.String("directory", "", "Directory to load images from")
	flag.Parse()

	wnd = giu.NewMasterWindow("Fast Mark Image Tagging", 1024, 768, 0)
	if *directory != "" {
		backend = Storage{*directory}
	}
	icon, _, err := image.Decode(bytes.NewReader(iconData))
	if err == nil {
		wnd.SetIcon(icon)
	} else {
		log.Printf("Error setting icon: %s", err)
	}

	updateFiles()
	wnd.Run(loop)
}
