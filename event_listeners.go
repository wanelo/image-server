package main

func initializeEventListeners(sc *ServerConfiguration) {
	go handleImageProcessed(sc)
	go handleOriginalDownloaded(sc)
}

func handleImageProcessed(sc *ServerConfiguration) {
	for {
		ic := <-sc.Events.ImageProcessed
		sc.DataStore.upload(ic.LocalResizedImagePath(), ic.MantaResizedImagePath())
	}
}

func handleOriginalDownloaded(sc *ServerConfiguration) {
	for {
		ic := <-sc.Events.OriginalDownloaded
		sc.DataStore.upload(ic.LocalOriginalImagePath(), ic.MantaOriginalImagePath())
	}
}
