package main

import (
	"github.com/opensource-f2f/open-podcasts/pkg/rss"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func main() {
	opt := &option{}
	cmd := &cobra.Command{
		Use:  "yaml-rss",
		RunE: opt.runE,
	}
	flags := cmd.Flags()
	flags.StringVarP(&opt.server, "server", "s", "",
		"The server of the RSS content")
	flags.StringVarP(&opt.showFile, "show-file", "", "show.yaml",
		"The show YAML file path")
	flags.StringVarP(&opt.itemsPattern, "items-pattern", "", "item-*.yaml",
		"The item files path pattern")
	cmd.SetOut(os.Stdout)
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func (o *option) runE(cmd *cobra.Command, args []string) (err error) {
	var showItemFiles []string
	if showItemFiles, err = filepath.Glob(o.itemsPattern); err != nil {
		return
	}

	var content string
	showFilePath := "show.yaml"
	if content, err = rss.GenerateRSSFromNonCRDFiles(o.server, showFilePath, showItemFiles); err != nil {
		return
	}
	cmd.Println(content)
	return
}

type option struct {
	server       string
	showFile     string
	itemsPattern string
}
