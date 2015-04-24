package main

import (
	"fmt"

	"github.com/wanelo/image-server/core"
)

var cmdVersion = &Command{
	Run:       runVersion,
	UsageLine: "version",
	Short:     "print images version",
	Long:      `Version prints the images version.`,
}

func runVersion(cmd *Command, args []string) {
	if len(args) != 0 {
		cmd.Usage()
	}
	fmt.Printf("images version %s\n", core.VERSION)
}
