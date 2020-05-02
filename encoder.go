// Package morse provides tools to convert morse code to and from latin alphabet.
package morse

import (
	"bufio"
	"bytes"
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

type state struct {
	notFirstWord  bool
	leadingSpace  bool
	trailingSpace bool
}

// Encoder provides basic buffered encoding.
type Encoder struct {
	buffIn  bytes.Buffer
	buffOut bytes.Buffer

	state
	convertable string
	// Punctuation defines whether punctuation characters will be encoded as well.
	// It must be set before first call to write.
	Punctuation bool
	// Extended defines whether to use extended Morse code.
	// It must be set before first call to write.
	Extended bool
}

func (e *Encoder) Write(b []byte) (int, error) {
	n, err := e.buffIn.Write(b)
	e.translateBuffer()
	return n, err
}

func (e *Encoder) Read(buff []byte) (int, error) {
	return e.buffOut.Read(buff)
}

func (e *Encoder) translateBuffer() {
	scanner := bufio.NewScanner(&e.buffIn)
	scanner.Split(e.scanWords)

	for scanner.Scan() {
		word := scanner.Text()
		if err := scanner.Err(); err != nil {
			fmt.Println(err)
		}

		// Skip entirely non-printable word.
		if !e.containsConvertable(word) {
			continue
		}

		// Current word is either continuation from previus write or entirely new word.
		if e.notFirstWord {
			if e.leadingSpace {
				e.buffOut.WriteString("//")
			} else {
				e.buffOut.WriteString("/")
			}
		}
		e.notFirstWord = true

		for i, r := range []rune(word) {
			if e.isConvertable(r) {
				if i > 0 {
					// Prepend letter delimeter.
					e.buffOut.WriteString("/")
				}
				if unicode.IsLetter(r) {
					r = unicode.ToUpper(r)
				}
				code, ok := itu[r]
				if !ok {
					panic(fmt.Sprintf("%q should be convertable", r))
				}
				e.buffOut.WriteString(code)
			}
		}
	}
}

func (e *Encoder) isConvertable(r rune) bool {
	if e.convertable == "" {
		e.initConvertable()
	}
	r = unicode.ToUpper(r)
	return strings.ContainsRune(e.convertable, r)
}

func (e *Encoder) containsConvertable(s string) bool {
	for _, r := range []rune(s) {
		if e.isConvertable(r) {
			return true
		}
	}
	return false
}

func (e *Encoder) initConvertable() {
	e.convertable = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	if e.Punctuation {
		e.convertable += ".,?'!/()&:;=+-_\"$@"
	}
	if e.Extended {
		e.convertable += "ÀÅÄÆĄĆÇĈĤĐÉĘĴŁÈŃÑÓÖØŚŜŠÜŬŹŽČ"
	}
}

func (e *Encoder) scanWords(data []byte, atEOF bool) (advance int, token []byte, err error) {
	// Skip leading spaces.
	start := 0
	for width := 0; start < len(data); start += width {
		var r rune
		r, width = utf8.DecodeRune(data[start:])
		if !unicode.IsSpace(r) {
			break
		}
	}
	// Check whether token is preceeded by space.
	if start > 0 || e.trailingSpace {
		e.leadingSpace = true
	} else {
		e.leadingSpace = false
	}

	// Scan until space, marking end of word.
	for width, i := 0, start; i < len(data); i += width {
		var r rune
		r, width = utf8.DecodeRune(data[i:])
		if unicode.IsSpace(r) {
			e.trailingSpace = true
			//fmt.Println("read with trailing:", i+width, data[start:i])
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word. Return it.
	if atEOF && len(data) > start {
		e.trailingSpace = false
		//fmt.Println("read:", len(data), data[start:], start, e.trailingSpace)
		return len(data), data[start:], nil
	}
	// Request more data.
	e.trailingSpace = e.leadingSpace
	//fmt.Println("no read, trailing:", e.trailingSpace)
	return start, nil, nil
}
