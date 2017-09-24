package ztex

import (
	"fmt"
)

// BoardType indicates the board type associated with the device.
type BoardType uint8

// String returns a human-readable description of a board type.
func (b BoardType) String() string {
	switch b {
	case 1:
		return "ZTEX FPGA Module"
	case 2:
		return "ZTEX USB-FPGA Module (Cypress CY7C68013A EZ-USB FX2)"
	case 3:
		return "ZTEX USB3-FPGA Module (Cypress CYUSB3033 EZ-USB FX3S)"
	default:
		return Unknown
	}
}

// Number returns the raw representation of a board type.
func (b BoardType) Number() uint8 { return uint8(b) }

// BoardSeries indicates an entire generation of boards.  Currently only
// Series 1 and Series 2 boards are supported.
type BoardSeries uint8

// String returns a human-readable description of a board series.
func (b BoardSeries) String() string {
	switch b {
	case 1:
		return "1"
	case 2:
		return "2"
	default:
		return Unknown
	}
}

// Number returns the raw representation of a board series.
func (b BoardSeries) Number() uint8 { return uint8(b) }

// BoardNumber indicates a board in a series.
type BoardNumber uint8

// String returns a human-readable description of a board number.
func (b BoardNumber) String() string {
	switch {
	case b == 255:
		return Unknown
	default:
		return fmt.Sprintf("%d", uint8(b))
	}
}

// Number returns the raw representation of a board number.
func (b BoardNumber) Number() uint8 { return uint8(b) }

// BoardVariant indicates a variation on a board series and number.
type BoardVariant [2]byte

// String returns a human-readable description of a board variant.
func (b BoardVariant) String() string { return string(b.Bytes()) }

// Bytes returns the raw representation of a board variant.
func (b BoardVariant) Bytes() []byte {
	switch {
	case b[0] == 0:
		return []byte{}
	case b[1] == 0:
		return []byte{b[0]}
	default:
		return []byte{b[0], b[1]}
	}
}

// BoardVersion indicates the series, number, and variant of a module.
type BoardVersion struct {
	BoardSeries
	BoardNumber
	BoardVariant
}

// String returns a human-readable representation of the board version.
func (b BoardVersion) String() string {
	return fmt.Sprintf("%v.%v%v", b.BoardSeries, b.BoardNumber, b.BoardVariant)
}

// BoardConfig indicates the type, series, number, and variant of a ZTEX
// USB-FPGA module.  For example, a ZTEX USB3-FPGA 2.18b module would be
// represented by
//
//   BoardConfig{
//     BoardType: BoardType(3),
//     BoardVersion{
//       BoardSeries: BoardSeries(2),
//       BoardNumber: BoardNumber(18),
//       BoardVariant: BoardVariant([2]byte{0x62, 0x00}]),
//     },
//   }
//
// as a BoardConfig structure.
type BoardConfig struct {
	BoardType
	BoardVersion
}

// String returns a human-readable representation of a board version.
func (b BoardConfig) String() string {
	return fmt.Sprintf("Type %v, Version %v.%v%v", b.BoardType, b.BoardSeries, b.BoardNumber, b.BoardVariant)
}
