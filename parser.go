package cfp

// Parse returns either Font object or nil if any error occur
func Parse(s string) *Font {
	font := New(s)

	if font.error != nil {
		return nil
	}

	return font
}
