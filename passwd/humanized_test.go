package passwd

import (
	"passman/word"
	"testing"
)

func TestHumanized(t *testing.T) {
	storage, err := word.NewStorageNamed("testdata/russian.dict")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer func() { _ = storage.Close() }()
	testWith(t, func() (Value, error) {
		return Humanized(storage, testPasswordLength, &Options{Digits: 1, UpperCaseLetters: 2, SpecialChars: 1})
	})
}

func TestHumanizedErr(t *testing.T) {
	storage, err := word.NewStorageNamed("testdata/russian.dict")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer func() { _ = storage.Close() }()
	_, _ = Humanized(storage, 166, &Options{Digits: 81, UpperCaseLetters: 1, SpecialChars: 84})
}

func FuzzHumanized(f *testing.F) {
	storage, err := word.NewStorageNamed("testdata/russian.dict")
	if err != nil {
		f.Fatalf("%s", err.Error())
	}
	defer func() { _ = storage.Close() }()
	f.Add(testPasswordLength, 1, 2, 1)
	f.Fuzz(func(t *testing.T, length, digits, upperCases, specialChars int) {
		options := Options{Digits: digits, UpperCaseLetters: upperCases, SpecialChars: specialChars}
		_, _ = Humanized(storage, length, &options)
	})
}
