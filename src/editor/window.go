package editor

import "github.com/gdamore/tcell/v2"

type Window struct {
	buffer *Buffer
	screen tcell.Screen
}
