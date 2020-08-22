package repl

import "strconv"

func isSpace(c byte) bool {
	return c == ' ' || c == '\r' || c == '\n'
}

func Lex(s string) []string {
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

func Parse(ss []string) (LisgValue, error) {
	head := ss[0]

	if head[0] == '"' && head[len(head) - 1] == '"' {
		return LisgString{
			value: head[1:len(head) - 1],
		}, nil
	} else if value, err := strconv.ParseFloat(head, 64); err == nil {
		return LisgNumber{
			value: value,
		}, nil
	} else if head == "(" {
		var rangeStart int
		depth := 0

		ranges := [][]int{}
		for i, item := range ss {
			if item == "(" {
				depth += 1
				if depth == 2 {
					rangeStart = i
				}
			} else if item == ")" {
				depth -= 1
				if depth == 1 {
					ranges = append(ranges, []int{rangeStart, i + 1})
				}
			} else if depth == 1 {
				// top-level expressions in the list
				ranges = append(ranges, []int{i, i + 1})
			}
		}

		children := make([]LisgValue, len(ranges))
		for i, r := range ranges {
			child, err := Parse(ss[r[0]:r[1]])
			if err != nil {
				return nil, err
			}
			children[i] = child
		}
		return LisgList{children: children}, nil
	} else {
		return LisgSymbol{
			value: head,
		}, nil
	}
}
