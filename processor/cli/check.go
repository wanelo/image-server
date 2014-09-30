package cli

import "os/exec"

var Available bool

func init() {
	Available = true
	cmd := exec.Command("convert", "--version")
	err := cmd.Run()
	if err != nil {
		Available = false
	}
}
