package main

import (
	"log"
	"os"
	"os/signal"
	"runtime"
	"strings"
	"syscall"
	"time"

	"net/http"
	_ "net/http/pprof"

	"github.com/codegangsta/cli"
	cliprocessor "github.com/wanelo/image-server/cli"
	"github.com/wanelo/image-server/core"
	fetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/logger"
	"github.com/wanelo/image-server/logger/graphite"
	"github.com/wanelo/image-server/paths"
	"github.com/wanelo/image-server/server"
	"github.com/wanelo/image-server/uploader"
)

func main() {

	app := cli.NewApp()
	app.Name = "images"
	app.Version = core.VERSION
	app.Usage = "Image server and CLI"
	app.Action = func(c *cli.Context) {
		println("Need to provide subcommand: server or process")
	}

	app.Flags = globalFlags()

	app.Commands = []cli.Command{
		{
			Name:      "server",
			ShortName: "s",
			Usage:     "image server",
			Action: func(c *cli.Context) {
				setGoMaxProcs(c.GlobalInt("gomaxprocs"))
				go handleShutdownSignals()

				if c.GlobalBool("profile") {
					go initializePprofServer()
				}

				sc, err := serverConfiguration(c)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				go initializeUploader(sc)

				port := c.GlobalString("port")
				server.InitializeServer(sc, c.GlobalString("listen"), port)
			},
		},
		{
			Name:      "process",
			ShortName: "p",
			Usage:     "process image dimensions",
			Action: func(c *cli.Context) {
				sc, err := serverConfiguration(c)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				// initializeUploader(sc)
				outputsStr := c.GlobalString("outputs")
				if outputsStr == "" {
					log.Println("Need to specify outputs: 'x300.jpg,x300.webp'")
					os.Exit(1)
				}

				// input := bufio.NewReader(os.Stdin)
				namespace := c.GlobalString("namespace")
				outputs := strings.Split(outputsStr, ",")
				path := c.Args().First()
				if path == "" {
					log.Println("Need to pass an image path ARG")
					os.Exit(1)
				}

				err = cliprocessor.Process(sc, namespace, outputs, path)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

			},
		},
		{
			Name:      "process_stream",
			ShortName: "ps",
			Usage:     "process image dimensions",
			Action: func(c *cli.Context) {
				sc, err := serverConfiguration(c)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

				// initializeUploader(sc)
				outputsStr := c.GlobalString("outputs")
				if outputsStr == "" {
					log.Println("Need to specify outputs: 'x300.jpg,x300.webp'")
					os.Exit(1)
				}

				// input := bufio.NewReader(os.Stdin)
				namespace := c.GlobalString("namespace")
				outputs := strings.Split(outputsStr, ",")
				err = cliprocessor.ProcessStream(sc, namespace, outputs, os.Stdin)
				if err != nil {
					log.Println(err)
					os.Exit(1)
				}

			},
		},
	}

	app.Run(os.Args)
}

// globalFlags returns flags. If the flags are not present, it will try
// extracting values from the environment, otherwise it will use default values
func globalFlags() []cli.Flag {
	return []cli.Flag{
		// HTTP Server settings
		cli.StringFlag{Name: "port", Value: "7000", Usage: "Specifies the server port."},
		cli.StringFlag{Name: "extensions", Value: "jpg,gif,webp", Usage: "Whitelisted extensions (separated by commas)"},
		cli.StringFlag{Name: "local_base_path", Value: "public", Usage: "Directory where the images will be saved"},

		// Uploader paths
		cli.StringFlag{Name: "remote_base_url", Value: "", Usage: "Source domain for images"},
		cli.StringFlag{Name: "remote_base_path", Value: "", Usage: "base path for cloud storage"},

		// For CLI
		cli.StringFlag{Name: "namespace", Value: "", Usage: "Namespace"},
		cli.StringFlag{Name: "outputs", Value: "", Usage: "Output files with dimension and compression: 'x300.jpg,x300.webp'"},
		cli.StringFlag{Name: "listen", Value: "127.0.0.1", Usage: "IP address the server listens to"},

		// S3 uploader
		cli.StringFlag{Name: "aws_access_key_id", Value: "", Usage: "S3 Access Key"},
		cli.StringFlag{Name: "aws_secret_key", Value: "", Usage: "S3 Secret"},
		cli.StringFlag{Name: "aws_bucket", Value: "", Usage: "S3 Bucket"},

		// Manta uploader
		cli.StringFlag{Name: "manta_url", Value: "", Usage: "URL of Manta endpoint. https://us-east.manta.joyent.com"},
		cli.StringFlag{Name: "manta_user", Value: "", Usage: "The account name"},
		cli.StringFlag{Name: "manta_key_id", Value: "", Usage: "The fingerprint of the account or user SSH public key. Example: $(ssh-keygen -l -f $HOME/.ssh/id_rsa.pub | awk '{print $2}')"},
		cli.StringFlag{Name: "sdc_identity", Value: "", Usage: "Example: $HOME/.ssh/id_rsa"},

		// Default image settings
		cli.IntFlag{Name: "maximum_width", Value: 1000, Usage: "Maximum image width"},
		cli.IntFlag{Name: "default_quality", Value: 75, Usage: "Default image compression quality"},

		// Settings
		cli.IntFlag{Name: "uploader_concurrency", Value: 10, Usage: "Uploader concurrency"},
		cli.IntFlag{Name: "processor_concurrency", Value: 4, Usage: "Processor concurrency"},
		cli.IntFlag{Name: "http_timeout", Value: 5, Usage: "HTTP request timeout in seconds"},
		cli.IntFlag{Name: "gomaxprocs", Value: 0, Usage: "It will use the default when set to 0"},

		// Monitoring and Profiling
		cli.StringFlag{Name: "graphite_host", Value: "127.0.0.1", Usage: "Graphite host"},
		cli.IntFlag{Name: "graphite_port", Value: 8125, Usage: "Graphite port"},
		cli.BoolFlag{Name: "profile", Usage: "Enable pprof"},
	}
}

