package passwd

import (
	"crypto/rand"
	"fmt"
	"math/big"
	"passman/word"
	"unicode"
)

type pair struct {
	first  int
	second int
}

type passwdBuffer struct {
	runes []rune
}

func (b *passwdBuffer) writeRune(r rune) {
	b.runes = append(b.runes, r)
}

func (b *passwdBuffer) writeString(str string) {
	b.runes = append(b.runes, []rune(str)...)
}

func (b *passwdBuffer) toUpper(count int) error {
	indices := make([]int, 0, len(b.runes))
	for i := range b.runes {
		if unicode.IsLetter(b.runes[i]) {
			indices = append(indices, i)
		}
	}
	if count >= len(indices) {
		for _, index := range indices {
			b.runes[index] = unicode.ToUpper(b.runes[index])
		}
		return nil
	}

	for i := 0; i < count; i++ {
		max := big.NewInt(int64(len(indices)))
		val, err := rand.Int(rand.Reader, max)
		if err != nil {
			return err
		}
		pos := int(val.Int64())
		index := indices[pos]
		b.runes[index] = unicode.ToUpper(b.runes[index])
		indices[pos], indices[len(indices)-1] = indices[len(indices)-1], indices[pos]
		indices = indices[:len(indices)-1]
	}
	return nil
}

func (b *passwdBuffer) apply(gen generator, count int) error {
	for i := 0; i < count; i++ {
		r, err := gen.next()
		if err != nil {
			return err
		}
		max := big.NewInt(int64(len(b.runes)))
		val, err := rand.Int(rand.Reader, max)
		if err != nil {
			return err
		}
		index := int(val.Int64())
		b.runes = append(b.runes[:index+1], b.runes[index:]...)
		b.runes[index] = r
	}
	return nil
}

func findValidPair(storage word.Storage, total int) (pair, error) {
	for {
		max := big.NewInt(int64(total))
		val, err := rand.Int(rand.Reader, max)
		if err != nil {
			return pair{}, err
		}
		length := int(val.Int64())
		if length == 0 {
			continue
		}
		count, err := storage.Words(length)
		if err != nil {
			return pair{}, err
		}
		other, err := storage.Words(total - length)
		if err != nil {
			return pair{}, err
		}
		if count != 0 && other != 0 {
			return pair{first: length, second: total - length}, nil
		}
	}
}

func randomWord(storage word.Storage, length int) (string, error) {
	count, err := storage.Words(length)
	if err != nil {
		return "", err
	}
	max := big.NewInt(int64(count))
	value, err := rand.Int(rand.Reader, max)
	if err != nil {
		return "", err
	}
	return storage.WordAt(length, int(value.Int64()))
}

func Humanized(storage word.Storage, length int, options *Options) (Value, error) {
	if err := options.Check(length); err != nil {
		return nil, err
	}
	wordTotalLength := length - options.Digits - options.SpecialChars
	switch {
	case wordTotalLength <= 0:
		return nil, fmt.Errorf("%w %d", errIncorrectLength, length)
	case wordTotalLength <= 2:
		return Random(length, options)
	}
	specialCharCount := options.SpecialChars
	p, err := findValidPair(storage, wordTotalLength)
	if err != nil {
		return nil, err
	}
	first, err := randomWord(storage, p.first)
	if err != nil {
		return nil, err
	}
	second, err := randomWord(storage, p.second)
	if err != nil {
		return nil, err
	}
	buffer := passwdBuffer{}
	buffer.writeString(first)
	if specialCharCount > 0 {
		r, err := specialChars.next()
		if err != nil {
			return nil, err
		}
		buffer.writeRune(r)
		specialCharCount--
	}
	buffer.writeString(second)
	if err := buffer.apply(digits, options.Digits); err != nil {
		return nil, err
	}
	if err := buffer.apply(specialChars, specialCharCount); err != nil {
		return nil, err
	}
	if err := buffer.toUpper(options.UpperCaseLetters); err != nil {
		return nil, err
	}
	return buffer.runes, nil
}
