package mappath

const (
	EscapeRune    = '\\'
	PathSeparator = '.'
)

func CurrentKey(path string) (curr string, left string) {
	if len(path) == 0 {
		return "", ""
	}

	var (
		result  = make([]rune, 0, len(path))
		escaped = false
	)

	for i, ch := range path {
		// append escaped rune as a regular rune
		if escaped {
			result = append(result, ch)
			escaped = false
			continue
		}

		if ch == EscapeRune {
			escaped = true
			continue
		}

		// unescaped dot found
		// break cycle and return found node name
		if ch == PathSeparator {
			if i+1 < len(path) {
				return string(result), path[i+1:]
			}
			return string(result), ""
		}

		// just regular rune
		result = append(result, ch)
	}

	// if the last rune was a backslash, append it as a regular rune
	if escaped {
		result = append(result, EscapeRune)
	}

	return string(result), ""
}
