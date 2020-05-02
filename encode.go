package morse

import (
	"bytes"
	"fmt"
	"strings"
)

// Basic wrapper function for Encoder type.
func Encode(s string) string {
	enc := &Encoder{}
	r := strings.NewReader(s)
	n, err := r.WriteTo(enc)
	if err != nil {
		panic(err)
	}
	if n != int64(len(s)) {
		panic(fmt.Sprintf("Couldn't encode whole string %q", s))
	}
	buf := bytes.Buffer{}
	buf.ReadFrom(enc)
	return buf.String()
}
