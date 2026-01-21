package main

import (
	"context"
	"flag"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/checker"
	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/config"
	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/output"
	"github.com/YOUR_GITHUB_USERNAME/go-api-monitor/internal/scheduler"
)

func main() {
	var (
		configPath   = flag.String("config", "configs/sample.yaml", "Path to YAML config")
		once         = flag.Bool("once", false, "Run one check cycle and exit")
		jsonOut      = flag.Bool("json", true, "Output results as JSON (default true)")
		writeToFile  = flag.String("out", "", "Optional output file path (append mode)")
		concurrency  = flag.Int("concurrency", 5, "Max concurrent endpoint checks")
	)
	flag.Parse()

	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(logger)

	cfg, err := config.Load(*configPath)
	if err != nil {
		slog.Error("failed to load config", "error", err)
		os.Exit(2)
	}

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	client := checker.NewHTTPClient(time.Duration(cfg.TimeoutSeconds) * time.Second)
	chk := checker.NewChecker(client, *concurrency)

	writer := output.NewWriter(output.Options{
		JSON:       *jsonOut,
		OutputPath: *writeToFile,
	})

	runOnce := func(runCtx context.Context) error {
		results := chk.CheckAll(runCtx, cfg.Endpoints)
		return writer.Write(results)
	}

	if *once {
		if err := runOnce(ctx); err != nil {
			slog.Error("run failed", "error", err)
			os.Exit(1)
		}
		return
	}

	s := scheduler.New(time.Duration(cfg.IntervalSeconds)*time.Second, runOnce)
	if err := s.Run(ctx); err != nil {
		slog.Error("scheduler stopped with error", "error", err)
		os.Exit(1)
	}

	slog.Info("shutdown complete")
}
