package graphite

import (
	"fmt"
	"log"

	"github.com/marpaia/graphite-golang"
	"github.com/wanelo/image-server/core"
)

type Logger struct {
	Host     string
	Port     int
	graphite *graphite.Graphite
}

func New(h string, p int) (l *Logger) {
	logger := &Logger{Host: h, Port: p}
	logger.initializeGraphite()
	return logger
}

func (l *Logger) ImagePosted() {
	l.track("new_image.request")
}

func (l *Logger) ImagePostingFailed() {
	l.track("new_image.request_failed")
}

func (l *Logger) ImageProcessed(ic *core.ImageConfiguration) {
	l.track("processed")
	l.track("processed." + ic.Format)
}

func (l *Logger) ImageProcessedWithErrors(ic *core.ImageConfiguration) {
	l.track("processed_failed")
	l.track("processed_failed." + ic.Format)
}

func (l *Logger) SourceDownloaded() {
	l.track("fetch.source_downloaded")
}

func (l *Logger) OriginalDownloaded(source string, destination string) {
	l.track("fetch.original_downloaded")
}

func (l *Logger) OriginalDownloadFailed(source string) {
	l.track("fetch.original_unavailable")
}

func (l *Logger) OriginalDownloadSkipped(source string) {
	l.track("fetch.original_download_skipped")
}

func (l *Logger) track(name string) {
	metric := fmt.Sprintf("stats.image_server.%s", name)
	l.graphite.SimpleSend(metric, "1")
}

func (l *Logger) initializeGraphite() {
	var err error

	l.graphite, err = graphite.NewGraphite(l.Host, l.Port)

	// if you couldn't connect to graphite, use a nop
	if err != nil {
		l.graphite = graphite.NewGraphiteNop(l.Host, l.Port)
	}

	log.Println("Loaded Graphite connection")
}
