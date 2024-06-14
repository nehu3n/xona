package highlight

import (
	"regexp"
	"strconv"
	"strings"

	"github.com/gdamore/tcell/v2"
)

type Highlighter struct {
	patterns *SyntaxPatterns
}

func NewHighlighter(patterns *SyntaxPatterns) *Highlighter {
	return &Highlighter{patterns: patterns}
}

func (h *Highlighter) Highlight(line string) []tcell.Style {
	styles := make([]tcell.Style, len(line))
	for i := range styles {
		styles[i] = tcell.StyleDefault
	}

	applyStyle := func(pattern string, color tcell.Color) {
		re := regexp.MustCompile(pattern)
		matches := re.FindAllStringIndex(line, -1)
		for _, match := range matches {
			for i := match[0]; i < match[1]; i++ {
				styles[i] = styles[i].Foreground(color)
			}
		}
	}

	const COLOR_SPLIT_SEPARATOR string = ", "

	for _, pattern := range h.patterns.Keywords.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Keywords.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))

		applyStyle(`\b`+regexp.QuoteMeta(pattern)+`\b`, styleColor) // Purple
	}
	for _, pattern := range h.patterns.Types.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Types.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))

		applyStyle(`\b`+regexp.QuoteMeta(pattern)+`\b`, styleColor) // Purple
	}
	for _, pattern := range h.patterns.Numbers.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Numbers.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))
		applyStyle(pattern, styleColor) // Orange
	}
	for _, pattern := range h.patterns.Methods.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Methods.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))
		applyStyle(pattern, styleColor) // Blue
	}
	for _, pattern := range h.patterns.Operators.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Operators.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))
		applyStyle(pattern, styleColor) // Purple
	}
	for _, pattern := range h.patterns.Comments.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Comments.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))
		applyStyle(pattern, styleColor) // Gray
	}
	for _, pattern := range h.patterns.Strings.Patterns {
		var firstColorRGB int
		var secondColorRGB int
		var threeColorRGB int

		var styleColor tcell.Color

		for i, color := range strings.Split(h.patterns.Strings.Color, COLOR_SPLIT_SEPARATOR) {
			if i == 0 {
				firstColorRGB, _ = strconv.Atoi(color)
			} else if i == 1 {
				secondColorRGB, _ = strconv.Atoi(color)
			} else if i == 2 {
				threeColorRGB, _ = strconv.Atoi(color)
			}
		}

		styleColor = tcell.NewRGBColor(int32(firstColorRGB), int32(secondColorRGB), int32(threeColorRGB))
		applyStyle(pattern, styleColor) // Orange
	}

	return styles
}
