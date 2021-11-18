package ui

import "github.com/gdamore/tcell"

// AudioPlayer represents an audio player
type AudioPlayer interface {
	Play()
	Seek(int) error
	Position() int

	Draw(screen tcell.Screen)
	Handle(event tcell.Event) (changed, quit bool)
}
