# morse
Go library for converting to morse code

Example converts standard input and prints it to standard output:
```go
package main

import (
	"bufio"
	"fmt"
	"io"
	"os"

	"github.com/kljunas2/morse"
)

func main() {
	flag.Parse()
	reader := bufio.NewReader(os.Stdin)
	enc := &morse.Encoder{Extended: true}
	reader.WriteTo(enc)
	io.Copy(os.Stdout, enc)
}
```
