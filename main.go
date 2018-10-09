package main

import (
  "os"
  "fmt"
  "flag"
  "errors"
  "github.com/mig-hub/nesrom/nesrom"
)

// Flags
// Each set corresponds to a subcommand

var checkCommand = flag.NewFlagSet("check", flag.ContinueOnError)
var headerCommand = flag.NewFlagSet("header", flag.ContinueOnError)
var hexOption = headerCommand.Bool("x", false, "display header as hexadecimal")
var tilesCommand = flag.NewFlagSet("tiles", flag.ContinueOnError)
var colorOption = tilesCommand.Bool("c", false, "display tiles in color")
var flagsets = []*flag.FlagSet{checkCommand,headerCommand,tilesCommand}

// Main wraps another function which returns an int
// This is so that it can exit and still execute defered calls

func main() {
  os.Exit(mainExitCode())
}

// The C-like main function which returns an int

func mainExitCode() int {

  // Make sure there is a command

  if len(os.Args) < 2 {
    fmt.Fprintln(os.Stderr, "Error: command missing")
    return 1
  }

  // Set flag, file and header struct

  fs, err := setFlagSet() 
  if err != nil {
    printErr(err)
    return 1
  }

  f, err := setFile(fs)
  if err != nil {
    printErr(err)
    return 1
  }
  defer f.Close()

  h, err := setHeader(f)
  if err != nil {
    printErr(err)
    return 1
  }

  // Handle command

  switch fs.Name() {
  case "check":
    // Fails before when invalid
    fmt.Println("Valid ROM")
  case "header":
    if *hexOption {
      printHexHeader(h)
    } else {
      printHeader(h)
    }
  case "tiles":
    f.Seek(h.CHROffset(), 0)
    printTiles(f)
  }

  return 0
}

func setFlagSet() (*flag.FlagSet, error) {

  var fs *flag.FlagSet

  // Find the flag with the subcommand name
  for _, current_fs := range flagsets {
    if current_fs.Name() == os.Args[1] {
      fs = current_fs
      break
    }
  }

  if fs == nil {
    return nil, errors.New(fmt.Sprint("Error: command unknown ", os.Args[1]))
  } else {
    fs.Parse(os.Args[2:])
    return fs, nil
  }

}

func setFile(fs *flag.FlagSet) (*os.File, error) {

  // Make sure there is a filename and only one

  if len(fs.Args()) != 1 {
    return nil, errors.New("Error: one ROM file should be passed as argument")
  }

  // Open ROM file

  f, err := os.Open(fs.Arg(0))
  if err != nil {
    return nil, errors.New(fmt.Sprint("Error: cannot open ROM file ", fs.Arg(0)))
  }

  return f, nil
}

func setHeader(f *os.File) (*nesrom.Header, error) {

  rawHeader := make([]byte, 16)
  _, err := f.Read(rawHeader)
  if err != nil {
    return nil, errors.New("Error: cannot read ROM file header")
  }

  h := nesrom.NewHeader(rawHeader)
  if ! h.IsValid() {
    return nil, errors.New("Invalid ROM, magic signature does not match [ 4e 45 53 1a ]")
  }

  return h, nil
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

// Chars used to represent pixels in the terminal

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

// Helpers

func printErr(s error) {
  fmt.Fprintln(os.Stderr, s)
}

func yesNo(b bool) string {
  if b {
    return "yes"
  } else {
    return "no"
  }
}

