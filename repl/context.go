package repl

import "fmt"

type LisgContext struct {
	bindings  map[LisgSymbol]LisgValue
	parentCtx *LisgContext
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
	if _, ok := c.bindings[symbol]; ok {
		c.bindings[symbol] = value
		return value, nil
	}

	return nil, fmt.Errorf("setting non-existant symbol %s", symbol)
}
