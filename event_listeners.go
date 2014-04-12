package main

import "github.com/marpaia/graphite-golang"

func initializeEventListeners(sc *ServerConfiguration, uwc chan *UploadWork, g *graphite.Graphite) {
	go handleImageProcessed(sc, uwc, g)
	go handleImageProcessedWithErrors(sc, g)
	go handleOriginalDownloaded(sc, uwc, g)
	go handleOriginalDownloadUnavailable(sc, g)
}

func handleImageProcessed(sc *ServerConfiguration, uwc chan *UploadWork, g *graphite.Graphite) {
	for {
		ic := <-sc.Events.ImageProcessed
		uwc <- &UploadWork{ic}
		g.SimpleSend("stats.image_server.image_request", "1")
		g.SimpleSend("stats.image_server.image_request."+ic.format, "1")
	}
}

func handleImageProcessedWithErrors(sc *ServerConfiguration, g *graphite.Graphite) {
	for {
		_ = <-sc.Events.ImageProcessedWithErrors
		g.SimpleSend("stats.image_server.image_request_fail", "1")
	}
}

func handleOriginalDownloaded(sc *ServerConfiguration, uwc chan *UploadWork, g *graphite.Graphite) {
	for {
		ic := <-sc.Events.OriginalDownloaded
		uwc <- &UploadWork{ic}
		g.SimpleSend("stats.image_server.original_downloaded", "1")
	}
}

func handleOriginalDownloadUnavailable(sc *ServerConfiguration, g *graphite.Graphite) {
	for {
		_ = <-sc.Events.OriginalDownloaded
		g.SimpleSend("stats.image_server.original_unavailable", "1")
	}
}
