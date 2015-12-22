package request

import (
	"fmt"
	"io"

	"github.com/image-server/image-server/core"
	"github.com/image-server/image-server/logger"
	"github.com/image-server/image-server/parser"
)

type Request struct {
	ServerConfiguration *core.ServerConfiguration
	Namespace           string
	Outputs             []string
	Uploader            core.Uploader
	Paths               core.Paths
	Hash                string
	SourceURL           string
	SourceData          io.ReadCloser
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
		logger.AllImagesAlreadyProcessed(r.Namespace, r.Hash, r.SourceURL)
		return nil
	}

	// Process all the outputs
	go func() {
		defer close(errorProcessingChannel)
		for _, filename := range missing {
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
		}
	}()

	// Upload all the outputs in parallel. This might be sequential if the
	// processing is slower than the uplaod
	for range missing {
		go func() {
			var errU error
			select {
			case ic := <-uploadQueue:
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
	for range missing {
		select {
		case err := <-uploadedChannel:
			if err != nil && firstErr == nil {
				firstErr = err
			}
		}
	}

	return firstErr
}

// CalculateMissingOutputs determine what versions need to be generated
func (r *Request) CalculateMissingOutputs() (itemOutputs []string, err error) {
	if r.Outputs == nil {
		return nil, nil
	}

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
