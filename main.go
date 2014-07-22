package main

import (
	"log"
	"os"
	"strings"

	"github.com/codegangsta/cli"
	"github.com/wanelo/image-server/core"
	fetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/logger"
	"github.com/wanelo/image-server/logger/graphite"
	"github.com/wanelo/image-server/paths"
	processor "github.com/wanelo/image-server/processor/cli"
	"github.com/wanelo/image-server/server"
	"github.com/wanelo/image-server/uploader/manta"
)

func main() {
	app := cli.NewApp()
	app.Name = "images"
	app.Version = "0.1.0"
	app.Usage = "Image server and CLI"
	app.Action = func(c *cli.Context) {
		println("boom! I say!")
	}

	app.Flags = []cli.Flag{
		cli.StringFlag{Name: "port", Value: "7000", Usage: "Specifies the server port."},
		cli.StringFlag{Name: "extensions", Value: "jpg,gif,webp", Usage: "Whitelisted extensions (separated by commas)"},
		cli.StringFlag{Name: "local_base_path", Value: "public", Usage: "Directory where the images will be saved"},
		cli.StringFlag{Name: "source_domain", Value: "http://wanelo.s3.amazonaws.com", Usage: "Source domain for images"},
		cli.StringFlag{Name: "remote_base_path", Value: "public/images/development", Usage: "base path for manta storage"},
		cli.StringFlag{Name: "graphite_host", Value: "127.0.0.1", Usage: "Graphite host"},
		cli.IntFlag{Name: "graphite_port", Value: 8125, Usage: "Graphite port"},
		cli.IntFlag{Name: "maximum_width", Value: 1000, Usage: "Maximum image width"},
		cli.IntFlag{Name: "default_quality", Value: 75, Usage: "Default image compression quality"},
		cli.IntFlag{Name: "uploader_concurrency", Value: 10, Usage: "Uploader concurrency"},
	}

	app.Commands = []cli.Command{
		{Name: "server",
			ShortName: "s",
			Usage:     "image server",
			Action: func(c *cli.Context) {
				sc, err := serverConfiguration(c)
				if err != nil {
					log.Panicln(err)
				}

				port := c.GlobalString("port")
				server.InitializeRouter(sc, port)
			},
		},
	}

	app.Run(os.Args)
}

func serverConfiguration(c *cli.Context) (*core.ServerConfiguration, error) {
	sc := serverConfigurationFromContext(c)

	loggers := []core.Logger{
		graphite.New(sc.GraphiteHost, sc.GraphitePort),
	}

	adapters := &core.Adapters{
		Fetcher:   &fetcher.Fetcher{},
		Processor: &processor.Processor{},
		Uploader:  manta.InitializeUploader(sc.RemoteBasePath),
		Paths:     &paths.Paths{sc.LocalBasePath, sc.RemoteBasePath},
		Logger:    &logger.Logger{loggers},
	}
	sc.Adapters = adapters

	return sc, nil
}

func serverConfigurationFromContext(c *cli.Context) *core.ServerConfiguration {
	return &core.ServerConfiguration{
		WhitelistedExtensions: strings.Split(c.GlobalString("extensions"), ","),
		LocalBasePath:         c.GlobalString("local_base_path"),
		GraphitePort:          c.GlobalInt("graphite_port"),
		GraphiteHost:          c.GlobalString("graphite_host"),
		MaximumWidth:          c.GlobalInt("maximum_width"),
		RemoteBasePath:        c.GlobalString("remote_base_path"),
		DefaultQuality:        uint(c.GlobalInt("default_quality")),
		SourceDomain:          c.GlobalString("source_domain"),
		UploaderConcurrency:   uint(c.GlobalInt("uploader_concurrency")),
	}
}
