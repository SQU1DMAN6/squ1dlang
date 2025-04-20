package evaluator

import (
	"squ1d/lexer"
	"squ1d/object"
	"squ1d/parser"
	"testing"
)

func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("Object is not Integer. Got %T (%+v)", obj, obj)
	}
	if result.Value != expected {
		t.Errorf("Object has wrong value. Got %d, expected %d", result.Value, expected)
		return false
	}

	return true
}

func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("Object is not Boolean. Got %T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("Object has wrong value. Got %t, expected %t", result.Value, expected)
		return false
	}
	return true
}

func TestEvalBooleanExpression(t *testing.T) {
	tests := []struct {
		input string
		expected bool
	}{
		{"true", true},
		{"false", false},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testBooleanObject(t, evaluated, tt.expected)
	}
}

//PARSING

func TestEvalIntegerParsing(t *testing.T) {
	tests := []struct {
		input string
		expected int64
	}{
		{"5", 5},
		{"10", 10},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerObject(t, evaluated, tt.expected)
	}
}

func testEval(input string) object.Object {
	l := lexer.New(input)
	p := parser.New(l)
	program := p.ParseProgram()

	return Eval(program)
}