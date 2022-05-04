package goblet

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
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
		n += v.writeByte(b)
		for i := 1; i < v.charLen; i++ {
			n += v.writeByte(v.padChar)
		}
	}
	return n, nil
}

// writeByte writes one or two bytes to buffer,
// taking into account that a blank character is inserted every 8 bytes for formatting
func (v *BinaryViewer) writeByte(b byte) (n int) {
	// bytes#Buffer.WriteByte() always returns nil error
	v.Buffer.WriteByte(b)
	n++
	if v.withoutBlankLen()%(v.lineBytes*8) == 0 {
		v.Buffer.WriteByte('\n')
		n++
	} else if v.withoutBlankLen()%8 == 0 {
		v.Buffer.WriteByte(' ')
		n++
	}
	return n
}

// withoutBlankLen returns the length excluding blank chars for the calculation of line break position
func (v *BinaryViewer) withoutBlankLen() int {
	return v.Buffer.Len() - (v.Buffer.Len() / 9)
}

func EncodeBit(src []byte) []byte {
	var s strings.Builder
	for _, b := range src {
		// strings#Builder.WriterString() always returns nil error
		s.WriteString(fmt.Sprintf("%08b", b))
	}
	return []byte(s.String())
}

func EncodeHex(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func EncodeBase64(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}
