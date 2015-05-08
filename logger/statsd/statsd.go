package statsd

import (
	"fmt"
	"log"
	"time"

	"github.com/quipo/statsd"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/logger"
)

type Logger struct {
	Host   string
	Port   int
	Prefix string
	statsd *statsd.StatsdBuffer
}

func Enable(host string, port int, prefix string) {
	l := &Logger{Host: host, Port: port, Prefix: prefix}
	l.initializeStatsd()

	logger.Loggers = append(logger.Loggers, l)
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
	metric := fmt.Sprintf("%s_count", name)
	l.statsd.Incr(metric, 1)
}

func (l *Logger) initializeStatsd() {
	server := fmt.Sprintf("%v:%v", l.Host, l.Port)
	statsdclient := statsd.NewStatsdClient(server, l.Prefix)
	statsdclient.CreateSocket()
	interval := time.Second * 2 // aggregate stats and flush every 2 seconds
	l.statsd = statsd.NewStatsdBuffer(interval, statsdclient)
	// defer stats.Close()

	log.Println("Loaded Statsd connection:", server)
}
