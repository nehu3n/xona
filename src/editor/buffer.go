package editor

type Buffer struct {
	content [][]rune
	cursorX int
	cursorY int
}

func NewBuffer() *Buffer {
	return &Buffer{
		content: [][]rune{{}},
		cursorX: 0,
		cursorY: 0,
	}
}

func (b *Buffer) Insert(char rune) {
	b.content[b.cursorY] = append(b.content[b.cursorY][:b.cursorX], append([]rune{char}, b.content[b.cursorY][b.cursorX:]...)...)
	b.cursorX++
}

func (b *Buffer) InsertString(text string) {
	for _, char := range text {
		if char == '\n' {
			b.InsertNewline()
		} else {
			b.Insert(char)
		}
	}
}

func (b *Buffer) Delete() {
	if b.cursorX > 0 {
		b.content[b.cursorY] = append(b.content[b.cursorY][:b.cursorX-1], b.content[b.cursorY][b.cursorX:]...)
		b.cursorX--
	} else if b.cursorY > 0 {
		b.cursorX = len(b.content[b.cursorY-1])
		b.content[b.cursorY-1] = append(b.content[b.cursorY-1], b.content[b.cursorY]...)
		b.content = append(b.content[:b.cursorY], b.content[b.cursorY+1:]...)
		b.cursorY--
	}
}

func (b *Buffer) MoveCursor(offsetX, offsetY int) {
	newCursorX := b.cursorX + offsetX
	newCursorY := b.cursorY + offsetY

	if newCursorY < 0 {
		newCursorY = 0
	} else if newCursorY >= len(b.content) {
		newCursorY = len(b.content) - 1
	}

	if newCursorX < 0 {
		if newCursorY > 0 {
			newCursorY--
			newCursorX = len(b.content[newCursorY])
		} else {
			newCursorX = 0
		}
	} else if newCursorX > len(b.content[newCursorY]) {
		if newCursorY < len(b.content)-1 {
			newCursorX = 0
			newCursorY++
		} else {
			newCursorX = len(b.content[newCursorY])
		}
	}

	b.cursorX = newCursorX
	b.cursorY = newCursorY
}

func (b *Buffer) GetContent() [][]rune {
	return b.content
}

func (b *Buffer) GetCursor() (int, int) {
	return b.cursorX, b.cursorY
}

func (b *Buffer) InsertNewline() {
	b.content = append(b.content[:b.cursorY+1], append([][]rune{{}}, b.content[b.cursorY+1:]...)...)
	b.content[b.cursorY+1] = append([]rune{}, b.content[b.cursorY][b.cursorX:]...)
	b.content[b.cursorY] = b.content[b.cursorY][:b.cursorX]
	b.cursorX = 0
	b.cursorY++
}
