package main

import (
	"flag"
	"fmt"
	"log"
	_ "net/http/pprof"
	"os"
	"strings"
	"time"

	"github.com/golang/glog"
	"github.com/wanelo/image-server/core"
	fetcher "github.com/wanelo/image-server/fetcher/http"
	"github.com/wanelo/image-server/logger"
	"github.com/wanelo/image-server/logger/graphite"
	"github.com/wanelo/image-server/paths"
	"github.com/wanelo/image-server/uploader"
)

// configT collects all the global state of the logging setup.
type configT struct {
	port          string
	extensions    string
	localBasePath string

	remoteBaseURL  string
	remoteBasePath string

	namespace string
	outputs   string
	listen    string

	awsAccessKeyID string
	awsSecretKey   string
	awsBucket      string

	mantaURL    string
	mantaUser   string
	mantaKeyID  string
	sdcIdentity string

	maximumWidth   int
	defaultQuality int

	uploaderConcurrency  int
	processorConcurrency int
	httpTimeout          int
	gomaxprocs           int

	graphiteHost string
	graphitePort int
	profile      bool
}

var config configT

var commands = []*Command{
	cmdServer,
	cmdCli,
}

// A Command is an implementation of a images command
// like images server or images cli.
type Command struct {
	// Run runs the command.
	// The args are the arguments after the command name.
	Run func(cmd *Command, args []string)

	// UsageLine is the one-line usage message.
	// The first word in the line is taken to be the command name.
	UsageLine string

	// Short is the short description shown in the 'go help' output.
	Short string

	// Long is the long message shown in the 'go help <this-command>' output.
	Long string

	// Flag is a set of flags specific to this command.
	Flag flag.FlagSet

	// CustomFlags indicates that the command will do its own
	// flag parsing.
	CustomFlags bool
}

// Name returns the command's name: the first word in the usage line.
func (c *Command) Name() string {
	name := c.UsageLine
	i := strings.Index(name, " ")
	if i >= 0 {
		name = name[:i]
	}
	return name
}

func (c *Command) Usage() {
	fmt.Fprintf(os.Stderr, "usage: %s\n\n", c.UsageLine)
	fmt.Fprintf(os.Stderr, "%s\n", strings.TrimSpace(c.Long))
	os.Exit(2)
}

func main() {
	registerFlags()
	flag.Parse()
	args := flag.Args()

	defer glog.Flush()

	for _, cmd := range commands {
		if cmd.Name() == args[0] && cmd.Run != nil {
			cmd.Flag.Usage = func() { cmd.Usage() }
			if cmd.CustomFlags {
				args = args[1:]
			} else {
				cmd.Flag.Parse(args[1:])
				args = cmd.Flag.Args()
			}
			cmd.Run(cmd, args)
			exit()
			return
		}
	}
}

func exit() {
	os.Exit(0)
}

// globalFlags returns flags. If the flags are not present, it will try
// extracting values from the environment, otherwise it will use default values
func registerFlags() {

	// HTTP Server settings
	flag.StringVar(&config.port, "port", "7000", "Specifies the server port.")
	flag.StringVar(&config.extensions, "extensions", "jpg,gif,webp", "Whitelisted extensions (separated by commas)")
	flag.StringVar(&config.localBasePath, "local_base_path", "public", "Directory where the images will be saved")

	// Uploader paths
	flag.StringVar(&config.remoteBaseURL, "remote_base_url", "", "Source domain for images")
	flag.StringVar(&config.remoteBasePath, "remote_base_path", "", "base path for cloud storage")

	// For CLI
	flag.StringVar(&config.namespace, "namespace", "", "Namespace")
	flag.StringVar(&config.outputs, "outputs", "", "Output files with dimension and compression: 'x300.jpg,x300.webp'")
	flag.StringVar(&config.listen, "listen", "127.0.0.1", "IP address the server listens to")

	// S3 uploader
	flag.StringVar(&config.awsAccessKeyID, "aws_access_key_id", "", "S3 Access Key")
	flag.StringVar(&config.awsSecretKey, "aws_secret_key", "", "S3 Secret")
	flag.StringVar(&config.awsBucket, "aws_bucket", "", "S3 Bucket")

	// Manta uploader
	flag.StringVar(&config.mantaURL, "manta_url", "", "URL of Manta endpoint. https://us-east.manta.joyent.com")
	flag.StringVar(&config.mantaUser, "manta_user", "", "The account name")
	flag.StringVar(&config.mantaKeyID, "manta_key_id", "", "The fingerprint of the account or user SSH public key. Example: $(ssh-keygen -l -f $HOME/.ssh/id_rsa.pub | awk '{print $2}')")
	flag.StringVar(&config.sdcIdentity, "sdc_identity", "", "Example: $HOME/.ssh/id_rsa")

	// Default image settings
	flag.IntVar(&config.maximumWidth, "maximum_width", 1000, "Maximum image width")
	flag.IntVar(&config.defaultQuality, "default_quality", 75, "Default image compression quality")

	// Settings
	flag.IntVar(&config.uploaderConcurrency, "uploader_concurrency", 10, "Uploader concurrency")
	flag.IntVar(&config.processorConcurrency, "processor_concurrency", 4, "Processor concurrency")
	flag.IntVar(&config.httpTimeout, "http_timeout", 5, "HTTP request timeout in seconds")
	flag.IntVar(&config.gomaxprocs, "gomaxprocs", 0, "It will use the default when set to 0")

	// Monitoring and Profiling
	flag.StringVar(&config.graphiteHost, "graphite_host", "127.0.0.1", "Graphite host")
	flag.IntVar(&config.graphitePort, "graphite_port", 8125, "Graphite port")
	flag.BoolVar(&config.profile, "profile", false, "Enable pprof")
}

// initializeUploader creates base path on destination server
func initializeUploader(sc *core.ServerConfiguration) {
	err := uploader.Initialize(sc)
	if err != nil {
		log.Println("EXITING: Unable to initialize uploader: ", err)
		os.Exit(2)
	}
}

func serverConfiguration() (*core.ServerConfiguration, error) {
	sc := serverConfigurationFromConfig()

	log.Println(sc)

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
func serverConfigurationFromConfig() *core.ServerConfiguration {
	httpTimeout := time.Duration(config.httpTimeout) * time.Second

	return &core.ServerConfiguration{
		WhitelistedExtensions: strings.Split(config.extensions, ","),
		LocalBasePath:         config.localBasePath,
		GraphitePort:          config.graphitePort,
		GraphiteHost:          config.graphiteHost,
		MaximumWidth:          config.maximumWidth,
		RemoteBasePath:        config.remoteBasePath,
		RemoteBaseURL:         config.remoteBaseURL,

		AWSAccessKeyID: config.awsAccessKeyID,
		AWSSecretKey:   config.awsSecretKey,
		AWSBucket:      config.awsBucket,

		// Manta uploader
		MantaURL:    config.mantaURL,
		MantaUser:   config.mantaUser,
		MantaKeyID:  config.mantaKeyID,
		SDCIdentity: config.sdcIdentity,

		Outputs:             config.outputs,
		DefaultQuality:      uint(config.defaultQuality),
		UploaderConcurrency: uint(config.uploaderConcurrency),
		HTTPTimeout:         httpTimeout,
	}
}
