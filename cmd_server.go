package main

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/wanelo/image-server/server"
)

var cmdServer = &Command{
	UsageLine: "server [flags]",
	Short:     "image server",
	Long: `
Flags:
    SERVER
    --listen '127.0.0.1'        IP address the server listens to
    --port '7000'               Specifies the server port.
    --local_base_path 'public'  Directory where the images will be saved
    --outputs                   i.e. 'full_size.jpg,full_size.webp,x110-q90.jpg,x200-q90.jpg'

  UPLOADERS:
  Images can be uploaded to Amazon S3 or Joyent's Manta if the following flags are provided.
    --remote_base_url   Source domain for images
    --remote_base_path  Base path for cloud storage

    S3
    --aws_access_key_id
    --aws_secret_key
    --aws_bucket

    MANTA
    --manta_url     URL of API. https://us-east.manta.joyent.com
    --manta_user    The account name
    --manta_key_id  The fingerprint of the account or user SSH public key. Example: $(ssh-keygen -l -f $HOME/.ssh/id_rsa.pub | awk '{print $2}')
    --sdc_identity  Example: $HOME/.ssh/id_rsa

  IMAGE CONFIGURATIONS
    --namespace                 Default namespace
    --extensions 'jpg,gif,webp' Whitelisted extensions (separated by commas)
    --maximum_width '1000'      Maximum image width
    --default_quality '75'      Default image compression quality

  SERVER TUNNING
    --uploader_concurrency '10' Uploader concurrency
    --processor_concurrency '4' Processor concurrency
    --http_timeout '5'          HTTP request timeout in seconds
    --gomaxprocs '0'            It will use the default when set to 0

  MONITORING & PROFILING
    --profile                   Enable pprof
    --graphite_host '127.0.0.1' Graphite host
    --graphite_port '8125'      Graphite port	`,
}

func init() {
	cmdServer.Run = runServer
}

func runServer(cmd *Command, args []string) {
	go handleShutdownSignals()

	if config.profile {
		go initializePprofServer()
	}

	sc, err := serverConfiguration()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	go initializeUploader(sc)

	port := config.port
	server.InitializeServer(sc, config.listen, port)
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
