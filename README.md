NESROM
======

The `go get` command will automatically compile the binary and place it in your 
`$GOPATH/bin` directory.

`go get github.com/mig-hub/nesrom`

This is still a work in progress but already implements the following commands 
for working with NES ROM files:

- `nesrom check zelda.nes` Prints a message and set the return status `0` if it 
  is a valid NES ROM, `2` if it is not.
- `nesrom header zelda.nes` Prints the info contained in the header, including 
  the size of PRG and CHR, the mirroring orientation, etc.
- `nesrom header -x zelda.nes` Prints the hexadecimal representation of the 
  header.
- `nesrom tiles zelda.nes` Prints the CHR banks in a text format which reveals 
  the tiles in the terminal.
- `nesrom tiles -c zelda.nes` Same but in color with terminal escape sequences.

