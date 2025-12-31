# FFwd UI Icon Assets

This directory contains the application icons for all platforms.

## Logo Design

The FFwd UI logo features:
- **Fast Forward Icon**: Two white triangular arrows pointing right (symbolizing video fast-forward)
- **"UI" Text**: Bold white text below the icon
- **Purple Gradient Background**: Modern gradient from #667eea to #764ba2
- **Rounded Corners**: 80px border radius for a modern app appearance

## Files

### Source
- **`logo.svg`** - Vector source file (512x512)
  - Scalable to any size without quality loss
  - Use this for creating new sizes or variants

### Flatpak/Linux/Wails
- **`appicon.png`** - 1024x1024 PNG for Flatpak package and Wails Linux builds
  - Referenced in `io.github.thesle.FFwdUI.yml` for Flatpak
  - Referenced in `main.go` for Wails Linux icon
  - Used by desktop environments and app launchers
  - **Note**: Wails requires exactly 1024x1024 pixels for proper rendering

### Windows
- **`windows/icon.ico`** - Multi-resolution .ico file
  - Contains sizes: 256x256, 128x128, 64x64, 48x48, 32x32, 16x16
  - Used for executable and taskbar icons

### macOS
- **`darwin/icons.iconset/`** - macOS icon set directory
  - Standard sizes: 16x16, 32x32, 64x64, 128x128, 256x256, 512x512
  - Retina @2x versions: 32x32@2x, 64x64@2x, 256x256@2x, 512x512@2x, 1024x1024@2x
  - Automatically compiled into .icns during Wails build

### Frontend
- **`../frontend/src/assets/images/logo-universal.png`** - 256x256 PNG
  - Used within the application UI if needed

## Regenerating Icons

If you need to modify the logo, edit `logo.svg` and regenerate:

```bash
# Wails/Flatpak icon (1024x1024 - required for Wails)
inkscape --export-type=png --export-filename=build/appicon.png --export-width=1024 --export-height=1024 build/logo.svg

# Windows .ico (multi-size)
convert build/logo.svg -background none -define icon:auto-resize=256,128,64,48,32,16 build/windows/icon.ico

# macOS iconset (run from project root)
mkdir -p build/darwin/icons.iconset
for size in 16 32 64 128 256 512; do 
  inkscape --export-type=png --export-filename=build/darwin/icons.iconset/icon_${size}x${size}.png --export-width=$size --export-height=$size build/logo.svg
done

# macOS @2x retina versions
for size in 16 32 128 256 512; do 
  double=$((size * 2))
  inkscape --export-type=png --export-filename=build/darwin/icons.iconset/icon_${size}x${size}@2x.png --export-width=$double --export-height=$double build/logo.svg
done

# Frontend logo
inkscape --export-type=png --export-filename=frontend/src/assets/images/logo-universal.png --export-width=256 --export-height=256 build/logo.svg
```

## Design Notes

- **Color Scheme**: Purple gradient (#667eea → #764ba2) for a modern, professional look
- **Icon Style**: Minimalist with clear symbolism (fast forward = video operations)
- **Platform Guidelines**: Follows standard icon sizing for Linux, Windows, and macOS
- **Accessibility**: High contrast white icons on colored background for visibility

## Tools Required

- **Inkscape** - For SVG to PNG conversion (recommended)
- **ImageMagick** - Alternative for .ico creation
- **rsvg-convert** - Alternative lightweight SVG converter

Install on Ubuntu:
```bash
sudo apt install inkscape imagemagick librsvg2-bin
```
