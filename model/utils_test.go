package model

import (
	"reflect"
	"testing"
)

func TestPairwise(t *testing.T) {
	// Test case 1: List of consecutive integers
	input1 := []int{1, 2, 3, 4}
	expected1 := []Edge{{Node1: 1, Node2: 2}, {Node1: 2, Node2: 3}, {Node1: 3, Node2: 4}}
	result1 := Pairwise(input1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, but got %v", expected1, result1)
	}

	// Test case 2: Empty input
	input2 := []int{}
	expected2 := []Edge{}
	result2 := Pairwise(input2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, but got %v", expected2, result2)
	}

	// Test case 3: Single element input
	input3 := []int{1}
	expected3 := []Edge{}
	result3 := Pairwise(input3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, but got %v", expected3, result3)
	}

	// Test case 4: Random input
	input4 := []int{5, 9, 12, 17}
	expected4 := []Edge{{Node1: 5, Node2: 9}, {Node1: 9, Node2: 12}, {Node1: 12, Node2: 17}}
	result4 := Pairwise(input4)
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed: Expected %v, but got %v", expected4, result4)
	}
}

func TestRange(t *testing.T) {
	// Test case 1: Positive range
	start1 := 1
	end1 := 5
	expected1 := []int{1, 2, 3, 4}
	result1 := Range(start1, end1)
	if !reflect.DeepEqual(result1, expected1) {
		t.Errorf("Test case 1 failed: Expected %v, but got %v", expected1, result1)
	}

	// Test case 2: Range starting from zero
	start2 := 0
	end2 := 3
	expected2 := []int{0, 1, 2}
	result2 := Range(start2, end2)
	if !reflect.DeepEqual(result2, expected2) {
		t.Errorf("Test case 2 failed: Expected %v, but got %v", expected2, result2)
	}

	// Test case 3: Empty range
	start3 := 5
	end3 := 5
	expected3 := []int{}
	result3 := Range(start3, end3)
	if !reflect.DeepEqual(result3, expected3) {
		t.Errorf("Test case 3 failed: Expected %v, but got %v", expected3, result3)
	}

	// Test case 4: Negative range
	start4 := -3
	end4 := 2
	expected4 := []int{-3, -2, -1, 0, 1}
	result4 := Range(start4, end4)
	if !reflect.DeepEqual(result4, expected4) {
		t.Errorf("Test case 4 failed: Expected %v, but got %v", expected4, result4)
	}
}
