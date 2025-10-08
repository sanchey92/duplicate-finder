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
