package main

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"os"
)

var (
	binFlag = flag.Bool("b", false, "use bits format")
	hexFlag = flag.Bool("hex", false, "use hex format")
	b64Flag = flag.Bool("b64", false, "use base64 format")
	widthFlag = flag.Int("w", 6, "width")
	separateFlag = flag.Int("s", 8, "separate")
)

func init() {
	flag.BoolVar(binFlag, "bits", false, "use bits format")
	flag.BoolVar(hexFlag, "hexadecimal", false, "use hex format")
	flag.BoolVar(b64Flag, "base64", false, "use base64 format")
	flag.IntVar(widthFlag, "width", 6, "width")
	flag.IntVar(separateFlag, "separate", 8, "separate")
}

func main() {
	flag.Parse()
	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("required filename")
		os.Exit(1)
	}

	f, err := os.Open(flag.Args()[0])
	if err != nil {
		fmt.Printf("failed to open file: %v", err)
		os.Exit(1)
	}
	defer f.Close()

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(f); err != nil {
		fmt.Printf("failed to read file: %v", err)
		os.Exit(1)
	}
	b := buf.Bytes()

	var encoded string
	switch {
	case *binFlag:
		for _, v := range b {
			encoded += fmt.Sprintf("%08b", v)
		}
	case *b64Flag:
		encoded = base64.StdEncoding.EncodeToString(b)
	default:
		encoded = hex.EncodeToString(b)
	}

	var s string
	for i, v := range encoded {
		s += string(v)
		if (i+1)%(*widthFlag**separateFlag) == 0 {
			s += "\n"
		} else if (i+1) % *separateFlag == 0 {
			s += " "
		}
	}
	fmt.Println(s)
}
