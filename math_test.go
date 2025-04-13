package main

import (
	"testing"
)

func validateTokens(target, expected []Token, t *testing.T) {
	for i, token := range expected {
		if token != target[i] {
			t.Error("Invalid token stream")
		}
	}
}

func TestTokenizer(t *testing.T) {
	exp := "8.5 * 10 - (14 / 10)"
	//exp := "8.5 * 10"
	result, _ := TokenizeMathExp(exp)
	expected := []Token{
		{Number, None, 8.5},
		{Multi, High, 0.0},
		{Number, None, 10.0},
		{Minus, Low, 0.0},
		{Open, VeryHigh, 0.0},
		{Number, None, 14.0},
		{Divi, High, 0.0},
		{Number, None, 10.0},
		{Close, VeryHigh, 0.0},
	}
	
	validateTokens(result, expected, t)
}

func TestParser(t *testing.T) {
	tokens := []Token{
		{Number, None, 2.0},
		{Multi, High, 0.0},
		{Number, None, 3.0},
		{Plus, Low, 0.0},
		{Number, None, 7.0},
		{Plus, Low, 0.0},
		{Number, None, 8.0},
		{Multi, High, 0.0},
		{Number, None, 4.0},
	}
	result, _, err := ParseMathTokens(tokens, 0)
	expected := 45.0

	if result != expected || err != nil {
		t.Errorf("Parsing went wrong. Expected: %f, Got: %f", expected, result)
	}
}

func TestExpEvaluator(t *testing.T) {
	exp := "10 - 15 + 30 * 1.5 / 2"
	//exp := "2 * 3 + 7 + 8 * 4"
	result, _ := EvaluateMathExp(exp)
	expected := 17.5
	//expected := 45
	
	if result != expected {
		t.Errorf("Evaluation faild. Expected: %f, Got: %f", expected, result)
	}
}
