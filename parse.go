package main

func isSpace(c byte) bool {
	return c == ' ' || c == '\r' || c == '\n'
}

func lex(s string) []string {
	tokens := []string{}

	startPos := -1
	escaped := false

	for i := range s {
		if s[i] == '(' {
			tokens = append(tokens, "(")
			continue
		}

		if s[i] == ')' {
			if startPos >= 0 {
				tokens = append(tokens, s[startPos:i])
				startPos = -1
			}
			tokens = append(tokens, ")")
			continue
		}

		if !isSpace(s[i]) && startPos < 0 {
			startPos = i
			continue
		}

		if startPos >= 0 {
			if s[startPos] == '"' {
				if s[i] == '\\' {
					escaped = true
					continue
				} else if s[i] == '"' && !escaped {
					tokens = append(tokens, s[startPos:i+1])
					startPos = -1
					continue
				}
			} else if isSpace(s[i]) {
				tokens = append(tokens, s[startPos:i])
				startPos = -1
				continue
			}
		}

		escaped = false
	}

	return tokens
}

func parse(ss []string) {
}
