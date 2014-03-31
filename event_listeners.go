package main

func initializeEventListeners(sc *ServerConfiguration) {
	go handleImageProcessed(sc)
	go handleOriginalDownloaded(sc)
}

func handleImageProcessed(sc *ServerConfiguration) {
	for {
		ic := <-sc.Events.ImageProcessed
		sc.DataStore.upload(ic.LocalResizedImagePath(), ic.MantaResizedImagePath())
		sc.Graphite.SimpleSend("stats.image_server.image_processed", "1")
	}
}

func handleOriginalDownloaded(sc *ServerConfiguration) {
	for {
		ic := <-sc.Events.OriginalDownloaded
		sc.DataStore.upload(ic.LocalOriginalImagePath(), ic.MantaOriginalImagePath())
		sc.Graphite.SimpleSend("stats.image_server.original_downloaded", "1")
	}
}
