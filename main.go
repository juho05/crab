package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"

	"github.com/Bananenpro/crab/interpreter"
)

func main() {
	rand.Seed(time.Now().UnixNano())

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
		fmt.Println(strings.Repeat("=", 50))
	}

	program, errs := interpreter.Parse(tokens, lines)
	for _, err := range errs {
		fmt.Fprintln(os.Stderr, err)
	}
	if len(errs) > 0 {
		os.Exit(1)
	}

	err = interpreter.Check(program, lines)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	if *verbose {
		for _, stmt := range program {
			fmt.Println(interpreter.PrintAST(stmt))
		}
		fmt.Println(strings.Repeat("=", 50))
	}

	err = interpreter.Interpret(program, lines)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
