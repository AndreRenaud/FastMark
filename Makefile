APP_NAME := Fastmark
BINARY := fastmark
OUTPUT := output
BUNDLE := $(OUTPUT)/$(APP_NAME).app
IDENTIFIER := com.github.andrerenaud.fastmark

GIT_VERSION := $(shell git describe --tags --always --dirty 2>/dev/null || echo 0.0.0)
VERSION := $(patsubst v%,%,$(GIT_VERSION))

GO_SOURCES := $(wildcard *.go) $(wildcard storage/*.go) go.mod go.sum

.PHONY: all app install clean

all: app

app: $(BUNDLE)

$(OUTPUT)/$(BINARY): $(GO_SOURCES) icon-128.png
	mkdir -p $(OUTPUT)
	go build -o $@ .

# macOS icon: scale the source PNG into an iconset, then compile to .icns
$(OUTPUT)/$(BINARY).icns: icon-128.png
	rm -rf $(OUTPUT)/$(BINARY).iconset
	mkdir -p $(OUTPUT)/$(BINARY).iconset
	sips -z 16 16 $< --out $(OUTPUT)/$(BINARY).iconset/icon_16x16.png >/dev/null
	sips -z 32 32 $< --out $(OUTPUT)/$(BINARY).iconset/icon_16x16@2x.png >/dev/null
	sips -z 32 32 $< --out $(OUTPUT)/$(BINARY).iconset/icon_32x32.png >/dev/null
	sips -z 64 64 $< --out $(OUTPUT)/$(BINARY).iconset/icon_32x32@2x.png >/dev/null
	cp $< $(OUTPUT)/$(BINARY).iconset/icon_128x128.png
	iconutil -c icns -o $@ $(OUTPUT)/$(BINARY).iconset
	rm -rf $(OUTPUT)/$(BINARY).iconset

$(BUNDLE): $(OUTPUT)/$(BINARY) $(OUTPUT)/$(BINARY).icns
	rm -rf $(BUNDLE)
	mkdir -p $(BUNDLE)/Contents/MacOS $(BUNDLE)/Contents/Resources
	cp $(OUTPUT)/$(BINARY) $(BUNDLE)/Contents/MacOS/$(BINARY)
	cp $(OUTPUT)/$(BINARY).icns $(BUNDLE)/Contents/Resources/$(BINARY).icns
	echo "$$INFO_PLIST" > $(BUNDLE)/Contents/Info.plist
	plutil -lint $(BUNDLE)/Contents/Info.plist
	codesign --force --deep --sign - $(BUNDLE)
	touch $(BUNDLE)
	@echo "Built $(BUNDLE) (version $(VERSION)) - copy it into /Applications to install"

install: $(BUNDLE)
	rm -rf /Applications/$(APP_NAME).app
	cp -R $(BUNDLE) /Applications/

clean:
	rm -rf $(OUTPUT)

define INFO_PLIST
<?xml version="1.0" encoding="UTF-8"?>
<!DOCTYPE plist PUBLIC "-//Apple//DTD PLIST 1.0//EN" "http://www.apple.com/DTDs/PropertyList-1.0.dtd">
<plist version="1.0">
<dict>
	<key>CFBundleName</key>
	<string>$(APP_NAME)</string>
	<key>CFBundleDisplayName</key>
	<string>$(APP_NAME)</string>
	<key>CFBundleExecutable</key>
	<string>$(BINARY)</string>
	<key>CFBundleIdentifier</key>
	<string>$(IDENTIFIER)</string>
	<key>CFBundleIconFile</key>
	<string>$(BINARY)</string>
	<key>CFBundlePackageType</key>
	<string>APPL</string>
	<key>CFBundleShortVersionString</key>
	<string>$(VERSION)</string>
	<key>CFBundleVersion</key>
	<string>$(VERSION)</string>
	<key>LSMinimumSystemVersion</key>
	<string>11.0</string>
	<key>NSHighResolutionCapable</key>
	<true/>
</dict>
</plist>
endef
export INFO_PLIST
