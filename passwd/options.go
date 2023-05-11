package passwd

import (
	"errors"
	"fmt"
)

type Options struct {
	Digits           int
	UpperCaseLetters int
	SpecialChars     int
}

var (
	errIncorrectLength  = errors.New("incorrect length")
	errIncorrectOptions = errors.New("incorrect options")
)

func (o *Options) Valid() error {
	if o.Digits < 0 || o.UpperCaseLetters < 0 || o.SpecialChars < 0 {
		return fmt.Errorf("%w %s", errIncorrectOptions, o.String())
	}
	return nil
}

func (o *Options) Check(length int) error {
	if err := o.Valid(); err != nil {
		return err
	}
	if length <= 0 || length > 256 {
		return fmt.Errorf("%w %d", errIncorrectLength, length)
	}
	if o.Other(length) < 0 {
		return fmt.Errorf("%w %d", errIncorrectLength, length)
	}
	return nil
}

func (o *Options) Other(length int) int {
	return length - o.Digits - o.UpperCaseLetters - o.SpecialChars
}

func (o *Options) String() string {
	return fmt.Sprintf("options[digits=%d,upperCaseLetters=%d,specialChars=%d]", o.Digits,
		o.UpperCaseLetters, o.SpecialChars)
}
