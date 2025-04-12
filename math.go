package main

import (
	"errors"
	"strconv"
	"unicode"
)

func EvaluateMathExp(exp string) (float64, error) {
	InvalidExp := errors.New("Invalid Expression")

	tokens, err := TokenizeMathExp(exp)
	if err != nil {
		return 0.0, InvalidExp
	}

	res := ParseMathTokens(tokens, 0)

	return res, nil
}

// func ParseMathTokens(tokens []Token, i int) float64 {
// 	/**** LOOP METHOD ****/
//
// 	acumulator := 0.0
//
// 	for i := 0; i < len(tokens); i += 1 {
// 		if tokens[i].kind == Number {
// 			acumulator += tokens[i].num
// 			continue
// 	}
//
// 		right := tokens[i + 1].num
//
// 		switch tokens[i].kind {
// 			default: return 0.0
// 			case Plus: acumulator += right
// 			case Minus: acumulator -= right
// 			case Multi: acumulator *= right
// 			case Divi: acumulator /= right
// 		}
//
// 		i += 1
// 	}
//
// 	return acumulator
//
//
// 	/**** RECURSION METHOD ****/
//
// 	left := 0.0
// 	right := 0.0
//
// 	if i >= len(tokens) - 1 {
// 		return tokens[i].num
// 	}
//
// 	if tokens[i].kind == Number {
// 		left = tokens[i].num
// 		i += 1
// 		right = ParseMathTokens(tokens, i + 1)
// 	}
//
// 	switch tokens[i].kind {
// 		default: return 0.0
// 		case Plus: left += right
// 		case Minus: left -= right
// 		case Multi: left *= right
// 		case Divi: left /= right
// 	}
//
// 	return left
// }

func ParseMathTokens(tokens []Token, i int) float64 {
	acumulator := 0.0

	for i = i; i < len(tokens); i += 1 {
		if tokens[i].kind == Number {
			acumulator += tokens[i].num
			continue
		}

		right := tokens[i + 1].num

		currentPrec := Low
		nextPrec := Low
		if i + 1 < len(tokens) - 1 {
			currentPrec = tokens[i].prec
			nextPrec = tokens[i + 2].prec
		}

		// If precedence increases we recurse
		recurse := currentPrec < nextPrec
		if recurse {
			right = ParseMathTokens(tokens, i + 1)
		}

		switch tokens[i].kind {
			case Plus: acumulator += right
			case Minus: acumulator -= right
			case Multi: acumulator *= right
			case Divi: acumulator /= right
		}

		if recurse {
			return acumulator
		}

		i += 1
	}

	return acumulator
}

func TokenizeMathExp(exp string) ([]Token, error) {
	tokens := []Token{}
	resErr := errors.New("invalid token")

	for i := 0; i < len(exp); i += 1 {
		nullToken := Token{Number, None, 0.0}
		token := nullToken
		switch exp[i] {
			case ' ': continue
			case '(': token = Token{Open, VeryHigh, 0.0}
			case ')': token = Token{Close, VeryHigh, 0.0}
			case '+': token = Token{Plus, Low, 0.0}
			case '-': token = Token{Minus, Low, 0.0}
			case '*': token = Token{Multi, High, 0.0}
			case '/': token = Token{Divi, High, 0.0}
		}
		if token != nullToken {
			tokens = append(tokens, token)
			continue
		}

		// Check for invalid tokens
		if !unicode.IsDigit(rune(exp[i])) {
			return nil, resErr
		}
		
		var numStr string
		for j := i; true; j += 1 {
			isPossibleNum := false
			if j < len(exp) {
				isPossibleNum = unicode.IsDigit(rune(exp[j])) || exp[j] == '.'
			}

			if isPossibleNum {
				numStr += string(exp[j])
				continue
			}
			
			// FIXME: Check for non space runes and return error.
			// Actually not just non space, but maybe merge this with
			// other tokens as well like +, -, * and so on

			if len(numStr) != 0 {
				numF64, err := strconv.ParseFloat(numStr, 64)
				if err != nil {
					return nil, resErr
				}

				tokens = append(tokens, Token{Number, None, numF64})
			}
			
			i += (j - i) - 1
			break
		}
	}

	return tokens, nil
}

type Token struct {
	kind TokenKind
	prec Precedence
	num float64
}

type TokenKind uint32
const (
	Number TokenKind = iota
	Open
	Close
	Plus
	Minus
	Multi
	Divi
)

type Precedence uint32
const (
	None Precedence = iota
	Low
	High
	VeryHigh
)
