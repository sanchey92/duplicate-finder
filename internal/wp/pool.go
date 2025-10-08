package wp

import (
	"context"
	"fmt"
	"sync"

	"github.com/sanchey92/duplicate-finder/internal/types"
)

type Hasher interface {
	Calculate(path string) (string, error)
}

type ProgressCallback func(processed, total int)

type WorkerPool struct {
	wg         sync.WaitGroup
	workers    int
	totalFiles int
	jobsCh     chan types.FileInfo
	resultCh   chan types.HashResult
	hasher     Hasher
}

func New(workersNum int, length int, hasher Hasher) *WorkerPool {
	return &WorkerPool{
		workers:    workersNum,
		totalFiles: length,
		jobsCh:     make(chan types.FileInfo, workersNum*2),
		resultCh:   make(chan types.HashResult, workersNum*2),
		hasher:     hasher,
	}
}

func (wp *WorkerPool) Process(ctx context.Context, files []types.FileInfo, cb ProgressCallback) ([]types.FileInfo, error) {
	for i := 0; i < wp.workers; i++ {
		wp.wg.Add(1)
		go wp.worker(ctx)
	}

	go wp.sendJobs(ctx, files)

	go func() {
		wp.wg.Wait()
		close(wp.resultCh)
	}()

	return wp.collectResults(ctx, cb)
}

func (wp *WorkerPool) sendJobs(ctx context.Context, files []types.FileInfo) {
	defer close(wp.jobsCh)

	for _, file := range files {
		select {
		case <-ctx.Done():
			return
		case wp.jobsCh <- file:
		}
	}
}

func (wp *WorkerPool) worker(ctx context.Context) {
	defer wp.wg.Done()
	for {
		select {
		case <-ctx.Done():
			return
		case file, ok := <-wp.jobsCh:
			if !ok {
				return
			}

			h, err := wp.hasher.Calculate(file.Path)
			select {
			case <-ctx.Done():
				return
			case wp.resultCh <- types.HashResult{
				File: file,
				Hash: h,
				Err:  err,
			}:
			}
		}
	}
}

func (wp *WorkerPool) collectResults(ctx context.Context, cb ProgressCallback) ([]types.FileInfo, error) {
	processedFiles := make([]types.FileInfo, 0, wp.totalFiles)
	processed := 0

	for {
		select {
		case <-ctx.Done():
			return processedFiles, fmt.Errorf("operation canceled")
		case result, ok := <-wp.resultCh:
			if !ok {
				fmt.Printf("Worker pool finished: processed %d/%d files\n", processed, wp.totalFiles)
				return processedFiles, nil
			}
			processed++

			wp.printProgress(processed, cb)

			if result.Err != nil {
				fmt.Printf("Warning: cannot calculate hash for %s: %v\n", result.File.Path, result.Err)
				continue
			}

			result.File.Hash = result.Hash
			processedFiles = append(processedFiles, result.File)
		}
	}
}

func (wp *WorkerPool) printProgress(processed int, cb ProgressCallback) {
	if cb != nil {
		cb(processed, wp.totalFiles)
	} else if processed%100 == 0 || processed == wp.totalFiles {
		fmt.Printf("Progress: %d/%d files processed\n", processed, wp.totalFiles)
	}
}
