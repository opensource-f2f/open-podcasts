package main

import (
	"github.com/opensource-f2f/open-podcasts/cli/cmd"
)

func main() {
	_ = cmd.NewPlayCommand().Execute()
}
