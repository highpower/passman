package translit

import (
	"bytes"
	"golang.org/x/text/unicode/norm"
	"unicode"
)

type transliterationProcessor struct {
	buffer   bytes.Buffer
	spaces   bool
	lastRune rune
}

const runeX = 'х'
const shortX = "h"
const specialX = "kh"

var specialRunes = map[rune]bool{
	'k': true, 'z': true, 'c': true, 's': true, 'e': true, 'h': true,
}

var replacements = map[rune]string{
	'а': "a",
	'б': "b",
	'в': "v",
	'г': "g",
	'д': "d",
	'е': "e",
	'ё': "yo",
	'ж': "zh",
	'з': "z",
	'и': "i",
	'й': "j",
	'к': "k",
	'л': "l",
	'м': "m",
	'н': "n",
	'о': "o",
	'п': "p",
	'р': "r",
	'с': "s",
	'т': "t",
	'у': "u",
	'ф': "f",
	'ц': "c",
	'ч': "ch",
	'ш': "sh",
	'щ': "shch",
	'ъ': "",
	'ы': "y",
	'ь': "",
	'э': "eh",
	'ю': "yu",
	'я': "ya"}

func (t *transliterationProcessor) transliterate(runes []rune) string {
	for _, r := range runes {
		switch {
		case unicode.IsSpace(r) || unicode.IsPunct(r):
			t.processSpace()
		case unicode.Is(unicode.Digit, r) || unicode.Is(unicode.Latin, r):
			t.processLatin(unicode.ToLower(r))
		case unicode.Is(unicode.Cyrillic, r):
			t.processCyrillic(r)
		}
	}
	return t.buffer.String()
}

func (t *transliterationProcessor) flushSpaces() {
	if !t.spaces {
		return
	}
	t.spaces = false
	if t.buffer.Len() != 0 {
		t.buffer.WriteRune('-')
	}
}

func (t *transliterationProcessor) processSpace() {
	t.spaces = true
}

func (t *transliterationProcessor) processLatin(r rune) {
	t.flushSpaces()
	t.buffer.WriteRune(r)
}

func (t *transliterationProcessor) processXRune() {
	if _, ok := specialRunes[t.lastRune]; !ok {
		t.buffer.WriteString(shortX)
	} else {
		t.buffer.WriteString(specialX)
	}
	t.lastRune = 'h'
}

func (t *transliterationProcessor) processCyrillic(r rune) {
	t.flushSpaces()
	if actual := unicode.ToLower(r); actual == runeX {
		t.processXRune()
	} else if replacement, ok := replacements[actual]; ok {
		runes := []rune(replacement)
		if len(runes) > 0 {
			t.lastRune = runes[len(runes)-1]
		}
		t.buffer.WriteString(replacement)
	}
}

func (t *transliterationProcessor) result() string {
	return t.buffer.String()
}

func ToLatin(value string) string {
	cp := transliterationProcessor{spaces: false, buffer: bytes.Buffer{}}
	cp.transliterate([]rune(norm.NFC.String(value)))
	return cp.result()
}

func Transliterated(value string) bool {
	for _, r := range value {
		if !unicode.IsDigit(r) && !(unicode.IsLetter(r) && unicode.Is(unicode.Latin, r)) && r != '-' {
			return false
		}
	}
	return true
}
