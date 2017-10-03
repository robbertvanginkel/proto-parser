package main

import (
    "os"
    "fmt"
	"text/scanner"
)

func main() {
        r, _ := os.Open(os.Args[1])
	s := new(scanner.Scanner)
	s.Init(r)
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		fmt.Printf("%s: \t%s\t%s\n", s.Position, scanner.TokenString(tok), s.TokenText())
	}
}
