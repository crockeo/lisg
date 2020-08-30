package repl

import "fmt"

func getBuiltin(name string) (LisgBuiltin, bool) {
	builtin, ok := map[string]LisgBuiltin{
		"eval":   eval,
		"defvar": defvar,
		"lambda": lambda,
		"funcall": funcall,
		// TODO: write these shorthands?
		// "defun": defun,
		// "apply": apply,
		"+": add,
		"*": mul,
		"-": sub,
		"/": div,
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
			return builtin(global, local, LisgList{children: list.children[1:]})
		}

		return funcall(global, local, list)
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

func lambda(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("lambda passed value %s that is not an argument list", value)
	}

	if len(args.children) < 2 {
		return nil, fmt.Errorf("insufficient args passed to lambda: %s", args)
	}

	argList, ok := args.children[0].(LisgList)
	if !ok {
		return nil, fmt.Errorf("") // TODO
	}

	fnArgs := make([]LisgSymbol, len(argList.children))
	for i, arg := range argList.children {
		symbol, ok := arg.(LisgSymbol)
		if !ok {
			return nil, fmt.Errorf("asdf") // TODO
		}
		fnArgs[i] = symbol
	}

	fnBody := make([]LisgValue, len(args.children)-1)
	for i, body := range args.children[1:] {
		fnBody[i] = body
	}

	return LisgFunction{
		args: fnArgs,
		body: fnBody,
	}, nil
}

func funcall(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	args, ok := value.(LisgList)
	if !ok {
		return nil, fmt.Errorf("funcall passed value %s that is not an argument list", value)
	}

	if len(args.children) < 1 {
		return nil, fmt.Errorf("insufficient args passed to funcall: %s", args)
	}

	head, err := eval(global, local, args.children[0])
	if err != nil {
		return nil, err
	}
	fn, ok := head.(LisgFunction)
	if !ok {
		return nil, fmt.Errorf("calling uncallable: %s", args.children[0])
	}

	if len(args.children) - 1 != len(fn.args) {
		return nil, fmt.Errorf("arg count mismatch, %d != %d", len(args.children) - 1, len(fn.args))
	}

	bindings := map[LisgSymbol]LisgValue{}
	for i := range args.children[1:] {
		bindings[fn.args[i]] = args.children[i + 1]
	}

	fnContext, err := local.PushContext(bindings)
	if err != nil {
		return nil, err
	}

	var retVal LisgValue = LisgList{children: []LisgValue{}}
	for _, value := range fn.body {
		retVal, err = eval(global, fnContext, value)
		if err != nil {
			return nil, err
		}
	}

	return retVal, nil
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
