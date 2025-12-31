package ffmpeg

import (
	"bufio"
	"context"
	"fmt"
	"os/exec"
	"regexp"
	"strconv"
	"sync"
	"time"
)

type Executor struct {
	ctx          context.Context
	ffmpegCtx    context.Context
	ffmpegCancel context.CancelFunc
	currentCmd   *exec.Cmd
	mu           sync.Mutex
	onProgress   func(percent float64, message string)
	onComplete   func()
	onError      func(error)
}

func NewExecutor(ctx context.Context) *Executor {
	return &Executor{
		ctx: ctx,
	}
}

func (e *Executor) SetProgressCallback(cb func(percent float64, message string)) {
	e.onProgress = cb
}

func (e *Executor) SetCompleteCallback(cb func()) {
	e.onComplete = cb
}

func (e *Executor) SetErrorCallback(cb func(error)) {
	e.onError = cb
}

func (e *Executor) Execute(args []string, duration float64) error {
	e.mu.Lock()

	ctx, cancel := context.WithTimeout(e.ctx, 30*time.Minute)
	e.ffmpegCtx = ctx
	e.ffmpegCancel = cancel

	cmd := exec.CommandContext(ctx, "ffmpeg", args...)
	e.currentCmd = cmd

	stderr, err := cmd.StderrPipe()
	if err != nil {
		e.mu.Unlock()
		return fmt.Errorf("failed to get stderr pipe: %w", err)
	}

	if err := cmd.Start(); err != nil {
		e.mu.Unlock()
		return fmt.Errorf("failed to start ffmpeg: %w", err)
	}

	e.mu.Unlock()

	go func() {
		scanner := bufio.NewScanner(stderr)
		timeRegex := regexp.MustCompile(`time=(\d{2}):(\d{2}):(\d{2}\.\d{2})`)

		for scanner.Scan() {
			line := scanner.Text()

			if matches := timeRegex.FindStringSubmatch(line); len(matches) == 4 {
				hours, _ := strconv.ParseFloat(matches[1], 64)
				minutes, _ := strconv.ParseFloat(matches[2], 64)
				seconds, _ := strconv.ParseFloat(matches[3], 64)

				currentTime := hours*3600 + minutes*60 + seconds

				if duration > 0 {
					percent := (currentTime / duration) * 100
					if percent > 100 {
						percent = 100
					}
					if e.onProgress != nil {
						e.onProgress(percent, fmt.Sprintf("Processing... %.1f%%", percent))
					}
				}
			}
		}
	}()

	go func() {
		err := cmd.Wait()

		e.mu.Lock()
		e.ffmpegCancel = nil
		e.currentCmd = nil
		e.mu.Unlock()

		if err != nil {
			if ctx.Err() == context.Canceled {
				if e.onError != nil {
					e.onError(fmt.Errorf("operation cancelled"))
				}
			} else {
				if e.onError != nil {
					e.onError(fmt.Errorf("ffmpeg error: %w", err))
				}
			}
		} else {
			if e.onComplete != nil {
				e.onComplete()
			}
		}
	}()

	return nil
}

func (e *Executor) Cancel() error {
	e.mu.Lock()
	defer e.mu.Unlock()

	if e.ffmpegCancel != nil {
		e.ffmpegCancel()
		return nil
	}

	return fmt.Errorf("no operation running")
}

func (e *Executor) IsRunning() bool {
	e.mu.Lock()
	defer e.mu.Unlock()
	return e.currentCmd != nil
}
