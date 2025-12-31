# FFwd UI - Quick Start Guide

## Prerequisites Check

Before running the application, ensure you have:

1. **FFmpeg installed**
   ```bash
   ffmpeg -version
   ```
   If not installed:
   - Ubuntu/Debian: `sudo apt install ffmpeg`
   - Arch: `sudo pacman -S ffmpeg`
   - macOS: `brew install ffmpeg`

2. **Go 1.21+**
   ```bash
   go version
   ```

3. **Node.js 16+**
   ```bash
   node --version
   ```

4. **Wails CLI**
   ```bash
   wails version
   ```
   If not installed:
   ```bash
   go install github.com/wailsapp/wails/v2/cmd/wails@latest
   ```

## First-Time Setup

1. **Install frontend dependencies:**
   ```bash
   cd frontend
   npm install
   cd ..
   ```

2. **Generate Wails bindings:**
   ```bash
   wails generate module
   ```

## Running the App

### Development Mode (Hot Reload)
```bash
wails dev
```

The app will launch with a development server. Frontend changes hot-reload automatically.

### Production Build
```bash
wails build
```

The binary will be in `build/bin/` directory.

## Testing the App

Once the app launches:

1. **Click the drop zone** to select a test video file
2. **Choose an operation:**
   - Trim Start: Remove 5 seconds from beginning
   - Trim to Length: Cut to 60 seconds
   - Extract Audio: Save as MP3
3. **Check the command preview** at the bottom
4. **Click Execute** to process
5. **Test Cancel** button during processing

## Common Operations

### Trim 10 seconds from start
1. Select operation: Trim Start
2. Set seconds: 10
3. Execute

### Extract audio as MP3
1. Select operation: Extract Audio
2. Choose format: MP3
3. Execute

### Cut video to 30 seconds
1. Select operation: Trim to Length
2. Set duration: 30
3. Execute

## Troubleshooting

**App won't start?**
- Check `wails doctor` for missing dependencies
- Ensure all prerequisites are installed

**FFmpeg errors?**
- Verify FFmpeg is in PATH: `which ffmpeg` (Linux/macOS) or `where ffmpeg` (Windows)
- Try a different input file (some codecs may not be supported)

**Build errors?**
- Run `wails build -clean` to clear cache
- Regenerate bindings: `wails generate module`

**Frontend issues?**
- Delete `node_modules` and reinstall: `cd frontend && rm -rf node_modules && npm install`

## Expected Behavior

### Successful Operation
1. Progress bar updates from 0% to 100%
2. Green success notification appears
3. Output file created at specified location

### Cancellation
1. Click Cancel during processing
2. FFmpeg process terminates immediately
3. No output file created (or partial file)
4. No zombie processes left running

## Performance Notes

- **Trim operations** (with `-c copy`) are fast - they don't re-encode
- **Audio extraction** re-encodes audio, takes longer
- **Progress tracking** updates in real-time based on FFmpeg output
- **Cancel** kills the process immediately via context cancellation

## Next Steps

Once familiar with basic operations, you can:
- Process batch files (future feature)
- Try different audio formats (AAC, WAV, FLAC)
- Experiment with various trim durations
- Monitor disk space before large operations
