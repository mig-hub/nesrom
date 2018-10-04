package main

import (
  "os"
  "fmt"
  "flag"
  "github.com/mig-hub/nesrom/nesrom"
)

// Flags

var checkCommand = flag.NewFlagSet("checkCommand", flag.ContinueOnError)
var headerCommand = flag.NewFlagSet("headerCommand", flag.ContinueOnError)
var hexOption = headerCommand.Bool("x", false, "display header as hexadecimal")
var tilesCommand = flag.NewFlagSet("tilesCommand", flag.ContinueOnError)
var colorOption = tilesCommand.Bool("c", false, "display tiles in color")

func main() {

  // Make sure there is a command

  if len(os.Args) < 2 {
    fmt.Println("Error: command missing")
    os.Exit(2)
  }

  // Handle command

  var (
    f *os.File
    h *nesrom.Header
  )

  switch os.Args[1] {
  case "check":
    f, h = initFileAndHeader(checkCommand)
    fmt.Println("Valid ROM")
  case "header":
    f, h = initFileAndHeader(headerCommand)
    if *hexOption {
      printHexHeader(h)
    } else {
      printHeader(h)
    }
  case "tiles":
    f, h = initFileAndHeader(tilesCommand)
    f.Seek(h.CHROffset(), 0)
    printTiles(f)
  default:
    fmt.Println("Error: command unknown", os.Args[1])
    f.Close()
    os.Exit(2)
  }

  if f != nil {
    f.Close()
  }

}

func initFileAndHeader(fs *flag.FlagSet) (*os.File, *nesrom.Header) {

  fs.Parse(os.Args[2:])

  // Make sure there is a filename and only one

  if len(fs.Args()) != 1 {
    fmt.Println("Error: one ROM file should be passed as argument")
    os.Exit(2)
  }

  // Open ROM file

  f, err := os.Open(fs.Arg(0))
  if err != nil {
    fmt.Println("Error: cannot open ROM file", fs.Arg(0))
    os.Exit(2)
  }

  // ROM Header

  rawHeader := make([]byte, 16)
  _, err = f.Read(rawHeader)
  if err != nil {
    fmt.Println("Error: cannot read ROM file header")
    f.Close()
    os.Exit(2)
  }

  h := nesrom.NewHeader(rawHeader)
  if ! h.IsValid() {
    fmt.Println("Invalid ROM, magic signature does not match [ 4e 45 53 1a ]")
    f.Close()
    os.Exit(2)
  }

  return f, h

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

// var drawPixelUnicode = []string{"  ", "\u2591\u2591", "\u2593\u2593", "\u2587\u2587" }
var drawPixel = []string{"  ", "--", "//", "##" }
var drawPixelColor = []string{"\033[30;40m  \033[0m", "\033[34;44m--\033[0m", "\033[31;41m//\033[0m", "\033[37;47m##\033[0m" }

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
        if *colorOption {
          fmt.Printf(drawPixelColor[pixel])
        } else {
          fmt.Printf(drawPixel[pixel])
        }
      }
      fmt.Printf("\n")
    }
    fmt.Println("")
  }
}

