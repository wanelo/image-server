package processor

type ImageProcessingResult struct {
	ResizedPath string
	Error       error
}

var ImageProcessings map[string][]chan ImageProcessingResult
