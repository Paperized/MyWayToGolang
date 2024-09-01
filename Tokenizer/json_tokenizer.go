package main

import (
	"errors"
	"fmt"
	"slices"
	"unicode"
)

type TokenType rune

const (
	CURLY_OPEN   = TokenType('{')
	CURLY_CLOSE  = TokenType('}')
	SQUARE_OPEN  = TokenType('[')
	SQUARE_CLOSE = TokenType(']')
	DOUBLEQUOTE  = TokenType('"')
	SPACE        = TokenType(' ')
	COLON        = TokenType(':')
	COMMA        = TokenType(',')
	BACKINLINE   = TokenType('\n')
	TRUE         = TokenType(0)
	FALSE        = TokenType(0)
	NULL         = TokenType(0)
	NUMBER       = TokenType(0)
	TEXT         = TokenType(0)
)

var allowedTokensAfterNumber []TokenType = []TokenType{CURLY_OPEN, CURLY_CLOSE, SQUARE_OPEN, SQUARE_CLOSE, SPACE, COLON, COMMA, BACKINLINE, DOUBLEQUOTE}

type Token struct {
	ttype   TokenType
	value   string
	startAt int
	endAt   int
}

func getCharacter(str string, index int) (rune, bool) {
	if index >= len(str) {
		return '-', false
	}

	return rune(str[index]), true
}

func safeSubstring(str string, from int, to int) string {
	if to > len(str) {
		to = len(str)
	}

	return str[from:to]
}

/*
token parole -> true, false, null | solo se dentro " -> anything-else**
token numerici -> num | num.num, solo uno 0 di fila all'inizio del numero intero



*/

func tokenizeJsonString(str string) ([]Token, error) {
	output := []Token{}

	strLen := len(str)

	for i := 0; i < strLen; i++ {
		isEOF := false
		c := rune(str[i])

		switch TokenType(c) {
		case CURLY_OPEN, CURLY_CLOSE, SQUARE_OPEN, SQUARE_CLOSE, SPACE, COLON, COMMA, BACKINLINE:
			output = append(output, Token{ttype: TokenType(c), value: string(c), startAt: i, endAt: i + 1})
		default:
			if c == 't' && safeSubstring(str, i, i+4) == "true" {
				output = append(output, Token{ttype: TRUE, value: str[i : i+4], startAt: i, endAt: i + 4})
				// we place ourselves at character 'e' so +3 (next iteration we will be moving in the new character +1)
				i += 3
				break
			}

			if c == 'f' && safeSubstring(str, i, i+5) == "false" {
				output = append(output, Token{ttype: FALSE, value: str[i : i+4], startAt: i, endAt: i + 5})
				// we place ourselves at character 'e' so +4 (next iteration we will be moving in the new character +1)
				i += 4
				break
			}

			if c == 'n' && safeSubstring(str, i, i+4) == "null" {
				output = append(output, Token{ttype: NULL, value: str[i : i+4], startAt: i, endAt: i + 4})
				// we place ourselves at character 'e' so +3 (next iteration we will be moving in the new character +1)
				i += 3
				break
			}

			if unicode.IsDigit(rune(c)) {
				var innerC rune
				innerC, isEOF = getCharacter(str, i+1)
				if isEOF {
					output = append(output, Token{ttype: NUMBER, value: string(c), startAt: i, endAt: i + 1})
					break
				}

				if c == '0' && unicode.IsDigit(innerC) && innerC == '0' {
					return output, errors.New(fmt.Sprintf("Unrecognized character: Numbers cannot start with two 0s at index %d: '%c'", i, str[i+1]))
				}

				var numberFound bool = false
				var foundNotAllowedAfterChar bool = false

				var j int
				for j = i + 1; j < strLen; j++ {
					innerC, isEOF = getCharacter(str, j)
					if isEOF {
						numberFound = true
						j--
						break
					}

					if innerC == '.' {
						break
					}

					if !unicode.IsDigit(innerC) {
						if j == i+1 && innerC == '-' {
							continue
						}

						numberFound = true
						foundNotAllowedAfterChar = !slices.Contains(allowedTokensAfterNumber, TokenType(innerC))
						break
					}
				}

				if numberFound {
					if foundNotAllowedAfterChar {
						return output, errors.New(fmt.Sprintf("Unrecognized character after a number at index %d: '%c'", i, str[j]))
					}

					output = append(output, Token{ttype: NUMBER, value: str[i : j+1], startAt: i, endAt: j + 1})
					i = j + 1
					break
				}

				for j = j + 1; j < strLen; j++ {
					innerC, isEOF = getCharacter(str, j)
					if isEOF {
						numberFound = true
						j--
						break
					}

					if !unicode.IsDigit(innerC) {
						numberFound = true
						foundNotAllowedAfterChar = !slices.Contains(allowedTokensAfterNumber, TokenType(innerC))
						break
					}
				}

				if foundNotAllowedAfterChar {
					return output, errors.New(fmt.Sprintf("Unrecognized character after a number at index %d: '%c'", i, str[j]))
				}

				output = append(output, Token{ttype: NUMBER, value: str[i : j+1], startAt: i, endAt: j + 1})
				i = j + 1
				break
			}

			if c == '"' {
				output = append(output, Token{ttype: DOUBLEQUOTE, value: str[i : i+1], startAt: i, endAt: i + 1})

				var nextChar rune
				var j int
				for j := j + 1; j < strLen; j++ {
					nextChar, isEOF = getCharacter(str, j)
					if isEOF {
						return output, errors.New(fmt.Sprintf("String was not closed before EOF at %d", i))
					}

					if nextChar == '"' {
						if str[j-1] == '\\' {
							continue
						}

						break
					}
				}

				output = append(output, Token{ttype: TEXT, value: str[i:j], startAt: i, endAt: j})
				output = append(output, Token{ttype: DOUBLEQUOTE, value: str[j : j+1], startAt: j, endAt: j + 1})
				i = j + 1
				break
			}

			return output, errors.New(fmt.Sprintf("Unrecognized character at index %d: '%c'", i, str[i]))
		}

		if isEOF {
			break
		}
	}

	return output, nil
}

/*
{
	"a": 1,
	"b": "c",
	"c": false,
	"d": [
		"a",
		"b",
		{
			"x": "y"
		}
	],
	"e": {
		...
	}
}
*/
