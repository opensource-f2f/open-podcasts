package advanced_ui

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
	"github.com/faiface/beep/speaker"
	playio "github.com/linuxsuren/goplay/pkg/io"
	"github.com/linuxsuren/goplay/pkg/rss"
	"github.com/linuxsuren/goplay/pkg/ui"
	"github.com/linuxsuren/goplay/pkg/util"
	"github.com/linuxsuren/http-downloader/pkg/net"
	"github.com/spf13/viper"
	"io"
	"os"
	"path"
	"time"
)

type TrackAudioPanel struct {
	ui.AudioPlayer

	audioUID string
}

func NewTrackAudioPanel(episode rss.Episode) (panel *TrackAudioPanel, err error) {
	_ = loadConfig()

	var buffer io.Reader
	if buffer, err = playWithLocalCache(episode.AudioLink); err == nil {
		var streamer beep.StreamSeekCloser
		var format beep.Format
		streamer, format, err = mp3.Decode(playio.SeekerWithoutCloser(buffer))
		if err != nil {
			return
		}

		speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/30))

		return &TrackAudioPanel{
			AudioPlayer: ui.NewAudioPanel(format.SampleRate, streamer, episode.Title),
			audioUID:    episode.UID,
		}, nil
	}
	return
}

func LoadAudioFile(rssURL string) (audioCacheFile string, err error) {
	cacheDir := os.ExpandEnv("$HOME/.config/goplay/cache")
	_ = os.MkdirAll(cacheDir, 0751)

	audioCacheFile = path.Join(cacheDir, fmt.Sprintf("%s.audio", util.HashCodeAsString(rssURL)))
	if _, err = os.Stat(audioCacheFile); err != nil {
		err = net.DownloadFileWithMultipleThread(rssURL, audioCacheFile, 4, true)
	}
	return
}

func saveOrGetCache(rssURL string) (reader io.Reader, err error) {
	var audioCacheFile string
	if audioCacheFile, err = LoadAudioFile(rssURL); err == nil {
		reader, err = os.Open(audioCacheFile)
	}
	return
}

func playWithLocalCache(trackURL string) (reader io.Reader, err error) {
	reader, err = saveOrGetCache(trackURL)
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
