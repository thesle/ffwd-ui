SHELL := /bin/bash

.DEFAULT_GOAL := help

.PHONY: dev dev-pr build help

# Use like:
#   make dev
#   make build

help:
	@echo "Usage:"
	@echo "  make [target]"
	@echo ""
	@echo "Targets:"
	@echo "  dev    - Run wails dev (auto-detects distro/version for webkit2_41 tag)"
	@echo "  build  - Build the app (injects BuildDate; auto-detects webkit2_41 tag)"
	@echo "  help   - Show this help message"
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
