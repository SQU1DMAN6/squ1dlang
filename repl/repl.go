package repl

import (
	"bufio"
	"fmt"
	"io"
	"os/user"
	"squ1d/evaluator"
	"squ1d/lexer"
	"squ1d/object"
	"squ1d/parser"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	env := object.NewEnvironment()
	_, err := user.Current()
	if err != nil {
		panic(err)
	}
	scanner := bufio.NewScanner(in)
	for {
		fmt.Print(PROMPT)
		scanned := scanner.Scan()
		if !scanned {
			return
		}
		line := scanner.Text()
		l := lexer.New(line)
		p := parser.New(l)
		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			printParserErrors(out, p.Errors())
			continue
		}
		evaluated := evaluator.Eval(program, env)
		if evaluated != nil {
			io.WriteString(out, evaluated.Inspect())
			io.WriteString(out, "\n")
		}
	}
}

func printParserErrors(out io.Writer, errors []string) {
	for _, msg := range errors {
		io.WriteString(out, "\t"+msg+"\n")
	}
}
