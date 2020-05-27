// Package morse provides tools to convert morse code to
// and from latin alphabet.
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
// It translates based on ITU code, skipping unconvertable characters.
type Encoder struct {
	buffIn  bytes.Buffer
	buffOut bytes.Buffer

	state
	convertable string
	// Punctuation defines whether punctuation characters
	// will be encoded as well.
	// It must be set before first call to write.
	Punctuation bool
	// Extended defines whether to use extended Morse code.
	// It must be set before first call to write.
	Extended bool
}

// Write appends the ocntents of to the encoder's buffer, translating it.
// The return value n is the length of p, err is always nil.
func (e *Encoder) Write(p []byte) (n int, err error) {
	// Call to bytes.Buffer.Write always returns nil error.
	n, _ = e.buffIn.Write(p)
	e.translateBuffer()
	return n, nil
}

// Read writes translated buffer to p.
func (e *Encoder) Read(p []byte) (int, error) {
	return e.buffOut.Read(p)
}

func (e *Encoder) translateBuffer() {
	scanner := bufio.NewScanner(&e.buffIn)
	scanner.Split(e.scanWords)

	for scanner.Scan() {
		word := scanner.Text()

		// Remove any non-convertable characters.
		cleanWord := e.sanitizeWord([]byte(word))

		// Skip entirely non-printable word.
		if len(cleanWord) == 0 {
			continue
		}

		// Current word is either continuation from previus write,
		// or entirely new word.
		if e.notFirstWord {
			if e.leadingSpace {
				e.buffOut.WriteString("//")
			} else {
				e.buffOut.WriteString("/")
			}
		}
		e.notFirstWord = true

		for i, r := range cleanWord {
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

func (e *Encoder) sanitizeWord(s []byte) []rune {
	var cleanWord []rune
	for i, width := 0, 0; i < len(s); i += width {
		var r rune
		r, width = utf8.DecodeRune(s[i:])
		if e.isConvertable(r) {
			cleanWord = append(cleanWord, unicode.ToUpper(r))
		}
	}
	return cleanWord
}

func (e *Encoder) isConvertable(r rune) bool {
	if e.convertable == "" {
		e.initConvertable()
	}
	r = unicode.ToUpper(r)
	return strings.ContainsRune(e.convertable, r)
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

func (e *Encoder) scanWords(data []byte, atEOF bool) (
	advance int, token []byte, err error) {
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
			return i + width, data[start:i], nil
		}
	}
	// If we're at EOF, we have a final, non-empty, non-terminated word.
	// Return it.
	if atEOF && len(data) > start {
		e.trailingSpace = false
		// Scanner has to advance to the end of the data slice.
		returnLen := len(data)
		// Check for runes split at the end.
		data, broken := e.checkSplit(data[start:])
		e.buffIn.Write(broken)
		return returnLen, data, nil
	}

	e.trailingSpace = e.leadingSpace
	// Request more data.
	return start, nil, nil
}

func (e *Encoder) checkSplit(data []byte) (full, broken []byte) {
	var lastStart int
	for i := len(data) - 1; i >= 0; i-- {
		// Return broken rune for next call to encodeBuffer.
		if utf8.RuneStart(data[i]) {
			lastStart = i
			break
		}
	}
	if !utf8.FullRune(data[lastStart:]) {
		return data[:lastStart], data[lastStart:]
	}
	return data, nil
}
