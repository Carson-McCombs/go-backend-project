package binarytree

import (
	transaction "go-fetch-backend/Types/Transaction"
	"testing"
	"time"
)

type testcase_Pop[T any] struct {
	sortedArray         []T
	index               uint64
	expectedArray       []T
	expectedPoppedValue T
}

func Test_Pop(t *testing.T) {
	testcases := testcase_Pop[int]{
		sortedArray:         []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9},
		index:               3,
		expectedArray:       []int{0, 1, 2, 4, 5, 6, 7, 8, 9},
		expectedPoppedValue: 3,
	}
	actualArray := testcases.sortedArray
	actualPoppedValue := RemoveAtIndex(&actualArray, testcases.index)
	t.Logf("Expected (%d) %v \n Actual (%d) %v", testcases.expectedPoppedValue, testcases.expectedArray, actualPoppedValue, actualArray)
	if actualPoppedValue != testcases.expectedPoppedValue {
		t.Fatalf("Error: Popping (Index %d) %v \n Expected (%d) %v \n Actual (%d) %v", testcases.index, testcases.sortedArray, testcases.expectedPoppedValue, testcases.expectedArray, actualPoppedValue, actualArray)
	}

}

type testcase_Insert[T any] struct {
	sortedArray   []T
	value         T
	expectedArray []T
}

func Test_InsertGeneric(t *testing.T) {
	testcases := []testcase_Insert[int]{
		{
			sortedArray:   []int{1, 2, 3, 4, 5, 6, 7, 8, 19, 20},
			value:         15,
			expectedArray: []int{1, 2, 3, 4, 5, 6, 7, 8, 15, 19, 20},
		},
	}
	compareFunc := func(valueA int, valueB int) int8 {
		if valueB < valueA {
			return -1
		} else if valueB > valueA {
			return 1
		} else {
			return 0
		}
	}
	for _, testcase := range testcases {
		actualArray := testcase.sortedArray
		InsertIntoSorted(&actualArray, testcase.value, compareFunc)
		if !isEqual(actualArray, testcase.expectedArray) {
			t.Fatalf("Error \n Expected: %v \n Actual: %v", testcase.expectedArray, actualArray)
		}
	}
}

func Test_InsertTransactions(t *testing.T) {
	testcases := []testcase_Insert[transaction.Transaction]{
		{
			sortedArray: []transaction.Transaction{
				{Payer: "DANNON", Points: 300, Timestamp: time.Date(2022, 10, 31, 10, 0, 0, 0, time.UTC)},
				{Payer: "UNILEVER", Points: 200, Timestamp: time.Date(2022, 10, 31, 11, 0, 0, 0, time.UTC)},
				{Payer: "MILLER COORS", Points: 10000, Timestamp: time.Date(2022, 11, 1, 14, 0, 0, 0, time.UTC)},
				{Payer: "DANNON", Points: 1000, Timestamp: time.Date(2022, 11, 2, 14, 0, 0, 0, time.UTC)},
			},
			value: transaction.Transaction{Payer: "DANNON", Points: -200, Timestamp: time.Date(2022, 10, 31, 15, 0, 0, 0, time.UTC)},
			expectedArray: []transaction.Transaction{
				{Payer: "DANNON", Points: 300, Timestamp: time.Date(2022, 10, 31, 10, 0, 0, 0, time.UTC)},
				{Payer: "UNILEVER", Points: 200, Timestamp: time.Date(2022, 10, 31, 11, 0, 0, 0, time.UTC)},
				{Payer: "DANNON", Points: -200, Timestamp: time.Date(2022, 10, 31, 15, 0, 0, 0, time.UTC)},
				{Payer: "MILLER COORS", Points: 10000, Timestamp: time.Date(2022, 11, 1, 14, 0, 0, 0, time.UTC)},
				{Payer: "DANNON", Points: 1000, Timestamp: time.Date(2022, 11, 2, 14, 0, 0, 0, time.UTC)},
			},
		},
	}
	for _, testcase := range testcases {
		actualArray := testcase.sortedArray
		InsertIntoSorted(&actualArray, testcase.value, transaction.Compare)
		if !isEqual(actualArray, testcase.expectedArray) {
			t.Fatalf("Error \n Expected: %v \n Actual: %v", testcase.expectedArray, actualArray)
		}
	}
}

func isEqual[T comparable](sliceA []T, sliceB []T) bool {
	if len(sliceA) != len(sliceB) {
		return false
	}
	for i := range len(sliceA) {
		if sliceA[i] != sliceB[i] {
			return false
		}
	}
	return true
}
