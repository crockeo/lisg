package repl

func Eval(global, local *LisgContext, value LisgValue) (LisgValue, error) {
	return eval(global, local, value)
}
