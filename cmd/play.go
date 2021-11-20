package cmd

import (
	"github.com/AlecAivazis/survey/v2"
	"github.com/gdamore/tcell"
	"github.com/linuxsuren/goplay/pkg/advanced_ui"
	"github.com/linuxsuren/goplay/pkg/rss"
	"github.com/spf13/cobra"
	"time"
)

func NewPlayCommand() (cmd *cobra.Command) {
	opt := &playOption{}

	cmd = &cobra.Command{
		Use:     "goplay",
		Example: "goplay opensource",
		Short:   "Play podcast",
		Args:    cobra.MinimumNArgs(1),
		RunE:    opt.runE,
	}
	return
}

type playOption struct {
}

func (o *playOption) runE(cmd *cobra.Command, args []string) (err error) {
	keyword := args[0]

	feeds, titles := rss.FindMapByKeyword(keyword)
	selector := &survey.Select{
		Message: "Select a podcast",
		Options: titles,
	}
	var choose string
	if err = survey.AskOne(selector, &choose); err != nil {
		return
	}

	feed := feeds[choose]
	var episodes map[string]rss.Episode
	episodes, titles = rss.ConvertEpisodeMap(feed.Items)
	selector = &survey.Select{
		Message: "Select an episode",
		Options: titles,
	}
	if err = survey.AskOne(selector, &choose); err != nil {
		return
	}

	play(episodes[choose])
	return
}

func play(episode rss.Episode) {
	screen, err := tcell.NewScreen()
	if err != nil {
		panic(err)
	}
	err = screen.Init()
	if err != nil {
		panic(err)
	}
	defer screen.Fini()

	var player *advanced_ui.TrackAudioPanel
	if player, err = advanced_ui.NewTrackAudioPanel(episode); err != nil {
		panic(err)
	}

	screen.Clear()
	player.Draw(screen)
	screen.Show()

	_ = player.Start()

	seconds := time.Tick(time.Second)
	events := make(chan tcell.Event)
	go func() {
		for {
			events <- screen.PollEvent()
		}
	}()

loop:
	for {
		select {
		case event := <-events:
			changed, quit := player.Handle(event)
			if quit {
				player.Stop()
				break loop
			}
			if changed {
				screen.Clear()
				player.Draw(screen)
				screen.Show()
			}
		case <-seconds:
			screen.Clear()
			player.Draw(screen)
			screen.Show()
		}
	}
}

