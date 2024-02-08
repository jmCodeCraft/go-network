package io

import (
	"testing"
)

func TestLinetoList(t *testing.T) {
	list := lineToList([]string{"1", "2", "3"})

	if len(list) != 3 {
		t.Errorf("Expected 3, got %d", len(list))
	}

	if list[0] != 1 {
		t.Errorf("Expected 1, got %d", list[0])
	}

	if list[1] != 2 {
		t.Errorf("Expected 2, got %d", list[1])
	}

	if list[2] != 3 {
		t.Errorf("Expected 3, got %d", list[2])
	}
}
