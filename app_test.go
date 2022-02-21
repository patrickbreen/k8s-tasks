package main

import (
	"testing"
)

// example test
func calculate() bool {
	return true
}

func TestWithTimeOut(t *testing.T) {
	if !calculate() {
		t.Errorf("The calculation was false")
	}
}
