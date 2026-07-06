package main

import (
	"fmt"
	"image"
	"slices"

	"github.com/guigui-gui/guigui"
	"github.com/guigui-gui/guigui/basicwidget"
)

// editorPane is the right-hand column: toolbar, info texts, the region
// editor and the metadata footer. It is hosted inside a basicwidget.Panel so
// it can grow taller than the window and scroll, since the editor always
// renders the image at the full pane width.
type editorPane struct {
	guigui.DefaultWidget

	model *appModel

	changeDirButton      basicwidget.Button
	contrastCheckbox     basicwidget.Checkbox
	contrastLabel        basicwidget.Text
	backendText          basicwidget.Text
	currentFileText      clickableText
	regionsText          basicwidget.Text
	drawingLabelText     basicwidget.Text
	editor               regionEditor
	helpText             basicwidget.Text
	summaryText          basicwidget.Text
	categoryText         basicwidget.Text
	updateMetadataButton basicwidget.Button

	colItems       []guigui.LinearLayoutItem
	toolbarItems   []guigui.LinearLayoutItem
	buttonRowItems []guigui.LinearLayoutItem
}

func (p *editorPane) SetModel(m *appModel) {
	p.model = m
}

func (p *editorPane) Build(context *guigui.Context, adder *guigui.ChildAdder) error {
	adder.AddWidget(&p.changeDirButton)
	adder.AddWidget(&p.contrastCheckbox)
	adder.AddWidget(&p.contrastLabel)
	adder.AddWidget(&p.backendText)
	adder.AddWidget(&p.currentFileText)
	adder.AddWidget(&p.regionsText)
	adder.AddWidget(&p.drawingLabelText)
	adder.AddWidget(&p.editor)
	adder.AddWidget(&p.helpText)
	adder.AddWidget(&p.summaryText)
	adder.AddWidget(&p.categoryText)
	adder.AddWidget(&p.updateMetadataButton)

	m := p.model
	if m == nil {
		return nil
	}

	p.changeDirButton.SetText("Change Directory")
	p.changeDirButton.OnDown(func(context *guigui.Context) {
		m.selectDirectory()
	})

	p.contrastCheckbox.SetValue(m.autoContrast)
	p.contrastCheckbox.OnValueChanged(func(context *guigui.Context, value bool) {
		m.autoContrast = value
		m.regenerateDisplayImage()
	})
	p.contrastLabel.SetValue("Auto contrast")
	p.contrastLabel.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	if m.backend != nil {
		p.backendText.SetValue(m.backend.Describe())
	}
	p.backendText.SetVerticalAlign(basicwidget.VerticalAlignMiddle)

	var file string
	if m.selectedIndex >= 0 && m.selectedIndex < len(m.files) {
		file = m.files[m.selectedIndex]
	}
	p.currentFileText.SetValue(fmt.Sprintf("Current file: %s", file))
	p.currentFileText.OnClick(func(context *guigui.Context) {
		if file != "" {
			copyToClipboard(file)
		}
	})

	var regionSummary string
	for i, region := range m.currentRegions.Regions {
		if i > 0 {
			regionSummary += ", "
		}
		regionSummary += fmt.Sprintf("%d: %s mid=%.3f,%.3f size=%.3fx%.3f", i, m.labelName(region.index), region.xMid, region.yMid, region.width, region.height)
	}
	p.regionsText.SetValue(fmt.Sprintf("Regions: %s", regionSummary))

	p.drawingLabelText.SetValue(fmt.Sprintf("Drawing label: %d %s (Press 1-9 to select new type)", m.drawingIndex, m.labelName(m.drawingIndex)))
	p.drawingLabelText.SetColor(RegionIndexColor(m.drawingIndex))

	p.editor.SetModel(m)

	p.helpText.SetValue("Press n to find next unlabeled image")

	meta := m.metadataSnapshot()
	p.summaryText.SetValue(meta.Summary())
	p.categoryText.SetValue(m.categorySummary(meta))
	p.categoryText.SetMultiline(true)

	p.updateMetadataButton.SetText("Update Metadata")
	p.updateMetadataButton.OnDown(func(context *guigui.Context) {
		m.startMetadataScan()
	})

	return nil
}

func (p *editorPane) layout(context *guigui.Context) guigui.LinearLayout {
	u := basicwidget.UnitSize(context)

	p.toolbarItems = slices.Delete(p.toolbarItems, 0, len(p.toolbarItems))
	p.toolbarItems = append(p.toolbarItems,
		guigui.LinearLayoutItem{Widget: &p.changeDirButton},
		guigui.LinearLayoutItem{Widget: &p.contrastCheckbox, Size: guigui.FixedSize(u)},
		guigui.LinearLayoutItem{Widget: &p.contrastLabel},
		guigui.LinearLayoutItem{Widget: &p.backendText, Size: guigui.FlexibleSize(1)},
	)
	toolbar := guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     p.toolbarItems,
		Gap:       u / 4,
	}

	p.buttonRowItems = slices.Delete(p.buttonRowItems, 0, len(p.buttonRowItems))
	p.buttonRowItems = append(p.buttonRowItems,
		guigui.LinearLayoutItem{Widget: &p.updateMetadataButton},
		guigui.LinearLayoutItem{Size: guigui.FlexibleSize(1)},
	)
	buttonRow := guigui.LinearLayout{
		Direction: guigui.LayoutDirectionHorizontal,
		Items:     p.buttonRowItems,
	}

	p.colItems = slices.Delete(p.colItems, 0, len(p.colItems))
	p.colItems = append(p.colItems,
		guigui.LinearLayoutItem{Layout: &toolbar},
		guigui.LinearLayoutItem{Widget: &p.currentFileText},
		guigui.LinearLayoutItem{Widget: &p.regionsText},
		guigui.LinearLayoutItem{Widget: &p.drawingLabelText},
		guigui.LinearLayoutItem{Widget: &p.editor},
		guigui.LinearLayoutItem{Widget: &p.helpText},
		guigui.LinearLayoutItem{Widget: &p.summaryText},
		guigui.LinearLayoutItem{Widget: &p.categoryText},
		guigui.LinearLayoutItem{Layout: &buttonRow},
	)
	return guigui.LinearLayout{
		Direction: guigui.LayoutDirectionVertical,
		Items:     p.colItems,
		Gap:       u / 4,
	}
}

func (p *editorPane) Layout(context *guigui.Context, widgetBounds *guigui.WidgetBounds, layouter *guigui.ChildLayouter) {
	p.layout(context).LayoutWidgets(context, widgetBounds.Bounds(), layouter)
}

func (p *editorPane) Measure(context *guigui.Context, constraints guigui.Constraints) image.Point {
	return p.layout(context).Measure(context, constraints)
}
