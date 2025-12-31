# FFwd UI

A minimalistic, user-friendly desktop application for common FFmpeg operations built with Go, Wails 2.11, and Svelte.

## Features

- **Trim Start**: Remove seconds from the beginning of a video file
- **Trim to Length**: Cut video to a specific duration
- **Extract Audio**: Save audio track as separate file (MP3, AAC, WAV, FLAC)
- **Disk Space Monitoring**: View available space on all mount points/drives
- **Command Preview**: See exact FFmpeg command before execution
- **Threaded Execution**: Non-blocking operations with cancellation support
- **Progress Tracking**: Real-time progress updates during encoding

## Requirements

### System Dependencies
- **FFmpeg**: Must be installed and available in your system PATH
  - Ubuntu/Debian: `sudo apt install ffmpeg`
  - macOS: `brew install ffmpeg`
  - Windows: Download from [ffmpeg.org](https://ffmpeg.org/download.html)

### Development Dependencies
- Go 1.21 or higher
- Node.js 16+ (for frontend build)
- Wails CLI 2.11
  ```bash
  go install github.com/wailsapp/wails/v2/cmd/wails@latest
  ```

## Installation

1. Clone the repository
2. Install frontend dependencies:
   ```bash
   cd frontend
   npm install
   cd ..
   ```
3. Generate Wails bindings:
   ```bash
   wails generate module
   ```

## Development

Run in live development mode with hot reload:

```bash
wails dev
```

The app will launch with a development server. Changes to frontend code will hot-reload automatically.

## Building

Build a production binary:

```bash
wails build
```

For optimized build with UPX compression:

```bash
wails build -clean -upx
```

The built application will be in the `build/bin` directory.

## Usage

1. **Select Input File**: Click the drop zone or use the file picker
2. **Choose Operation**:
   - Trim Start: Remove N seconds from beginning
   - Trim to Length: Cut video to specific duration
   - Extract Audio: Save audio track in various formats
3. **Configure Parameters**: Adjust seconds, duration, or audio format
4. **Review Command**: Check the FFmpeg command at the bottom
5. **Execute**: Click Execute to start processing
6. **Monitor Progress**: Watch the progress bar and cancel if needed

## Architecture

### Backend (Go)
- **FFmpeg Executor**: Context-based execution with timeout and cancellation
- **Operations**: Command builders for trim and audio extraction
- **File Probe**: FFprobe integration for file information
- **Disk Space**: Platform-specific utilities for Linux, macOS, and Windows

### Frontend (Svelte + Bulma)
- **Reactive UI**: Svelte components with Bulma CSS framework
- **Event System**: Real-time updates via Wails runtime events
- **State Management**: Svelte stores for application state

## Key Technical Features

### Context-Based Process Management
FFmpeg processes run with context cancellation, ensuring:
- Immediate termination when Cancel is clicked
- No zombie processes consuming CPU
- Proper cleanup of system resources

### Cross-Platform Disk Space Detection
- Linux: Uses `syscall.Statfs` on `/proc/mounts`
- macOS: Uses `syscall.Statfs` with `df` parsing
- Windows: Uses `GetDiskFreeSpaceExW` for drive letters

## Project Structure

```
ffwd-ui/
├── app.go                 # Main app with Wails bindings
├── main.go                # Application entry point
├── ffmpeg/
│   ├── executor.go        # Context-based FFmpeg execution
│   ├── operations.go      # Command builders
│   └── probe.go          # FFprobe wrapper
├── system/
│   ├── disk_linux.go     # Linux disk space
│   ├── disk_windows.go   # Windows disk space
│   ├── mounts_*.go       # Platform-specific mount detection
├── models/
│   └── types.go          # Shared data types
└── frontend/
    └── src/
        ├── App.svelte    # Main UI component
        └── style.css     # Bulma + custom styles
```

## Troubleshooting

**FFmpeg not found**
- Ensure FFmpeg is installed and in your system PATH
- Test: `ffmpeg -version` in terminal

**Build errors**
- Run `wails generate module` to regenerate bindings
- Clear build cache: `wails build -clean`

**Frontend errors**
- Reinstall dependencies: `cd frontend && npm install`
- Check Node.js version: `node --version` (requires 16+)

## License

MIT License

## Contributing

Contributions welcome! Please feel free to submit a Pull Request.
