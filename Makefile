SHELL := /bin/bash

.DEFAULT_GOAL := help

.PHONY: dev dev-pr build help flatpak flatpak-install flatpak-bundle flatpak-check

# Use like:
#   make dev
#   make build

help:
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  dev              - Run wails dev (auto-detects distro/version for webkit2_41 tag)"
	@echo "  build            - Build the app (injects BuildDate; auto-detects webkit2_41 tag)"
	@echo "  flatpak          - Build Flatpak package"
	@echo "  flatpak-install  - Build and install Flatpak locally"
	@echo "  flatpak-bundle   - Create distributable Flatpak bundle"
	@echo "  flatpak-check    - Check Flatpak prerequisites"
	@echo "  help             - Show this help message"
	@echo ""
	@echo "Default target: help"

dev:
	@DISTRO=$$(lsb_release -is 2>/dev/null || echo "Unknown"); \
	VERSION=$$(lsb_release -rs 2>/dev/null || echo "0"); \
	# Use webkit2_41 for Ubuntu >= 24 and Linux Mint >= 22 \
	if [[ ("$$DISTRO" == "Ubuntu" && "$$(printf '%s\n' "$$VERSION" "24.00" | sort -V | head -n1)" == "24.00") || \
	  ("$$DISTRO" == "Elementary" && "$$(printf '%s\n' "$$VERSION" "8" | sort -V | head -n1)" == "8") || \
      ("$$DISTRO" == "Linuxmint" && "$$(printf '%s\n' "$$VERSION" "22.0" | sort -V | head -n1)" == "22.0") ]]; then \
	  echo "Detected Ubuntu >= 24 or Linux Mint >= 22. Using webkit2_41 build command..."; \
	  wails dev -tags webkit2_41; \
	else \
	  echo "Detected Linux Mint < 22, Ubuntu < 24, or another distro. Using normal build command..."; \
	  wails dev; \
	fi

build:
	@sudo rm -rf build-dir repo .flatpak-builder
	@DISTRO=$$(lsb_release -is 2>/dev/null || echo "Unknown"); \
	VERSION=$$(lsb_release -rs 2>/dev/null || echo "0"); \
	# Use webkit2_41 for Ubuntu >= 24 and Linux Mint >= 22 \
	if [[ ("$$DISTRO" == "Ubuntu" && "$$(printf '%s\n' "$$VERSION" "24.00" | sort -V | head -n1)" == "24.00") || \
	  ("$$DISTRO" == "Elementary" && "$$(printf '%s\n' "$$VERSION" "8" | sort -V | head -n1)" == "8") || \
      ("$$DISTRO" == "Linuxmint" && "$$(printf '%s\n' "$$VERSION" "22.0" | sort -V | head -n1)" == "22.0") ]]; then \
	  echo "Detected Ubuntu >= 24 or Linux Mint >= 22. Using webkit2_41 build command..."; \
	  wails build -tags webkit2_41 --ldflags="-X main.BuildDate=$$(date -u '+%Y-%m-%d_%H:%M:%S')"; \
	else \
	  echo "Detected Linux Mint < 22, Ubuntu < 24, or another distro. Using normal build command..."; \
	  wails build --ldflags="-X main.BuildDate=$$(date -u '+%Y-%m-%d_%H:%M:%S')"; \
	fi

flatpak-check:
	@echo "Checking Flatpak prerequisites..."
	@command -v flatpak-builder >/dev/null 2>&1 || { echo "Error: flatpak-builder is not installed. Install with: sudo apt install flatpak-builder"; exit 1; }
	@echo "✓ flatpak-builder found"
	@flatpak list | grep -q "org.gnome.Platform.*49" || { echo "Installing org.gnome.Platform//49..."; flatpak install -y flathub org.gnome.Platform//49; }
	@flatpak list | grep -q "org.gnome.Sdk.*49" || { echo "Installing org.gnome.Sdk//49..."; flatpak install -y flathub org.gnome.Sdk//49; }
	@echo "✓ All prerequisites installed"

flatpak: build flatpak-check
	@echo "Building Flatpak package..."
	flatpak-builder --force-clean build-dir io.github.thesle.FFwdUI.yml
	@echo ""
	@echo "Build complete! To install locally, run: make flatpak-install"

flatpak-install: build flatpak-check
	@echo "Building and installing Flatpak locally..."
	flatpak-builder --user --install --force-clean build-dir io.github.thesle.FFwdUI.yml
	@echo ""
	@echo "Installation complete! Run with: flatpak run io.github.thesle.FFwdUI"

flatpak-bundle: build flatpak-check
	@echo "Creating distributable Flatpak bundle..."
	flatpak-builder --repo=repo --force-clean build-dir io.github.thesle.FFwdUI.yml
	flatpak build-bundle repo ffwd-ui.flatpak io.github.thesle.FFwdUI
	@echo ""
	@echo "Bundle created: ffwd-ui.flatpak"
	@echo "Users can install with: flatpak install ffwd-ui.flatpak"
