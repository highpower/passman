package common

import (
	"fmt"
	"os"
)

func Halt(err error) {
	_, _ = fmt.Fprintf(os.Stderr, "%s\n", err.Error())
	os.Exit(1)
}
