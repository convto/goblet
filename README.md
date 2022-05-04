# goblet

This is a binary viewer made by Go.

It reads binaries from stdin or file and outputs them in bit, hex, or base64 format.

## Installation

```
$ go install github.com/convto/goblet/cmd/goblet
```

## Usage

goblet reads binaries from a stdin or file and outputs them in any text encoding.  
Here is a simple example

```
$ goblet -bit test.bin 
00001000 10111001 01100000 00010000 10110010 10010010
00000100 00011000 00000001 
```

For multi-bit per character encodings such as base64 and hex, the output is padded.

```
$ goblet -hex test.bin
0___8___ b___9___ 6___0___ 1___0___ b___2___ 9___2___
0___4___ 1___8___ 0___1___ 
```

```
$ goblet -base64 test.bin
C_____L_ ____l___ __g_____ E_____L_ ____K___ __S_____
B_____B_ ____g___ __B_____
```

Currently supported options are

```
$ goblet -h
Usage of goblet:
  -base64
    	use base64 format
  -bit
    	use bit format
  -hex
    	use hex format
  -w int
    	width (default 6)
```