package repl

type LisgValue interface {}

type LisgSymbol struct {
	value string
}

type LisgString struct {
	value string
}

type LisgNumber struct {
	value float64
}

type LisgList struct {
	children []LisgValue
}

type LisgFunction func(global, local *LisgContext, value LisgValue) (LisgValue, error)
