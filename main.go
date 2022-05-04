package main

import (
	"flag"
	"fmt"
	"io"
	"os"
)

var (
	binFlag   = flag.Bool("bit", false, "use bit format")
	b64Flag   = flag.Bool("base64", false, "use base64 format")
	hexFlag   = flag.Bool("hex", false, "use hex format")
	widthFlag = flag.Int("w", 6, "width")
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

	in, err := io.ReadAll(r)
	if err != nil {
		fmt.Printf("failed to read input: %v", err)
		os.Exit(1)
	}

	var viewer *BinaryViewer
	switch {
	case *binFlag:
		buf := make([]byte, 0, ViewerCap(len(in), PadSizeBit))
		viewer = NewBinaryViewer(buf, PadSizeBit, *widthFlag, DefaultPadChar)
		if _, err := viewer.Write(encodeBits(in)); err != nil {
			fmt.Printf("failed to write bits to binary viewer: %v", err)
			os.Exit(1)
		}
	case *b64Flag:
		buf := make([]byte, 0, ViewerCap(len(in), PadSizeBase64))
		viewer = NewBinaryViewer(buf, PadSizeBase64, *widthFlag, DefaultPadChar)
		if _, err := viewer.Write(encodeBase64(in)); err != nil {
			fmt.Printf("failet to write base64 to binary viewer: %v", err)
			os.Exit(1)
		}
	case *hexFlag:
		fallthrough
	default:
		buf := make([]byte, 0, ViewerCap(len(in), PadSizeHex))
		viewer = NewBinaryViewer(buf, PadSizeHex, *widthFlag, DefaultPadChar)
		if _, err := viewer.Write(encodeHex(in)); err != nil {
			fmt.Printf("failet to write hex to binary viewer: %v", err)
			os.Exit(1)
		}
	}

	io.Copy(os.Stdout, viewer)
}
