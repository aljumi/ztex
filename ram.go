package ztex

import (
	"fmt"
	"strings"
)

// RAMSize indicates the amount of RAM available on the module.
type RAMSize uint8

// String returns a human-readable representation of the RAM size.
func (r RAMSize) String() string {
	return binaryPrefix((uint64(r&0xf0))<<((uint(r&0xf))+16), "B")
}

// Number returns a raw numeric representation of the RAM size.
func (r RAMSize) Number() uint8 { return uint8(r) }

// RAMType indicates the type of RAM available on the module.
type RAMType uint8

// String returns a human-readable representation of the RAM type.
func (r RAMType) String() string {
	switch r {
	case 1:
		return "DDR-200 SDRAM"
	case 2:
		return "DDR-266 SDRAM"
	case 3:
		return "DDR-333 SDRAM"
	case 4:
		return "DDR-400 SDRAM"
	case 5:
		return "DDR2-400 SDRAM"
	case 6:
		return "DDR2-533 SDRAM"
	case 7:
		return "DDR2-667 SDRAM"
	case 8:
		return "DDR2-800 SDRAM"
	case 9:
		return "DDR2-1066 SDRAM"
	case 10:
		return "DDR3-800 SDRAM"
	default:
		return "Unknown"
	}
}

// RAMConfig indicates the size and type of the RAM in the module.
type RAMConfig struct {
	RAMSize
	RAMType
}

// String returns a human-readable representation of the RAM configuration.
func (r RAMConfig) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Size(%v)", r.RAMSize))
	x = append(x, fmt.Sprintf("Type(%v)", r.RAMType))
	return strings.Join(x, ", ")
}
