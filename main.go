package main

import (
  "os"
  "fmt"
  "github.com/mig-hub/nesrom/nesrom"
)

func main() {

  if len(os.Args) < 2 {
    fmt.Println("Error: command missing")
    os.Exit(1)
  }

  if len(os.Args) < 3 {
    fmt.Println("Error: rom file missing")
    os.Exit(1)
  }

  f, err := os.Open(os.Args[2])
  if err != nil {
    fmt.Println("Error: cannot open ROM file", os.Args[2])
    os.Exit(1)
  }
  defer f.Close()

  rawHeader := make([]byte, 16)
  _, err = f.Read(rawHeader)
  if err != nil {
    fmt.Println("Error: cannot read ROM file header")
    f.Close()
    os.Exit(1)
  }

  h := nesrom.NewHeader(rawHeader)
  if ! h.IsValid() {
    fmt.Println("Invalid ROM, magic signature does not match [ 4e 45 53 1a ]")
    f.Close()
    os.Exit(1)
  }

  switch os.Args[1] {
  case "check":
    // Failed further up if not
    fmt.Println("Valid ROM")
  case "hexheader":
    printHexHeader(h)
  case "header":
    printHeader(h)
  case "tiles":
    f.Seek(h.CHROffset(), 0)
    printTiles(f)
  default:
    fmt.Println("Error: command unknown", os.Args[1])
    f.Close()
    os.Exit(1)
  }

}

func printHexHeader(h *nesrom.Header) {
  fmt.Printf("% x", h.Raw)
  fmt.Printf("\n")
}

func printHeader(h *nesrom.Header) {
  h.ParseFlags()
  fmt.Println("Format:            ", h.FormatName())
  fmt.Println("PGR-ROM banks:     ", h.PRGBanks(), "* 16KB")
  fmt.Println("CHR-ROM banks:     ", h.CHRBanks(), "* 8KB")
  fmt.Println("RAM banks:         ", h.RAMBanks(), "* 8KB")
  fmt.Println("Mirroring:         ", h.MirroringName())
  fmt.Println("Battery-backed RAM:", yesNo(h.BatteryBackedRAM))
  fmt.Println("Trainer:           ", yesNo(h.HasTrainer))
  fmt.Println("4 screen mirroring:", yesNo(h.FourScreenMirroring))
  fmt.Println("VS Unisystem      :", yesNo(h.VSUnisystem))
  fmt.Println("Playchoice-10 8KB :", yesNo(h.Playchoice10))
  fmt.Println("Mapper:            ", h.MapperName())
}

func yesNo(b bool) string {
  if b {
    return "yes"
  } else {
    return "no"
  }
}

var drawPixel = []string{"  ", "\u2591\u2591", "\u2593\u2593", "\u2587\u2587" }

func printTiles(f *os.File) {
  raw := make([]byte, 16)
  for {
    n, err := f.Read(raw)
    if n != 16 {
      break
    }
    if err != nil {
      panic("Ahhhhh")
    }
    tile := nesrom.NewTile(raw)
    for row:=0; row<8; row++ {
      for _, pixel := range tile.Pixels[8*row:8*row+8] {
        fmt.Printf(drawPixel[pixel])
      }
      fmt.Printf("\n")
    }
    fmt.Println("")
  }
}