// initializeUploader creates base path on destination server
func initializeUploader(sc *core.ServerConfiguration) {
	err := uploader.Initialize(sc)
	if err != nil {
		log.Println("EXITING: Unable to initialize uploader: ", err)
		os.Exit(2)
	}
}

func serverConfiguration(c *cli.Context) (*core.ServerConfiguration, error) {
	sc := serverConfigurationFromContext(c)

	loggers := []core.Logger{
		graphite.New(sc.GraphiteHost, sc.GraphitePort),
	}

	adapters := &core.Adapters{
		Fetcher: &fetcher.Fetcher{},
		Paths:   &paths.Paths{LocalBasePath: sc.LocalBasePath, RemoteBasePath: sc.RemoteBasePath, RemoteBaseURL: sc.RemoteBaseURL},
		Logger:  &logger.Logger{Loggers: loggers},
	}
	sc.Adapters = adapters

	return sc, nil
}

// serverConfigurationFromContext returns a core.ServerConfiguration initialized
// from command line flags or defaults.
// Command line flags preceding the Command (server, process, etc) are registered
// as globals. Flags succeeding the Command are not globals.
func serverConfigurationFromContext(c *cli.Context) *core.ServerConfiguration {
	httpTimeout := time.Duration(c.GlobalInt("http_timeout")) * time.Second

	return &core.ServerConfiguration{
		WhitelistedExtensions: strings.Split(c.GlobalString("extensions"), ","),
		LocalBasePath:         c.GlobalString("local_base_path"),
		GraphitePort:          c.GlobalInt("graphite_port"),
		GraphiteHost:          c.GlobalString("graphite_host"),
		MaximumWidth:          c.GlobalInt("maximum_width"),
		RemoteBasePath:        c.GlobalString("remote_base_path"),
		RemoteBaseURL:         c.GlobalString("remote_base_url"),

		AWSAccessKeyID: c.GlobalString("aws_access_key_id"),
		AWSSecretKey:   c.GlobalString("aws_secret_key"),
		AWSBucket:      c.GlobalString("aws_bucket"),

		// Manta uploader
		MantaURL:    c.GlobalString("manta_url"),
		MantaUser:   c.GlobalString("manta_user"),
		MantaKeyID:  c.GlobalString("manta_key_id"),
		SDCIdentity: c.GlobalString("sdc_identity"),

		Outputs:             c.GlobalString("outputs"),
		DefaultQuality:      uint(c.GlobalInt("default_quality")),
		UploaderConcurrency: uint(c.GlobalInt("uploader_concurrency")),
		HTTPTimeout:         httpTimeout,
	}
}

func handleShutdownSignals() {
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT)

	<-shutdown
	server.ShuttingDown = true
	log.Println("Shutting down. Allowing requests to finish within 30 seconds. Interrupt again to quit immediately.")

	go func() {
		shutdown := make(chan os.Signal, 1)
		signal.Notify(shutdown, syscall.SIGHUP, syscall.SIGINT)

		<-shutdown
		log.Println("Forced to shutdown.")
		os.Exit(0)
	}()
}

func initializePprofServer() {
	log.Println(http.ListenAndServe("localhost:6060", nil))
}

func setGoMaxProcs(maxprocs int) {
	if maxprocs != 0 {
		runtime.GOMAXPROCS(maxprocs)
	}
}
