package main

import (
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	binFlag      = flag.Bool("b", false, "use bits format")
	b64Flag      = flag.Bool("b64", false, "use base64 format")
	widthFlag    = flag.Int("w", 6, "width")
	separateFlag = flag.Int("s", 8, "separate")
)

func main() {
	flag.Parse()
	args := flag.Args()

	var r io.Reader
	if len(args) == 1 {
		f, err := os.Open(flag.Args()[0])
		if err != nil {
			fmt.Printf("failed to open file: %v", err)
			os.Exit(1)
		}
		defer f.Close()
		r = f
	} else {
		r = os.Stdin
	}

	var buf bytes.Buffer
	if _, err := buf.ReadFrom(r); err != nil {
		fmt.Printf("failed to read input: %v", err)
		os.Exit(1)
	}

	var encoded string
	switch {
	case *binFlag:
		encoded = encodeBits(buf.Bytes())
	case *b64Flag:
		encoded = encodeBase64(buf.Bytes())
	default:
		encoded = encodeHex(buf.Bytes())
	}

	var s string
	for i, v := range encoded {
		s += string(v)
		if (i+1)%(*widthFlag**separateFlag) == 0 {
			s += "\n"
		} else if (i+1)%*separateFlag == 0 {
			s += " "
		}
	}
	fmt.Println(s)
}
