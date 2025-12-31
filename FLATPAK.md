# Flatpak Build Guide

This document provides detailed instructions for building and distributing FFwd UI as a Flatpak.

## Quick Start

Use the Makefile for the easiest build experience:

```bash
# Build and install locally
make flatpak-install

# Or just build without installing
make flatpak

# Or create a distributable bundle
make flatpak-bundle
```

The Makefile targets will automatically:
- Check for required dependencies
- Install necessary runtimes
- Build the Flatpak
- Provide next steps

## Manual Build Process

### 1. Install Prerequisites

```bash
# Ubuntu/Debian
sudo apt install flatpak flatpak-builder

# Fedora
sudo dnf install flatpak flatpak-builder

# Arch Linux
sudo pacman -S flatpak flatpak-builder
```

### 2. Add Flathub Repository

```bash
flatpak remote-add --if-not-exists flathub https://flathub.org/repo/flathub.flatpakrepo
```

### 3. Install Required Runtimes

```bash
flatpak install flathub org.gnome.Platform//47
flatpak install flathub org.gnome.Sdk//47
flatpak install flathub org.freedesktop.Sdk.Extension.golang
```

### 4. Build the Flatpak

Using Make (recommended):

```bash
make flatpak-install
```

Or manually:

```bash
flatpak-builder --force-clean build-dir io.github.thesle.FFwdUI.yml
flatpak-builder --user --install --force-clean build-dir io.github.thesle.FFwdUI.yml
```

### 5. Run the Application

```bash
flatpak run io.github.thesle.FFwdUI
```

## Creating a Distributable Bundle

To create a single `.flatpak` file that users can install:

Using Make (recommended):

```bash
make flatpak-bundle
```

Or manually:

```bash
# Build and create repository
flatpak-builder --repo=repo --force-clean build-dir io.github.thesle.FFwdUI.yml

# Create bundle
flatpak build-bundle repo ffwd-ui.flatpak io.github.thesle.FFwdUI
```

Users can then install with:

```bash
flatpak install ffwd-ui.flatpak
```

## Configuration Files

### Manifest (`io.github.thesle.FFwdUI.yml`)

The main Flatpak manifest that defines:
- **Runtime**: org.gnome.Platform 47 (includes WebKit2GTK 4.1)
- **SDK**: org.gnome.Sdk 47 with Golang extension
- **Modules**:
  - FFmpeg 6.1.1 with common codecs (x264, x265, VP9, AV1, MP3, Vorbis, Opus)
  - FFwd UI application

### Desktop Entry (`io.github.thesle.FFwdUI.desktop`)

Defines how the app appears in application menus and launchers.

### AppData/Metainfo (`io.github.thesle.FFwdUI.metainfo.xml`)

Provides metadata for software centers and app stores. Includes:
- Application description
- Feature list
- Screenshots (to be added)
- Release information
- Developer information

## Important Notes

### Before Building

1. **Update Git URL**: Edit `io.github.thesle.FFwdUI.yml` and change the git source URL and tag to match your repository:
   ```yaml
   sources:
     - type: git
       url: https://github.com/YOUR_USERNAME/ffwd-ui.git
       tag: v1.0.0  # or use 'branch: main' for development
   ```

2. **Icon Requirement**: Ensure `build/appicon.png` exists and is a 512x512 PNG file.

3. **Application ID**: The app ID `io.github.thesle.FFwdUI` follows the reverse-DNS naming convention. Update if you use a different domain.

### Permissions

The Flatpak has access to:
- Home directory (for selecting input/output files)
- XDG directories (Documents, Videos, Music, Downloads)
- Display (X11 and Wayland)
- GPU acceleration (for hardware encoding)

These are defined in the `finish-args` section of the manifest.

### Dependencies Included

The Flatpak bundles everything needed:
- **FFmpeg 6.1.1** with codecs for H.264, H.265, VP8/VP9, AV1, MP3, AAC, Vorbis, Opus, FLAC
- **WebKit2GTK 4.1** (included in GNOME Platform 47)
- **Go runtime** for building the application

Users don't need to install FFmpeg separately!

## Publishing to Flathub

To make your app available on Flathub:

1. Create a GitHub repository for your Flatpak manifest
2. Fork [flathub/flathub](https://github.com/flathub/flathub)
3. Add your app manifest
4. Submit a pull request
5. Follow the [Flathub submission guidelines](https://docs.flathub.org/docs/for-app-authors/submission)

Requirements for Flathub:
- Public Git repository
- AppData with screenshots
- Valid desktop file
- Appropriate license
- Working build

## Troubleshooting

### Build Fails with "No such file or directory"

Ensure all source files exist:
- `build/appicon.png` (512x512 PNG)
- `io.github.thesle.FFwdUI.desktop`
- `io.github.thesle.FFwdUI.metainfo.xml`

### FFmpeg Build Errors

If FFmpeg fails to build, check that you have enough disk space (minimum 5GB free).

### WebKit Build Errors

WebKit is a large dependency. Ensure you have:
- At least 8GB RAM
- 10GB+ free disk space
- Build time: 30-60 minutes depending on hardware

### Runtime Not Found

Install the required runtime:
```bash
flatpak install flathub org.freedesktop.Platform//23.08
```

### Permission Issues

If the app can't access files, check the permissions in `finish-args`. You may need to add additional filesystem access.

## Testing

After installation, verify:

1. **Application launches**: `flatpak run io.github.thesle.FFwdUI`
2. **File selection works**: Try selecting input/output files
3. **FFmpeg operations work**: Test a simple trim operation
4. **Command preview displays**: Check the command preview at the bottom

## Updating the Flatpak

When releasing a new version:

1. Update the version in `io.github.thesle.FFwdUI.metainfo.xml`
2. Add release notes
3. Update the git tag in `io.github.thesle.FFwdUI.yml`
4. Rebuild and redistribute

## Resources

- [Flatpak Documentation](https://docs.flatpak.org/)
- [Flatpak Builder Documentation](https://docs.flatpak.org/en/latest/flatpak-builder.html)
- [Flathub Submission Guidelines](https://docs.flathub.org/docs/for-app-authors/submission)
- [Desktop Entry Specification](https://specifications.freedesktop.org/desktop-entry-spec/latest/)
- [AppData Guidelines](https://www.freedesktop.org/software/appstream/docs/chap-Metadata.html)
