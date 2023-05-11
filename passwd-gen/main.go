package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"passman/common"
	"passman/passwd"
	"passman/word"
)

func bindOptions(options *passwd.Options) {
	flag.IntVar(&options.Digits, "digits", 1,
		"required number of digits in password")
	flag.IntVar(&options.UpperCaseLetters, "uppercase", 2,
		"required number of uppercase letters in password")
	flag.IntVar(&options.SpecialChars, "special", 1,
		"required number of special chars in password")
}

func runRandom(length int, options *passwd.Options) {
	p, err := passwd.Random(length, options)
	if err != nil {
		common.Halt(err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", string(p))
}

func runHuman(file string, length int, options *passwd.Options) {
	storage, err := word.NewStorageNamed(file)
	if err != nil {
		common.Halt(err)
	}
	defer func() { _ = storage.Close() }()
	p, err := passwd.Humanized(storage, length, options)
	if err != nil {
		common.Halt(err)
	}
	_, _ = fmt.Fprintf(os.Stdout, "%s\n", string(p))
}

func configured(humanized bool, file string) bool {
	return (humanized && file != "") || (!humanized && file == "")
}

func usage(writer io.Writer) {
	flag.CommandLine.SetOutput(writer)
	flag.PrintDefaults()
}

func main() {

	var file string
	flag.StringVar(&file, "file", "", "dictionary file to generate human passwords")

	var human bool
	flag.BoolVar(&human, "human", false, "generate password which is easy to remember")

	var length int
	flag.IntVar(&length, "length", 10, "required password length")

	options := passwd.Options{}
	bindOptions(&options)

	flag.Usage = func() { usage(os.Stdout) }
	flag.Parse()

	if !flag.Parsed() || !configured(human, file) || length == 0 {
		usage(os.Stderr)
		os.Exit(1)
	}
	if human {
		runHuman(file, length, &options)
		return
	}
	runRandom(length, &options)
}
