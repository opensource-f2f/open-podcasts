package advanced_ui

import (
	"bytes"
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	"github.com/linuxsuren/goplay/pkg/broadcast"
	playio "github.com/linuxsuren/goplay/pkg/io"
	"github.com/linuxsuren/goplay/pkg/ui"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

type TrackAudioPanel struct {
	ui.AudioPlayer

	audioUID string
}

func NewTrackAudioPanel(trackInfo *broadcast.TrackInfo) (panel *TrackAudioPanel, err error) {
	_ = loadConfig()

	var resp *http.Response
	if resp, err = http.Get(trackInfo.PlayURL64); err == nil {
		data, _ := io.ReadAll(resp.Body)
		buffer := bytes.NewReader(data)

		var streamer beep.StreamSeekCloser
		var format beep.Format
		streamer, format, err = mp3.Decode(playio.SeekerWithoutCloser(buffer))
		if err != nil {
			return
		}

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

		return &TrackAudioPanel{
			AudioPlayer: ui.NewAudioPanel(format.SampleRate, streamer),
			audioUID:    strconv.Itoa(trackInfo.UID),
		}, nil
	}
	return
}

func loadConfig() (err error) {
	viper.SetConfigName("goplay")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("$HOME/.config")
	viper.AddConfigPath(".")
	if err = viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			// Config file not found; ignore error if desired
			err = nil
		} else {
			err = fmt.Errorf("failed to load config: %s, error: %v", os.ExpandEnv("$HOME/.config/goplay.yaml"), err)
		}
	}
	return
}

// Start starts to play and try to seek the last position
func (t *TrackAudioPanel) Start() error {
	t.Play()
	if position := viper.GetInt(t.audioUID); position > 0 {
		return t.Seek(position)
	}
	return nil
}

// Stop stops the audio player
func (t *TrackAudioPanel) Stop() {
	viper.Set(t.audioUID, t.Position())
	_ = viper.WriteConfigAs(os.ExpandEnv("$HOME/.config/goplay.yaml")) // TODO print a warning log
}
