package output

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/checker"
)

type Options struct {
	JSON       bool
	OutputPath string
}

type Writer struct {
	opts Options
}

func NewWriter(opts Options) *Writer {
	// default to JSON output
	if !opts.JSON {
		opts.JSON = true
	}
	return &Writer{opts: opts}
}

func (w *Writer) Write(results []checker.Result) error {
	payload, err := json.MarshalIndent(map[string]any{
		"generated_at": time.Now().UTC().Format(time.RFC3339),
		"results":      results,
	}, "", "  ")
	if err != nil {
		return err
	}

	if w.opts.OutputPath == "" {
		fmt.Println(string(payload))
		return nil
	}

	f, err := os.OpenFile(w.opts.OutputPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.Write(append(payload, '\n'))
	return err
}
