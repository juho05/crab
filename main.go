package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/Bananenpro/crab/interpreter"
)

func main() {
	verbose := flag.Bool("verbose", false, "Print verbose output.")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <file>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		flag.PrintDefaults()
	}

	flag.Parse()

	if flag.NArg() != 1 {
		flag.Usage()
		os.Exit(1)
	}

	sourceFile, err := os.Open(flag.Arg(0))
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to open source file: %s\n", err)
		os.Exit(1)
	}

	tokens, lines, err := interpreter.Scan(sourceFile)
	sourceFile.Close()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Tokens:", tokens)
	}

	program, err := interpreter.Parse(tokens, lines)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *verbose {
		fmt.Println("Program:", interpreter.PrintAST(program))
	}
}
