package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	FibCode = `
(define fib (n)
  (let loop ((n n)
             (a 1)
             (b 2))
    (loop (- n 1)
          b
          (+ a b))))`

	StringCode = `
(define prepend-hello (s)
  (string-append "hello " s))`

	EscapedString = "\"hello \\\" world\""
)

func TestLexFibCode(t *testing.T) {
	value := lex(FibCode)
	assert.Equal(
		t,
		[]string{
			"(", "define", "fib", "(",
			"n", ")", "(", "let",
			"loop", "(", "(", "n",
			"n", ")", "(", "a",
			"1", ")", "(", "b",
			"2", ")", ")", "(",
			"loop", "(", "-", "n",
			"1", ")", "b", "(",
			"+", "a", "b", ")",
			")", ")", ")",
		},
		value,
	)
}

func TestLexStringStuff(t *testing.T) {
	value := lex(StringCode)
	assert.Equal(
		t,
		[]string{
			"(", "define", "prepend-hello", "(",
			"s", ")", "(", "string-append",
			"\"hello \"", "s", ")", ")",
		},
		value,
	)
}

func TestEscapedString(t *testing.T) {
	value := lex(EscapedString)
	assert.Equal(
		t,
		[]string{
			"\"hello \\\" world\"",
		},
		value,
	)
}

// func TestParseSymbol(t *testing.T) {
// 	value, err := parse("asdf", 0, 4)
// 	require.NoError(t, err)

// 	assert.Equal(
// 		t,
// 		value,
// 		&LisgSymbol{
// 			value: "asdf",
// 		},
// 	)
// }

// func TestParseString(t *testing.T) {

// }

// func TestParseNumber(t *testing.T) {
// }

// func TestParseList(t *testing.T) {
// }
