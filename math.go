package main

import (
	"errors"
	"math"
	"strconv"
	"unicode"
)

func EvaluateMathExp(exp string) (float64, error) {
	InvalidExp := errors.New("Invalid Expression")

	tokens, err := TokenizeMathExp(exp)
	if err != nil {
		return 0.0, InvalidExp
	}

	res, _, err := ParseMathTokens(tokens, 0, false)
	if err != nil {
		return 0.0, InvalidExp
	}

	return res, nil
}

func ParseMathTokens(tokens []Token, i int, isRecursive bool) (float64, int, error) {
	InvalidSyntax := errors.New("Invalid Syntex")
	acumulator := 0.0

	for i = i; i < len(tokens); i += 1 {
		previous := Token{Invalid, None, 0.0}
		if i != 0 {
			previous = tokens[i - 1]
		}

		// Check for invalid Syntax, like double operator and no operator between numbers
		if previous.kind == tokens[i].kind || previous.kind >= Plus && tokens[i].kind >= Plus {
			return 0.0, 0, InvalidSyntax
		}

		if tokens[i].kind == Number {
			acumulator += tokens[i].num
			continue
		}
		
		j := i
		right := tokens[i + 1].num

		currentPrec := tokens[i].prec
		nextPrec := Low
		if i + 1 < len(tokens) - 1 {
			nextPrec = tokens[i + 2].prec
		}

		// If precedence increases we recurse
		recurse := currentPrec < nextPrec
		if recurse {
			var err error
			right, j, err = ParseMathTokens(tokens, i + 1, true)
			if err != nil {
				return 0.0, 0, err
			}
		}
		
		switch tokens[i].kind {
			case Plus: acumulator += right
			case Minus: acumulator -= right
			case Multi: acumulator *= right
			case Divi: acumulator /= right
			case Powe: acumulator = math.Pow(acumulator, right)
		}
		
		if j != i {
			i = j
		} else if recurse || nextPrec < currentPrec && isRecursive {
			// If precedence decreases we return to preserve the order of the operations
			return acumulator, j, nil
		}

		i += 1
	}

	return acumulator, 0, nil
}

func TokenizeMathExp(exp string) ([]Token, error) {
	InvalidToken := errors.New("Invalid token")
	tokens := []Token{}

	for i := 0; i < len(exp); i += 1 {
		nullToken := Token{Invalid, None, 0.0}
		token := nullToken
		switch exp[i] {
			case ' ': continue
			case '(': token = Token{Open, VeryHigh, 0.0}
			case ')': token = Token{Close, VeryHigh, 0.0}
			case '+': token = Token{Plus, Low, 0.0}
			case '-': {
				last := Token{Plus, None, 0.0}
				if len(tokens) != 0 {
					last = tokens[len(tokens) - 1]
				}

				if last.kind < Plus {
					token = Token{Minus, Low, 0.0}
				}
			}	
			case '*': token = Token{Multi, High, 0.0}
			case '/': token = Token{Divi, High, 0.0}
			case '^': token = Token{Powe, VeryHigh, 0.0}
		}
		if token != nullToken {
			tokens = append(tokens, token)
			continue
		}

		var numStr string
		for j := i; true; j += 1 {
			isPossibleNum := false
			if j < len(exp) {
				isPossibleNum = unicode.IsDigit(rune(exp[j])) || exp[j] == '.' || exp[j] == '-'
			}

			if isPossibleNum {
				numStr += string(exp[j])
				continue
			}
			
			// NOTE: This actually handles all the possible errors of the tokenizer
			// since everything that's not valid is treated as a possible number and 
			// then checked here
			numF64, err := strconv.ParseFloat(numStr, 64)
			if err != nil {
				return nil, InvalidToken
			}

			tokens = append(tokens, Token{Number, None, numF64})
			
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
	Invalid TokenKind = iota 
	Number
	Open
	Close
	Plus
	Minus
	Multi
	Divi
	Powe
)

type Precedence uint32
const (
	None Precedence = iota
	Low
	High
	VeryHigh
)
