package main

import (
	"bufio"
	"flag"
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

	scanner := bufio.NewScanner(os.Stdin)
	storage, err := word.UpdateStorage(file, scanner)
	if err != nil {
		common.Halt(err)
	}
	_ = storage.Close()
}
