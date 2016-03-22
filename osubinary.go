// Package osubinary allows for binary reading of data encoded with an OsuBinaryReader. Or something like that.
package osubinary

import (
	"encoding/binary"
	"io"
	"reflect"
)

// BinaryReader is the data to export into various variables.
type BinaryReader struct {
	OsuReader
}

// OsuReader is a wrapped io.Reader that can read files in the "osu! format".
type OsuReader struct {
	io.Reader
}

// New creates a new OsuReader.
func New(r io.Reader) OsuReader {
	return OsuReader{
		Reader: r,
	}
}

// Unmarshal is a wrapper for b.OsuRead. Do not use this - it's deprecated.
func (b OsuReader) Unmarshal(out ...interface{}) error {
	return b.OsuRead(out...)
}

// OsuRead decodes some OsuReader into various data structers passed as arguments.
func (b OsuReader) OsuRead(out ...interface{}) error {
	for _, vOriginal := range out {
		switch v := vOriginal.(type) {
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

		case *[]int32:
			var arrlen uint16
			err := binary.Read(b.Reader, binary.LittleEndian, &arrlen)
			if err != nil {
				return err
			}
			finalArr := make([]int32, arrlen)
			for i := 0; i < int(arrlen); i++ {
				err = binary.Read(b.Reader, binary.LittleEndian, &finalArr[i])
				if err != nil {
					return err
				}
			}
			*v = finalArr

		case *[]uint32:
			var arrlen uint16
			err := binary.Read(b.Reader, binary.LittleEndian, &arrlen)
			if err != nil {
				return err
			}
			finalArr := make([]uint32, arrlen)
			for i := 0; i < int(arrlen); i++ {
				err = binary.Read(b.Reader, binary.LittleEndian, &finalArr[i])
				if err != nil {
					return err
				}
			}
			*v = finalArr

		default:
			err := binary.Read(b.Reader, binary.LittleEndian, vOriginal)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// ReadString takes a string out of the reader.
func (b OsuReader) ReadString() ([]byte, error) {
	return ReadString(b.Reader)
}

// OsuWriter is an io.Writer for binary data, and the custom datatypes of osu! databases.
type OsuWriter struct {
	io.Writer
}

// NewWriter creates a new OsuWriter.
func NewWriter(w io.Writer) OsuWriter {
	return OsuWriter{
		Writer: w,
	}
}

// OsuWrite writes data to the writer in binary using the "osu format".
func (o OsuWriter) OsuWrite(data ...interface{}) error {
	for _, vOriginal := range data {
		switch v := vOriginal.(type) {
		case *string, *[]byte, *[]int32, *[]uint32:
			o.OsuWrite(reflect.TypeOf(v).Elem())

		case string:
			_, err := o.Writer.Write(MakeString(v))
			if err != nil {
				return err
			}
		case []byte:
			_, err := o.Writer.Write(MakeString(string(v)))
			if err != nil {
				return err
			}

		case []int32:
			err := binary.Write(o.Writer, binary.LittleEndian, uint16(len(v)))
			if err != nil {
				return err
			}
			for _, el := range v {
				err = binary.Write(o.Writer, binary.LittleEndian, el)
				if err != nil {
					return err
				}
			}

		case []uint32:
			err := binary.Write(o.Writer, binary.LittleEndian, uint16(len(v)))
			if err != nil {
				return err
			}
			for _, el := range v {
				err = binary.Write(o.Writer, binary.LittleEndian, el)
				if err != nil {
					return err
				}
			}

		default:
			err := binary.Write(o.Writer, binary.LittleEndian, vOriginal)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
