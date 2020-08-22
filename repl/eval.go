package repl

import "fmt"

type LisgBuiltin func(args LisgList) (LisgValue, error)

var (
	builtins = map[string]LisgBuiltin{
		"+": add,
	}
)

func add(args LisgList) (LisgValue, error) {
	var acc float64
	for _, value := range args.children[1:] {
		number, ok := value.(LisgNumber)
		if !ok {
			return nil, fmt.Errorf("invalid type, not LisgNumber: %s", value)
		}
		acc += number.value
	}
	return LisgNumber{value: acc}, nil
}

func evalList(list LisgList) (LisgValue, error) {
	if len(list.children) == 0 {
		// this is the nil type for the language, so it's atomic
		return list, nil
	}

	head, ok := list.children[0].(LisgSymbol)
	if !ok {
		// TODO raise
	}

	if builtin, ok := builtins[head.value]; ok {
		return builtin(list)
	} else {
		// TODO raise
	}

	return nil, nil
}

func Eval(v LisgValue) (LisgValue, error) {
	if list, ok := v.(LisgList); ok {
		return evalList(list)
	} else {
		return v, nil
	}
}
