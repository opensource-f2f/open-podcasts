package main

import (
	"github.com/opensource-f2f/open-podcasts/pkg/rss"
	"github.com/spf13/cobra"
	"os"
	"path/filepath"
)

func main() {
	cmd := &cobra.Command{
		Use: "yaml-rss",
		RunE: func(cmd *cobra.Command, args []string) (err error) {
			var showItemFiles []string
			if showItemFiles, err = findDefaultItems(); err != nil {
				return
			}

			var content string
			showFilePath := "show.yaml"
			if content, err = rss.GenerateRSSFromNonCRDFiles("", showFilePath, showItemFiles); err != nil {
				return
			}
			cmd.Println(content)
			return
		},
	}
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func findDefaultItems() (files []string, err error) {
	files, err = filepath.Glob("item-*.")
	return
}
