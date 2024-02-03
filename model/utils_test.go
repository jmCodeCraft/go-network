package model

import (
	"reflect"
	"testing"
)

func TestPairwise(t *testing.T) {
	// Test case 1: Pairwise connections for a sequence of 1, 2, 3
	result := Pairwise([]int{1, 2, 3})
	expected := []Edge{
		{Node(1), Node(2)},
		{Node(2), Node(3)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 2: Pairwise connections for an empty sequence
	result = Pairwise([]int{})
	expected = []Edge{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 3: Pairwise connections for a sequence of 5, 10, 15
	result = Pairwise([]int{5, 10, 15})
	expected = []Edge{
		{Node(5), Node(10)},
		{Node(10), Node(15)},
	}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Add more test cases as needed...
}

func TestRange(t *testing.T) {
	// Test case 1: Range from 0 to 5
	result := Range(0, 5)
	expected := []int{0, 1, 2, 3, 4}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 2: Range from 3 to 8
	result = Range(3, 8)
	expected = []int{3, 4, 5, 6, 7}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 3: Range from -2 to 3
	result = Range(-2, 3)
	expected = []int{-2, -1, 0, 1, 2}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Test case 4: Range from 5 to 5 (empty range)
	result = Range(5, 5)
	expected = []int{}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("Expected %v, but got %v", expected, result)
	}

	// Add more test cases as needed...
}
