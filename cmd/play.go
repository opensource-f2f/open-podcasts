package cmd

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gdamore/tcell"
	"github.com/linuxsuren/goplay/pkg/advanced_ui"
	"github.com/linuxsuren/goplay/pkg/rss"
	exec2 "github.com/linuxsuren/http-downloader/pkg/exec"
	"github.com/spf13/cobra"
	"os/exec"
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
	if len(titles) == 0 {
		err = fmt.Errorf("no podcast found by keyword: %s", keyword)
		return
	}
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
	if len(titles) == 0 {
		err = fmt.Errorf("no episode found by keyword: %s", keyword)
		return
	}
	selector = &survey.Select{
		Message: "Select an episode",
		Options: titles,
	}
	if err = survey.AskOne(selector, &choose); err != nil {
		return
	}

	episode := episodes[choose]
	if !isSupport(episode) {
		var ok bool
		if ok, err = playWithPotentialTools(episode); !ok || err != nil {
			err = fmt.Errorf("currently, only support mp3")
		}
	} else {
		err = play(episode)
	}
	return
}

func isSupport(episode rss.Episode) bool {
	return episode.Type == "audio/mpeg"
}

func playWithPotentialTools(episode rss.Episode) (ok bool, err error) {
	var mplayer string
	if mplayer, err = exec.LookPath("mplayer"); err == nil {
		selector := &survey.Confirm{
			Message: fmt.Sprintf("Do you want to play '%s' with mplayer?", episode.Title),
		}

		err = survey.AskOne(selector, &ok)
		if err == nil && ok {
			var audioCacheFile string
			if audioCacheFile, err = advanced_ui.LoadAudioFile(episode.AudioLink); err == nil {
				err = exec2.RunCommand(mplayer, audioCacheFile)
			}
		}
	}
	return
}

func play(episode rss.Episode) (err error) {
	if _, err = advanced_ui.LoadAudioFile(episode.AudioLink); err != nil {
		err = fmt.Errorf("failed to load audio file from: %s, error: %v", episode.AudioLink, err)
		return
	}

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
		err = fmt.Errorf("failed to play '%s', error: %v", episode.Title, err)
		return
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
	return
}
