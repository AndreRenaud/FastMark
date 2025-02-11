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

	directory string
	splitPos  = float32(200)
)

type Region struct {
	xMid   float64
	yMid   float64
	width  float64
	height float64

	index int
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
			log.Printf("Invalid line: %s", scanner.Text())
			continue
		}
		region := Region{}
		region.index, err = strconv.Atoi(columns[0])
		region.xMid, err = strconv.ParseFloat(columns[1], 64)
		region.yMid, err = strconv.ParseFloat(columns[2], 64)
		region.width, err = strconv.ParseFloat(columns[3], 64)
		region.height, err = strconv.ParseFloat(columns[4], 64)
		retval = append(retval, region)
	}

	return retval
}

func (r Region) Color() color.RGBA {
	switch r.index {
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
		return color.RGBA{128, 128, 128, 255}
	}
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
			if strings.HasSuffix(m, ".png") || strings.HasSuffix(m, ".jpg") {
				files = append(files, filepath.Base(m))
			}
		}
	}
	fileRows = make([]*giu.TableRowWidget, len(files))
	fileLabels = make([]*giu.SelectableWidget, len(files))
	for i, file := range files {
		fileLabels[i] = giu.Selectable(file)
		fileLabels[i].OnClick(func() {
			selectFile(i)
		})
		fileRows[i] = giu.TableRow(fileLabels[i])
	}
	selectFile(0)
}

func loop() {
	window := giu.SingleWindow()
	var file string
	if selectedIndex >= 0 && selectedIndex < len(files) {
		file = files[selectedIndex]
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
				giu.Label(file),
				giu.Custom(func() {
					canvas := giu.GetCanvas()
					pos := giu.GetCursorScreenPos()
					imageWidth := 800
					imageHeight := 600

					if currentImage != nil {
						max := pos.Add(image.Point{X: imageWidth, Y: imageHeight})
						canvas.AddImage(currentImage, pos, max)
					}
					// Check if the user has stopped drawing
					if drawingRect {
						end := giu.GetMousePos().Sub(pos)
						if !giu.IsMouseDown(giu.MouseButtonLeft) {
							newRegion := Region{
								xMid:   float64(drawingStart.X+end.X) / 2 / float64(imageWidth),
								yMid:   float64(drawingStart.Y+end.Y) / 2 / float64(imageHeight),
								width:  float64(end.X-drawingStart.X) / float64(imageWidth),
								height: float64(end.Y-drawingStart.Y) / float64(imageHeight),
								index:  1, // TODO: Make this configurable
							}
							currentRegions = append(currentRegions, newRegion)
							log.Printf("Added new region %#v", newRegion)
							drawingRect = false
						}
						canvas.AddRect(pos.Add(drawingStart), pos.Add(end), color.RGBA{255, 0, 0, 255}, 0, 0, 2)
					} else if giu.IsMouseDown(giu.MouseButtonLeft) {
						// Get the current screen position of the cursor
						scr := giu.GetMousePos().Sub(pos)
						drawingRect = true
						drawingStart = image.Point{X: int(scr.X), Y: int(scr.Y)}
					}

					for _, region := range currentRegions {
						w := int(float32(region.width) * float32(imageWidth))
						h := int(float32(region.height) * float32(imageHeight))
						x := int(float32(region.xMid)*float32(imageWidth)) - w/2
						y := int(float32(region.yMid)*float32(imageHeight)) - h/2
						color := region.Color()
						canvas.AddRect(pos.Add(image.Point{X: x, Y: y}), pos.Add(image.Point{X: x + w, Y: y + h}), color, 0, 0, 2)
					}
				}),
				giu.Event().OnClick(giu.MouseButtonLeft, func() {
					log.Printf("Click")
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

}

func main() {
	flag.StringVar(&directory, "directory", "", "Directory to load images from")
	flag.Parse()

	wnd := giu.NewMasterWindow("Fast Mark Image Tagging", 1024, 768, 0)

	updateFiles()
	wnd.Run(loop)
}
