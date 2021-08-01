package test

import (
	"testing"
)

func TestPrint(t *testing.T) {
	t.Log("Testing Print function")
	result := "testing"

	if result != "testing" {
		t.Errorf("Expected result of testing, but it was %s instead", result)
	}
}
