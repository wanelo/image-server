package core

import (
	"flag"
	"strings"
)

// ServerConfiguration struct
// Most of this configuration comes from json config
type ServerConfiguration struct {
	SourceDomain          string
	WhitelistedExtensions []string
	MaximumWidth          int
	LocalBasePath         string
	MantaBasePath         string
	DefaultQuality        uint
	GraphiteHost          string
	GraphitePort          int
	Environment           string
	NamespaceMappings     map[string]string
	Events                *EventChannels
	Adapters              *Adapters
}

// EventChannels struct
// Available image processing/downloading events
type EventChannels struct {
	ImageProcessed              chan *ImageConfiguration
	ImageProcessedWithErrors    chan *ImageConfiguration
	OriginalDownloaded          chan *ImageConfiguration
	OriginalDownloadUnavailable chan *ImageConfiguration
}

// NamespaceMapping Maps a url namespace with a source path i.e 'p' => 'product/images'
type NamespaceMapping struct {
	Namespace string
	Source    string
}

// ServerConfigurationFromFlags initializes a ServerConfiguration from flags
func ServerConfigurationFromFlags() (*ServerConfiguration, error) {
	whitelistedExtensions := flag.String("extensions", "jpg,gif,webp", "Whitelisted extensions (separated by commas)")
	localBasePath := flag.String("local_base_path", "public", "Directory where the images will be saved")
	graphitePort := flag.Int("graphite_port", 8125, "Graphite port")
	graphiteHost := flag.String("graphite_host", "127.0.0.1", "Graphite Host")
	maximumWidth := flag.Int("maximum_width", 1000, "maximum image width")
	mantaBasePath := flag.String("manta_base_path", "public/images/development", "base path for manta storage")
	defaultQuality := flag.Uint("default_quality", 75, "Default image compression quality")
	sourceDomain := flag.String("source_domain", "http://wanelo.s3.amazonaws.com", "Source domain for images")
	flag.Parse()

	sc := &ServerConfiguration{
		WhitelistedExtensions: strings.Split(*whitelistedExtensions, ","),
		LocalBasePath:         *localBasePath,
		GraphitePort:          *graphitePort,
		GraphiteHost:          *graphiteHost,
		MaximumWidth:          *maximumWidth,
		MantaBasePath:         *mantaBasePath,
		DefaultQuality:        *defaultQuality,
		SourceDomain:          *sourceDomain,
	}
	return sc, nil
}
