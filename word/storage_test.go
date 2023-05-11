package word

import (
	"bufio"
	"os"
	"sort"
	"strings"
	"testing"
)

type testSource struct {
	current int
	words   []string
}

func (s *testSource) Scan() bool {
	s.current++
	return s.current < len(s.words)
}

func (s *testSource) Text() string {
	return s.words[s.current]
}

func TestStorage(t *testing.T) {
	file, err := os.Open("testdata/words.txt")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	defer func() { _ = file.Close() }()
	source := testSource{current: -1}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := prepareLine(scanner.Text())
		if line != "" {
			source.words = append(source.words, line)
		}
	}
	sort.Strings(source.words)
	updated, err := UpdateStorage("testdata/storage.dat", &source)
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	_ = updated.Close()
	storage, err := NewStorageNamed("testdata/storage.dat")
	if err != nil {
		t.Fatalf("%s", err.Error())
	}
	words := []string(nil)
	for i := 1; i < MaxLength; i++ {
		count, err := storage.Words(i)
		if err != nil {
			t.Fatalf("%s", err.Error())
		}
		for j := 0; j < count; j++ {
			word, err := storage.WordAt(i, j)
			if err != nil {
				t.Fatalf("%s", err.Error())
			}
			words = append(words, word)
		}
	}
	sort.Strings(words)
	if len(words) != len(source.words) {
		t.Errorf("incorrect result: expected [%s], got [%s]", strings.Join(source.words, ","),
			strings.Join(words, ","))
	}
}

func TestBadFile(t *testing.T) {
	if _, err := NewStorageNamed("testdata/words.txt"); err == nil {
		t.Errorf("error expected but did not happen")
	}
}
