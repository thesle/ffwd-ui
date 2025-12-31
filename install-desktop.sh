#!/bin/bash
# Install FFwd UI desktop file and icon for testing

set -e

APP_NAME="ffwd-ui"
BINARY_PATH="$(pwd)/build/bin/ffwd-ui"
ICON_PATH="$(pwd)/build/appicon.png"
DESKTOP_FILE="$HOME/.local/share/applications/ffwd-ui.desktop"
ICON_DEST="$HOME/.local/share/icons/hicolor/1024x1024/apps/ffwd-ui.png"

echo "Installing FFwd UI desktop file..."

# Create directories if they don't exist
mkdir -p "$HOME/.local/share/applications"
mkdir -p "$HOME/.local/share/icons/hicolor/1024x1024/apps"

# Copy icon
cp "$ICON_PATH" "$ICON_DEST"
echo "✓ Icon installed to: $ICON_DEST"

# Create desktop file
cat > "$DESKTOP_FILE" << EOF
[Desktop Entry]
Name=FFwd UI
Comment=A minimalistic FFmpeg GUI for common video operations
Exec=$BINARY_PATH
Icon=ffwd-ui
Type=Application
Categories=AudioVideo;Video;AudioVideoEditing;
Keywords=ffmpeg;video;audio;convert;trim;extract;
StartupNotify=true
Terminal=false
EOF

chmod +x "$DESKTOP_FILE"
echo "✓ Desktop file installed to: $DESKTOP_FILE"

# Update desktop database
if command -v update-desktop-database &> /dev/null; then
    update-desktop-database "$HOME/.local/share/applications"
    echo "✓ Desktop database updated"
fi

# Update icon cache
if command -v gtk-update-icon-cache &> /dev/null; then
    gtk-update-icon-cache -f -t "$HOME/.local/share/icons/hicolor" 2>/dev/null || true
    echo "✓ Icon cache updated"
fi

echo ""
echo "Installation complete!"
echo "You can now:"
echo "  1. Search for 'FFwd UI' in your application menu"
echo "  2. Launch it from there - the icon should display correctly"
echo ""
echo "To uninstall, run: rm $DESKTOP_FILE $ICON_DEST"
