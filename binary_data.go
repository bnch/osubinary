// Package osubinary allows for binary reading of data encoded with an OsuBinaryReader. Or something like that.
package osubinary

import (
	"io"
	"fmt"
	"encoding/binary"
	"github.com/bnch/uleb128"
	"errors"
)

// BinaryData is the data to export into various variables.
type BinaryData struct {
	reader io.Reader
}

// New creates a new BinaryData.
func New(r io.Reader) BinaryData {
	return BinaryData{
		reader: r,
	}
}

// Unmarshal decodes some BinaryData into various data structers passed as arguments.
func (b BinaryData) Unmarshal(out ...interface{}) error {
	for _, v := range out {
		switch v := v.(type) {

		case *string:
			d, err := b.ReadString()
			if err != nil {
				return err
			}
			*v = string(d)
		case *[]byte:
			d, err := b.ReadString()
			if err != nil {
				return err
			}
			*v = d

		case *int8:
			var finVal int8
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal
		case *uint8:
			var finVal uint8
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal

		case *int16:
			var finVal int16
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal
		case *uint16:
			var finVal uint16
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal

		case *int32:
			var finVal int32
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal
		case *uint32:
			var finVal uint32
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal

		case *int64:
			var finVal int64
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal
		case *uint64:
			var finVal uint64
			err := binary.Read(b.reader, binary.LittleEndian, &finVal)
			if err != nil {
				return err
			}
			*v = finVal

		case *[]uint32:
			var arrlen uint16
			err := binary.Read(b.reader, binary.LittleEndian, &arrlen)
			if err != nil {
				return err
			}
			finalArr := make([]uint32, arrlen)
			for i := 0; i < int(arrlen); i++ {
				err = binary.Read(b.reader, binary.LittleEndian, &finalArr[i])
				if err != nil {
					return err
				}
			}
			*v = finalArr

		case *[]int32:
			var arrlen uint16
			err := binary.Read(b.reader, binary.LittleEndian, &arrlen)
			if err != nil {
				return err
			}
			finalArr := make([]int32, arrlen)
			for i := 0; i < int(arrlen); i++ {
				err = binary.Read(b.reader, binary.LittleEndian, &finalArr[i])
				if err != nil {
					return err
				}
			}
			*v = finalArr

		default:
			return fmt.Errorf("osubinary.Unmarshal: type not supported (%T)", v)
		}
	}
	return nil
}

// ReadString takes a string out of the reader.
func (b BinaryData) ReadString() ([]byte, error) {
	bslice := make([]byte, 1)
	b.reader.Read(bslice)
	if bslice[0] != 11 {
		return []byte{}, errors.New("osubinary.ReadString: was expecting string, does not begin with byte 11")
	}
	strlen := uleb128.UnmarshalReader(b.reader)
	bslice = make([]byte, strlen)
	d, err := b.reader.Read(bslice)
	if d < strlen {
		return []byte{}, fmt.Errorf("osubinary.ReadString: unexpected end of string (expected to read %d bytes, read %d)", strlen, d)
	}
	return bslice, err
}
