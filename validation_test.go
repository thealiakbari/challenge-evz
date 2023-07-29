package main

import (
	"testing"
)

func TestConcatWithPlusOperator(t *testing.T) {
	result := concatStringPlusOperation("Ali", "Akbari")
	expected := "Ali Akbari"
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestConcatWithSprintf(t *testing.T) {
	result := concatWhitFmt("Ali", "Akbari")
	expected := "Ali Akbari"
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}

func TestConcatWithJoin(t *testing.T) {
	result := concatStringJoin("Ali", "Akbari")
	expected := "Ali Akbari"
	if result != expected {
		t.Errorf("Expected '%s', but got '%s'", expected, result)
	}
}
