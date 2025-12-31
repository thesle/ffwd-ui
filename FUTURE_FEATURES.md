# Future Feature Roadmap

This document tracks potential enhancements for FFwd UI that have been identified but not yet implemented.

## Batch Processing

**Priority:** High  
**Complexity:** Medium-High

Allow users to queue multiple files with the same operation settings and process them sequentially.

### Implementation Notes
- Add file list UI component
- Queue system with processing state
- Progress tracking per file
- Overall progress indicator
- Ability to reorder/remove queued items

### Use Cases
- Converting entire folder of videos to same format
- Batch trimming multiple clips
- Bulk resolution changes

---

## Recent Files

**Priority:** Medium  
**Complexity:** Low

Remember last 10-15 input and output file paths for quick access.

### Implementation Notes
- Store in localStorage
- Dropdown/autocomplete on file input
- Clear history option
- Group by operation type

### Benefits
- Faster workflow for repeated operations
- Easy access to common locations

---

## File Size Estimation

**Priority:** Medium  
**Complexity:** Medium

Estimate output file size before processing based on operation parameters.

### Implementation Notes
- Calculate from duration × bitrate
- Adjust for codec overhead
- Show "Estimated output: ~250 MB"
- Warn if insufficient disk space

### Calculations
- Copy operations: similar to input
- Bitrate changes: duration × target bitrate
- Resolution changes: scale by pixel ratio

---

## Operation History

**Priority:** Medium  
**Complexity:** Medium

Log completed operations with ability to view and rerun previous commands.

### Implementation Notes
- Store command history in file or database
- UI panel showing recent operations
- Click to rerun with same parameters
- Export history as shell script

### Storage Format
```json
{
  "timestamp": "2025-12-31T13:00:00Z",
  "operation": "trim_start",
  "input": "/path/to/input.mp4",
  "output": "/path/to/output.mp4",
  "command": "ffmpeg ...",
  "duration": 45.2,
  "success": true
}
```

---

## Crop Validation

**Priority:** High  
**Complexity:** Low

Validate crop parameters against actual video dimensions to prevent errors.

### Implementation Notes
- Check crop_width + crop_x <= video_width
- Check crop_height + crop_y <= video_height
- Show validation errors in UI
- Suggest maximum values

### Additional Validations
- Trim range: end > start
- Resolution: positive dimensions
- Volume: 0-500% reasonable range
- Bitrate: valid format (e.g., "2M", "500k")

---

## Progress Accuracy Improvement

**Priority:** Medium  
**Complexity:** Medium

Improve progress tracking accuracy beyond time-based estimation.

### Implementation Notes
- Parse FFmpeg frame count output
- Track bytes processed for file size operations
- Use multiple metrics for better estimation
- Handle operations with unknown duration

### FFmpeg Parsing
```
frame= 1234 fps= 45 q=-1.0 size=   12345kB time=00:00:41.23 bitrate=2450.1kbits/s
```
Extract frame count relative to total frames for better accuracy.

---

## Custom FFmpeg Path

**Priority:** Low  
**Complexity:** Low

Allow users to specify custom FFmpeg binary location.

### Implementation Notes
- Settings/preferences dialog
- Path validation (check if executable exists)
- Test run to verify FFmpeg version
- Default to system PATH if not specified

### UI
- Settings button in title bar or menu
- Text input with file picker
- "Test FFmpeg" button
- Display detected version

---

## Implementation Priority

Based on user value and complexity:

1. **Crop Validation** - Quick win, prevents frustrating errors
2. **Recent Files** - Low complexity, high convenience
3. **File Size Estimation** - Useful, moderate complexity
4. **Batch Processing** - High value but significant work
5. **Operation History** - Nice to have, moderate work
6. **Progress Accuracy** - Improvement to existing feature
7. **Custom FFmpeg Path** - Edge case, implement if users request

## Notes

These features were identified on 2025-12-31 during initial development. Priority and complexity assessments should be revisited as the application evolves and user feedback is gathered.
