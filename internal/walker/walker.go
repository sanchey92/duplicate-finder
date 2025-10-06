package walker

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/sanchey92/duplicate-finder/internal/types"
)

type Walker struct {
	rootPath string
}

type FileFilter func(info types.FileInfo, osInfo os.FileInfo) bool

func New(path string) (*Walker, error) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		return nil, fmt.Errorf("directory does not exist: %s", path)
	}

	info, err := os.Stat(path)
	if err != nil {
		return nil, fmt.Errorf("cannot access path: %s: %w", path, err)
	}

	if !info.IsDir() {
		return nil, fmt.Errorf("path is not directory: %s", path)
	}

	return &Walker{rootPath: path}, nil
}

func (w *Walker) WalkFiles(filter FileFilter) ([]types.FileInfo, error) {
	var files []types.FileInfo

	if filter == nil {
		filter = defaultFilter()
	}

	err := filepath.WalkDir(w.rootPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			fmt.Printf("Warning: cannot access %s: %v\n", path, err)
			return nil
		}

		if d.IsDir() {
			return nil
		}

		info, err := d.Info()
		if err != nil {
			fmt.Printf("Warning: cannot get info %s: %v\n", path, err)
		}

		fileInfo := types.FileInfo{
			Path: path,
			Size: info.Size(),
		}

		if !filter(fileInfo, info) {
			return nil
		}

		files = append(files, fileInfo)
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to walk dirrectory %s: %w", w.rootPath, err)
	}

	return files, nil
}

func defaultFilter() FileFilter {
	return func(info types.FileInfo, osInfo os.FileInfo) bool {
		if osInfo.Size() == 0 {
			return false
		}

		if osInfo.Mode()&os.ModeSymlink != 0 {
			return false
		}

		if !osInfo.Mode().IsRegular() {
			return false
		}
		return true
	}
}

func (w *Walker) GetRootPath() string {
	return w.rootPath
}
