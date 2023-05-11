package word

import (
	"fmt"
	"testing"
)

func TestActual(t *testing.T) {

	type testCase struct {
		num            int
		resultExpected int
	}

	testCases := []testCase{
		{num: 0, resultExpected: 8},
		{num: 3, resultExpected: 8},
		{num: 7, resultExpected: 8},
		{num: 8, resultExpected: 8},
		{num: 12, resultExpected: 16},
		{num: 16, resultExpected: 16},
		{num: 17, resultExpected: 32},
	}

	for _, tc := range testCases {
		t.Run(fmt.Sprintf("Actual/%d", tc.num), func(t *testing.T) {
			if res := alignedSize(tc.num); res != tc.resultExpected {
				t.Errorf("incorrect result: expected %d, got %d", tc.resultExpected, res)
			}
		})
	}
}
