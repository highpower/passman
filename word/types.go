package word

import (
	"io"
	"passman/translit"
)

type Source interface {
	Scan() bool
	Text() string
}

type Storage interface {
	Words(len int) (int, error)
	WordAt(len, index int) (string, error)
}

type CloseableStorage interface {
	Storage
	io.Closer
}

const MaxLength = 48

func ReadSource(source Source, f func(str string) error) error {
	for source.Scan() {
		word := translit.ToLatin(prepareLine(source.Text()))
		if word == "" {
			continue
		}
		if err := f(word); err != nil {
			return err
		}
	}
	return nil
}
