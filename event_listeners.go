package main

func initializeEventListeners(sc *ServerConfiguration) {
	go handleImageProcessed(sc)
	go handleImageProcessedWithErrors(sc)
	go handleOriginalDownloaded(sc)
	go handleOriginalDownloadUnavailable(sc)
}

func handleImageProcessed(sc *ServerConfiguration) {
	for {
		ic := <-sc.Events.ImageProcessed
		sc.DataStore.upload(ic.LocalResizedImagePath(), ic.MantaResizedImagePath())
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

func handleOriginalDownloaded(sc *ServerConfiguration) {
	for {
		ic := <-sc.Events.OriginalDownloaded
		sc.DataStore.upload(ic.LocalOriginalImagePath(), ic.MantaOriginalImagePath())
		sc.Graphite.SimpleSend("stats.image_server.original_downloaded", "1")
	}
}

func handleOriginalDownloadUnavailable(sc *ServerConfiguration) {
	for {
		_ = <-sc.Events.OriginalDownloaded
		sc.Graphite.SimpleSend("stats.image_server.original_unavailable", "1")
	}
}
