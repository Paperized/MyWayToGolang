package parser

import (
	"fmt"
	"strconv"
	"strings"
)

/*
PARSER RESULT: <output_value>
<output_value> ::= <object> | <array>
<array> ::= "[" "]" | "[" <elements> "]"
<elements> ::= <all_values> | <elements> "," <all_values>
<all_values> ::= true | false | null | NUMBER | TEXT | <output_value>
<object> ::= "{" "}" | "{" <properties> "}" OK
<properties> ::= <property> | <properties> "," <property> OK
<property> ::= TEXT ":" <all_values> OK
*/

func JsonStringToMap(input string, alloc ...bool) (any, error) {
	tokens, err := TokenizeJsonString(input, alloc...)
	if err != nil {
		return nil, err
	}
	if len(tokens) == 0 {
		return nil, fmt.Errorf("parsing error json not provided")
	}

	var result any

	switch tokens[0].ttype {
	case CURLY_OPEN:
		result, _, err = parseJsonObject(0, tokens)
	case SQUARE_OPEN:
		result, _, err = parseJsonArray(0, tokens)
	default:
		return nil, fmt.Errorf("parsing error: starting character not valid")
	}

	return result, err
}

func nextToken(at int, tokens []Token) *Token {
	if at >= len(tokens) {
		return nil
	}

	return &tokens[at]
}

func isTokenTypeAt(ttype TokenType, at int, tokens []Token) (bool, error) {
	if at >= len(tokens) {
		return false, fmt.Errorf("parsing error expected %c at index %d", ttype, tokens[at-1].startAt)
	}

	return tokens[at].ttype == ttype, nil
}

func jsonValue(fromIndex int, tokens []Token) (any, int, error) {
	valueToken := nextToken(fromIndex, tokens)
	if valueToken == nil {
		return nil, fromIndex, fmt.Errorf("parser error: expected json value but EOF at index: %d", tokens[fromIndex-1].startAt)
	}

	var value any
	var err error
	switch valueToken.ttype {
	case TRUE:
		return true, fromIndex + 1, nil
	case FALSE:
		return false, fromIndex + 1, nil
	case NULL:
		return nil, fromIndex + 1, nil
	case TEXT:
		return valueToken.value, fromIndex + 1, nil
	case NUMBER:
		isFloat := strings.Contains(valueToken.value, ".")
		if isFloat {
			value, err = strconv.ParseFloat(valueToken.value, 64)
		} else {
			value, err = strconv.Atoi(valueToken.value)
		}

		if err != nil {
			err = fmt.Errorf("parsing error property number malformatted at index: %d", tokens[fromIndex-1].startAt)
		}
		return value, fromIndex + 1, err
	case CURLY_OPEN:
		value, fromIndex, err = parseJsonObject(fromIndex, tokens)
		return value, fromIndex, err
	case SQUARE_OPEN:
		value, fromIndex, err = parseJsonArray(fromIndex, tokens)
		return value, fromIndex, err
	default:
		return nil, fromIndex, fmt.Errorf("parsing error unrecognized value at index: %d", tokens[fromIndex-1].startAt)
	}
}

func parseJsonProperty(fromIndex int, tokens []Token) (string, any, int, error) {
	key := ""
	var value any = nil

	isType, err := isTokenTypeAt(TEXT, fromIndex, tokens)
	if err != nil {
		return key, value, fromIndex, err
	}
	if !isType {
		return key, value, fromIndex, fmt.Errorf("parsing error expected key at index: %d", tokens[fromIndex].startAt)
	}
	key = tokens[fromIndex].value
	fromIndex += 1

	isType, err = isTokenTypeAt(COLON, fromIndex, tokens)
	if err != nil {
		return key, value, fromIndex, err
	}
	if !isType {
		return key, value, fromIndex, fmt.Errorf("parsing error missing colon at index: %d", tokens[fromIndex].startAt)
	}
	fromIndex += 1

	valueToken := nextToken(fromIndex, tokens)
	if valueToken == nil {
		return key, value, fromIndex, fmt.Errorf("parsing error property value not defined at index: %d", tokens[fromIndex-1].startAt)
	}

	value, fromIndex, err = jsonValue(fromIndex, tokens)
	return key, value, fromIndex, err
}

