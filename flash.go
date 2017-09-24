package ztex

import (
	"fmt"
	"strings"
)

// FlashEnabled indicates whether or not the flash is enabled.
type FlashEnabled uint8

// String returns a human-readable description of whether or not the
// flash is enabled.
func (f FlashEnabled) String() string {
	switch f {
	case 0:
		return "Disabled"
	case 1:
		return "Enabled"
	default:
		return "Unknown"
	}
}

// FlashSector represents the size of a sector in the flash.
type FlashSector [2]uint8

// String returns a human-readable description of the size of a sector
// in the flash.
func (f FlashSector) String() string { return binaryPrefix(f.Number(), "B") }

// Number returns the size of a sector in the flash (in bytes).
func (f FlashSector) Number() uint64 {
	z := uint64(bytesToUint16(f))
	if z&0x8000 != 0 {
		z = 1 << (z & 0x7fff)
	}
	return z
}

// FlashCount represents the number of sectors in the flash.
type FlashCount [4]uint8

// String returns a human-readable description of the number of sectors
// in the flash.
func (f FlashCount) String() string { return fmt.Sprintf("%v", f.Number()) }

// Number returns the number of sectors in the flash.
func (f FlashCount) Number() uint32 { return bytesToUint32(f) }

// FlashError represents the error code in the flash.
type FlashError uint8

// String returns a human-readable description of the flash error code.
func (f FlashError) String() string {
	switch f {
	case 0:
		return "None"
	case 1:
		return "Command Error"
	case 2:
		return "Timeout Error"
	case 3:
		return "Busy Error"
	case 4:
		return "Pending Error"
	case 5:
		return "Read Error"
	case 6:
		return "Write Error"
	case 7:
		return "Unsupported Error"
	case 8:
		return "Runtime Error"
	default:
		return fmt.Sprintf("Unknown Error [%v]", uint8(f))
	}
}

// FlashStatus indicates the current status of the flash.
type FlashStatus struct {
	FlashEnabled
	FlashSector
	FlashCount
	FlashError
}

// String returns a human-readable description of the flash status.
func (f FlashStatus) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Enabled(%v)", f.FlashEnabled))
	x = append(x, fmt.Sprintf("Sector(%v)", f.FlashSector))
	x = append(x, fmt.Sprintf("Count(%v)", f.FlashCount))
	x = append(x, fmt.Sprintf("Error(%v)", f.FlashError))
	return strings.Join(x, ", ")
}
