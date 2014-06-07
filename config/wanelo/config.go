package config

import (
	"flag"
	"strings"

	"github.com/wanelo/image-server/core"
	"github.com/wanelo/image-server/events"
	"github.com/wanelo/image-server/processor/cli"
	sm "github.com/wanelo/image-server/source_mapper/waneloS3"
	"github.com/wanelo/image-server/uploader/manta"
)

// ServerConfiguration initializes a ServerConfiguration from flags
func ServerConfiguration() (*core.ServerConfiguration, error) {
	sc := configurationFromFlags()

	sc.Events = &core.EventChannels{
		ImageProcessed:     make(chan *core.ImageConfiguration),
		OriginalDownloaded: make(chan *core.ImageConfiguration),
	}

	adapters := &core.Adapters{
		Processor:    &cli.Processor{sc},
		SourceMapper: &sm.SourceMapper{},
		Uploader:     manta.InitializeUploader(sc),
	}
	sc.Adapters = adapters

	go events.InitializeEventListeners(sc)

	return sc, nil
}

func configurationFromFlags() *core.ServerConfiguration {
	var (
		whitelistedExtensions = flag.String("extensions", "jpg,gif,webp", "Whitelisted extensions (separated by commas)")
		localBasePath         = flag.String("local_base_path", "public", "Directory where the images will be saved")
		sourceDomain          = flag.String("source_domain", "http://wanelo.s3.amazonaws.com", "Source domain for images")
		mantaBasePath         = flag.String("manta_base_path", "public/images/development", "base path for manta storage")
		graphiteHost          = flag.String("graphite_host", "127.0.0.1", "Graphite Host")
		graphitePort          = flag.Int("graphite_port", 8125, "Graphite port")
		maximumWidth          = flag.Int("maximum_width", 1000, "Maximum image width")
		defaultQuality        = flag.Uint("default_quality", 75, "Default image compression quality")
		uploaderConcurrency   = flag.Uint("uploader_concurrency", 10, "Uploader concurrency")
	)
	flag.Parse()

	sc := &core.ServerConfiguration{
		WhitelistedExtensions: strings.Split(*whitelistedExtensions, ","),
		LocalBasePath:         *localBasePath,
		GraphitePort:          *graphitePort,
		GraphiteHost:          *graphiteHost,
		MaximumWidth:          *maximumWidth,
		MantaBasePath:         *mantaBasePath,
		DefaultQuality:        *defaultQuality,
		SourceDomain:          *sourceDomain,
		UploaderConcurrency:   *uploaderConcurrency,
	}

	return sc
}
