package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	DefaultByteWidth = 6
	PadSizeBit       = 0
	PadSizeHex       = 3
	PadSizeBase64    = 5
	DefaultPadChar   = '_'
)

// BinaryViewer receives string encoded binary and handles the formatted I/O.
type BinaryViewer struct {
	*bytes.Buffer
	lineBytes int // used to determine line breaks when writing
	padSize   int
	padChar   byte // supports only one ascii char
}

// ViewerCap returns cap of bytes required by BinaryViewer.
func ViewerCap(n, pad int) int {
	padded := n + n*pad
	return padded + padded/8 // to insert a blank character every 8 bytes.
}

// NewBinaryViewer creates binary viewer with options.
// buf and padSize are required, byteWidth and padChar are optional.
// TODO: refactor to functional option pattern
func NewBinaryViewer(buf []byte, padSize, byteWidth int, padChar byte) *BinaryViewer {
	if byteWidth == 0 {
		byteWidth = DefaultByteWidth
	}
	if padChar == 0 {
		padChar = DefaultPadChar
	}
	return &BinaryViewer{
		Buffer:    bytes.NewBuffer(buf),
		lineBytes: byteWidth,
		padSize:   padSize,
		padChar:   padChar,
	}
}

func (v *BinaryViewer) Write(p []byte) (n int, err error) {
	for _, b := range p {
		n += v.writeByte(b)
		for i := 1; i <= v.padSize; i++ {
			n += v.writeByte(v.padChar)
		}
	}
	return n, nil
}

// writeByte writes one or two bytes to buffer,
// taking into account that a blank character is inserted every 8 bytes for formatting.
func (v *BinaryViewer) writeByte(b byte) (n int) {
	// bytes#Buffer.WriteByte() always returns nil error
	v.Buffer.WriteByte(b)
	n++
	if v.PaddedLen()%(v.lineBytes*8) == 0 {
		v.Buffer.WriteByte('\n')
		n++
	} else if v.PaddedLen()%8 == 0 {
		v.Buffer.WriteByte(' ')
		n++
	}
	return n
}

func (v *BinaryViewer) PaddedLen() int {
	return v.Buffer.Len() - v.BlankLen()
}

func (v *BinaryViewer) BlankLen() int {
	return v.Buffer.Len() / 9
}

func encodeBits(src []byte) []byte {
	var s strings.Builder
	for _, b := range src {
		// strings#Builder.WriterString() always returns nil error
		s.WriteString(fmt.Sprintf("%08b", b))
	}
	return []byte(s.String())
}

func encodeHex(src []byte) []byte {
	dst := make([]byte, hex.EncodedLen(len(src)))
	hex.Encode(dst, src)
	return dst
}

func encodeBase64(src []byte) []byte {
	dst := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
	base64.StdEncoding.Encode(dst, src)
	return dst
}
