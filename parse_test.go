package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
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

func TestLexEscapedString(t *testing.T) {
	value := lex(EscapedString)
	assert.Equal(
		t,
		[]string{
			"\"hello \\\" world\"",
		},
		value,
	)
}

func TestParseString(t *testing.T) {
	lisgValue, err := parse([]string{"\"hello\""})
	require.NoError(t, err)

	lisgString, ok := lisgValue.(LisgString)
	require.True(t, ok)
	assert.Equal(t, "hello", lisgString.value)
}

func TestParseNumbers(t *testing.T) {
	numberStrings := []string{"1234", "1234.5"}
	lisgNumbers := []LisgNumber{
		{1234},
		{1234.5},
	}

	for i, numberString := range numberStrings {
		lisgValue, err := parse([]string{numberString})
		require.NoError(t, err)

		lisgNumber, ok := lisgValue.(LisgNumber)
		require.True(t, ok)
		assert.Equal(t, lisgNumbers[i], lisgNumber)
	}
}

func TestParseSymbol(t *testing.T) {
	lisgValue, err := parse([]string{"symbol-name"})
	require.NoError(t, err)

	lisgSymbol, ok := lisgValue.(LisgSymbol)
	require.True(t, ok)
	assert.Equal(t, LisgSymbol{value: "symbol-name"}, lisgSymbol)
}

func TestParseSimpleList(t *testing.T) {
	lisgValue, err := parse([]string{
		"(",
		"string-append",
		"\"hello\"",
		"1234.5",
		")",
	})
	require.NoError(t, err)

	lisgList, ok := lisgValue.(LisgList)
	require.True(t, ok)

	assert.Equal(
		t,
		LisgSymbol{value: "string-append"},
		lisgList.children[0].(LisgSymbol),
	)
	assert.Equal(
		t,
		LisgString{value: "hello"},
		lisgList.children[1].(LisgString),
	)
	assert.Equal(
		t,
		LisgNumber{value: 1234.5},
		lisgList.children[2].(LisgNumber),
	)
}

func TestLexParseEmbedded(t *testing.T) {
	symbols := lex(StringCode)

	lisgValue, err := parse(symbols)
	require.NoError(t, err)

	lisgList, ok := lisgValue.(LisgList)
	require.True(t, ok)

	assert.Equal(
		t,
		LisgList{
			children: []LisgValue{
				LisgSymbol{value: "define"},
				LisgSymbol{value: "prepend-hello"},
				LisgList{
					children: []LisgValue{
						LisgSymbol{value: "s"},
					},
				},
				LisgList{
					children: []LisgValue{
						LisgSymbol{value: "string-append"},
						LisgString{value: "hello "},
						LisgSymbol{value: "s"},
					},
				},
			},
		},
		lisgList,
	)
}
