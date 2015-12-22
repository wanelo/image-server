package logfile

import (
	"github.com/golang/glog"
	"github.com/image-server/image-server/core"
	"github.com/image-server/image-server/logger"
)

type Logger struct {
}

func Enable() {
	l := &Logger{}
	logger.Loggers = append(logger.Loggers, l)
}

func (l *Logger) ImagePosted() {
}

func (l *Logger) ImagePostingFailed() {
}

func (l *Logger) ImageProcessed(ic *core.ImageConfiguration) {
}

func (l *Logger) ImageAlreadyProcessed(ic *core.ImageConfiguration) {
}

func (l *Logger) ImageProcessedWithErrors(ic *core.ImageConfiguration) {
}

func (l *Logger) AllImagesAlreadyProcessed(namespace string, hash string, sourceURL string) {
	glog.Warningf("All images already processed: namespace=%v hash=%v source=%v", namespace, hash, sourceURL)
}

func (l *Logger) SourceDownloaded() {
}

func (l *Logger) OriginalDownloaded(source string, destination string) {
}

func (l *Logger) OriginalDownloadFailed(source string) {
}

func (l *Logger) OriginalDownloadSkipped(source string) {
}
