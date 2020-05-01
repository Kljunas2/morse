package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/kljunas2/morse"
)

var extended = flag.Bool("e", true, "use extended Morse code")
var punctuation = flag.Bool("p", false, "convert punctuation")
var noNewline = flag.Bool("n", false, "do not append newline")

func main() {
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	enc := &morse.Encoder{Extended: *extended, Punctuation: *punctuation}
	reader.WriteTo(enc)
	io.Copy(os.Stdout, enc)
	if !*noNewline {
		fmt.Println()
	}
}
