package goblet

import (
	"bytes"
)

const (
	DefaultByteWidth = 6
	DefaultPadChar   = '_'
	CharLenBit       = 1
	CharLenHex       = 4
	CharLenBase64    = 6
)

// BinaryViewer provides I/O for viewing text-encoded binaries,
// write encoded binaries, you can read the formatted contents
type BinaryViewer struct {
	*bytes.Buffer
	lineBytes int  // used to determine line breaks when writing
	charLen   int  // number of bits per encoded character
	padChar   byte // supports only one ascii char
}

// viewerCap returns a convenient cap to BinaryViewer
func viewerCap(n, pad int) int {
	padded := n + n*pad
	return padded + padded/8 // to insert a blank character every 8 bytes
}

// NewBinaryViewer creates binary viewer with options
// charLen is required, other parameter are optional
// TODO: refactor to functional option pattern
func NewBinaryViewer(srcLen, charLen, byteWidth int, padChar byte) *BinaryViewer {
	if byteWidth == 0 {
		byteWidth = DefaultByteWidth
	}
	if padChar == 0 {
		padChar = DefaultPadChar
	}
	buf := make([]byte, 0, viewerCap(srcLen, charLen))
	return &BinaryViewer{
		Buffer:    bytes.NewBuffer(buf),
		lineBytes: byteWidth,
		charLen:   charLen,
		padChar:   padChar,
	}
}

func (v *BinaryViewer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		v.writeByte(b)
		n++
		for i := 1; i < v.charLen; i++ {
			v.writeByte(v.padChar)
		}
	}
	return n, nil
}

// writeByte writes one or two bytes to buffer,
// taking into account that a blank character is inserted every 8 bytes for formatting
func (v *BinaryViewer) writeByte(b byte) {
	// bytes#Buffer.WriteByte() always returns nil error
	v.Buffer.WriteByte(b)
	if v.withoutBlankLen()%(v.lineBytes*8) == 0 {
		v.Buffer.WriteByte('\n')
	} else if v.withoutBlankLen()%8 == 0 {
		v.Buffer.WriteByte(' ')
	}
}

// withoutBlankLen returns the length excluding blank chars for the calculation of line break position
func (v *BinaryViewer) withoutBlankLen() int {
	return v.Buffer.Len() - (v.Buffer.Len() / 9)
}
