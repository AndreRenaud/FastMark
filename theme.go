package main

import (
	"github.com/AllenDang/cimgui-go/imgui"
	"github.com/AllenDang/giu"
)

func FastMarkTheme() *giu.StyleSetter {
	return giu.Style().
		SetStyleFloat(giu.StyleVarAlpha, 1.0).
		SetStyleFloat(giu.StyleVarFrameRounding, 3.0).
		SetStyleFloat(giu.StyleVarScrollbarSize, 16.0).
		SetStyleFloat(giu.StyleVarScrollbarRounding, 2.0).
		SetStyleFloat(giu.StyleVarFrameBorderSize, 1.0).
		SetFontSize(15).
		SetColorVec4(giu.StyleColorText, imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 1.00}).
		SetColorVec4(giu.StyleColorTextDisabled, imgui.Vec4{X: 0.60, Y: 0.60, Z: 0.60, W: 1.00}).
		SetColorVec4(giu.StyleColorWindowBg, imgui.Vec4{X: 0.94, Y: 0.94, Z: 0.94, W: 0.94}).
		SetColorVec4(giu.StyleColorChildBg, imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 0.00}).
		SetColorVec4(giu.StyleColorPopupBg, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.94}).
		SetColorVec4(giu.StyleColorBorder, imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 0.39}).
		SetColorVec4(giu.StyleColorBorderShadow, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.10}).
		SetColorVec4(giu.StyleColorFrameBg, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.94}).
		SetColorVec4(giu.StyleColorFrameBgHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.40}).
		SetColorVec4(giu.StyleColorFrameBgActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.67}).
		SetColorVec4(giu.StyleColorTitleBg, imgui.Vec4{X: 0.96, Y: 0.96, Z: 0.96, W: 1.00}).
		SetColorVec4(giu.StyleColorTitleBgCollapsed, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.51}).
		SetColorVec4(giu.StyleColorTitleBgActive, imgui.Vec4{X: 0.82, Y: 0.82, Z: 0.82, W: 1.00}).
		SetColorVec4(giu.StyleColorMenuBarBg, imgui.Vec4{X: 0.86, Y: 0.86, Z: 0.86, W: 1.00}).
		SetColorVec4(giu.StyleColorScrollbarBg, imgui.Vec4{X: 0.98, Y: 0.98, Z: 0.98, W: 0.53}).
		SetColorVec4(giu.StyleColorScrollbarGrab, imgui.Vec4{X: 0.69, Y: 0.69, Z: 0.69, W: 1.00}).
		SetColorVec4(giu.StyleColorScrollbarGrabHovered, imgui.Vec4{X: 0.59, Y: 0.59, Z: 0.59, W: 1.00}).
		SetColorVec4(giu.StyleColorScrollbarGrabActive, imgui.Vec4{X: 0.49, Y: 0.49, Z: 0.49, W: 1.00}).
		//SetColorVec4(giu.StyleColorComboBg, imgui.Vec4{X: 0.86, Y: 0.86, Z: 0.86, W: 0.99}).
		SetColorVec4(giu.StyleColorCheckMark, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorSliderGrab, imgui.Vec4{X: 0.24, Y: 0.52, Z: 0.88, W: 1.00}).
		SetColorVec4(giu.StyleColorSliderGrabActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorButton, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.40}).
		SetColorVec4(giu.StyleColorButtonHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorButtonActive, imgui.Vec4{X: 0.06, Y: 0.53, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorHeader, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.31}).
		SetColorVec4(giu.StyleColorHeaderHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.80}).
		SetColorVec4(giu.StyleColorHeaderActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		//SetColorVec4(giu.StyleColorColumn, imgui.Vec4{X: 0.39, Y: 0.39, Z: 0.39, W: 1.00}).
		//SetColorVec4(giu.StyleColorColumnHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.78}).
		//SetColorVec4(giu.StyleColorColumnActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorResizeGrip, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.50}).
		SetColorVec4(giu.StyleColorResizeGripHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.67}).
		SetColorVec4(giu.StyleColorResizeGripActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.95}).
		//SetColorVec4(giu.StyleColorCloseButton, imgui.Vec4{X: 0.59, Y: 0.59, Z: 0.59, W: 0.50}).
		//SetColorVec4(giu.StyleColorCloseButtonHovered, imgui.Vec4{X: 0.98, Y: 0.39, Z: 0.36, W: 1.00}).
		//SetColorVec4(giu.StyleColorCloseButtonActive, imgui.Vec4{X: 0.98, Y: 0.39, Z: 0.36, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotLines, imgui.Vec4{X: 0.39, Y: 0.39, Z: 0.39, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotLinesHovered, imgui.Vec4{X: 1.00, Y: 0.43, Z: 0.35, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotHistogram, imgui.Vec4{X: 0.90, Y: 0.70, Z: 0.00, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotHistogramHovered, imgui.Vec4{X: 1.00, Y: 0.60, Z: 0.00, W: 1.00}).
		SetColorVec4(giu.StyleColorTextSelectedBg, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.35}).
		SetColorVec4(giu.StyleColorModalWindowDimBg, imgui.Vec4{X: 0.20, Y: 0.20, Z: 0.20, W: 0.35})
}
