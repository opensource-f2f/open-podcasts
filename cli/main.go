package main

import (
	"github.com/linuxsuren/open-podcasts/cli/cmd"
)

func main() {
	_ = cmd.NewPlayCommand().Execute()
}
