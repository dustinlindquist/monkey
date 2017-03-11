package repl

import (
	"bufio"
	"fmt"
	"io"

	"github.com/dustinlindquist/monkey/lexer"
	"github.com/dustinlindquist/monkey/token"
)

const PROMPT = ">> "

func Start(in io.Reader, out io.Writer) {
	var scanner = bufio.NewScanner(in)

	for {
		fmt.Println(PROMPT)
		var scanned = scanner.Scan()
		if !scanned {
			return
		}

		var line = scanner.Text()
		var l = lexer.New(line)

		for tok := l.NextToken(); tok.Type != token.EOF; tok = l.NextToken() {
			fmt.Println("%+v\n", tok)
		}
	}
}
