package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher"
	"github.com/wanelo/image-server/parser"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/uploader"
)

func ResizeHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	vars := mux.Vars(req)
	filename := vars["filename"]

	ic, err := parser.NameToConfiguration(sc, filename)
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound)
		return
	}

	ic.ID = unpartitionHash(vars["id1"], vars["id2"], vars["id3"], vars["id4"])
	ic.Namespace = vars["namespace"]

	localResizedPath := sc.Adapters.Paths.LocalImagePath(ic.Namespace, ic.ID, ic.Filename)
	localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(ic.Namespace, ic.ID)
	remoteOriginalPath := sc.Adapters.Paths.RemoteOriginalURL(ic.Namespace, ic.ID)

	// download original image
	f := fetcher.NewUniqueFetcher(remoteOriginalPath, localOriginalPath)
	_, err = f.Fetch()
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound)
		return
	}

	err = processAndUpload(sc, ic)
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound)
		return
	}

	http.ServeFile(w, req, localResizedPath)
}

func unpartitionHash(p1, p2, p3, p4 string) string {
	return fmt.Sprintf("%s%s%s%s", p1, p2, p3, p4)
}

func processAndUpload(sc *core.ServerConfiguration, ic *core.ImageConfiguration) error {
	localResizedPath := sc.Adapters.Paths.LocalImagePath(ic.Namespace, ic.ID, ic.Filename)
	localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(ic.Namespace, ic.ID)

	// process image
	pchan := &processor.ProcessorChannels{
		ImageProcessed: make(chan *core.ImageConfiguration),
		Skipped:        make(chan string),
	}

	p := processor.Processor{
		Source:             localOriginalPath,
		Destination:        localResizedPath,
		ImageConfiguration: ic,
		Channels:           pchan,
	}

	err := p.CreateImage()
	if err != nil {
		return err
	}

	select {
	case <-pchan.ImageProcessed:
		uploader := uploader.DefaultUploader(sc)
		remoteResizedPath := sc.Adapters.Paths.RemoteImagePath(ic.Namespace, ic.ID, ic.Filename)
		err = uploader.Upload(localResizedPath, remoteResizedPath, ic.ToContentType())
	case path := <-pchan.Skipped:
		log.Printf("Skipped processing %s", path)
	}
	return err
}
