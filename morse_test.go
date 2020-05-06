package morse

import (
	"bytes"
	"fmt"
	"testing"
)

func TestStrings(t *testing.T) {
	tests := []string{"aaa", "AaA", "!AaA", "aaa aaa aaa", "avc ! cba", "12305", "če čebula nebi še imela."}
	for _, s := range tests {
		testEncoder(s, t)
		testPunctuation(s, t)
		testExtended(s, t)
		testFunc(s, t)
	}
}

func testEncoder(s string, t *testing.T) {
	encoder := &Encoder{}
	n, err := encoder.Write([]byte(s))
	if n != len(s) || err != nil {
		t.Errorf("encoder faied to convert %q", s)
	}
	buf := bytes.Buffer{}
	buf.ReadFrom(encoder)
	fmt.Printf("encoder encoded %q to %q\n", s, buf.String())
}

func testPunctuation(s string, t *testing.T) {
	encoder := &Encoder{}
	encoder.Punctuation = true
	n, err := encoder.Write([]byte(s))
	if n != len(s) || err != nil {
		t.Errorf("encoder(p) faied to convert %q", s)
	}
	buf := bytes.Buffer{}
	buf.ReadFrom(encoder)
	fmt.Printf("encoder(p) encoded %q to %q\n", s, buf.String())
}

func testExtended(s string, t *testing.T) {
	encoder := &Encoder{}
	encoder.Extended = true
	n, err := encoder.Write([]byte(s))
	if n != len(s) || err != nil {
		t.Errorf("encoder(e) faied to convert %q", s)
	}
	buf := bytes.Buffer{}
	buf.ReadFrom(encoder)
	fmt.Printf("encoder(e) encoded %q to %q\n", s, buf.String())
}

func testFunc(s string, t *testing.T) {
	fmt.Printf("func encoded %q to %q\n", s, Encode(s))
}

func TestMultiEncode(t *testing.T) {
	encoder := &Encoder{}
	encoder.Write([]byte("aaaa 22"))
	encoder.Write([]byte("bbb "))
	encoder.Write([]byte("ccccc"))
	encoder.Write([]byte(" dddddddd"))
	buf := bytes.Buffer{}
	buf.ReadFrom(encoder)
	fmt.Printf("encoded to %q\n", buf.String())
}

func ExampleEncode() {
	fmt.Println(Encode("Encode this into Morse code."))
	// Output:
	//./-./-.-./---/-.././/-/..../../...//../-./-/---//--/---/.-./..././/-.-./---/-../.
}

func TestBrokenWrite(t *testing.T) {
	encoder := &Encoder{}
	encoder.Extended = true
	str := "čač bcčdef"
	byt := []byte(str)
	encoder.Write(byt[:4])
	encoder.Write(byt[4:])
	buf := bytes.Buffer{}
	buf.ReadFrom(encoder)
	fmt.Printf("%qencoded to %q\n", str, buf.String())
}
