package main

import (
	"encoding/base64"
	"encoding/hex"
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/convto/bit"
	"github.com/convto/goblet"
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

	var viewer *goblet.BinaryViewer
	switch {
	case *binFlag:
		viewer = goblet.NewBinaryViewer(len(in), goblet.CharLenBit, *widthFlag, goblet.DefaultPadChar)
		if _, err := bit.NewEncoder(viewer).Write(in); err != nil {
			fmt.Printf("failed to write bits to binary viewer: %v", err)
			os.Exit(1)
		}
	case *b64Flag:
		viewer = goblet.NewBinaryViewer(len(in), goblet.CharLenBase64, *widthFlag, goblet.DefaultPadChar)
		enc := base64.NewEncoder(base64.StdEncoding, viewer)
		defer enc.Close()
		if _, err := enc.Write(in); err != nil {
			fmt.Printf("failet to write base64 to binary viewer: %v", err)
			os.Exit(1)
		}
	case *hexFlag:
		fallthrough
	default:
		viewer = goblet.NewBinaryViewer(len(in), goblet.CharLenHex, *widthFlag, goblet.DefaultPadChar)
		if _, err := hex.NewEncoder(viewer).Write(in); err != nil {
			fmt.Printf("failet to write hex to binary viewer: %v", err)
			os.Exit(1)
		}
	}

	io.Copy(os.Stdout, viewer)
}
