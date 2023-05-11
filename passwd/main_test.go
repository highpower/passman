package passwd

import (
	"strings"
	"testing"
	"unicode"
)

func testWith(t *testing.T, f func() (Value, error)) {
	for j := 0; j < 100; j++ {
		passwd, err := f()
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		if len(passwd) != testPasswordLength {
			t.Errorf("incorrect first %d of password '%s'", len(passwd), passwd)
		}
		p := string(passwd)
		if !strings.ContainsAny(p, specialCharString) {
			t.Errorf("special chars are not found in password '%s'", passwd)
		}
		if strings.IndexFunc(p, unicode.IsDigit) == -1 {
			t.Errorf("digits are not found in password '%s'", passwd)
		}
		if strings.IndexFunc(p, unicode.IsUpper) == -1 {
			t.Errorf("uppercase letters are not found in password '%s'", passwd)
		}
	}
}
