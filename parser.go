package cfp

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

const (
	StateVariation        = iota // 0
	StateLineHeight              // 1
	StateFontFamily              // 2
	StateBeforeFontFamily        // 3
	StateAfterOblique            // 4

	FontFamily  = "font-family"
	FontSize    = "font-size"
	FontStyle   = "font-style"
	FontVariant = "font-variant"
	FontWeight  = "font-weight"
	FontStretch = "font-stretch"
	LineHeight  = "line-height"
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

func Parse(s string) map[string]any {
	state := StateVariation
	buffer := ""
	result := map[string]any{FontFamily: []string{}}

	rns := []rune(s)
	for i := 0; i < len(rns); i++ {
		char := rns[i]

		switch {
		case state == StateBeforeFontFamily && (char == '"' || char == '\''):
			index := i + 1

			ind := strings.IndexRune(s[index:], char)
			if ind == -1 {
				return nil
			}
			index += ind + 1
			for s[index-2] == '\\' {
				ind = strings.IndexRune(s[index:], char)
				if ind == -1 {
					return nil
				}
				index += ind + 1
			}

			result[FontFamily] = append(result[FontFamily].([]string), s[i:index])

			i = index - 1
			state = StateFontFamily
			buffer = ""
		case state == StateFontFamily && char == ',':
			state = StateBeforeFontFamily
			buffer = ""
		case state == StateBeforeFontFamily && char == ',':
			identifier := parseIdentifier(buffer)
			if identifier != "" {
				result[FontFamily] = append(result[FontFamily].([]string), identifier)
			}
			buffer = ""
		case state == StateAfterOblique && char == ' ':
			if obliqueRegexp.MatchString(buffer) {
				result[FontStyle] = fmt.Sprintf("%v %s", result[FontStyle], buffer)
				buffer = ""
			} else {
				// The 'oblique' token was not followed by an angle.
				// Backtrack to allow the token to be parsed as VARIATION
				i--
			}
			state = StateVariation
		case state == StateVariation && (char == ' ' || char == '/'):
			switch {
			case variantSizeTShortRegexp.MatchString(buffer) || variantCompareRegexp.MatchString(buffer) || variantSizeRegexp.MatchString(buffer):
				if char == '/' {
					state = StateLineHeight
				} else {
					state = StateBeforeFontFamily
				}
				result[FontSize] = buffer
			case styleItalicRegexp.MatchString(buffer):
				result[FontStyle] = buffer
			case styleObliqueRegexp.MatchString(buffer):
				result[FontStyle] = buffer
				state = StateAfterOblique
			case variantSmallCapsRegexp.MatchString(buffer):
				result[FontVariant] = buffer
			case weightRegexp.MatchString(buffer):
				result[FontWeight] = buffer
			case weightSizeRegexp.MatchString(buffer):
				n, _ := strconv.ParseFloat(buffer, 64)
				if n >= 1 && n <= 1000 {
					result[FontWeight] = buffer
				}
			case stretchRegexp.MatchString(buffer):
				result[FontStretch] = buffer
			}

			buffer = ""
		case state == StateLineHeight && char == ' ':
			if lineHeightRegexp.MatchString(buffer) {
				result[LineHeight] = buffer
			}
			state = StateBeforeFontFamily
			buffer = ""
		default:
			buffer = fmt.Sprintf("%s%c", buffer, char)
		}

	}

	// This is for the case where a string was specified followed by
	// an identifier, but without a separating comma.
	if state == StateFontFamily && !onlyWhiteSpacesRegexp.MatchString(buffer) {
		return nil
	}

	if state == StateBeforeFontFamily {
		identifier := parseIdentifier(buffer)
		if identifier != "" {
			result[FontFamily] = append(result[FontFamily].([]string), identifier)
		}
	}

	ff, ok := result[FontFamily].([]string)

	if result[FontSize] != nil && ok && len(ff) > 0 {
		return result
	}

	return nil
}

func parseIdentifier(s string) string {
	parts := strings.Split(
		moreThanOneSpaceRegexp.ReplaceAllString(strings.TrimSpace(s), " "),
		" ",
	)

	for _, part := range parts {
		if identifierRegexp.MatchString(part) || !noneIdentifierRegexp.MatchString(part) {
			return ""
		}
	}

	return strings.Join(parts, " ")
}
