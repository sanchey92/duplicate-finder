package finder

import (
	"context"
	"fmt"

	"github.com/sanchey92/duplicate-finder/internal/types"
	"github.com/sanchey92/duplicate-finder/internal/walker"
	"github.com/sanchey92/duplicate-finder/internal/wp"
)

type Pool interface {
	Process(ctx context.Context, files []types.FileInfo, cb wp.ProgressCallback) ([]types.FileInfo, error)
}
type Walker interface {
	WalkFiles(filter walker.FileFilter) ([]types.FileInfo, error)
}

type Finder struct {
	pool   Pool
	walker Walker
}

func New(wp Pool, w Walker) (*Finder, error) {
	if wp == nil || w == nil {
		return nil, fmt.Errorf("worker pool and walker must be initialized")
	}

	return &Finder{
		pool:   wp,
		walker: w,
	}, nil
}

func (f *Finder) Start(ctx context.Context, cb wp.ProgressCallback) (types.DuplicateGroup, *types.ScanStats, error) {
	files, err := f.collectFiles()
	if err != nil {
		return nil, nil, err
	}

	totalFiles := len(files)
	if totalFiles == 0 {
		return make(types.DuplicateGroup), &types.ScanStats{}, nil
	}

	processedFiles, err := f.pool.Process(ctx, files, cb)
	if err != nil {
		return nil, nil, err
	}

	hashGroup := f.groupByHash(processedFiles)
	duplicated := f.filterDuplicates(hashGroup)
	stats := f.stats(totalFiles, len(processedFiles), duplicated)

	return duplicated, stats, nil
}

func (f *Finder) collectFiles() ([]types.FileInfo, error) {
	files, err := f.walker.WalkFiles(nil)
	if err != nil {
		return nil, fmt.Errorf("failed to walk files: %w", err)
	}

	return files, nil
}

func (f *Finder) groupByHash(files []types.FileInfo) types.DuplicateGroup {
	hashGroups := make(types.DuplicateGroup)

	for _, file := range files {
		if file.Hash == "" {
			continue
		}
		hashGroups[file.Hash] = append(hashGroups[file.Hash], file)
	}

	return hashGroups
}

func (f *Finder) filterDuplicates(hashGroups types.DuplicateGroup) types.DuplicateGroup {
	count := 0
	for _, g := range hashGroups {
		if len(g) > 1 {
			count++
		}
	}

	duplicates := make(types.DuplicateGroup, count)

	for h, g := range hashGroups {
		if len(g) > 1 {
			duplicates[h] = g
		}
	}
	return duplicates

}

func (f *Finder) stats(totalFiles, processed int, duplicates types.DuplicateGroup) *types.ScanStats {
	stats := &types.ScanStats{
		TotalFiles:     totalFiles,
		ProcessedFiles: processed,
		DuplicateGroup: len(duplicates),
	}

	for _, group := range duplicates {
		if len(group) < 2 {
			continue
		}

		var totalSize int64
		for _, file := range group {
			totalSize += file.Size
		}

		wastedSize := totalSize - group[0].Size
		stats.TotalWastedSpace += wastedSize
	}

	return stats
}
