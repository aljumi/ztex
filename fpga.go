package ztex

import (
	"fmt"
	"strings"
)

// FPGAType indicates which FPGA device is present.
type FPGAType [2]byte

// String returns a human-readable representation of an FPGA type.
func (f FPGAType) String() string {
	switch f.Number() {
	case 1:
		return "Xilinx Spartan-6 XC6SLX9"
	case 2:
		return "Xilinx Spartan-6 XC6SLX16"
	case 3:
		return "Xilinx Spartan-6 XC6SLX25"
	case 4:
		return "Xilinx Spartan-6 XC6SLX45"
	case 5:
		return "Xilinx Spartan-6 XC6SLX75"
	case 6:
		return "Xilinx Spartan-6 XC6SLX100"
	case 7:
		return "Xilinx Spartan-6 XC6SLX150"
	case 8:
		return "Xilinx Artix-7 XC7A35T"
	case 9:
		return "Xilinx Artix-7 XC7A50T"
	case 10:
		return "Xilinx Artix-7 XC7A75T"
	case 11:
		return "Xilinx Artix-7 XC7A100T"
	case 12:
		return "Xilinx Artix-7 XC7A200T"
	case 13:
		return "Xilinx Spartan-6 XC6SLX150 [x4]"
	case 14:
		return "Xilinx Artix-7 XC7A15T"
	default:
		return "Unknown"
	}
}

// Bytes returns a raw representation of an FPGA type.
func (f FPGAType) Bytes() []byte { return []byte{f[0], f[1]} }

// Number returns a numeric representation of an FPGA type.
func (f FPGAType) Number() uint16 { return (uint16(f[0]) << 0) | (uint16(f[1]) << 8) }

// FPGAPackage indicates the mechanical packaging of the FPGA.
type FPGAPackage uint8

// String returns a human-readable representation of the FPGA package.
func (f FPGAPackage) String() string {
	switch f {
	case 1:
		return "FTG256"
	case 2:
		return "CSG324"
	case 3:
		return "CSG484"
	case 4:
		return "FBG484"
	default:
		return "Unknown"
	}
}

// Number returns the raw numeric representation of an FPGA package.
func (f FPGAPackage) Number() uint8 { return uint8(f) }

// FPGAGrade indicates the speed grade, operating voltages, and
// temperature range of the FPGA.
type FPGAGrade [3]byte

// String returns a human-readable representation of the FPGA grade.
func (f FPGAGrade) String() string { return string(f.Bytes()) }

// Bytes returns the raw representation of the FPGA grade.
func (f FPGAGrade) Bytes() []byte {
	switch {
	case f[0] == 0:
		return []byte{}
	case f[1] == 0:
		return []byte{f[0]}
	case f[2] == 0:
		return []byte{f[0], f[1]}
	default:
		return []byte{f[0], f[1], f[2]}
	}
}

// FPGAConfig indicates the type, package, speed grade, etc. of the FPGA
// present in a device.
type FPGAConfig struct {
	FPGAType
	FPGAPackage
	FPGAGrade
}

// String returns a human-readable representation of the FPGA version.
func (f FPGAConfig) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Type(%v)", f.FPGAType))
	x = append(x, fmt.Sprintf("Package(%v)", f.FPGAPackage))
	x = append(x, fmt.Sprintf("Grade(%v)", f.FPGAGrade))
	return strings.Join(x, ", ")
}

// FPGAConfigured indicates whether or not the FPGA is configured.
type FPGAConfigured uint8

// String returns a human-readable description of the FPGA configuration
// indicator.
func (f FPGAConfigured) String() string {
	switch f {
	case 0:
		return "Configured"
	case 1:
		return "Unconfigured"
	default:
		return "Unknown"
	}
}

// Number returns the raw numeric representation of the FPGA configuration
// indicator.
func (f FPGAConfigured) Number() uint8 { return uint8(f) }

// Bool returns true if and only if the FPGA is configured.
func (f FPGAConfigured) Bool() bool { return f == 0 }

// FPGAChecksum represents the number of bytes
type FPGAChecksum uint8

// FPGATransferred represents the number of bytes transferred.
type FPGATransferred [4]uint8

// String returns a human-readable description of the number of bytes
// transferred.
func (f FPGATransferred) String() string {
	return fmt.Sprintf("%v", binaryPrefix(uint64(f.Number()), "B"))
}

// Number returns the number of bytes transferred.
func (f FPGATransferred) Number() uint32 { return bytesToUint32(f) }

// FPGAInit represents the number of INIT_B states.
type FPGAInit uint8

// FPGAResult represents the result of previous FPGA configuration.
type FPGAResult uint8

// String returns a human-readable description of the FPGA configuration
// result.
func (f FPGAResult) String() string {
	switch f {
	case 0:
		return "No Error"
	case 1:
		return "Already Configured Error"
	case 2:
		return "Flash Error"
	case 3:
		return "No Bitstream Error"
	case 4:
		return "Configuration Error"
	default:
		return "Unknown"
	}
}

// Bool returns true if and only if the result indicates that
// configuration was successful.
func (f FPGAResult) Bool() bool { return f == 0 }

// FPGASwapped represents the bit order of the FPGA bitstream.
type FPGASwapped uint8

// String returns a human-readable description of the bitstream bit order.
func (f FPGASwapped) String() string {
	switch f {
	case 0:
		return "Unswapped"
	case 1:
		return "Swapped"
	default:
		return "Unknown"
	}
}

// Number returns the raw numeric representation of the bitstream bit order.
func (f FPGASwapped) Number() uint8 { return uint8(f) }

// Bool returns true if and only if the bitstream bit order is swapped.
func (f FPGASwapped) Bool() bool { return f == 1 }

// FPGAStatus indicates the status of the FPGA.
type FPGAStatus struct {
	FPGAConfigured
	FPGAChecksum
	FPGATransferred
	FPGAInit
	FPGAResult
	FPGASwapped
}

// String returns a human-readable description of the FPGA status.
func (f FPGAStatus) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Configured(%v)", f.FPGAConfigured))
	x = append(x, fmt.Sprintf("Checksum(%v)", f.FPGAChecksum))
	x = append(x, fmt.Sprintf("Transferred(%v)", f.FPGATransferred))
	x = append(x, fmt.Sprintf("Init(%v)", f.FPGAInit))
	x = append(x, fmt.Sprintf("Result(%v)", f.FPGAResult))
	x = append(x, fmt.Sprintf("Swapped(%v)", f.FPGASwapped))
	return strings.Join(x, ", ")
}
