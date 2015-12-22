package main

import (
	"log"
	"os"
	"strings"

	cliprocessor "github.com/image-server/image-server/cli"
)

var cmdCli = &Command{
	UsageLine: "cli",
}

func runCli(cmd *Command, args []string) {
	sc, err := serverConfiguration()
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	// initializeUploader(sc)
	outputsStr := sc.Outputs
	if outputsStr == "" {
		log.Println("Need to specify outputs: 'x300.jpg,x300.webp'")
		os.Exit(1)
	}

	// input := bufio.NewReader(os.Stdin)
	namespace := config.namespace
	outputs := strings.Split(outputsStr, ",")
	path := args[0]
	if path != "" {
		err = cliprocessor.Process(sc, namespace, outputs, path)
	} else {
		err = cliprocessor.ProcessStream(sc, namespace, outputs, os.Stdin)
	}
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
}

func init() {
	cmdCli.Run = runCli
}