func parseJsonObject(fromIndex int, tokens []Token) (map[string]any, int, error) {
	tokensLen := len(tokens)
	result := map[string]any{}

	isOfType, err := isTokenTypeAt(CURLY_OPEN, fromIndex, tokens)
	if err != nil {
		return nil, fromIndex, err
	}
	if !isOfType {
		return nil, fromIndex, fmt.Errorf("parsing error: object must start with curly open brackets at index: %d", tokens[fromIndex].startAt)
	}
	fromIndex += 1

	var keyProp string
	var valueProp any
	for ; fromIndex < tokensLen; fromIndex++ {
		token := tokens[fromIndex]
		switch token.ttype {
		case CURLY_CLOSE:
			return result, fromIndex + 1, nil
		case TEXT:
			if len(result) > 0 && tokens[fromIndex-1].ttype != COMMA {
				return result, fromIndex + 1, fmt.Errorf("parsing error: property must be separated by comma at index: %d", tokens[fromIndex].startAt)
			}

			keyProp, valueProp, fromIndex, err = parseJsonProperty(fromIndex, tokens)
			if err != nil {
				return result, fromIndex, err
			}

			result[keyProp] = valueProp
			// adjust to -1 since forloop will increase it by 1
			fromIndex -= 1
		case COMMA:
			if tokens[fromIndex-1].ttype == COMMA {
				return result, fromIndex + 1, fmt.Errorf("parsing error: propery expected but found comma at index: %d", tokens[fromIndex].startAt)
			}
		default:
			return result, fromIndex + 1, fmt.Errorf("parsing error unrecognized value at index: %d", tokens[fromIndex].startAt)
		}
	}

	return result, fromIndex, fmt.Errorf("unreachable code at index: %d", tokens[fromIndex].startAt)
}

func parseJsonArray(fromIndex int, tokens []Token) ([]any, int, error) {
	tokensLen := len(tokens)
	result := []any{}

	isOfType, err := isTokenTypeAt(SQUARE_OPEN, fromIndex, tokens)
	if err != nil {
		return nil, fromIndex, err
	}
	if !isOfType {
		return nil, fromIndex, fmt.Errorf("parsing error: object must start with square open brackets at index: %d", tokens[fromIndex].startAt)
	}
	fromIndex += 1

	var value any
	for ; fromIndex < tokensLen; fromIndex++ {
		token := tokens[fromIndex]
		switch token.ttype {
		case SQUARE_CLOSE:
			return result, fromIndex + 1, nil
		case TRUE, FALSE, NULL, TEXT, NUMBER, CURLY_OPEN, SQUARE_OPEN:
			if len(result) > 0 && tokens[fromIndex-1].ttype != COMMA {
				return result, fromIndex + 1, fmt.Errorf("parsing error: value must be separated by comma at index: %d", tokens[fromIndex].startAt)
			}

			value, fromIndex, err = jsonValue(fromIndex, tokens)
			if err != nil {
				return result, fromIndex, err
			}

			result = append(result, value)
			// adjust to -1 since forloop will increase it by 1
			fromIndex -= 1
		case COMMA:
			if tokens[fromIndex-1].ttype == COMMA {
				return result, fromIndex + 1, fmt.Errorf("parsing error: value expected but found comma at index: %d", tokens[fromIndex].startAt)
			}
		default:
			return result, fromIndex + 1, fmt.Errorf("parsing error unrecognized value at index: %d", tokens[fromIndex].startAt)
		}
	}

	return result, fromIndex, fmt.Errorf("unreachable code at index: %d", tokens[fromIndex].startAt)
}
