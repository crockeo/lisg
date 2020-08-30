package repl

import "fmt"

type LisgContext struct {
	bindings  map[LisgSymbol]LisgValue
	parentCtx *LisgContext
}

func BaseContext() *LisgContext {
	return &LisgContext{
		bindings: map[LisgSymbol]LisgValue{},
		parentCtx: nil,
	}
}

// GetValue retrieves the value that a symbol maps onto, if it exists in the current or any parent
// contexts.
func (c *LisgContext) GetValue(symbol LisgSymbol) (LisgValue, error) {
	if value, ok := c.bindings[symbol]; ok {
		return value, nil
	}

	if c.parentCtx != nil {
		return c.parentCtx.GetValue(symbol)
	}

	return nil, fmt.Errorf("symbol has no value %s", symbol)
}

// SetValue sets a symbol's value, if it exists in only the current context.
func (c *LisgContext) SetValue(symbol LisgSymbol, value LisgValue) (LisgValue, error) {
	c.bindings[symbol] = value
	return value, nil
	// if _, ok := c.bindings[symbol]; ok {
	// 	c.bindings[symbol] = value
	// 	return value, nil
	// }

	// return nil, fmt.Errorf("setting non-existant symbol %s", symbol)
}

// PushContext pushes a new context onto the stack of contexts.
func (c *LisgContext) PushContext(bindings map[LisgSymbol]LisgValue) (*LisgContext, error) {
	return &LisgContext{
		bindings: bindings,
		parentCtx: c,
	}, nil
}

// PopContext pops a context from the stack of contexts.
func (c *LisgContext) PopContext() (*LisgContext, error) {
	if c.parentCtx == nil {
		return nil, fmt.Errorf("attempting to pop to nil context")
	}

	return c.parentCtx, nil
}
