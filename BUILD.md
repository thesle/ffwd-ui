# FFmpeg Frontend UI - Build Plan

## Project Overview

A minimalistic, user-friendly desktop application for common FFmpeg operations built with Go, Wails 2.11, and Svelte. The app will be named `FFwd UI`.

## Technology Stack

- **Backend**: Go 1.25+
- **Frontend Framework**: Svelte 4
- **Desktop Framework**: Wails 2.11
- **CSS Library**: Bulma 0.9.4
- **FFmpeg**: External dependency (must be installed on system)

## Core Features

### 1. File Operations
- **Input File Selection**: Drag-and-drop or file picker
- **Output File Naming**: Manual naming with validation
- **File Information Display**: Show input file details (codec, duration, size)

### 2. Video Operations
- **Trim Start**: Remove N seconds from beginning
- **Trim to Length**: Cut video to specific duration
- **Extract Audio**: Save audio track as separate file (MP3, AAC, WAV)

### 3. System Information
- **Disk Space Display**: Show available space on all mount points/drives
- **Real-time Updates**: Refresh disk space before operations

### 4. Process Management
- **Threaded Execution**: Non-blocking FFmpeg operations
- **Context-based Cancellation**: Proper process termination
- **Progress Tracking**: Show encoding progress (if possible)
- **Command Preview**: Display exact FFmpeg command before execution

### 5. UI/UX Features
- **Minimalistic Design**: Clean, uncluttered interface
- **Responsive Feedback**: Loading states, error messages
- **Command Display**: Show FFmpeg commands at bottom of UI
- **Cancel Button**: Immediate process termination

## Architecture

### Backend (Go)

#### Package Structure
```
/
├── main.go                 # Wails app initialization
├── app.go                  # Main app struct
├── ffmpeg/
│   ├── executor.go        # FFmpeg command execution with context
│   ├── operations.go      # Operation builders (trim, extract, etc.)
│   └── probe.go          # FFprobe for file info
├── system/
│   ├── disk.go           # Disk space utilities
│   └── platform.go       # Platform-specific implementations
└── models/
    ├── config.go         # App configuration
    └── types.go          # Shared types
```

#### Key Backend Components

**1. FFmpeg Executor**
```go
type FFmpegExecutor struct {
    ctx    context.Context
    cancel context.CancelFunc
    cmd    *exec.Cmd
    mu     sync.Mutex
}

// Methods:
// - Execute(args []string) error
// - Cancel() error
// - GetProgress() (float64, error)
```

**2. Operation Builders**
```go
// TrimStart removes N seconds from beginning
func BuildTrimStartCommand(input, output string, seconds float64) []string

// TrimToLength cuts video to specific duration
func BuildTrimToLengthCommand(input, output string, duration float64) []string

// ExtractAudio saves audio track
func BuildExtractAudioCommand(input, output, format string) []string
```

**3. Disk Space Service**
```go
type MountPoint struct {
    Path      string
    Total     uint64
    Available uint64
    Used      uint64
}

func GetMountPoints() ([]MountPoint, error)
```

### Frontend (Svelte)

#### Component Structure
```
frontend/src/
├── App.svelte              # Main app component
├── components/
│   ├── FileInput.svelte   # Drag-drop file input
│   ├── Operations.svelte  # Operation selection panel
│   ├── DiskSpace.svelte   # Disk space display
│   ├── CommandPreview.svelte  # FFmpeg command display
│   └── ProgressBar.svelte # Operation progress
├── stores/
│   └── app.js            # Svelte stores for state
└── styles/
    └── main.css          # Custom styles + Bulma
```

#### State Management
Using Svelte stores for:
- Input/output file paths
- Selected operation and parameters
- FFmpeg command preview
- Disk space information
- Operation status (idle, running, complete, error)

### Wails Integration

**Exposed Go Methods** (called from frontend):
```go
// File operations
func (a *App) SelectInputFile() (string, error)
func (a *App) SelectOutputFile(defaultName string) (string, error)
func (a *App) GetFileInfo(path string) (*FileInfo, error)

// FFmpeg operations
func (a *App) TrimStart(input, output string, seconds float64) error
func (a *App) TrimToLength(input, output string, duration float64) error
func (a *App) ExtractAudio(input, output, format string) error
func (a *App) CancelOperation() error

// System info
func (a *App) GetDiskSpace() ([]MountPoint, error)

// Command preview
func (a *App) PreviewCommand(operation string, params map[string]interface{}) (string, error)
```

**Runtime Events** (Go → Frontend):
```go
// Emit progress updates
runtime.EventsEmit(ctx, "ffmpeg:progress", progress)

// Emit completion
runtime.EventsEmit(ctx, "ffmpeg:complete", result)

// Emit errors
runtime.EventsEmit(ctx, "ffmpeg:error", errorMsg)
```

## Implementation Phases

### Phase 1: Project Setup
- [ ] Initialize Wails project with Svelte template
- [ ] Install Bulma CSS framework
- [ ] Set up Go module dependencies
- [ ] Configure Wails build settings
- [ ] Create basic project structure

### Phase 2: Backend Core
- [ ] Implement FFmpeg executor with context cancellation
- [ ] Build operation command generators
- [ ] Create FFprobe wrapper for file info
- [ ] Implement disk space utilities (Linux, Windows, macOS)
- [ ] Add error handling and validation

