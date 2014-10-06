package request

import (
	"fmt"
	"log"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/info"
	"github.com/wanelo/image-server/parser"
	"github.com/wanelo/image-server/processor"
)

type Request struct {
	ServerConfiguration *core.ServerConfiguration
	Namespace           string
	Outputs             []string
	Uploader            core.Uploader
	Paths               core.Paths
	Hash                string
	SourceURL           string
	directoryListing    map[string]string
}

func (r *Request) ProcessMultiple() error {
	missing, err := r.CalculateMissingOutputs()
	if err != nil {
		return err
	}

	uploadQueue := make(chan *core.ImageConfiguration)
	errorProcessingChannel := make(chan error)
	uploadedChannel := make(chan error)
	defer close(uploadedChannel)

	if missing == nil {
		// All the files are already uploaded. Nothing do do!
		return nil
	}

	// Process all the outputs
	go func() {
		defer close(errorProcessingChannel)
		for _, filename := range missing {
			log.Println(filename, "started")
			ic, err := parser.NameToConfiguration(r.ServerConfiguration, filename)
			if err != nil {
				errorProcessingChannel <- err
				return
			}
			ic.Namespace = r.Namespace
			ic.ID = r.Hash

			err = r.Process(ic)
			if err != nil {
				errorProcessingChannel <- err
				return
			}
			uploadQueue <- ic
			log.Println(filename, "complete")
		}
	}()

	// Upload all the outputs in parallel. This might be sequential if the
	// processing is slower than the uplaod
	for _, _ = range missing {
		go func() {
			var errU error
			select {
			case ic := <-uploadQueue:
				log.Println("about to upload!")
				localResizedPath := r.Paths.LocalImagePath(r.Namespace, r.Hash, ic.Filename)
				remoteResizedPath := r.Paths.RemoteImagePath(ic.Namespace, ic.ID, ic.Filename)
				errU = r.Uploader.Upload(localResizedPath, remoteResizedPath, ic.ToContentType())
			case errP := <-errorProcessingChannel:
				errU = errP
			}
			uploadedChannel <- errU
		}()
	}

	var firstErr error
	// wait till everything finishes, return on first error
	for _, _ = range missing {
		select {
		case err := <-uploadedChannel:
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

func (r *Request) Process(ic *core.ImageConfiguration) error {
	// The original file will be downloaded only once, even when every dimension requests it
	err := r.DownloadOriginal()
	if err != nil {
		return err
	}

	localResizedPath := r.Paths.LocalImagePath(r.Namespace, r.Hash, ic.Filename)
	localOriginalPath := r.Paths.LocalOriginalPath(r.Namespace, r.Hash)

	// process image
	pchan := &processor.ProcessorChannels{
		ImageProcessed: make(chan *core.ImageConfiguration),
		Skipped:        make(chan string),
	}

	info := &info.Info{
		Path: localOriginalPath,
	}
	id, err := info.ImageDetails()
	if err != nil {
		return err
	}

	p := processor.Processor{
		Source:             localOriginalPath,
		Destination:        localResizedPath,
		ImageConfiguration: ic,
		ImageDetails:       id,
		Channels:           pchan,
	}

	err = p.CreateImage()
	if err != nil {
		return err
	}

	select {
	case <-pchan.ImageProcessed:
		log.Printf("Processed (resize handler) %s", localResizedPath)
	case <-pchan.Skipped:
		log.Printf("Skipped processing (resize handler) %s", localResizedPath)
	}

	return nil
}

// CalculateMissingOutputs determine what versions need to be generated
func (r *Request) CalculateMissingOutputs() (itemOutputs []string, err error) {
	err = r.FetchRemoteFileListing()

	if err == nil {
		for _, output := range r.Outputs {
			if r.RemoteMissesFile(output) {
				itemOutputs = append(itemOutputs, output)
			}
		}

	} else {
		return nil, err
	}

	return itemOutputs, nil
}

func (r *Request) RemoteMissesFile(filename string) bool {
	_, ok := r.directoryListing[filename]
	return !ok
}

func (r *Request) FetchRemoteFileListing() error {
	if r.directoryListing == nil {
		r.directoryListing = make(map[string]string)
	} else {
		// Already fetched the listing
		return nil
	}
	if r.Uploader == nil {
		return fmt.Errorf("missing uploader")
	}

	remoteDirectory := r.Paths.RemoteImageDirectory(r.Namespace, r.Hash)
	entries, err := r.Uploader.ListDirectory(remoteDirectory)

	if err != nil {
		return err
	}

	for _, entry := range entries {
		r.directoryListing[entry] = entry
	}
	return nil
}
