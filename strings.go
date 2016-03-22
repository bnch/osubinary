package osubinary

import (
	"errors"
	"fmt"
	"io"

	"github.com/bnch/uleb128"
)

// ReadString takes a string in the "osu! format" out of an io.Reader.
func ReadString(r io.Reader) ([]byte, error) {
	bslice := make([]byte, 1)
	r.Read(bslice)
	if bslice[0] != 11 {
		return []byte{}, errors.New("osubinary.ReadString: was expecting string, does not begin with byte 11")
	}
	strlen := uleb128.UnmarshalReader(r)
	bslice = make([]byte, strlen)
	d, err := r.Read(bslice)
	if d < strlen {
		return []byte{}, fmt.Errorf("osubinary.ReadString: unexpected end of string (expected to read %d bytes, read %d)", strlen, d)
	}
	return bslice, err
}

// MakeString generates a string in the "osu! format".
// It is essentially made up by the 11 byte, an uleb int indicating the size of the string, and then the raw string.
func MakeString(s string) []byte {
	b := []byte(s)
	lenUleb := uleb128.Marshal(len(b))
	end := make([]byte, 1+len(lenUleb)+len(b))
	end[0] = 11
	copied := copy(end[1:], lenUleb)
	copy(end[copied+1:], b)
	return end
}
