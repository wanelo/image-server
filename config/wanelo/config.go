package config

import (
	"flag"
	"strings"

	"github.com/wanelo/image-server/core"
	fetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/logger"
	"github.com/wanelo/image-server/logger/graphite"
	"github.com/wanelo/image-server/paths"
	"github.com/wanelo/image-server/processor/cli"
	"github.com/wanelo/image-server/uploader/manta"
)

// ServerConfiguration initializes a ServerConfiguration from flags
func ServerConfiguration() (*core.ServerConfiguration, error) {
	sc := configurationFromFlags()

	loggers := []core.Logger{
		graphite.New(sc.GraphiteHost, sc.GraphitePort),
	}

	adapters := &core.Adapters{
		Fetcher:   &fetcher.Fetcher{},
		Processor: &cli.Processor{},
		Uploader:  manta.InitializeUploader(sc.RemoteBasePath),
		Paths:     &paths.Paths{sc.LocalBasePath, sc.RemoteBasePath},
		Logger:    &logger.Logger{loggers},
	}
	sc.Adapters = adapters

	return sc, nil
}

func configurationFromFlags() *core.ServerConfiguration {
	var (
		whitelistedExtensions = flag.String("extensions", "jpg,gif,webp", "Whitelisted extensions (separated by commas)")
		localBasePath         = flag.String("local_base_path", "public", "Directory where the images will be saved")
		sourceDomain          = flag.String("source_domain", "http://wanelo.s3.amazonaws.com", "Source domain for images")
		remoteBasePath        = flag.String("remote_base_path", "public/images/development", "base path for manta storage")
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
		RemoteBasePath:        *remoteBasePath,
		DefaultQuality:        *defaultQuality,
		SourceDomain:          *sourceDomain,
		UploaderConcurrency:   *uploaderConcurrency,
	}

	return sc
}
