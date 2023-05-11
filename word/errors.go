package word

import "errors"

var (
	ErrNoSuchWord      = errors.New("no such word")
	ErrIncorrectLength = errors.New("incorrect length")
	ErrIncorrectFile   = errors.New("incorrect file")
)

var (
	errIncorrectRead  = errors.New("incorrect read")
	errIncorrectWrite = errors.New("incorrect write")
)
