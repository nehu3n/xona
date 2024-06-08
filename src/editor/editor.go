package editor

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"time"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

const (
	NOTIFICATION_TYPE_SUCCESS = "success"
	NOTIFICATION_TYPE_ERROR   = "error"
	NOTIFICATION_TYPE_WARN    = "warn"
	NOTIFICATION_TYPE_INFO    = "info"
)

func Editor(filePath string) {
	screen, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("Failed to create screen: %v", err)
	}
	if err := screen.Init(); err != nil {
		log.Fatalf("Failed to initialize screen: %v", err)
	}
	defer screen.Fini()

	screen.EnableMouse()
	var mouseScrollActive bool = true

	var quitAppConfirmation bool = false
	var unsavedChanges bool = false

	buffer := NewBuffer()

	var notificationMessage string
	var notificationType string
	var notificationEnd time.Time

	if filePath != "" {
		file, err := os.Open(filePath)
		if err != nil {
			log.Fatalf("Failed to open file: %v", err)
		}
		defer file.Close()

		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			buffer.InsertString(scanner.Text())
			buffer.InsertNewline()
		}

		if err := scanner.Err(); err != nil {
			log.Fatalf("Failed to read file: %v", err)
		}
	}

	quit := func() {
		screen.Fini()
		log.Println("Xona quited")
	}

	showNotification := func(message string, _type string) {
		notificationMessage = message
		notificationType = _type
		notificationEnd = time.Now().Add(1 * time.Second)
	}

	copyLineInClipboard := func() {
		_, cursorY := buffer.GetCursor()
		line := buffer.content[cursorY]
		clipboard.WriteAll(string(line))

		showNotification("Line copied to clipboard", NOTIFICATION_TYPE_SUCCESS)
	}

	saveFile := func() {
		if filePath == "" {
			return
		}
		file, err := os.Create(filePath)
		if err != nil {
			log.Fatalf("Failed to save file: %v", err)
		}
		defer file.Close()

		writer := bufio.NewWriter(file)
		for _, line := range buffer.GetContent() {
			writer.WriteString(string(line) + "\n")
		}
		writer.Flush()
		unsavedChanges = false
		showNotification("File saved", NOTIFICATION_TYPE_SUCCESS)
	}

	_ = unsavedChanges

	draw := func() {
		screen.Clear()

		screenWidth, screenHeight := screen.Size()

		cursorX, cursorY := buffer.GetCursor()
		if !mouseScrollActive {
			if cursorY < buffer.viewTop {
				buffer.viewTop = cursorY
			} else if cursorY >= buffer.viewTop+screenHeight {
				buffer.viewTop = cursorY - screenHeight + 1
			}
		}

		for y := 0; y < screenHeight; y++ {
			lineIdx := buffer.viewTop + y
			if lineIdx >= len(buffer.content) {
				break
			}
			line := buffer.content[lineIdx]
			lineNumber := strconv.Itoa(lineIdx + 1)

			for x, r := range lineNumber {
				screen.SetContent(x, y, r, nil, tcell.StyleDefault)
			}

			for x, r := range line {
				screen.SetContent(len(lineNumber)+1+x, y, r, nil, tcell.StyleDefault)
			}
		}

		screen.ShowCursor(len(strconv.Itoa(cursorY+1))+1+cursorX, cursorY-buffer.viewTop)

		if time.Now().Before(notificationEnd) {
			msgX := screenWidth - len(notificationMessage)
			msgY := screenHeight - 1
			for i, r := range notificationMessage {
				var color tcell.Color
				if notificationType == NOTIFICATION_TYPE_SUCCESS {
					color = tcell.ColorGreen
				} else if notificationType == NOTIFICATION_TYPE_ERROR {
					color = tcell.ColorRed
				} else if notificationType == NOTIFICATION_TYPE_WARN {
					color = tcell.ColorYellow
				} else if notificationType == NOTIFICATION_TYPE_INFO {
					color = tcell.ColorBlue
				}

				screen.SetContent(msgX+i, msgY, r, nil, tcell.StyleDefault.Foreground(color))
			}
		}

		screen.Show()
	}

	adjustViewTop := func() {
		cursorY, _ := buffer.GetCursor()
		_, screenHeight := screen.Size()
		if cursorY < buffer.viewTop {
			buffer.viewTop = cursorY
		} else if cursorY >= buffer.viewTop+screenHeight {
			buffer.viewTop = cursorY - screenHeight + 1
		} else if cursorY >= buffer.viewTop+screenHeight-1 {
			buffer.viewTop++
		}
	}

	handleKey := func(key *tcell.EventKey) {
		mouseScrollActive = false

		switch key.Key() {
		case tcell.KeyCtrlQ:
			quitAppConfirmation = true
			/*
				if unsavedChanges {
					saveFile()
				}
			*/
			quit()
			return
		case tcell.KeyCtrlS:
			if unsavedChanges {
				saveFile()
			}
		case tcell.KeyEnter:
			buffer.InsertNewline()
			unsavedChanges = true
			adjustViewTop()
		case tcell.KeyRune:
			buffer.Insert(key.Rune())
			unsavedChanges = true
			adjustViewTop()
		case tcell.KeyBackspace, tcell.KeyBackspace2:
			buffer.Delete()
			unsavedChanges = true
			adjustViewTop()
		case tcell.KeyLeft:
			buffer.MoveCursor(-1, 0)
			adjustViewTop()
		case tcell.KeyRight:
			buffer.MoveCursor(1, 0)
			adjustViewTop()
		case tcell.KeyUp:
			buffer.MoveCursor(0, -1)
			adjustViewTop()
		case tcell.KeyDown:
			buffer.MoveCursor(0, 1)
			adjustViewTop()
		case tcell.KeyCtrlL:
			copyLineInClipboard()
		}
	}

	handleMouse := func(mouse *tcell.EventMouse) {
		_, screenHeight := screen.Size()
		switch mouse.Buttons() {
		case tcell.Button1:
			mouseScrollActive = false
			mouseX, mouseY := mouse.Position()
			newCursorY := buffer.viewTop + mouseY
			if newCursorY < len(buffer.content) {
				lineNumberWidth := len(strconv.Itoa(newCursorY+1)) + 1
				if mouseX >= lineNumberWidth {
					newCursorX := mouseX - lineNumberWidth
					buffer.SetCursor(newCursorX, newCursorY)
				} else {
					buffer.SetCursor(0, newCursorY)
				}
				adjustViewTop()
			}
		case tcell.WheelDown:
			mouseScrollActive = true
			if buffer.viewTop < len(buffer.content)-screenHeight {
				buffer.viewTop++
			}
		case tcell.WheelUp:
			mouseScrollActive = true
			if buffer.viewTop > 0 {
				buffer.viewTop--
			}
		}
	}

	for {
		if quitAppConfirmation {
			return
		}

		draw()
		ev := screen.PollEvent()
		switch ev := ev.(type) {
		case *tcell.EventKey:
			handleKey(ev)
		case *tcell.EventMouse:
			handleMouse(ev)
		case *tcell.EventResize:
			screen.Sync()
		}
	}
}
