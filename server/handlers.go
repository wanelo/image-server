package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/unrolled/render"
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

	namespace := vars["namespace"]
	id1 := vars["id1"]
	id2 := vars["id2"]
	id3 := vars["id3"]
	id4 := vars["id4"]
	hash := fmt.Sprintf("%s%s%s%s", id1, id2, id3, id4)

	ic.ID = hash
	ic.Namespace = namespace

	localPath := sc.Adapters.Paths.LocalImagePath(namespace, hash, filename)
	localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, hash)

	// download original image
	remoteOriginalPath := sc.Adapters.Paths.RemoteOriginalURL(namespace, hash)
	log.Println(remoteOriginalPath)
	f := fetcher.NewUniqueFetcher(remoteOriginalPath, localOriginalPath)
	_, err = f.Fetch()
	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound)
		return
	}

	// process image
	pchan := &processor.ProcessorChannels{
		ImageProcessed: make(chan *core.ImageConfiguration),
		Skipped:        make(chan string),
	}
	defer close(pchan.ImageProcessed)
	defer close(pchan.Skipped)

	p := processor.Processor{
		Source:             localOriginalPath,
		Destination:        localPath,
		ImageConfiguration: ic,
		Channels:           pchan,
	}

	resizedPath, err := p.CreateImage()

	if err != nil {
		errorHandler(err, w, req, http.StatusNotFound)
		return
	}

	select {
	case <-pchan.ImageProcessed:
		log.Println("about to upload to manta")
		uploader := uploader.DefaultUploader(sc)
		remoteResizedPath := sc.Adapters.Paths.RemoteImagePath(namespace, hash, filename)
		err = uploader.Upload(localPath, remoteResizedPath, ic.ToContentType())
		if err != nil {
			log.Println(err)
		}
	case path := <-pchan.Skipped:
		log.Printf("Skipped processing %s", path)
	}

	http.ServeFile(w, req, resizedPath)
}

func errorHandlerJSON(err error, w http.ResponseWriter, status int) {
	r := render.New(render.Options{
		IndentJSON: true,
	})

	log.Println(err)
	json := map[string]string{
		"error": fmt.Sprintf("%s", err),
	}
	r.JSON(w, status, json)
}

func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int) {
	w.WriteHeader(status)
	if status == http.StatusNotFound {
		fmt.Fprint(w, "404 image not available. ", err)
	}
}

// func genericImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *core.ServerConfiguration) {
// 	ic, err := parser.NameToConfiguration(sc, params["filename"])
// 	if err != nil {
// 		errorHandler(err, w, r, http.StatusNotFound, ic)
// 	}
// 	qs := r.URL.Query()
//
// 	ic.ServerConfiguration = sc
// 	ic.Namespace = params["namespace"]
// 	ic.ID = params["id1"] + params["id2"] + params["id3"]
// 	ic.Source = qs.Get("source")
// 	imageHandler(ic, w, r)
//
// }

// func multiImageHandler(params martini.Params, r *http.Request, w http.ResponseWriter, sc *core.ServerConfiguration) {
// 	qs := r.URL.Query()
//
// 	go func() {
// 		outputs := strings.Split(qs.Get("outputs"), ",")
// 		fmt.Println("multiImageHandler")
// 		fmt.Println(outputs)
// 		for _, filename := range outputs {
// 			fmt.Println(filename)
// 			ic, err := parser.NameToConfiguration(sc, filename)
// 			if err != nil {
// 				continue
// 			}
// 			qs := r.URL.Query()
//
// 			ic.ServerConfiguration = sc
// 			ic.Namespace = params["namespace"]
// 			ic.ID = params["id1"] + params["id2"] + params["id3"]
// 			ic.Source = qs.Get("source")
//
// 			allowed, _ := allowedImage(ic)
// 			if allowed {
// 				err := sc.Adapters.Fetcher.FetchOriginal(ic)
// 				if err != nil {
// 					return
// 				}
// 				sc.Adapters.Processor.CreateImage(ic)
// 			}
// 		}
// 	}()
//
// 	ic := &core.ImageConfiguration{
// 		ServerConfiguration: sc,
// 		Namespace:           params["namespace"],
// 		ID:                  params["id1"] + params["id2"] + params["id3"],
// 		Source:              qs.Get("source"),
// 	}
// 	err := sc.Adapters.Fetcher.FetchOriginal(ic)
// 	if err != nil {
// 		errorHandler(err, w, r, http.StatusNotFound, ic)
// 		return
// 	}
// }
//
// func imageHandler(ic *core.ImageConfiguration, w http.ResponseWriter, r *http.Request) {
// 	allowed, err := allowedImage(ic)
// 	if !allowed {
// 		errorHandler(err, w, r, http.StatusNotAcceptable, ic)
// 		return
// 	}
//
// 	sc := ic.ServerConfiguration
// 	err = sc.Adapters.Fetcher.FetchOriginal(ic)
// 	if err != nil {
// 		errorHandler(err, w, r, http.StatusNotFound, ic)
// 		return
// 	}
// 	resizedPath, err := sc.Adapters.Processor.CreateImage(ic)
// 	if err != nil {
// 		errorHandler(err, w, r, http.StatusNotFound, ic)
// 		return
// 	}
// 	http.ServeFile(w, r, resizedPath)
// }
//
// func errorHandler(err error, w http.ResponseWriter, r *http.Request, status int, ic *core.ImageConfiguration) {
// 	w.WriteHeader(status)
// 	if status == http.StatusNotFound {
// 		fmt.Fprint(w, "404 image not available. ", err)
// 	}
//
// 	go func() {
// 		sc := ic.ServerConfiguration
// 		sc.Events.ImageProcessedWithErrors <- ic
// 	}()
// }
//
// func allowedImage(ic *core.ImageConfiguration) (bool, error) {
// 	sc := ic.ServerConfiguration
// 	// verify maximum width
// 	if ic.Width > sc.MaximumWidth {
// 		err := fmt.Errorf("maximum width is: %v\n", sc.MaximumWidth)
// 		return false, err
// 	}
//
// 	// verify image format
// 	for _, ext := range sc.WhitelistedExtensions {
// 		if ext == ic.Format {
// 			return true, nil
// 		}
// 	}
// 	return false, fmt.Errorf("format not allowed %s", ic.Format)
// }
