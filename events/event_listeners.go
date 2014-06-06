package events

import (
	"log"

	graphite "github.com/marpaia/graphite-golang"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/uploader"
)

// InitializeEventListeners activates all workers that listen to events
func InitializeEventListeners(sc *core.ServerConfiguration) {
	uwc := uploader.UploadWorkers(sc.UploaderConcurrency)
	g := initializeGraphite(sc)
	go handleImageProcessed(sc, uwc, g)
	go handleImageProcessedWithErrors(sc, g)
	go handleOriginalDownloaded(sc, uwc, g)
	go handleOriginalDownloadUnavailable(sc, g)
}

func handleImageProcessed(sc *core.ServerConfiguration, uwc chan *uploader.UploadWork, g *graphite.Graphite) {
	for {
		ic := <-sc.Events.ImageProcessed
		resizedPath := ic.LocalResizedImagePath()
		log.Printf("Processed image: %s", resizedPath)

		uwc <- &uploader.UploadWork{ic, sc.Adapters.Uploader.Upload}
		g.SimpleSend("stats.image_server.image_request", "1")
		g.SimpleSend("stats.image_server.image_request."+ic.Format, "1")
	}
}

func handleImageProcessedWithErrors(sc *core.ServerConfiguration, g *graphite.Graphite) {
	for {
		_ = <-sc.Events.ImageProcessedWithErrors
		g.SimpleSend("stats.image_server.image_request_fail", "1")
	}
}

func handleOriginalDownloaded(sc *core.ServerConfiguration, uwc chan *uploader.UploadWork, g *graphite.Graphite) {
	for {
		ic := <-sc.Events.OriginalDownloaded
		uwc <- &uploader.UploadWork{ic, sc.Adapters.Uploader.UploadOriginal}
		g.SimpleSend("stats.image_server.original_downloaded", "1")
	}
}

func handleOriginalDownloadUnavailable(sc *core.ServerConfiguration, g *graphite.Graphite) {
	for {
		_ = <-sc.Events.OriginalDownloaded
		g.SimpleSend("stats.image_server.original_unavailable", "1")
	}
}
