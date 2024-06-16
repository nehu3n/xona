package editor

import (
	"log"

	"github.com/gdamore/tcell/v2"
)

type Window struct {
	buffer *Buffer
	screen tcell.Screen
}

func (e *Editor) openNewWindow() {
	if len(e.windows) == 4 {
		return
	}

	windowScreen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}

	e.windows = append(e.windows, &Window{buffer: NewBuffer(), screen: windowScreen})
	currentWindowIndex = len(e.windows) - 1

	if err := windowScreen.Init(); err != nil {
		log.Fatalf("Failed to initialize new window screen: %v", err)
	}
}
func (e *Editor) switchWindow() {
	numWindows := len(e.windows)
	if numWindows > 1 {
		currentWindowIndex = (currentWindowIndex + 1) % numWindows
	}

}
