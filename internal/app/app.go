package app

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"

	"github.com/sanchey92/duplicate-finder/internal/config"
	"github.com/sanchey92/duplicate-finder/internal/finder"
	"github.com/sanchey92/duplicate-finder/internal/hash"
	"github.com/sanchey92/duplicate-finder/internal/walker"
	"github.com/sanchey92/duplicate-finder/internal/wp"
	"github.com/sanchey92/duplicate-finder/pkg/formatter"
)

type App struct {
	cfg       *config.Config
	finder    *finder.Finder
	formatter *formatter.Formatter
	output    io.WriteCloser
}

func New(cfg *config.Config) (*App, error) {
	if cfg.Path == "" {
		return nil, fmt.Errorf("path is required")
	}

	w, err := walker.New(cfg.Path)
	if err != nil {
		return nil, fmt.Errorf("walker init: %w", err)
	}

	h, err := hash.New(cfg.Algorithm)
	if err != nil {
		return nil, fmt.Errorf("hasher init: %w", err)
	}

	pool := wp.New(cfg.Workers, 0, h)

	f, err := finder.New(pool, w)
	if err != nil {
		return nil, fmt.Errorf("finder init: %w", err)
	}

	out, err := setupOutput(cfg.OutputFile)
	if err != nil {
		return nil, fmt.Errorf("output init: %w", err)
	}

	form, err := formatter.New(formatter.OutputFormat(cfg.Format), out)
	if err != nil {
		return nil, fmt.Errorf("formatter init: %w", err)
	}

	return &App{
		cfg:       cfg,
		finder:    f,
		formatter: form,
		output:    out,
	}, nil
}

func (a *App) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigCh := make(chan os.Signal, 1)
	errCh := make(chan error, 1)

	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		errCh <- a.run(ctx)
	}()

	select {
	case <-sigCh:
		_, _ = fmt.Fprintln(os.Stderr, "\n shutdown initiated...")
		cancel()
		<-errCh
		return fmt.Errorf("interrupted by signal")
	case err := <-errCh:
		return err
	}
}

func (a *App) run(ctx context.Context) error {
	if !a.cfg.Quiet {
		a.formatter.PrintHeader(a.cfg.Path, a.cfg.Algorithm, a.cfg.Workers)
	}

	var callback wp.ProgressCallback
	if !a.cfg.Quiet && a.cfg.ShowProgress {
		callback = a.formatter.PrintProgressBar
	}

	duplicates, stats, err := a.finder.Start(ctx, callback)
	if err != nil {
		return fmt.Errorf("scan failed: %w", err)
	}

	if err = a.formatter.PrintResults(duplicates, stats); err != nil {
		return fmt.Errorf("output failed: %w", err)
	}

	return nil
}

func setupOutput(path string) (io.WriteCloser, error) {
	if path == "" {
		return nopCloser{os.Stdout}, nil
	}

	f, err := os.Create(path)
	if err != nil {
		return nil, fmt.Errorf("cannot create output file: %w", err)
	}

	return f, nil
}

type nopCloser struct {
	io.Writer
}

func (nopCloser) Close() error {
	return nil
}
