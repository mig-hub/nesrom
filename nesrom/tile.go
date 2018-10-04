package nesrom

import (
  "math/bits"
)

type Tile struct {
  Raw []byte
  Pixels []byte
}

var offset = []uint8{1,2,4,8,16,32,64,128}

func NewTile(raw []byte) *Tile {
  t := &Tile{}
  t.Raw = raw
  t.Pixels = make([]byte, 64)
  var reversedLow, reversedHigh uint8
  for row := uint8(0); row<8; row++ {
    reversedLow = bits.Reverse8(uint8(raw[row]))
    reversedHigh = bits.Reverse8(uint8(raw[row+8]))
    for pixel := uint8(0); pixel<8; pixel++ {
      t.Pixels[8*row+pixel] = (reversedLow & offset[pixel]) >> pixel
      t.Pixels[8*row+pixel] = t.Pixels[8*row+pixel] | ((reversedHigh & offset[pixel]) >> (pixel-1))
    }
  }
  return t
}

