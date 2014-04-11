package main

func initializeEventListeners(sc *ServerConfiguration, uwc chan *UploadWork) {
	go handleImageProcessed(sc, uwc)
	go handleImageProcessedWithErrors(sc)
	go handleOriginalDownloaded(sc, uwc)
	go handleOriginalDownloadUnavailable(sc)
}

func handleImageProcessed(sc *ServerConfiguration, uwc chan *UploadWork) {
	for {
		ic := <-sc.Events.ImageProcessed
		uwc <- &UploadWork{ic}
		sc.Graphite.SimpleSend("stats.image_server.image_request", "1")
		sc.Graphite.SimpleSend("stats.image_server.image_request."+ic.format, "1")
	}
}

func handleImageProcessedWithErrors(sc *ServerConfiguration) {
	for {
		_ = <-sc.Events.ImageProcessedWithErrors
		sc.Graphite.SimpleSend("stats.image_server.image_request_fail", "1")
	}
}

func handleOriginalDownloaded(sc *ServerConfiguration, uwc chan *UploadWork) {
	for {
		ic := <-sc.Events.OriginalDownloaded
		uwc <- &UploadWork{ic}
		sc.Graphite.SimpleSend("stats.image_server.original_downloaded", "1")
	}
}

func handleOriginalDownloadUnavailable(sc *ServerConfiguration) {
	for {
		_ = <-sc.Events.OriginalDownloaded
		sc.Graphite.SimpleSend("stats.image_server.original_unavailable", "1")
	}
}
