package main

import (
	_ "embed"

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

//go:embed micross.ttf
var micross []byte

func Win98Theme() *giu.StyleSetter {
	ctx := giu.Context
	fi := ctx.FontAtlas.AddFontFromBytes("micross.ttf", micross)
	//hspacing := 8
	//vspacing := 6
	/*
			ImGuiIO& io = ImGui::GetIO();
		io.Fonts->Clear();
		ImFont* font = io.Fonts->AddFontFromFileTTF("C:\\Windows\\Fonts\\segoeui.ttf", 18.0f);
		if (font != NULL) {
			font->DisplayOffset.y -= 1;
			io.FontDefault = font;
		} else {
			io.Fonts->AddFontDefault();
		}
		io.Fonts->Build();
	*/

	white := imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 1.00}
	transparent := imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 0.00}
	dark := imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 0.20}
	darker := imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 0.50}
	background := imgui.Vec4{X: 0.95, Y: 0.95, Z: 0.95, W: 1.00}
	text := imgui.Vec4{X: 0.10, Y: 0.10, Z: 0.10, W: 1.00}
	border := imgui.Vec4{X: 0.60, Y: 0.60, Z: 0.60, W: 1.00}
	grab := imgui.Vec4{X: 0.69, Y: 0.69, Z: 0.69, W: 1.00}
	header := imgui.Vec4{X: 0.86, Y: 0.86, Z: 0.86, W: 1.00}
	active := imgui.Vec4{X: 0.00, Y: 0.47, Z: 0.84, W: 1.00}
	hover := imgui.Vec4{X: 0.00, Y: 0.47, Z: 0.84, W: 0.20}
	return giu.Style().
		SetFont(fi).
		//DisplaySafeAreaPadding(0, 0).
		SetStyleFloat(giu.StyleVarWindowPadding, 4). //WindowPadding(hspacing/2, vspacing)
		SetStyleFloat(giu.StyleVarFramePadding, 6).
		SetStyleFloat(giu.StyleVarItemSpacing, 6).
		SetStyleFloat(giu.StyleVarItemInnerSpacing, 6).
		SetStyleFloat(giu.StyleVarIndentSpacing, 20.0).
		SetStyleFloat(giu.StyleVarWindowRounding, 0.0).
		SetStyleFloat(giu.StyleVarFrameRounding, 0.0).
		SetStyleFloat(giu.StyleVarWindowBorderSize, 0.0).
		SetStyleFloat(giu.StyleVarFrameBorderSize, 1.0).
		SetStyleFloat(giu.StyleVarPopupBorderSize, 1.0).
		SetStyleFloat(giu.StyleVarScrollbarSize, 20.0).
		SetStyleFloat(giu.StyleVarScrollbarRounding, 0.0).
		SetStyleFloat(giu.StyleVarGrabMinSize, 5.0).
		SetStyleFloat(giu.StyleVarGrabRounding, 0.0).
		SetColor(giu.StyleColorText, giu.Vec4ToRGBA(text)).
		SetColor(giu.StyleColorWindowBg, giu.Vec4ToRGBA(background)).
		SetColor(giu.StyleColorChildBg, giu.Vec4ToRGBA(background)).
		SetColor(giu.StyleColorPopupBg, giu.Vec4ToRGBA(white)).
		SetColor(giu.StyleColorBorder, giu.Vec4ToRGBA(border)).
		SetColor(giu.StyleColorBorderShadow, giu.Vec4ToRGBA(transparent)).
		SetColor(giu.StyleColorButton, giu.Vec4ToRGBA(header)).
		SetColor(giu.StyleColorButtonHovered, giu.Vec4ToRGBA(hover)).
		SetColor(giu.StyleColorButtonActive, giu.Vec4ToRGBA(active)).
		SetColor(giu.StyleColorFrameBg, giu.Vec4ToRGBA(white)).
		SetColor(giu.StyleColorFrameBgHovered, giu.Vec4ToRGBA(hover)).
		SetColor(giu.StyleColorFrameBgActive, giu.Vec4ToRGBA(active)).
		SetColor(giu.StyleColorMenuBarBg, giu.Vec4ToRGBA(header)).
		SetColor(giu.StyleColorHeader, giu.Vec4ToRGBA(header)).
		SetColor(giu.StyleColorHeaderHovered, giu.Vec4ToRGBA(hover)).
		SetColor(giu.StyleColorHeaderActive, giu.Vec4ToRGBA(active)).
		SetColor(giu.StyleColorCheckMark, giu.Vec4ToRGBA(text)).
		SetColor(giu.StyleColorSliderGrab, giu.Vec4ToRGBA(grab)).
		SetColor(giu.StyleColorSliderGrabActive, giu.Vec4ToRGBA(active)).
		//SetColor(giu.StyleColorCloseButton, giu.Vec4ToRGBA(transparent)).
		//SetColor(giu.StyleColorCloseButtonHovered, giu.Vec4ToRGBA(transparent)).
		//SetColor(giu.StyleColorCloseButtonActive, giu.Vec4ToRGBA(transparent)).
		SetColor(giu.StyleColorScrollbarBg, giu.Vec4ToRGBA(header)).
		SetColor(giu.StyleColorScrollbarGrab, giu.Vec4ToRGBA(grab)).
		SetColor(giu.StyleColorScrollbarGrabHovered, giu.Vec4ToRGBA(dark)).
		SetColor(giu.StyleColorScrollbarGrabActive, giu.Vec4ToRGBA(darker))

	/*
		ImGuiStyle* style = &ImGui::GetStyle();
		int hspacing = 8;
		int vspacing = 6;
		style->DisplaySafeAreaPadding = ImVec2(0, 0);
		style->WindowPadding = ImVec2(hspacing/2, vspacing);
		style->FramePadding = ImVec2(hspacing, vspacing);
		style->ItemSpacing = ImVec2(hspacing, vspacing);
		style->ItemInnerSpacing = ImVec2(hspacing, vspacing);
		style->IndentSpacing = 20.0f;

	*/
}

