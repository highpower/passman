package passwd

import (
	"bytes"
	"crypto/rand"
	"fmt"
	"math/big"
	"unicode"
)

type Value []rune

type generator interface {
	next() (rune, error)
}

type listGenerator struct {
	length *big.Int
	runes  []rune
}

type uppercaseGenerator struct{}

type roundRobinGenerator struct {
	totalWeight *big.Int
	items       []generatorItem
}

type generatorItem struct {
	weight int
	gen    generator
}

const specialCharString = "!@#$%^&*;?"

//goland:noinspection SpellCheckingInspection

var (
	digits       = newListGenerator("0123456789")
	letters      = newListGenerator("abcdefghijklmnopqrstuvwxyz")
	specialChars = newListGenerator(specialCharString)
	roundRobin   = newRoundRobinGenerator()
)

func (v Value) String() string {
	return string(v)
}

func (v Value) MarshalJSON() ([]byte, error) {
	str := string(v)
	return []byte(str), nil
}

func (g *listGenerator) next() (rune, error) {
	val, err := rand.Int(rand.Reader, g.length)
	if err != nil {
		return 0, err
	}
	return g.runes[int(val.Int64())], nil
}

func (u uppercaseGenerator) next() (rune, error) {
	r, err := letters.next()
	if err != nil {
		return 0, err
	}
	return unicode.ToUpper(r), nil
}

func (g *roundRobinGenerator) find(random int) generator {
	sum := 0
	for _, item := range g.items {
		sum += item.weight * 100
		if random < sum {
			return item.gen
		}
	}
	panic(fmt.Errorf("incorrect algorithm"))
}

func (g *roundRobinGenerator) next() (rune, error) {
	x, err := rand.Int(rand.Reader, g.totalWeight)
	if err != nil {
		return 0, err
	}
	gen := g.find(int(x.Int64()))
	return gen.next()
}

func newListGenerator(chars string) generator {
	runes := []rune(chars)
	return &listGenerator{length: big.NewInt(int64(len(runes))), runes: runes}
}

func newRoundRobinGenerator() generator {
	totalWeight := 0
	items := []generatorItem{
		{gen: digits, weight: 2},
		{gen: uppercaseGenerator{}, weight: 4},
		{gen: specialChars, weight: 1},
		{gen: letters, weight: 18},
	}
	for _, i := range items {
		totalWeight += i.weight
	}
	return &roundRobinGenerator{totalWeight: big.NewInt(int64(totalWeight * 100)), items: items}
}

func generatorList(length int, options *Options) ([]generator, error) {
	result := make([]generator, 0, length)
	for i := 0; i < options.Digits; i++ {
		result = append(result, digits)
	}
	for i := 0; i < options.UpperCaseLetters; i++ {
		result = append(result, uppercaseGenerator{})
	}
	for i := 0; i < options.SpecialChars; i++ {
		result = append(result, specialChars)
	}
	for i := 0; i < options.Other(length); i++ {
		result = append(result, roundRobin)
	}
	return result, nil
}

func Random(length int, options *Options) (Value, error) {
	if err := options.Check(length); err != nil {
		return nil, err
	}
	generators, err := generatorList(length, options)
	if err != nil {
		return nil, err
	}
	buffer := bytes.Buffer{}
	for len(generators) > 0 {
		index := 0
		if len(generators) != 1 {
			max := big.NewInt(int64(len(generators)))
			x, err := rand.Int(rand.Reader, max)
			if err != nil {
				return nil, err
			}
			index = int(x.Int64())
		}
		char, err := generators[index].next()
		if err != nil {
			return nil, err
		}
		buffer.WriteRune(char)
		generators[index], generators[len(generators)-1] = generators[len(generators)-1], generators[index]
		generators = generators[:len(generators)-1]
	}
	return Value(buffer.String()), nil
}
