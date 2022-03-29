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

func Parse(s string) map[string]any {
	state := StateVariation
	buffer := ""
	result := map[string]any{FontFamily: []string{}}

	rns := []rune(s)
	for i := 0; i < len(rns); i++ {
		char := rns[i]

		if state == StateBeforeFontFamily && (char == '"' || char == '\'') {
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

		} else if state == StateFontFamily && char == ',' {
			state = StateBeforeFontFamily
			buffer = ""
		} else if state == StateBeforeFontFamily && char == ',' {
			identifier := parseIdentifier(buffer)
			if identifier != "" {
				result[FontFamily] = append(result[FontFamily].([]string), identifier)
			}
			buffer = ""
		} else if state == StateAfterOblique && char == ' ' {
			if regexp.MustCompile(`^(?:\+|-)?(?:[0-9]*\.)?[0-9]+(?:deg|grad|rad|turn)$`).MatchString(buffer) {
				result[FontStyle] = fmt.Sprintf("%v %s", result[FontStyle], buffer)
				buffer = ""
			} else {
				// The 'oblique' token was not followed by an angle.
				// Backtrack to allow the token to be parsed as VARIATION
				i--
			}
			state = StateVariation
		} else if state == StateVariation && (char == ' ' || char == '/') {
			if regexp.MustCompile(`^(?:(?:xx|x)-large|(?:xx|s)-small|small|large|medium)$`).MatchString(buffer) ||
				regexp.MustCompile(`^(?:larg|small)er$`).MatchString(buffer) ||
				regexp.MustCompile(`^(?:\+|-)?(?:[0-9]*\.)?[0-9]+(?:em|ex|ch|rem|vh|vw|vmin|vmax|px|mm|cm|in|pt|pc|%)$`).MatchString(buffer) {
				if char == '/' {
					state = StateLineHeight
				} else {
					state = StateBeforeFontFamily
				}
				result[FontSize] = buffer
			} else if regexp.MustCompile(`^italic$`).MatchString(buffer) {
				result[FontStyle] = buffer
			} else if regexp.MustCompile(`^oblique$`).MatchString(buffer) {
				result[FontStyle] = buffer
				state = StateAfterOblique
			} else if regexp.MustCompile(`^small-caps$`).MatchString(buffer) {
				result[FontVariant] = buffer
			} else if regexp.MustCompile(`^(?:bold(?:er)?|lighter)$`).MatchString(buffer) {
				result[FontWeight] = buffer
			} else if regexp.MustCompile(`^[+-]?(?:[0-9]*\.)?[0-9]+(?:e[+-]?(?:0|[1-9][0-9]*))?$`).MatchString(buffer) {
				n, _ := strconv.ParseFloat(buffer, 64)
				if n >= 1 && n <= 1000 {
					result[FontWeight] = buffer
				}
			} else if regexp.MustCompile(`^(?:(?:ultra|extra|semi)-)?(?:condensed|expanded)$`).MatchString(buffer) {
				result[FontStretch] = buffer
			}

			buffer = ""
		} else if state == StateLineHeight && char == ' ' {
			if regexp.MustCompile(`^(?:\+|-)?([0-9]*\.)?[0-9]+(?:em|ex|ch|rem|vh|vw|vmin|vmax|px|mm|cm|in|pt|pc|%)?$`).MatchString(buffer) {
				result[LineHeight] = buffer
			}
			state = StateBeforeFontFamily
			buffer = ""
		} else {
			buffer = fmt.Sprintf("%s%c", buffer, char)
		}

	}

	// This is for the case where a string was specified followed by
	// an identifier, but without a separating comma.
	if state == StateFontFamily && !regexp.MustCompile(`^\s*$`).MatchString(buffer) {
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
		regexp.MustCompile(`\s+`).ReplaceAllString(strings.TrimSpace(s), " "),
		" ",
	)

	r1 := regexp.MustCompile(`^(?:-?\d|--)`)
	r2 := regexp.MustCompile(`^(?:[_a-zA-Z0-9-]|[^\0-\237]|(?:\\[0-9a-f]{1,6}(?:\r\n|[ \n\r\t\f])?|\\[^\n\r\f0-9a-f]))+$`)

	for _, part := range parts {
		if r1.MatchString(part) || !r2.MatchString(part) {
			return ""
		}
	}

	return strings.Join(parts, " ")
}
