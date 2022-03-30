package cfp

import (
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	stateVariation        = iota // 0
	stateLineHeight              // 1
	stateFontFamily              // 2
	stateBeforeFontFamily        // 3
	stateAfterOblique            // 4
)

var (
	obliqueRegexp = regexp.MustCompile(`^(?:\+|-)?(?:[0-9]*\.)?[0-9]+(?:deg|grad|rad|turn)$`)

	variantSizeTShortRegexp = regexp.MustCompile(`^(?:(?:xx|x)-large|(?:xx|s)-small|small|large|medium)$`)
	variantCompareRegexp    = regexp.MustCompile(`^(?:larg|small)er$`)
	variantSizeRegexp       = regexp.MustCompile(`^(?:\+|-)?(?:[0-9]*\.)?[0-9]+(?:em|ex|ch|rem|vh|vw|vmin|vmax|px|mm|cm|in|pt|pc|%)$`)
	variantSmallCapsRegexp  = regexp.MustCompile(`^small-caps$`)

	styleItalicRegexp  = regexp.MustCompile(`^italic$`)
	styleObliqueRegexp = regexp.MustCompile(`^oblique$`)

	weightRegexp     = regexp.MustCompile(`^(?:bold(?:er)?|lighter)$`)
	weightSizeRegexp = regexp.MustCompile(`^[+-]?(?:[0-9]*\.)?[0-9]+(?:e[+-]?(?:0|[1-9][0-9]*))?$`)

	stretchRegexp = regexp.MustCompile(`^(?:(?:ultra|extra|semi)-)?(?:condensed|expanded)$`)

	lineHeightRegexp = regexp.MustCompile(`^(?:\+|-)?([0-9]*\.)?[0-9]+(?:em|ex|ch|rem|vh|vw|vmin|vmax|px|mm|cm|in|pt|pc|%)?$`)

	onlyWhiteSpacesRegexp  = regexp.MustCompile(`^\s*$`)
	moreThanOneSpaceRegexp = regexp.MustCompile(`\s+`)

	identifierRegexp     = regexp.MustCompile(`^(?:-?\d|--)`)
	noneIdentifierRegexp = regexp.MustCompile(`^(?:[_a-zA-Z0-9-]|[^\0-\237]|(?:\\[0-9a-f]{1,6}(?:\r\n|[ \n\r\t\f])?|\\[^\n\r\f0-9a-f]))+$`)
)

// Font represents a CSS font
type Font struct {
	Family     []string `json:"font-family,omitempty"`
	Size       string   `json:"font-size,omitempty"`
	Style      string   `json:"font-style,omitempty"`
	Variant    string   `json:"font-variant,omitempty"`
	Weight     string   `json:"font-weight,omitempty"`
	Stretch    string   `json:"font-stretch,omitempty"`
	LineHeight string   `json:"line-height,omitempty"`

	error  error
	raw    string
	runes  []rune
	state  int
	itr    int
	buffer string
}

// New creates a new font from the string input
func New(s string) *Font {
	font := &Font{raw: s}

	font.build()

	return font
}

// Error returns an error from parsing
func (f *Font) Error() error {
	return f.error
}

func (f *Font) build() {
	f.runes = []rune(f.raw)
	f.state = stateVariation
	f.Family = make([]string, 0)

	for ; f.itr < len(f.runes); f.itr++ {
		switch f.state {
		case stateBeforeFontFamily:
			err := f.beforeFontFamily()
			if err != nil {
				f.error = err
				return
			}
		case stateFontFamily:
			_ = f.fontFamily()
		case stateAfterOblique:
			_ = f.afterOblique()
		case stateVariation:
			_ = f.variation()
		case stateLineHeight:
			_ = f.lineHeight()
		}
	}

	// This is for the case where a string was specified followed by
	// an identifier, but without a separating comma.
	if f.state == stateFontFamily && !onlyWhiteSpacesRegexp.MatchString(f.buffer) {
		f.error = errors.New("incorrect format of font-family")
		return
	}

	if f.state == stateBeforeFontFamily {
		identifier := f.parseIdentifier()
		if identifier != "" {
			f.Family = append(f.Family, identifier)
		}
	}

	if f.Size != "" && len(f.Family) > 0 {
		return
	}

	f.error = errors.New("wasn't able to parse Font")
}

