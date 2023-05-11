package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"passman/common"
	"passman/word"
)

func usage(writer io.Writer) {
	flag.CommandLine.SetOutput(writer)
	flag.PrintDefaults()
}

func main() {

	var file string
	flag.StringVar(&file, "file", "", "file to store words")

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()
	if !flag.Parsed() || file == "" {
		usage(os.Stderr)
		os.Exit(1)
	}

	storage, err := word.NewStorageNamed(file)
	if err != nil {
		common.Halt(err)
	}
	defer func() { _ = storage.Close() }()
	for length := 1; length < word.MaxLength; length++ {
		count, err := storage.Words(length)
		if err != nil {
			common.Halt(err)
		}
		for i := 0; i < count; i++ {
			w, err := storage.WordAt(length, i)
			if err != nil {
				common.Halt(err)
			}
			_, _ = fmt.Fprintf(os.Stdout, "%s\n", w)
		}
	}
}
