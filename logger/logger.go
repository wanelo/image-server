package logger

import "github.com/wanelo/image-server/core"

var Loggers []core.Logger

func initialize() {
	Loggers = []core.Logger{}
}

func ImagePosted() {
	for _, logger := range Loggers {
		go logger.ImagePosted()
	}
}

func ImagePostingFailed() {
	for _, logger := range Loggers {
		go logger.ImagePostingFailed()
	}
}

func ImageProcessed(ic *core.ImageConfiguration) {
	for _, logger := range Loggers {
		go logger.ImageProcessed(ic)
	}
}

func ImageProcessedWithErrors(ic *core.ImageConfiguration) {
	for _, logger := range Loggers {
		go logger.ImageProcessedWithErrors(ic)
	}
}

func SourceDownloaded() {
	for _, logger := range Loggers {
		go logger.SourceDownloaded()
	}
}

func OriginalDownloaded(source string, destination string) {
	for _, logger := range Loggers {
		go logger.OriginalDownloaded(source, destination)
	}
}

func OriginalDownloadFailed(source string) {
	for _, logger := range Loggers {
		go logger.OriginalDownloadFailed(source)
	}
}

func OriginalDownloadSkipped(source string) {
	for _, logger := range Loggers {
		go logger.OriginalDownloadSkipped(source)
	}
}
