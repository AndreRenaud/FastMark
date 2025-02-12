package main

import (
	"bufio"
	"flag"
	"fmt"
	"image"
	"image/color"
	"log"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"

	"github.com/sqweek/dialog"

	"github.com/AllenDang/giu"
)

var (
	files      []string
	fileRows   []*giu.TableRowWidget
	fileLabels []*giu.SelectableWidget

	currentImage   *giu.Texture
	selectedIndex  int
	currentRegions []Region

	drawingRect  bool
	drawingStart image.Point
	drawingIndex int

	directory string
	splitPos  = float32(300)

	labels []string

	wnd *giu.MasterWindow
)

type Region struct {
	xMid   float64
	yMid   float64
	width  float64
	height float64

	index int
}

func labelName(index int) string {
	if index >= 0 && index < len(labels) {
		return labels[index]
	}
	return "unknown"
}

func parseLabelData(filename string) []Region {
	file, err := os.Open(filename)
	if err != nil {
		log.Printf("Error opening file %s: %s", filename, err)
		return nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	var retval []Region
	for scanner.Scan() {
		columns := strings.Fields(scanner.Text())
		if len(columns) != 5 {
			log.Printf("Invalid line: %s in %s", scanner.Text(), filename)
			continue
		}
		region := Region{}
		region.index, err = strconv.Atoi(columns[0])
		region.xMid, err = strconv.ParseFloat(columns[1], 64)
		region.yMid, err = strconv.ParseFloat(columns[2], 64)
		region.width, err = strconv.ParseFloat(columns[3], 64)
		region.height, err = strconv.ParseFloat(columns[4], 64)

		if !region.IsValid() {
			log.Printf("Invalid region: %#v in %s", columns, filename)
			continue
		}
		retval = append(retval, region)
	}

	return retval
}

func saveRegions() {
	if selectedIndex < 0 || selectedIndex >= len(files) {
		log.Printf("Invalid file index %d", selectedIndex)
		return
	}
	filename := files[selectedIndex]
	ext := filepath.Ext(filename)
	labelFile := filepath.Join(directory, "../labels/", strings.TrimSuffix(filename, ext)+".txt")
	log.Printf("Saving regions to %s", labelFile)
	file, err := os.OpenFile(labelFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Printf("Error creating file %s: %s", filename, err)
	}
	defer file.Close()
	for _, region := range currentRegions {
		fmt.Fprintf(file, "%d %f %f %f %f\n", region.index, region.xMid, region.yMid, region.width, region.height)
	}
}

func (r Region) Color() color.Color {
	switch r.index {
	case 0:
		return color.RGBA{255, 128, 64, 255}
	case 1:
		return color.RGBA{255, 0, 0, 255}
	case 2:
		return color.RGBA{0, 255, 0, 255}
	case 3:
		return color.RGBA{0, 0, 255, 255}
	case 4:
		return color.RGBA{255, 255, 0, 255}
	case 5:
		return color.RGBA{255, 0, 255, 255}
	case 6:
		return color.RGBA{0, 255, 255, 255}
	default:
		return color.YCbCr{255, uint8(r.index * 16), uint8(r.index * 16)}
	}
}

func (r Region) IsValid() bool {
	if r.width <= 0 || r.height <= 0 || r.width > 1 || r.height > 1 {
		log.Printf("Invalid width/height: %#v", r)
		return false
	}
	if r.xMid < 0 || r.xMid > 1 || r.yMid < 0 || r.yMid > 1 {
		log.Printf("Invalid x/y mid: %#v", r)
		return false
	}
	if r.xMid-r.width/2 < 0 || r.xMid+r.width/2 > 1 {
		log.Printf("Invalid x range: %#v", r)
		return false
	}
	if r.yMid-r.height/2 < 0 || r.yMid+r.height/2 > 1 {
		log.Printf("Invalid y range: %#v %f %f", r, r.yMid-r.height/2, r.yMid+r.height/2)
		return false
	}

	// TODO: Is this legit? These are too small to be useful
	if r.width < 0.001 || r.height < 0.001 {
		return false
	}
	return true
}

func drawFile(filename string) {
	currentImage = nil
	fullFilename := filepath.Join(directory, filename)
	if rgba, err := giu.LoadImage(fullFilename); err != nil {
		log.Printf("Error loading image %s: %s", filename, err)
	} else {
		giu.EnqueueNewTextureFromRgba(rgba, func(t *giu.Texture) {
			currentImage = t
		})
	}

	ext := filepath.Ext(filename)

	labelFile := filepath.Join(directory, "../labels/", strings.TrimSuffix(filename, ext)+".txt")

	currentRegions = parseLabelData(labelFile)
}

func selectDirectory() {
	newDirectory, err := dialog.Directory().Title("Load images").Browse()
	if err != nil {
		log.Printf("Error selecting directory: %s", err)
		return
	}
	directory = newDirectory
	updateFiles()
}

func selectFile(i int) {
	if i < 0 || i >= len(files) {
		log.Printf("Invalid file index %d", i)
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

	match, err := filepath.Glob(filepath.Join(directory, "*"))
	if err != nil {
		log.Printf("Error listing files in %s: %s", directory, err)
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

	labelsFile := filepath.Join(directory, "../labels.txt")
	file, err := os.Open(labelsFile)
	if err != nil {
		log.Printf("Error opening labels file %s: %s", labelsFile, err)
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
	for i, region := range currentRegions {
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
	for i, region := range currentRegions {
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
	for i, region := range currentRegions {
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
					giu.Label(directory),
				),
				giu.Labelf("Current file: %s", file),
				giu.Labelf("Regions: %s", regionSummary),
				giu.Labelf("Drawing label: %d %s (Press 1-9 to select new type)\n", drawingIndex, labelName(drawingIndex)),
				giu.Custom(func() {
					canvas := giu.GetCanvas()
					pos := giu.GetCursorScreenPos()
					wpos := giu.GetCursorPos()

					// Draw the image as big as the remaining space
					imageWidth := windowWidth - wpos.X
					imageHeight := windowHeight - wpos.Y

					if currentImage != nil {
						max := pos.Add(image.Point{X: imageWidth, Y: imageHeight})
						canvas.AddImage(currentImage, pos, max)
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
							if !newRegion.IsValid() {
								log.Printf("Invalid region: %#v", newRegion)
							} else {
								currentRegions = append(currentRegions, newRegion)
								log.Printf("Added new region %#v", newRegion)
								saveRegions()
							}
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
						// Find the region that was clicked
						click := giu.GetMousePos().Sub(pos)
						index := getClosestRegion(click, imageWidth, imageHeight)
						if index >= 0 {
							currentRegions = append(currentRegions[:index], currentRegions[index+1:]...)
							log.Printf("Removed region %d: %#v", index, currentRegions)
							saveRegions()
						}
					}

					for _, region := range currentRegions {
						w := int(float32(region.width) * float32(imageWidth))
						h := int(float32(region.height) * float32(imageHeight))
						x := int(float32(region.xMid)*float32(imageWidth)) - w/2
						y := int(float32(region.yMid)*float32(imageHeight)) - h/2
						color := region.Color()
						canvas.AddRect(pos.Add(image.Point{X: x, Y: y}), pos.Add(image.Point{X: x + w, Y: y + h}), color, 0, 0, 2)

						if !giu.IsMouseDown(giu.MouseButtonLeft) {
							mouse := giu.GetMousePos()
							relMouse := mouse.Sub(pos)
							if relMouse.X >= x && relMouse.X <= x+w && relMouse.Y >= y && relMouse.Y <= y+h {
								canvas.AddText(mouse.Add(image.Point{0, -20}), color, fmt.Sprintf("%s - %d", labelName(region.index), region.index))
							}
						}
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
	for i := 0; i < 9; i++ {
		if giu.IsKeyPressed(giu.Key(int(giu.Key0) + i)) {
			drawingIndex = i
		}
	}
}

func main() {
	flag.StringVar(&directory, "directory", "", "Directory to load images from")
	flag.Parse()

	wnd = giu.NewMasterWindow("Fast Mark Image Tagging", 1024, 768, 0)

	updateFiles()
	wnd.Run(loop)
}
