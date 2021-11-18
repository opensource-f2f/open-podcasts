package main

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/gdamore/tcell"
	"github.com/linuxsuren/goplay/pkg/advanced_ui"
	"github.com/linuxsuren/goplay/pkg/broadcast"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: goplay keyword")
		return
	}
	keyword := os.Args[1]

	var albumID int
	if result, err := broadcast.Search(keyword); err != nil {
		panic(err)
	} else if len(result.Data.AlbumViews.Albums) == 0 {
		fmt.Println("not found any albums by keyword:", keyword)
		return
	} else {
		albumTitle := make([]string, 0)
		albumMap := make(map[string]int, 0)
		for _, v := range result.Data.AlbumViews.Albums {
			albumTitle = append(albumTitle, v.AlbumInfo.Title)
			albumMap[v.AlbumInfo.Title] = v.AlbumInfo.ID
		}

		selector := &survey.Select{
			Message: "Select a album",
			Options: albumTitle,
		}
		var choose string
		if err = survey.AskOne(selector, &choose); err != nil {
			return
		}
		albumID = albumMap[choose]
	}

	tracks, err := broadcast.GetTrackList(albumID, 1, false)
	if err != nil {
		panic(err)
	}

	trackTitle := make([]string, 0)
	trackMap := make(map[string]*broadcast.TrackInfo, 0)
	for _, v := range tracks.Data.List {
		trackTitle = append(trackTitle, v.Title)
		trackMap[v.Title] = v
	}

	selector := &survey.Select{
		Message: "Select a track",
		Options: trackTitle,
	}

	var choose string
	if err = survey.AskOne(selector, &choose); err == nil {
		fmt.Println("start to play", choose)
		trackInfo := trackMap[choose]
		play(trackInfo)
	}
}

func play(trackInfo *broadcast.TrackInfo) {
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
	if player, err = advanced_ui.NewTrackAudioPanel(trackInfo); err != nil {
		panic(err)
	}

	//streamer, format, err := mp3.Decode(reader)
	//if err != nil {
	//	panic(err)
	//}
	//defer streamer.Close()
	//
	//speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

	//ap := ui.NewAudioPanel(format.SampleRate, streamer)

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
