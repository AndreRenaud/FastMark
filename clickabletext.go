package main

import (
	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

// clickableText is a basicwidget.Text that invokes a callback when clicked,
// showing a pointer cursor while hovered.
type clickableText struct {
	basicwidget.Text

	onClick func(context *guigui.Context)
}

func (c *clickableText) OnClick(f func(context *guigui.Context)) {
	c.onClick = f
}

func (c *clickableText) HandlePointingInput(context *guigui.Context, widgetBounds *guigui.WidgetBounds) guigui.HandleInputResult {
	if c.onClick != nil && widgetBounds.IsHitAtCursor() && inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		c.onClick(context)
		return guigui.HandleInputByWidget(c)
	}
	return c.Text.HandlePointingInput(context, widgetBounds)
}

func (c *clickableText) CursorShape(context *guigui.Context, widgetBounds *guigui.WidgetBounds) (ebiten.CursorShapeType, bool) {
	if c.onClick != nil && widgetBounds.IsHitAtCursor() {
		return ebiten.CursorShapePointer, true
	}
	return c.Text.CursorShape(context, widgetBounds)
}
