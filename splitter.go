package main

import (
	"image/color"

	"github.com/guigui-gui/guigui"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/vector"
)

// splitter is a draggable vertical divider between the file list and the
// editor pane. It reports drags to its parent, which owns the sidebar width.
type splitter struct {
	guigui.DefaultWidget

	dragging   bool
	dragStartX int

	onDragStart func()
	onDrag      func(deltaX int)
}

func (s *splitter) OnDragStart(f func()) {
	s.onDragStart = f
}

// OnDrag registers a handler called during a drag with the horizontal
// distance from the drag's start.
func (s *splitter) OnDrag(f func(deltaX int)) {
	s.onDrag = f
}

func (s *splitter) HandlePointingInput(context *guigui.Context, widgetBounds *guigui.WidgetBounds) guigui.HandleInputResult {
	cursorX, _ := ebiten.CursorPosition()

	if s.dragging {
		if !ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) {
			s.dragging = false
		} else if s.onDrag != nil {
			s.onDrag(cursorX - s.dragStartX)
		}
		return guigui.HandleInputByWidget(s)
	}

	if widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		s.dragging = true
		s.dragStartX = cursorX
		if s.onDragStart != nil {
			s.onDragStart()
		}
		return guigui.HandleInputByWidget(s)
	}

	return guigui.HandleInputResult{}
}

func (s *splitter) CursorShape(context *guigui.Context, widgetBounds *guigui.WidgetBounds) (ebiten.CursorShapeType, bool) {
	return ebiten.CursorShapeEWResize, true
}

func (s *splitter) Draw(context *guigui.Context, widgetBounds *guigui.WidgetBounds, dst *ebiten.Image) {
	b := widgetBounds.Bounds()
	x := float32(b.Min.X+b.Max.X) / 2
	width := float32(1 * context.Scale())
	vector.StrokeLine(dst, x, float32(b.Min.Y), x, float32(b.Max.Y), width, color.RGBA{0x80, 0x80, 0x80, 0xff}, false)
}
