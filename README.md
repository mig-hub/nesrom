NESROM
======

The `go get` command will automatically fetch the dependencies listed above, 
compile the binary and place it in your `$GOPATH/bin` directory.

`go get github.com/mig-hub/nesrom`

This is still a work in progress but already implements the following commands 
for working with NES ROM files:

`nesrom check zelda.nes` Prints a message and set the return status `0` if it is 
a valid NES ROM, `1` if it is not.

`nesrom hexheader zelda.nes` Prints the hexadecimal representation of the 
header.

`nesrom header zelda.nes` Prints the info contained in the header, including the 
size of PRG and CHR, the mirroring orientation, etc.

`nesrom tiles zelda.nes` Prints the CHR banks in a text format which reveals the 
tiles in the terminal using Unicode shading characters.

