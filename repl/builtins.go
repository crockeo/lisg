package repl

import "fmt"

var (
)

func getBuiltin(name string) (LisgFunction, bool) {
	builtin, ok := map[string]LisgFunction{
		"eval":   eval,
		"defvar": defvar,
		"let":    let,
		"car":    car,
		"cdr":    cdr,
		"+":      add,
		"*":      mul,
		"-":      sub,
		"/":      div,
	}[name]
	return builtin, ok
}

func eval(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	if list, ok := value.(LisgList); ok {
		if len(list.children) == 0 {
			// nil -> nil
			return list, nil
		}

		headSymbol, ok := list.children[0].(LisgSymbol)
		if !ok {
			return nil, fmt.Errorf("list headed by uncallable %s", list.children[0])
		}

		if builtin, ok := getBuiltin(headSymbol.value); ok {
			return builtin(global, local, LisgList{children:list.children[1:]})
		}

		_, err := eval(global, local, headSymbol)
		if err != nil {
			return nil, err
		}

		// TODO: use this when we call functions
		// evaluatedBody := LisgList{children: make([]LisgValue, len(list.children) - 1)}
		// for i, child := range list.children[1:] {
		// 	evaluatedChild, err := eval(global, local, child)
		// 	if err != nil {
		// 		return nil, err
		// 	}
		// 	evaluatedBody.children[i] = evaluatedChild
		// }

		// TODO: try to convert headValue to a function once that exists


		return nil, nil
	} else if symbol, ok := value.(LisgSymbol); ok {
		return local.GetValue(symbol)
	} else {
		return value, nil
	}
}

func defvar(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("defvar passed value %s that is not an argument list", value)
	}

	if len(args.children) != 2 {
		return nil, fmt.Errorf("can only accept tuple, not %s", args)
	}

	symbol, ok := args.children[0].(LisgSymbol)
	if !ok {
		return nil, fmt.Errorf("cannot redefine non-symbol %s", args.children[0])
	}

	assignment, err := eval(global, local, args.children[1])
	if err != nil {
		return nil, err
	}

	newValue, err := global.SetValue(symbol, assignment)
	if err != nil {
		return nil, err
	}

	return newValue, nil
}

func let(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	return nil, fmt.Errorf("let unimplemented")
}

func car(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	return nil, fmt.Errorf("car unimplemented")
}

func cdr(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	return nil, fmt.Errorf("cdr unimplemented")
}

func add(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("add passed value %s that is not an argument list", value)
	}

	var acc float64 = 0
	for _, value := range args.children {
		evalValue, err := eval(global, local, value)
		if err != nil {
			return nil, err
		}

		number, ok := evalValue.(LisgNumber)
		if !ok {
			return nil, fmt.Errorf("%s is not a number", value)
		}
		acc += number.value
	}
	return LisgNumber{value: acc}, nil
}

func mul(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("mul passed value %s that is not an argument list", value)
	}

	var acc float64 = 1
	for _, value := range args.children {
		evalValue, err := eval(global, local, value)
		if err != nil {
			return nil, err
		}

		number, ok := evalValue.(LisgNumber)
		if !ok {
			return nil, fmt.Errorf("%s is not a number", value)
		}
		acc += number.value
	}
	return LisgNumber{value: acc}, nil
}

func sub(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("sub passed value %s that is not an argument list", value)
	}

	if len(args.children) < 1 {
		return LisgNumber{value: 0}, nil
	}

	evalAcc, err := eval(global, local, args.children[0])
	if err != nil {
		return nil, err
	}

	accValue, ok := evalAcc.(LisgNumber)
	if !ok {
		return nil, fmt.Errorf("%s is not a number", value)
	}

	acc := accValue.value
	for _, value := range args.children[1:] {
		evalValue, err := eval(global, local, value)
		if err != nil {
			return nil, err
		}

		number, ok := evalValue.(LisgNumber)
		if !ok {
			return nil, fmt.Errorf("%s is not a number", value)
		}
		acc -= number.value
	}

	return LisgNumber{value: acc}, nil
}

func div(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("div passed value %s that is not an argument list", value)
	}

	if len(args.children) < 1 {
		return LisgNumber{value: 0}, nil
	}

	evalAcc, err := eval(global, local, args.children[0])
	if err != nil {
		return nil, err
	}

	accValue, ok := evalAcc.(LisgNumber)
	if !ok {
		return nil, fmt.Errorf("%s is not a number", value)
	}

	acc := accValue.value
	for _, value := range args.children[1:] {
		evalValue, err := eval(global, local, value)
		if err != nil {
			return nil, err
		}

		number, ok := evalValue.(LisgNumber)
		if !ok {
			return nil, fmt.Errorf("%s is not a number", value)
		}
		acc /= number.value
	}

	return LisgNumber{value: acc}, nil
}
