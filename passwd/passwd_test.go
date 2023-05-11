package passwd

import (
	"strings"
	"testing"
	"unicode"
)

const testPasswordLength = 14

func TestRoundRobin(t *testing.T) {

	type testResults struct {
		digitCount       int
		letterCount      int
		uppercaseCount   int
		specialCharCount int
	}

	results := testResults{}
	for i := 0; i < 10000; i++ {
		c, err := roundRobin.next()
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		switch {
		case unicode.IsDigit(c):
			results.digitCount++
		case unicode.IsLower(c):
			results.letterCount++
		case unicode.IsUpper(c):
			results.uppercaseCount++
		case strings.ContainsRune(specialCharString, c):
			results.specialCharCount++
		}
	}
}

func TestRandom(t *testing.T) {
	testWith(t, func() (Value, error) {
		return Random(testPasswordLength, &Options{Digits: 1, UpperCaseLetters: 2, SpecialChars: 1})
	})
}

func BenchmarkRandom(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, err := Random(testPasswordLength, &Options{Digits: 1, UpperCaseLetters: 2, SpecialChars: 1})
		if err != nil {
			b.Fatalf("%s", err.Error())
		}
	}
}

func FuzzRandom(f *testing.F) {
	f.Add(testPasswordLength, 1, 2, 1)
	f.Fuzz(func(t *testing.T, length, digits, upperCases, specialChars int) {
		options := Options{Digits: digits, UpperCaseLetters: upperCases, SpecialChars: specialChars}
		_, _ = Random(length, &options)
	})
}
