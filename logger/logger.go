package logger

import "github.com/wanelo/image-server/core"

type Logger struct {
	Loggers []core.Logger
}

func (l *Logger) ImageProcessed(ic *core.ImageConfiguration) {
	for _, logger := range l.Loggers {
		go logger.ImageProcessed(ic)
	}
}

func (l *Logger) ImageProcessedWithErrors(ic *core.ImageConfiguration) {
	for _, logger := range l.Loggers {
		go logger.ImageProcessedWithErrors(ic)
	}
}

func (l *Logger) OriginalDownloaded(source string, destination string) {
	for _, logger := range l.Loggers {
		go logger.OriginalDownloaded(source, destination)
	}
}

func (l *Logger) OriginalDownloadFailed(source string) {
	for _, logger := range l.Loggers {
		go logger.OriginalDownloadFailed(source)
	}
}

func (l *Logger) OriginalDownloadSkipped(source string) {
	for _, logger := range l.Loggers {
		go logger.OriginalDownloadSkipped(source)
	}
}