func (f *Font) char() rune {
	return f.runes[f.itr]
}

func (f *Font) isQuoteChar() bool {
	char := f.char()

	return char == '"' || char == '\''
}

func (f *Font) isCommaChar() bool {
	return f.char() == ','
}

func (f *Font) isWhiteSpaceChar() bool {
	return f.char() == ' '
}

func (f *Font) isSlashChar() bool {
	return f.char() == '/'
}

func (f *Font) clearBuffer() {
	f.buffer = ""
}

func (f *Font) saveToBuffer() {
	f.buffer = fmt.Sprintf("%s%c", f.buffer, f.char())
}

func (f *Font) beforeFontFamily() error {
	switch {
	case f.isQuoteChar():
		char := f.char()
		index := f.itr + 1

		ind := strings.IndexRune(f.raw[index:], char)
		if ind == -1 {
			return errors.New("unclosed quote")
		}
		index += ind + 1
		for f.raw[index-2] == '\\' {
			ind = strings.IndexRune(f.raw[index:], char)
			if ind == -1 {
				return errors.New("unclosed quote")
			}
			index += ind + 1
		}

		f.Family = append(f.Family, f.raw[f.itr:index])

		f.itr = index - 1
		f.state = stateFontFamily
		f.clearBuffer()
	case f.isCommaChar():
		identifier := f.parseIdentifier()
		if identifier != "" {
			f.Family = append(f.Family, identifier)
		}
		f.clearBuffer()
	default:
		f.saveToBuffer()
	}

	return nil
}

func (f *Font) fontFamily() error {
	switch {
	case f.isCommaChar():
		f.state = stateBeforeFontFamily
		f.clearBuffer()
	default:
		f.saveToBuffer()
	}

	return nil
}

func (f *Font) afterOblique() error {
	switch {
	case f.isWhiteSpaceChar():
		if obliqueRegexp.MatchString(f.buffer) {
			f.Style = fmt.Sprintf("%s %s", f.Style, f.buffer)
			f.clearBuffer()
		} else {
			// The 'oblique' token was not followed by an angle.
			// Backtrack to allow the token to be parsed as VARIATION
			f.itr--
		}
		f.state = stateVariation
	default:
		f.saveToBuffer()
	}

	return nil
}

func (f *Font) variation() error {
	switch {
	case f.isWhiteSpaceChar() || f.isSlashChar():
		switch {
		case variantSizeTShortRegexp.MatchString(f.buffer) || variantCompareRegexp.MatchString(f.buffer) || variantSizeRegexp.MatchString(f.buffer):
			if f.char() == '/' {
				f.state = stateLineHeight
			} else {
				f.state = stateBeforeFontFamily
			}
			f.Size = f.buffer
		case styleItalicRegexp.MatchString(f.buffer):
			f.Style = f.buffer
		case styleObliqueRegexp.MatchString(f.buffer):
			f.Style = f.buffer
			f.state = stateAfterOblique
		case variantSmallCapsRegexp.MatchString(f.buffer):
			f.Variant = f.buffer
		case weightRegexp.MatchString(f.buffer):
			f.Weight = f.buffer
		case weightSizeRegexp.MatchString(f.buffer):
			n, _ := strconv.ParseFloat(f.buffer, 64)
			if n >= 1 && n <= 1000 {
				f.Weight = f.buffer
			}
		case stretchRegexp.MatchString(f.buffer):
			f.Stretch = f.buffer
		}

		f.clearBuffer()
	default:
		f.saveToBuffer()
	}

	return nil
}

func (f *Font) lineHeight() error {
	switch {
	case f.isWhiteSpaceChar():
		if lineHeightRegexp.MatchString(f.buffer) {
			f.LineHeight = f.buffer
		}
		f.state = stateBeforeFontFamily
		f.clearBuffer()
	default:
		f.saveToBuffer()
	}
	return nil
}

func (f *Font) parseIdentifier() string {
	parts := strings.Split(
		moreThanOneSpaceRegexp.ReplaceAllString(strings.TrimSpace(f.buffer), " "),
		" ",
	)

	for _, part := range parts {
		if identifierRegexp.MatchString(part) || !noneIdentifierRegexp.MatchString(part) {
			return ""
		}
	}

	return strings.Join(parts, " ")
}
