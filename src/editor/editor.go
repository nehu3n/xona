package editor

import (
	"log"
	"strconv"

	"github.com/gdamore/tcell/v2"
)

func Editor() {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	buffer := NewBuffer()

	quit := func() {
		screen.Fini()
		log.Println("Xona quited")
	}

	draw := func() {
		screen.Clear()
		for y, line := range buffer.GetContent() {
			lineNumber := strconv.Itoa(y + 1)
			for x, r := range lineNumber {
				screen.SetContent(x, y, r, nil, tcell.StyleDefault)
			}
			for x, r := range line {
				screen.SetContent(len(lineNumber)+1+x, y, r, nil, tcell.StyleDefault)
			}
		}
		cursorX, cursorY := buffer.GetCursor()
		screen.ShowCursor(len(strconv.Itoa(cursorY+1))+1+cursorX, cursorY)
		screen.Show()
	}

	for {
		draw()
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			switch ev.Key() {
			case tcell.KeyCtrlC:
				quit()
				return
			case tcell.KeyEnter:
				buffer.InsertNewline()
			case tcell.KeyRune:
				buffer.Insert(ev.Rune())
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				buffer.Delete()
			case tcell.KeyLeft:
				buffer.MoveCursor(-1, 0)
			case tcell.KeyRight:
				buffer.MoveCursor(1, 0)
			case tcell.KeyUp:
				buffer.MoveCursor(0, -1)
			case tcell.KeyDown:
				buffer.MoveCursor(0, 1)
			}
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}
