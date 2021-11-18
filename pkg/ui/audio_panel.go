package ui

import (
	"fmt"
	"github.com/faiface/beep"
	"github.com/faiface/beep/effects"
	"github.com/faiface/beep/speaker"
	"github.com/gdamore/tcell"
	"os"
	"time"
	"unicode"
)

type audioPanel struct {
	sampleRate beep.SampleRate
	streamer   beep.StreamSeeker
	ctrl       *beep.Ctrl
	resampler  *beep.Resampler
	volume     *effects.Volume
}

// NewAudioPanel creates a audio panel
func NewAudioPanel(sampleRate beep.SampleRate, streamer beep.StreamSeeker) AudioPlayer {
	ctrl := &beep.Ctrl{Streamer: beep.Loop(-1, streamer)}
	resampler := beep.ResampleRatio(4, 1, ctrl)
	volume := &effects.Volume{Streamer: resampler, Base: 2}
	return &audioPanel{sampleRate, streamer, ctrl, resampler, volume}
}

// Play starts to play a audio
func (ap *audioPanel) Play() {
	speaker.Play(ap.volume)
}

// Seek seeks the track
func (ap *audioPanel) Seek(position int) error {
	return ap.streamer.Seek(position)
}

// Position returns the position
func (ap *audioPanel) Position() int {
	return ap.streamer.Position()
}

// Draw draws the infomation to a screen
func (ap *audioPanel) Draw(screen tcell.Screen) {
	mainStyle := tcell.StyleDefault.
		Background(tcell.NewHexColor(0x473437)).
		Foreground(tcell.NewHexColor(0xD7D8A2))
	statusStyle := mainStyle.
		Foreground(tcell.NewHexColor(0xDDC074)).
		Bold(true)

	screen.Fill(' ', mainStyle)

	drawTextLine(screen, 0, 0, "Welcome to the Speedy Player!", mainStyle)
	drawTextLine(screen, 0, 1, "Press [ESC] to quit.", mainStyle)
	drawTextLine(screen, 0, 2, "Press [SPACE] to pause/resume.", mainStyle)
	drawTextLine(screen, 0, 3, "Use keys in (?/?) to turn the buttons.", mainStyle)

	speaker.Lock()
	position := ap.sampleRate.D(ap.streamer.Position())
	length := ap.sampleRate.D(ap.streamer.Len())
	volume := ap.volume.Volume
	speed := ap.resampler.Ratio()
	speaker.Unlock()

	positionStatus := fmt.Sprintf("%v / %v", position.Round(time.Second), length.Round(time.Second))
	volumeStatus := fmt.Sprintf("%.1f", volume)
	speedStatus := fmt.Sprintf("%.3fx", speed)

	drawTextLine(screen, 0, 5, "Position (Q/W):", mainStyle)
	drawTextLine(screen, 16, 5, positionStatus, statusStyle)

	drawTextLine(screen, 0, 6, "Volume   (A/S):", mainStyle)
	drawTextLine(screen, 16, 6, volumeStatus, statusStyle)

	drawTextLine(screen, 0, 7, "Speed    (Z/X):", mainStyle)
	drawTextLine(screen, 16, 7, speedStatus, statusStyle)
}

// Handle handles the input events
func (ap *audioPanel) Handle(event tcell.Event) (changed, quit bool) {
	switch event := event.(type) {
	case *tcell.EventKey:
		if event.Key() == tcell.KeyESC {
			return false, true
		}

		if event.Key() != tcell.KeyRune {
			return false, false
		}

		switch unicode.ToLower(event.Rune()) {
		case ' ':
			speaker.Lock()
			ap.ctrl.Paused = !ap.ctrl.Paused
			speaker.Unlock()
			return false, false

		case 'q', 'w':
			speaker.Lock()
			newPos := ap.streamer.Position()
			if event.Rune() == 'q' {
				newPos -= ap.sampleRate.N(time.Second)
			}
			if event.Rune() == 'w' {
				newPos += ap.sampleRate.N(time.Second)
			}
			if newPos < 0 {
				newPos = 0
			}
			if newPos >= ap.streamer.Len() {
				newPos = ap.streamer.Len() - 1
			}
			if err := ap.streamer.Seek(newPos); err != nil {
				report(err)
			}
			speaker.Unlock()
			return true, false

		case 'a':
			speaker.Lock()
			ap.volume.Volume -= 0.1
			speaker.Unlock()
			return true, false

		case 's':
			speaker.Lock()
			ap.volume.Volume += 0.1
			speaker.Unlock()
			return true, false

		case 'z':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 15 / 16)
			speaker.Unlock()
			return true, false

		case 'x':
			speaker.Lock()
			ap.resampler.SetRatio(ap.resampler.Ratio() * 16 / 15)
			speaker.Unlock()
			return true, false
		}
	}
	return false, false
}

func drawTextLine(screen tcell.Screen, x, y int, s string, style tcell.Style) {
	for _, r := range s {
		screen.SetContent(x, y, r, nil, style)
		x++
	}
}

func report(err error) {
	fmt.Fprintln(os.Stderr, err)
	os.Exit(1)
}
