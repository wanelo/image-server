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
	mantaclient "github.com/wanelo/image-server/uploader/manta/client"
)

// Item represents all image properties needed for the result of the processing
type Item struct {
	Hash   string
	URL    string
	Width  int
	Height int
}

// ToTabDelimited creates a tab delimited text representation of an Item
func (i Item) ToTabDelimited() string {
	return fmt.Sprintf("%s\t%s\t%d\t%d\n", i.Hash, i.URL, i.Width, i.Height)
}

// Process instanciates image processing based on the tab delimited input that
// contains source image urls and hashes. Each image is processed by a pool of
// of digesters
func Process(sc *core.ServerConfiguration, namespace string, outputs []string, input io.Reader) error {
	done := make(chan struct{})
	defer close(done)

	idsc := enqueueAll(done, input)

	// Start a fixed number of goroutines to read and digest images.
	c := make(chan result) // HLc
	var wg sync.WaitGroup

	numDigesters := int(sc.ProcessorConcurrency)
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

// digester processes image items till done is received.
func digester(sc *core.ServerConfiguration, namespace string, outputs []string, done <-chan struct{}, items <-chan *Item, c chan<- result) {
	for item := range items {
		var itemOutputs []string
		fetchedExistingOutputs := false
		var existingFiles map[string]mantaclient.Entry

		if item.Hash != "" {
			itemOutputs, existingFiles, _ = calculateMissingOutputs(sc, namespace, item.Hash, outputs)
			fetchedExistingOutputs = true
		}

		imageDetails, err := downloadOriginal(sc, namespace, item)
		if err != nil {
			continue
		}

		// have not tried to retrieve existing outputs
		if !fetchedExistingOutputs {
			itemOutputs, existingFiles, err = calculateMissingOutputs(sc, namespace, item.Hash, outputs)
			log.Println(itemOutputs, err)
			if err != nil {
				// process all outputs
				copy(itemOutputs, outputs)
			}
		}

		if _, ok := existingFiles["original"]; !ok {
			err = uploadOriginal(sc, namespace, item, imageDetails)
			if err != nil {
				continue
			}
		}

		log.Printf("About to process image: %s", item.Hash)
		localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, item.Hash)

		for _, filename := range itemOutputs {
			err := processImage(sc, namespace, item.Hash, localOriginalPath, filename)
			if err != nil {
				log.Println(err)
				continue
			}
		}

		select {
		case c <- result{item.Hash, err}:
		case <-done:
			return
		}

		fmt.Fprint(os.Stdout, item.ToTabDelimited())
	}
}

func calculateMissingOutputs(sc *core.ServerConfiguration, namespace string, imageHash string, outputs []string) ([]string, map[string]mantaclient.Entry, error) {
	// Determine what versions need to be generated
	var itemOutputs []string
	c := mantaclient.DefaultClient()
	m := make(map[string]mantaclient.Entry)
	remoteDirectory := sc.Adapters.Paths.RemoteImageDirectory(namespace, imageHash)
	entries, err := c.ListDirectory(remoteDirectory)
	if err == nil {

		for _, entry := range entries {
			if entry.Type == "object" {
				m[entry.Name] = entry
			} else {
				// got a directory
				invalidDirectoryPath := fmt.Sprintf("%s/%s", remoteDirectory, entry.Name)
				err2 := c.DeleteDirectory(invalidDirectoryPath)
				if err2 != nil {
					log.Println("Unable to delete directory", invalidDirectoryPath)
				} else {
					log.Println("Deleted invalid directory", invalidDirectoryPath)
				}
			}
		}

		for _, output := range outputs {
			if _, ok := m[output]; ok {
				log.Printf("Skipping %s/%s", remoteDirectory, output)
			} else {
				itemOutputs = append(itemOutputs, output)
			}
		}

	} else {
		return nil, nil, err
	}

	return itemOutputs, m, nil
}

func downloadOriginal(sc *core.ServerConfiguration, namespace string, item *Item) (*info.ImageDetails, error) {
	// Image does not have a hash, need to upload source and get image hash
	f := fetcher.OriginalFetcher{Paths: sc.Adapters.Paths}
	imageDetails, _, err := f.Fetch(namespace, item.URL, item.Hash)

	if err != nil {
		return nil, err
	}

	hash := imageDetails.Hash
	item.Width = imageDetails.Width
	item.Height = imageDetails.Height
	item.Hash = hash

	return imageDetails, nil
}

func uploadOriginal(sc *core.ServerConfiguration, namespace string, item *Item, imageDetails *info.ImageDetails) error {

	localOriginalPath := sc.Adapters.Paths.LocalOriginalPath(namespace, item.Hash)
	uploader := uploader.DefaultUploader(sc.RemoteBasePath)

	err := uploader.CreateDirectory(sc.Adapters.Paths.RemoteImageDirectory(namespace, item.Hash))
	if err != nil {
		return err
	}

	destination := sc.Adapters.Paths.RemoteOriginalPath(namespace, item.Hash)
	localInfoPath := sc.Adapters.Paths.LocalInfoPath(namespace, item.Hash)
	remoteInfoPath := sc.Adapters.Paths.RemoteInfoPath(namespace, item.Hash)

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
		uploader := uploader.DefaultUploader(sc.RemoteBasePath)
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
