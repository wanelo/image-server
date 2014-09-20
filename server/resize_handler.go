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

// ResizeHandler asumes the original image is either stores locally or on the remote server
// it returns the processed image in given dimension and format.
// When an image is requested more than once, only one will do the processing,
// and both requests will return the same output
func ResizeHandler(w http.ResponseWriter, req *http.Request, sc *core.ServerConfiguration) {
	vars := mux.Vars(req)
	filename := vars["filename"]

	ic, err := parser.NameToConfiguration(sc, filename)
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound)
		return
	}

	ic.ID = varsToHash(vars)
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

func varsToHash(vars map[string]string) string {
	return fmt.Sprintf("%s%s%s%s", vars["id1"], vars["id2"], vars["id3"], vars["id4"])
}

func processAndUploadFromOutputs(sc *core.ServerConfiguration, localOriginalPath string, namespace string, hash string, outputs []string) error {
	for _, filename := range outputs {
		ic, err := parser.NameToConfiguration(sc, filename)
		if err != nil {
			return err
		}
		ic.Namespace = namespace
		ic.ID = hash

		err = processAndUpload(sc, ic)
		if err != nil {
			return err
		}
	}
	return nil
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
