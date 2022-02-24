module github.com/linuxsuren/open-podcasts/cli

go 1.16

require (
	github.com/SlyMarbo/rss v1.0.1
	github.com/faiface/beep v1.1.0
	github.com/gdamore/tcell v1.3.0
	github.com/ghodss/yaml v1.0.0
	github.com/linuxsuren/cobra-extension v0.0.11
	github.com/linuxsuren/http-downloader v0.0.50
	github.com/spf13/cobra v1.2.1
	github.com/spf13/viper v1.9.0
	github.com/stretchr/testify v1.7.0
)

replace github.com/SlyMarbo/rss v1.0.1 => github.com/LinuxSuRen/rss v1.0.2-0.20211120161457-3f8efe372d7a
