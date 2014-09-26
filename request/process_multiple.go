package request

import (
	"fmt"
	"log"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher"
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
	directoryListing    map[string]string
}

func (r *Request) ProcessMultiple() error {
	missing, err := r.CalculateMissingOutputs()
	if err != nil {
		return err
	}

	if missing == nil {
		// All the files are already uploaded. Nothing do do!
		return nil
	}

	for _, filename := range missing {
		ic, err := parser.NameToConfiguration(r.ServerConfiguration, filename)
		if err != nil {
			return err
		}
		ic.Namespace = r.Namespace
		ic.ID = r.Hash

		err = r.Process(ic)
		if err != nil {
			return err
		}
	}
	return err
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

	p := processor.Processor{
		Source:             localOriginalPath,
		Destination:        localResizedPath,
		ImageConfiguration: ic,
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

	remoteResizedPath := r.Paths.RemoteImagePath(ic.Namespace, ic.ID, ic.Filename)
	err = r.Uploader.Upload(localResizedPath, remoteResizedPath, ic.ToContentType())
	return err
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

func (r *Request) DownloadOriginal() error {
	localOriginalPath := r.Paths.LocalOriginalPath(r.Namespace, r.Hash)
	remoteOriginalPath := r.Paths.RemoteOriginalURL(r.Namespace, r.Hash)

	// download original image
	f := fetcher.NewUniqueFetcher(remoteOriginalPath, localOriginalPath)
	_, err := f.Fetch()
	return err
}
