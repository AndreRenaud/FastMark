package main

import (
	"fmt"
	"image"
	"image/color"
	"log"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// regionEditor displays the current image aspect-fit and lets the user draw,
// delete and re-tag regions with the mouse.
type regionEditor struct {
	guigui.DefaultWidget

	model *appModel

	labelTexts guigui.WidgetSlice[*basicwidget.Text]

	drawingRect  bool
	drawingStart image.Point // relative to the displayed image's origin

	// lastAspect is the height/width ratio of the most recent image, kept so
	// the layout stays stable while the next image is decoding.
	lastAspect float64
}

func (e *regionEditor) SetModel(m *appModel) {
	e.model = m
}

func (e *regionEditor) cancelDrawing() {
	if e.drawingRect {
		e.drawingRect = false
		guigui.RequestRedraw(e)
	}
}

// imageRect returns the rectangle the image is displayed in: the full width
// of the widget, with the height following from the image's aspect ratio.
func (e *regionEditor) imageRect(bounds image.Rectangle) image.Rectangle {
	if e.model == nil || e.model.displayImage == nil {
		return image.Rectangle{}
	}
	size := e.model.displayImage.Bounds().Size()
	if size.X <= 0 || size.Y <= 0 {
		return image.Rectangle{}
	}
	w := bounds.Dx()
	h := size.Y * w / size.X
	return image.Rectangle{Min: bounds.Min, Max: bounds.Min.Add(image.Pt(w, h))}
}

// Measure sizes the editor to the full available width, with the height
// derived from the current image's aspect ratio.
func (e *regionEditor) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	aspect := e.lastAspect
	if aspect == 0 {
		aspect = 0.75 // default before the first image has loaded
	}
	w, ok := constraints.FixedWidth()
	if !ok {
		w = basicwidget.UnitSize(context) * 16
	}
	return image.Pt(w, int(float64(w)*aspect))
}

// regionRect converts a normalized region to display coordinates within ir.
func regionRect(region Region, ir image.Rectangle) image.Rectangle {
	w := int(region.width * float64(ir.Dx()))
	h := int(region.height * float64(ir.Dy()))
	x := int(region.xMid*float64(ir.Dx())) - w/2
	y := int(region.yMid*float64(ir.Dy())) - h/2
	return image.Rect(x, y, x+w, y+h).Add(ir.Min)
}

func (e *regionEditor) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	m := e.model
	n := 0
	if m != nil && m.displayImage != nil {
		n = len(m.currentRegions.Regions)
		size := m.displayImage.Bounds().Size()
		if size.X > 0 {
			e.lastAspect = float64(size.Y) / float64(size.X)
		}
	}
	e.labelTexts.SetLen(n)
	for i := range n {
		adder.AddWidget(e.labelTexts.At(i))
	}
	for i := range n {
		region := m.currentRegions.Regions[i]
		t := e.labelTexts.At(i)
		t.SetValue(fmt.Sprintf("%s - %d", m.labelName(region.index), region.index))
		t.SetColor(region.Color())
	}
	return nil
}

func (e *regionEditor) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	m := e.model
	if m == nil {
		return
	}
	ir := e.imageRect(widgetBounds.Bounds())
	lh := basicwidget.LineHeight(context)
	u := basicwidget.UnitSize(context)
	for i := range e.labelTexts.Len() {
		if i >= len(m.currentRegions.Regions) {
			break
		}
		rr := regionRect(m.currentRegions.Regions[i], ir)
		pos := image.Pt(rr.Min.X+rr.Dx()/2, rr.Min.Y-lh)
		layouter.LayoutWidget(e.labelTexts.At(i), image.Rectangle{Min: pos, Max: pos.Add(image.Pt(u*8, lh))})
	}
}

func (e *regionEditor) HandlePointingInput(context *guigui.Context, widgetBounds *guigui.WidgetBounds) guigui.HandleInputResult {
	m := e.model
	if m == nil || m.displayImage == nil {
		return guigui.HandleInputResult{}
	}
	ir := e.imageRect(widgetBounds.Bounds())
	cursor := image.Pt(ebiten.CursorPosition())

	if e.drawingRect {
		// Keep the widget repainting so the in-progress rectangle tracks the cursor.
		guigui.RequestRedraw(e)
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			end := cursor.Sub(ir.Min)
			// Create a new well formed region clamped within the image
			newRect := image.Rect(e.drawingStart.X, e.drawingStart.Y, end.X, end.Y)
			newRect = newRect.Intersect(image.Rect(0, 0, ir.Dx(), ir.Dy())).Canon()
			log.Printf("New rect: %v", newRect)
			newRegion := Region{
				xMid:   (float64(newRect.Dx())/2 + float64(newRect.Min.X)) / float64(ir.Dx()),
				yMid:   (float64(newRect.Dy())/2 + float64(newRect.Min.Y)) / float64(ir.Dy()),
				width:  float64(newRect.Dx()) / float64(ir.Dx()),
				height: float64(newRect.Dy()) / float64(ir.Dy()),
				index:  m.drawingIndex,
			}
			m.currentRegions.AddRegion(newRegion)
			e.drawingRect = false
		}
		return guigui.HandleInputByWidget(e)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && cursor.In(ir) {
		e.drawingRect = true
		e.drawingStart = cursor.Sub(ir.Min)
		return guigui.HandleInputByWidget(e)
	}

	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonRight) && cursor.In(ir) {
		changeRegion := -1
		// If we're pressing a number key, change the region type, otherwise delete it
		for d := range 10 {
			if ebiten.IsKeyPressed(ebiten.KeyDigit0 + ebiten.Key(d)) {
				changeRegion = d
			}
		}
		// Find the region that was clicked
		click := cursor.Sub(ir.Min)
		index := m.getClosestRegion(click, ir.Dx(), ir.Dy())
		if index >= 0 {
			if changeRegion >= 0 {
				m.currentRegions.Regions[index].index = changeRegion
				// Do this async so we don't block the UI
				go m.currentRegions.Save()
			} else {
				m.currentRegions.Remove(index)
			}
		}
		return guigui.HandleInputByWidget(e)
	}

	return guigui.HandleInputResult{}
}

func (e *regionEditor) Draw(context *guigui.Context, widgetBounds *guigui.WidgetBounds, dst *ebiten.Image) {
	m := e.model
	if m == nil || m.displayImage == nil {
		return
	}
	ir := e.imageRect(widgetBounds.Bounds())
	size := m.displayImage.Bounds().Size()

	op := &ebiten.DrawImageOptions{}
	op.GeoM.Scale(float64(ir.Dx())/float64(size.X), float64(ir.Dy())/float64(size.Y))
	op.GeoM.Translate(float64(ir.Min.X), float64(ir.Min.Y))
	op.Filter = ebiten.FilterLinear
	dst.DrawImage(m.displayImage, op)

	for _, region := range m.currentRegions.Regions {
		strokeRect(dst, regionRect(region, ir), region.Color())
	}

	if e.drawingRect {
		cursor := image.Pt(ebiten.CursorPosition())
		start := ir.Min.Add(e.drawingStart)
		strokeRect(dst, image.Rect(start.X, start.Y, cursor.X, cursor.Y), color.RGBA{255, 0, 0, 255})
	}
}

func strokeRect(dst *ebiten.Image, r image.Rectangle, clr color.Color) {
	vector.StrokeRect(dst, float32(r.Min.X), float32(r.Min.Y), float32(r.Dx()), float32(r.Dy()), 2, clr, false)
}