// LightTheme generates a default GIU theme's StyleSetter.
func LightTheme() *giu.StyleSetter {
	return giu.Style().
		SetStyleFloat(giu.StyleVarWindowRounding, 2).
		SetStyleFloat(giu.StyleVarFrameRounding, 4).
		SetStyleFloat(giu.StyleVarGrabRounding, 4).
		SetStyleFloat(giu.StyleVarFrameBorderSize, 1).
		SetColorVec4(giu.StyleColorText, imgui.Vec4{X: 0.10, Y: 0.10, Z: 0.10, W: 1.00}).
		SetColorVec4(giu.StyleColorTextDisabled, imgui.Vec4{X: 0.60, Y: 0.60, Z: 0.60, W: 1.00}).
		SetColorVec4(giu.StyleColorWindowBg, imgui.Vec4{X: 0.94, Y: 0.94, Z: 0.94, W: 1.00}).
		SetColorVec4(giu.StyleColorChildBg, imgui.Vec4{X: 0.97, Y: 0.97, Z: 0.97, W: 1.00}).
		SetColorVec4(giu.StyleColorPopupBg, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.98}).
		SetColorVec4(giu.StyleColorBorder, imgui.Vec4{X: 0.70, Y: 0.70, Z: 0.70, W: 1.00}).
		SetColorVec4(giu.StyleColorBorderShadow, imgui.Vec4{X: 0.00, Y: 0.00, Z: 0.00, W: 0.00}).
		SetColorVec4(giu.StyleColorFrameBg, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 1.00}).
		SetColorVec4(giu.StyleColorFrameBgHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.40}).
		SetColorVec4(giu.StyleColorFrameBgActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.67}).
		SetColorVec4(giu.StyleColorTitleBg, imgui.Vec4{X: 0.96, Y: 0.96, Z: 0.96, W: 1.00}).
		SetColorVec4(giu.StyleColorTitleBgActive, imgui.Vec4{X: 0.82, Y: 0.82, Z: 0.82, W: 1.00}).
		SetColorVec4(giu.StyleColorTitleBgCollapsed, imgui.Vec4{X: 1.00, Y: 1.00, Z: 1.00, W: 0.51}).
		SetColorVec4(giu.StyleColorMenuBarBg, imgui.Vec4{X: 0.86, Y: 0.86, Z: 0.86, W: 1.00}).
		SetColorVec4(giu.StyleColorScrollbarBg, imgui.Vec4{X: 0.98, Y: 0.98, Z: 0.98, W: 0.53}).
		SetColorVec4(giu.StyleColorScrollbarGrab, imgui.Vec4{X: 0.69, Y: 0.69, Z: 0.69, W: 1.00}).
		SetColorVec4(giu.StyleColorScrollbarGrabHovered, imgui.Vec4{X: 0.59, Y: 0.59, Z: 0.59, W: 1.00}).
		SetColorVec4(giu.StyleColorScrollbarGrabActive, imgui.Vec4{X: 0.49, Y: 0.49, Z: 0.49, W: 1.00}).
		SetColorVec4(giu.StyleColorCheckMark, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorSliderGrab, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorSliderGrabActive, imgui.Vec4{X: 0.06, Y: 0.53, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorButton, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.40}).
		SetColorVec4(giu.StyleColorButtonHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorButtonActive, imgui.Vec4{X: 0.06, Y: 0.53, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorHeader, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.31}).
		SetColorVec4(giu.StyleColorHeaderHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.80}).
		SetColorVec4(giu.StyleColorHeaderActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorSeparator, imgui.Vec4{X: 0.39, Y: 0.39, Z: 0.39, W: 0.62}).
		SetColorVec4(giu.StyleColorSeparatorHovered, imgui.Vec4{X: 0.14, Y: 0.44, Z: 0.80, W: 0.78}).
		SetColorVec4(giu.StyleColorSeparatorActive, imgui.Vec4{X: 0.14, Y: 0.44, Z: 0.80, W: 1.00}).
		SetColorVec4(giu.StyleColorResizeGrip, imgui.Vec4{X: 0.35, Y: 0.35, Z: 0.35, W: 0.17}).
		SetColorVec4(giu.StyleColorResizeGripHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.67}).
		SetColorVec4(giu.StyleColorResizeGripActive, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.95}).
		SetColorVec4(giu.StyleColorTab, imgui.Vec4{X: 0.76, Y: 0.80, Z: 0.84, W: 0.93}).
		SetColorVec4(giu.StyleColorTabHovered, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.80}).
		SetColorVec4(giu.StyleColorTabActive, imgui.Vec4{X: 0.60, Y: 0.73, Z: 0.88, W: 1.00}).
		SetColorVec4(giu.StyleColorTabUnfocused, imgui.Vec4{X: 0.92, Y: 0.93, Z: 0.94, W: 1.00}).
		SetColorVec4(giu.StyleColorTabUnfocusedActive, imgui.Vec4{X: 0.74, Y: 0.82, Z: 0.91, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotLines, imgui.Vec4{X: 0.39, Y: 0.39, Z: 0.39, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotLinesHovered, imgui.Vec4{X: 1.00, Y: 0.43, Z: 0.35, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotHistogram, imgui.Vec4{X: 0.90, Y: 0.70, Z: 0.00, W: 1.00}).
		SetColorVec4(giu.StyleColorPlotHistogramHovered, imgui.Vec4{X: 1.00, Y: 0.45, Z: 0.00, W: 1.00}).
		SetColorVec4(giu.StyleColorTextSelectedBg, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.35}).
		SetColorVec4(giu.StyleColorDragDropTarget, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.95}).
		SetColorVec4(giu.StyleColorNavWindowingHighlight, imgui.Vec4{X: 0.26, Y: 0.59, Z: 0.98, W: 0.80}).
		SetColorVec4(giu.StyleColorNavWindowingHighlight, imgui.Vec4{X: 0.70, Y: 0.70, Z: 0.70, W: 0.70}).
		SetColorVec4(giu.StyleColorTableHeaderBg, imgui.Vec4{X: 0.78, Y: 0.87, Z: 0.98, W: 1.00}).
		SetColorVec4(giu.StyleColorTableBorderStrong, imgui.Vec4{X: 0.57, Y: 0.57, Z: 0.64, W: 1.00}).
		SetColorVec4(giu.StyleColorTableBorderLight, imgui.Vec4{X: 0.68, Y: 0.68, Z: 0.74, W: 1.00})
}
