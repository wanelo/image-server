package cli

import (
	"bufio"
	"io"
	"log"
	"strings"
	"sync"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher"
	"github.com/wanelo/image-server/parser"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/uploader"
)

func Process(sc *core.ServerConfiguration, namespace string, outputs []string, input io.Reader) error {
	done := make(chan struct{})
	defer close(done)

	idsc := enqueueAll(done, input)

	// Start a fixed number of goroutines to read and digest files.
	c := make(chan result) // HLc
	var wg sync.WaitGroup

	numDigesters := 10
	wg.Add(numDigesters)

	for i := 0; i < numDigesters; i++ {
		go func() {
			digester(sc, namespace, outputs, done, idsc, c) // HLc
			wg.Done()
		}()
	}
	go func() {
		wg.Wait()
		close(c) // HLc
	}()
	// End of pipeline. OMIT

	for r := range c {
		log.Printf("Completed processing image %v", r.ID)
	}

	return nil
}

func enqueueAll(done <-chan struct{}, input io.Reader) <-chan string {
	idsc := make(chan string)
	go func() { // HL
		// Close the ids channel after Walk returns.
		defer close(idsc) // HL

		reader := bufio.NewReader(input)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}
			idsc <- strings.TrimSpace(line)
		}
	}()
	return idsc
}

// A result is the product of reading and summing a file using MD5.
type result struct {
	ID  string
	Err error
}

// digester reads path names from paths and sends digests of the corresponding
// files on c until either paths or done is closed.
func digester(sc *core.ServerConfiguration, namespace string, outputs []string, done <-chan struct{}, ids <-chan string, c chan<- result) {
	for hash := range ids { // HLpaths
		log.Printf("About to process image: %s", hash)
		// download original image
		localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, hash)
		remoteOriginalPath := sc.Adapters.Paths.RemoteOriginalURL(namespace, hash)
		log.Println(remoteOriginalPath)
		f := fetcher.NewUniqueFetcher(remoteOriginalPath, localOriginalPath)
		err := f.Fetch()
		if err != nil {
			log.Printf("Unable to download image for %s: %s", hash, err)
			continue
		}

		for _, filename := range outputs {
			ic, err := parser.NameToConfiguration(sc, filename)
			if err != nil {
				log.Printf("Error parsing name: %v\n", err)
				continue
			}

			ic.Namespace = namespace
			ic.ID = hash

			// process image
			pchan := &processor.ProcessorChannels{
				ImageProcessed: make(chan *core.ImageConfiguration),
				Skipped:        make(chan string),
			}

			localPath := sc.Adapters.Paths.LocalImagePath(namespace, hash, filename)

			p := processor.Processor{
				Source:             localOriginalPath,
				Destination:        localPath,
				ImageConfiguration: ic,
				Channels:           pchan,
			}

			_, err = p.CreateImage()

			if err != nil {
				continue
			}

			select {
			case <-pchan.ImageProcessed:
				log.Println("about to upload to manta")
				uploader := &uploader.Uploader{sc.RemoteBasePath}
				remoteResizedPath := sc.Adapters.Paths.RemoteImagePath(namespace, hash, filename)
				err = uploader.Upload(localPath, remoteResizedPath)
				if err != nil {
					log.Println(err)
				}
			case path := <-pchan.Skipped:
				log.Printf("Skipped processing %s", path)
			}

		}

		select {
		case c <- result{hash, err}:
		case <-done:
			return
		}

	}
}
