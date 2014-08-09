package cli

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sync"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/fetcher"
	"github.com/wanelo/image-server/info"
	"github.com/wanelo/image-server/parser"
	"github.com/wanelo/image-server/processor"
	"github.com/wanelo/image-server/uploader"
)

type Item struct {
	Hash   string
	URL    string
	Width  int
	Height int
}

func (i Item) ToTabDelimited() string {
	return fmt.Sprintf("%s\t%s\t%d\t%d\n", i.Hash, i.URL, i.Width, i.Height)
}

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

func enqueueAll(done <-chan struct{}, input io.Reader) <-chan *Item {
	idsc := make(chan *Item)
	go func() { // HL
		// Close the ids channel after Walk returns.
		defer close(idsc) // HL

		reader := bufio.NewReader(input)
		for {
			line, err := reader.ReadString('\n')
			if err != nil {
				break
			}

			item, err := lineToItem(line)
			if err == nil {
				idsc <- item
			}

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
func digester(sc *core.ServerConfiguration, namespace string, outputs []string, done <-chan struct{}, items <-chan *Item, c chan<- result) {
	for item := range items { // HLpaths
		hash := item.Hash

		if hash == "" {
			err := downloadOriginal(sc, namespace, item)
			if err != nil {
				continue
			}
			if item.Hash == "" {
				log.Panic("It should have created an image hash")
			}
			hash = item.Hash
		}

		log.Printf("About to process image: %s", hash)

		localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, hash)
		remoteOriginalPath := sc.Adapters.Paths.RemoteOriginalURL(namespace, hash)
		log.Println(remoteOriginalPath)
		f := fetcher.NewUniqueFetcher(remoteOriginalPath, localOriginalPath)
		_, err := f.Fetch()
		if err != nil {
			log.Printf("Unable to download image for %s: %s", hash, err)
			continue
		}

		for _, filename := range outputs {
			err := processImage(sc, namespace, hash, localOriginalPath, filename)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		select {
		case c <- result{hash, err}:
		case <-done:
			return
		}

		fmt.Fprint(os.Stdout, item.ToTabDelimited())
	}
}

func downloadOriginal(sc *core.ServerConfiguration, namespace string, item *Item) error {
	if item.URL == "" {
		return fmt.Errorf("Missing Hash & URL")
	}

	// Image does not have a hash, need to upload source and get image hash
	f := fetcher.NewSourceFetcher(sc.Adapters.Paths)
	imageDetails, downloaded, err := f.Fetch(item.URL, namespace)

	if err != nil {
		return err
	}

	hash := imageDetails.Hash
	if downloaded {
		localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, hash)
		uploader := &uploader.Uploader{sc.RemoteBasePath}

		err := uploader.CreateDirectory(sc.Adapters.Paths.RemoteImageDirectory(namespace, hash))
		if err != nil {
			return err
		}

		destination := sc.Adapters.Paths.RemoteOriginalPath(namespace, hash)
		localInfoPath := sc.Adapters.Paths.LocalInfoPath(namespace, hash)
		remoteInfoPath := sc.Adapters.Paths.RemoteInfoPath(namespace, hash)

		err = info.SaveImageDetail(imageDetails, localInfoPath)
		if err != nil {
			return err
		}

		// upload info
		err = uploader.Upload(localInfoPath, remoteInfoPath)
		if err != nil {
			return err
		}

		// upload original image
		err = uploader.Upload(localOriginalPath, destination)
		if err != nil {
			return err
		}
	}

  item.Width = imageDetails.Width
	item.Height = imageDetails.Height
	item.Hash = hash
	return nil
}

func processImage(sc *core.ServerConfiguration, namespace string, hash string, localOriginalPath string, filename string) error {
	ic, err := parser.NameToConfiguration(sc, filename)
	if err != nil {
		return fmt.Errorf("Error parsing name: %v\n", err)
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
		return err
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

	return nil
}

func lineToItem(line string) (*Item, error) {
	hr, _ := regexp.Compile("([a-z0-9]{32})")
	ur, _ := regexp.Compile("(htt[^\t\n\f\r ]+)")

	hash := hr.FindString(line)
	url := ur.FindString(line)
	return &Item{hash, url, 0, 0}, nil
}
