package main

import (
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"strings"
)

const (
	bitLen  = 1
	hexLen  = 4
	b64Len  = 6
	sepChar = "_"
)

func encodeBits(b []byte) string {
	var s string
	for _, v := range b {
		s += fmt.Sprintf("%08b", v)
	}
	return s
}

func encodeHex(b []byte) string {
	var s string
	for _, r := range hex.EncodeToString(b) {
		s += string(r)
		s += strings.Repeat(sepChar, hexLen-1)
	}
	return s
}

func encodeBase64(b []byte) string {
	var s string
	for _, r := range base64.StdEncoding.EncodeToString(b) {
		s += string(r)
		s += strings.Repeat(sepChar, b64Len-1)
	}
	return s
}
