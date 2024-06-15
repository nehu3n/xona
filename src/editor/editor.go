package editor

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"time"
	"xona/src/editor/highlight"

	"github.com/atotto/clipboard"
	"github.com/gdamore/tcell/v2"
)

const (
	NOTIFICATION_TYPE_SUCCESS = "success"
	NOTIFICATION_TYPE_ERROR   = "error"
	NOTIFICATION_TYPE_WARN    = "warn"
	NOTIFICATION_TYPE_INFO    = "info"
)

var (
	searchMode    bool
	searchQuery   string
	searchResults []int
)

var extensionsHighlight map[string]string = map[string]string{
	".go":   "go.toml",
	".py":   "python.toml",
	".js":   "javascript.toml",
	".ts":   "typescript.toml",
	".rs":   "rust.toml",
	".rb":   "ruby.toml",
	".sh":   "bash.toml",
	".cs":   "csharp.toml",
	".md":   "markdown.toml",
	".php":  "php.toml",
	".html": "html.toml",
	".css":  "css.toml",
	".toml": "toml.toml",
	".json": "json.toml",
	".yaml": "yaml.toml",
	".xml":  "xml.toml",
	".sql":  "sql.toml",
}

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

	var highlightSearch bool = false

	buffer := NewBuffer()

	var notificationMessage string
	var notificationType string
	var notificationEnd time.Time

	patternsMap := highlight.LoadAllPatterns()
	var highlighter *highlight.Highlighter = highlight.NewHighlighter(patternsMap["txt.toml"])

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

		maxLineNumber := len(buffer.content)
		maxLineNumberLen := len(strconv.Itoa(maxLineNumber))

		for y := 0; y < screenHeight; y++ {
			lineIdx := buffer.viewTop + y
			if lineIdx >= len(buffer.content) {
				break
			}
			line := buffer.content[lineIdx]
			lineNumber := strconv.Itoa(lineIdx + 1)

			styles := highlighter.Highlight(string(line))

			lineNumberSpace := maxLineNumberLen + 1
			lineNumberXOffset := maxLineNumberLen - len(lineNumber)

			for x, r := range lineNumber {
				screen.SetContent(lineNumberXOffset+x, y, r, nil, tcell.StyleDefault.Foreground(tcell.ColorGray).Background(tcell.ColorDefault))
			}

			highlightSearch = false
			for _, match := range searchResults {
				if match == lineIdx {
					highlightSearch = true
					break
				}
			}

			for x, r := range line {
				var style tcell.Style
				style = styles[x]

				if highlightSearch {
					style = tcell.StyleDefault.Foreground(tcell.ColorBlack).Background(tcell.ColorYellow)
				}

				screen.SetContent(lineNumberSpace+x, y, r, nil, style)
			}
		}

		screen.ShowCursor(maxLineNumberLen+1+cursorX, cursorY-buffer.viewTop)

		if searchMode {
			searchPrompt := "Search: "

			for i, r := range searchPrompt {
				screen.SetContent(screenWidth-len(searchPrompt)-len(searchQuery)+i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}

			for i, r := range searchQuery {
				screen.SetContent(screenWidth-len(searchQuery)+i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
		}

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

	askForNewFilePath := func(screen tcell.Screen) string {
		screenWidth, _ := screen.Size()

		prompt := "Enter file path: "
		x := screenWidth - len(prompt) - 1

		input := ""
		for {
			draw()
			for i, r := range prompt {
				screen.SetContent(x+i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
			for i, r := range input {
				screen.SetContent(x+len(prompt)+i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
			screen.Show()

			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlQ:
					return ""
				case tcell.KeyESC:
					return ""
				case tcell.KeyEnter:
					return input
				case tcell.KeyRune:
					input += string(ev.Rune())
				case tcell.KeyDelete:
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				}

				x = screenWidth - len(prompt) - len(input) - 1
				if x < 0 {
					x = 0
				}
			}
		}
	}

	confirmQuitWithoutSaving := func(screen tcell.Screen) string {
		screenWidth, _ := screen.Size()

		prompt := "Do you want to save the file before leaving? (y/n): "
		x := screenWidth - len(prompt) - 1

		input := ""
		for {
			draw()
			for i, r := range prompt {
				screen.SetContent(x+i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorWhite))
			}
			for i, r := range input {
				screen.SetContent(x+len(prompt)+i, 0, r, nil, tcell.StyleDefault.Foreground(tcell.ColorYellow))
			}
			screen.Show()

			ev := screen.PollEvent()
			switch ev := ev.(type) {
			case *tcell.EventKey:
				switch ev.Key() {
				case tcell.KeyCtrlQ:
					return ""
				case tcell.KeyESC:
					return ""
				case tcell.KeyEnter:
					return input
				case tcell.KeyRune:
					if len(input) != 3 {
						input += string(ev.Rune())
					}
				case tcell.KeyDelete:
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				case tcell.KeyBackspace, tcell.KeyBackspace2:
					if len(input) > 0 {
						input = input[:len(input)-1]
					}
				}

				x = screenWidth - len(prompt) - len(input) - 1
				if x < 0 {
					x = 0
				}
			}
		}
	}

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

		ext := filepath.Ext(filePath)

		tomlFileLang, ok := extensionsHighlight[ext]
		if ok {
			highlighter = highlight.NewHighlighter(patternsMap[tomlFileLang])
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
			newFilePath := askForNewFilePath(screen)

			if newFilePath == "" {
				return
			} else {
				filePath = newFilePath
			}
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

	performSearch := func(query string) []int {
		var results []int
		for i, line := range buffer.GetContent() {
			matched, _ := regexp.MatchString(query, string(line))
			if matched {
				results = append(results, i)
			}
		}
		return results
	}

	handleKey := func(key *tcell.EventKey) {
		mouseScrollActive = false

		if searchMode {
			switch key.Key() {
			case tcell.KeyCtrlQ:
				searchMode = false
				searchQuery = ""
				searchResults = nil

				if highlightSearch {
					highlightSearch = false
				}
			case tcell.KeyEscape:
				searchMode = false
				searchQuery = ""
				searchResults = nil

				if highlightSearch {
					highlightSearch = false
				}
			case tcell.KeyEnter:
				searchMode = false
				searchResults = performSearch(searchQuery)
				if len(searchResults) > 0 {
					buffer.SetCursor(0, searchResults[0])
					adjustViewTop()
				}

				searchQuery = ""
			case tcell.KeyBackspace, tcell.KeyBackspace2:
				if len(searchQuery) > 0 {
					searchQuery = searchQuery[:len(searchQuery)-1]
				}
			case tcell.KeyRune:
				searchQuery += string(key.Rune())
			case tcell.KeyDelete:
				searchQuery = string(searchQuery[len(searchQuery)-1])
			}
			return
		}

		switch key.Key() {
		case tcell.KeyCtrlQ:
			quitAppConfirmation = true

			confirmReady := false

			for !confirmReady {
				if unsavedChanges {
					confirm := confirmQuitWithoutSaving(screen)

					if confirm == "y" || confirm == "yes" || confirm == "Y" || confirm == "YES" {
						saveFile()
						confirmReady = true
						quit()
						return
					} else if confirm == "n" || confirm == "no" || confirm == "N" || confirm == "NO" {
						confirmReady = true
						quit()
						return
					} else {
						showNotification("Invalid confirm answer", NOTIFICATION_TYPE_ERROR)
					}
				}
			}
		case tcell.KeyCtrlS:
			if unsavedChanges {
				saveFile()
			}
		case tcell.KeyCtrlF:
			searchMode = true
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
