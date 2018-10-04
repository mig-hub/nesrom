package nesrom

import (
  "bytes"
)

type Header struct {
  Raw []byte
  VerticalMirroring bool
  BatteryBackedRAM bool
  HasTrainer bool
  FourScreenMirroring bool
  VSUnisystem bool
  Playchoice10 bool
  MapperNumber byte
  IsNES20 bool
  IsiNES bool
  IsArchaiciNES bool
}

func NewHeader(raw []byte) *Header {
  h := &Header{}
  h.Raw = raw
  return h
}

func (h *Header) IsValid() bool {
  return bytes.Equal(h.Raw[:4], []byte{0x4e, 0x45, 0x53, 0x1a}[:])
}

func (h *Header) ParseFlags() {
  // Flag 6
  h.VerticalMirroring = (h.Raw[6] & 1) != 0
  h.BatteryBackedRAM = (h.Raw[6] & 2) != 0
  h.HasTrainer = (h.Raw[6] & 4) != 0
  h.FourScreenMirroring = (h.Raw[6] & 8) != 0
  h.MapperNumber = (h.Raw[6] & 0xF0) >> 4
  // Flag 7
  h.VSUnisystem = (h.Raw[7] & 1) != 0
  h.Playchoice10 = (h.Raw[7] & 2) != 0
  h.MapperNumber = h.MapperNumber | (h.Raw[7] & 0xF0)
  // Format detection
  masked := h.Raw[7] & 0x0C
  if masked == 0x08 { // Something missing about size but I don't entirely get it
    h.IsNES20 = true
  } else if masked == 0x00 && h.Raw[12]==0 && h.Raw[13]==0 && h.Raw[14]==0 && h.Raw[15]==0 {
    h.IsiNES = true
  } else {
    h.IsArchaiciNES = true
  }
}

func (h *Header) PRGBanks() uint8 {
  return uint8(h.Raw[4])
}

func (h *Header) PRGOffset() int64 {
  var a int64 = 16 // header
  if h.HasTrainer {
    a += 256
  }
  return a
}

func (h *Header) CHRBanks() uint8 {
  return uint8(h.Raw[5])
}

func (h *Header) CHROffset() int64 {
  a := h.PRGOffset()
  a += int64(h.PRGBanks()) * 16 * 1024
  return a
}

func (h *Header) RAMBanks() uint8 {
  n := uint8(h.Raw[8])
  if n >= 1 {
    return n
  } else {
    return 1
  }
}

func (h *Header) FormatName() string {
  if h.IsNES20 {
    return "NES 2.0"
  }
  if h.IsiNES {
    return "iNES"
  }
  if h.IsArchaiciNES {
    return "Archaic iNES"
  }
  return "Not set, flags where probably not parsed"
}

func (h *Header) MirroringName() string {
  if h.VerticalMirroring { 
    return "vertical" 
  } else {
    return "horizontal"
  }
}

func (h *Header) MapperName() string {
  switch h.MapperNumber {
  case 0:
    return "NROM, no mapper"
  case 1:
    return "Nintendo MMC1"
  case 2:
    return "UNROM switch"
  case 3:
    return "CNROM switch"
  case 4:
    return "Nintendo MMC3"
  case 5:
    return "Nintendo MMC5"
  case 6:
    return "FFE F4xxx"
  case 7:
    return "AOROM switch"
  case 8:
    return "FFE F3xxx"
  case 9:
    return "Nintendo MMC2"
  case 10:
    return "Nintendo MMC4"
  case 11:
    return "ColorDreams chip"
  case 12:
    return "FFE F6xxx"
  case 15:
    return "100-in-1 switch"
  case 16:
    return "Bandai chip"
  case 17:
    return "FFE F8xxx"
  case 18:
    return "Jaleco SS8806 chip"
  case 19:
    return "Namcot 106 chip"
  case 20:
    return "Nintendo DiskSystem"
  case 21:
    return "Konami VRC4a"
  case 22:
    return "Konami VRC2a"
  case 23:
    return "Konami VRC2a"
  case 24:
    return "Konami VRC6"
  case 25:
    return "Konami VRC4b"
  case 32:
    return "Irem G-101 chip"
  case 33:
    return "Taito TC0190/TC0350"
  case 34:
    return "32 KB ROM switch"
  case 64:
    return "Tengen RAMBO-1 chip"
  case 65:
    return "Irem H-3001 chip"
  case 66:
    return "GNROM switch"
  case 67:
    return "SunSoft3 chip"
  case 68:
    return "SunSoft4 chip"
  case 69:
    return "SunSoft5 FME-7 chip"
  case 71:
    return "Camerica chip"
  case 78:
    return "Irem 74HC161/32-based"
  case 91:
    return "Pirate HK-SF3 chip"
  default:
    return "Unknown mapper"
  }
}