### Phase 3: Frontend UI
- [ ] Create file input component (drag-drop + picker)
- [ ] Build operation selection interface
- [ ] Design output file naming form
- [ ] Implement disk space display component
- [ ] Add command preview panel at bottom
- [ ] Style with Bulma for clean, minimal look

### Phase 4: Integration
- [ ] Wire up Wails bindings
- [ ] Connect UI to backend methods
- [ ] Implement event listeners for progress/status
- [ ] Add cancel button functionality
- [ ] Handle errors and display user feedback

### Phase 5: Polish & Testing
- [ ] Test all operations with various file types
- [ ] Test cancellation behavior (no zombie processes)
- [ ] Verify disk space accuracy on all platforms
- [ ] Add input validation and error messages
- [ ] Performance testing with large files
- [ ] Cross-platform testing (Linux, Windows, macOS)

## Technical Implementation Details

### Context-Based FFmpeg Execution

```go
type App struct {
    ctx       context.Context
    ffmpegCtx context.Context
    ffmpegCancel context.CancelFunc
    currentCmd *exec.Cmd
    mu sync.Mutex
}

func (a *App) ExecuteFFmpeg(args []string) error {
    a.mu.Lock()
    
    // Create cancellable context with timeout
    ctx, cancel := context.WithTimeout(a.ctx, 30*time.Minute)
    a.ffmpegCtx = ctx
    a.ffmpegCancel = cancel
    
    // Create command with context
    cmd := exec.CommandContext(ctx, "ffmpeg", args...)
    a.currentCmd = cmd
    a.mu.Unlock()
    
    // Run in goroutine
    go func() {
        err := cmd.Run()
        if err != nil {
            if ctx.Err() == context.Canceled {
                runtime.EventsEmit(a.ctx, "ffmpeg:cancelled")
            } else {
                runtime.EventsEmit(a.ctx, "ffmpeg:error", err.Error())
            }
        } else {
            runtime.EventsEmit(a.ctx, "ffmpeg:complete")
        }
    }()
    
    return nil
}

func (a *App) CancelOperation() error {
    a.mu.Lock()
    defer a.mu.Unlock()
    
    if a.ffmpegCancel != nil {
        a.ffmpegCancel() // This kills the process immediately
    }
    return nil
}
```

### Disk Space Detection

**Linux**: Use `syscall.Statfs` on mount points
**Windows**: Use `GetDiskFreeSpaceEx` via syscall
**macOS**: Use `syscall.Statfs` similar to Linux

```go
// +build linux darwin

func getDiskSpace(path string) (uint64, uint64, error) {
    var stat syscall.Statfs_t
    err := syscall.Statfs(path, &stat)
    if err != nil {
        return 0, 0, err
    }
    
    total := stat.Blocks * uint64(stat.Bsize)
    available := stat.Bavail * uint64(stat.Bsize)
    
    return total, available, nil
}
```

### UI Layout (Minimalistic)

```
┌─────────────────────────────────────────────────┐
│  FFmpeg Frontend                                │
├─────────────────────────────────────────────────┤
│                                                 │
│  Input File:  [Drop file here or click]         │
│               /path/to/video.mp4 (1.2 GB)       │
│                                                 │
│  Operation:   ○ Trim Start (seconds: [5   ])    │
│               ○ Trim to Length (seconds: [60 ]) │
│               ● Extract Audio (format: [MP3▼])  │
│                                                 │
│  Output File: [video_audio.mp3           ]      │
│                                                 │
│  Disk Space:  / (root): 45.2 GB free            │
│               /home: 120.5 GB free              │
│                                                 │
│               [Execute] [Cancel]                │
│                                                 │
│  Progress:    ████████░░░░░░░░░░ 40%            │
│                                                 │
├─────────────────────────────────────────────────┤
│  Command: ffmpeg -i video.mp4 -vn -acodec       │
│           libmp3lame -q:a 2 video_audio.mp3     │
└─────────────────────────────────────────────────┘
```

## Dependencies

### Go Dependencies
```go
require (
    github.com/wailsapp/wails/v2 v2.11.0
)
```

### Frontend Dependencies
```json
{
  "dependencies": {
    "svelte": "^4.0.0",
    "bulma": "^0.9.4"
  }
}
```

### System Requirements
- FFmpeg installed and in PATH
- Go 1.21 or higher
- Node.js 16+ (for frontend build)
- Operating System: Linux, Windows, or macOS

## Build Commands

```bash
# Development
wails dev

# Production build
wails build

# Production build with optimization
wails build -clean -upx
```

## Error Handling Strategy

1. **Validation Errors**: Show in UI immediately (missing files, invalid parameters)
2. **FFmpeg Errors**: Parse stderr, display user-friendly messages
3. **System Errors**: Disk full, permissions, missing FFmpeg binary
4. **Cancellation**: Clean process termination, clear UI state

## Future Enhancements (Post-MVP)

- Batch processing multiple files
- Video format conversion
- Custom FFmpeg arguments (advanced mode)
- Preset configurations
- Drag timeline for visual trimming
- Video preview thumbnails
- Progress with ETA calculation
- Save/load operation templates

## Testing Strategy

### Unit Tests
- FFmpeg command builders
- Disk space calculations
- File path validation

### Integration Tests
- Full workflow tests with sample files
- Cancellation behavior verification
- Error handling scenarios

### Manual Testing
- Cross-platform UI/UX testing
- Large file handling
- Edge cases (invalid files, disk full, etc.)

---

**Ready to proceed?** Review this plan and let me know if you'd like any adjustments before we start building!
