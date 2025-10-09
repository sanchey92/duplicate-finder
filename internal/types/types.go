package types

type FileInfo struct {
	Path string
	Size int64
	Hash string
}

type HashResult struct {
	File FileInfo
	Hash string
	Err  error
}

type ScanStats struct {
	TotalFiles       int
	ProcessedFiles   int
	DuplicateGroup   int
	TotalWastedSpace int64
}

type DuplicateGroup map[string][]FileInfo
