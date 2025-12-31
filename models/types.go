package models

type FileInfo struct {
	Path     string  `json:"path"`
	Size     int64   `json:"size"`
	Duration float64 `json:"duration"`
	Format   string  `json:"format"`
	Codec    string  `json:"codec"`
	Width    int     `json:"width"`
	Height   int     `json:"height"`
}

type MountPoint struct {
	Path      string `json:"path"`
	Total     uint64 `json:"total"`
	Available uint64 `json:"available"`
	Used      uint64 `json:"used"`
}

type OperationParams struct {
	Operation string                 `json:"operation"`
	Input     string                 `json:"input"`
	Output    string                 `json:"output"`
	Params    map[string]interface{} `json:"params"`
}

type ProgressUpdate struct {
	Percent float64 `json:"percent"`
	Message string  `json:"message"`
}
