package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"
	"path/filepath"
	"squ1d/evaluator"
	"squ1d/lexer"
	"squ1d/object"
	"squ1d/parser"
	"squ1d/repl"
	"strings"
)

func main() {
	if len(os.Args) > 1 {
		// File mode
		filename := os.Args[1]
		runFile(filename)
	} else {
		// REPL mode
		user, err := user.Current()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Hello %s! This is the SQU1D programming language!\n", user.Username)
		fmt.Printf("Feel free to type in commands\n")
		repl.Start(os.Stdin, os.Stdout)
	}
}

func runFile(filename string) {
	// Check file extension
	expectedFormat := ".sqd"
	actualFormat := strings.ToLower(filepath.Ext(filename))

	if actualFormat != expectedFormat {
		fmt.Printf("Incorrect file format: Got %s, expected %s\n", actualFormat, expectedFormat)
		return
	}

	// Read file
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("Failed to read file %s: %s\n", filename, err)
		return
	}

	code := string(data)

	// Run code through the ususal lexer > parser > evaluator pipeline
	l := lexer.New(code)
	p := parser.New(l)
	program := p.ParseProgram()

	if len(p.Errors()) != 0 {
		for _, msg := range p.Errors() {
			fmt.Println("Parser error: ", msg)
		}
		return
	}

	env := object.NewEnvironment()
	evaluated := evaluator.Eval(program, env)

	if evaluated != nil {
		fmt.Println(evaluated.Inspect())
	}
}