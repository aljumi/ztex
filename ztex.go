// Package ztex manages ZTEX USB-FPGA modules.
package ztex

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gousb"
)

const (
	// VendorID is the ZTEX USB vendor ID (VID).
	VendorID = gousb.ID(0x221A)

	// ProductID is the standard ZTEX USB product ID (PID)
	ProductID = gousb.ID(0x0100)

	// Unknown is the string used to describe unknown values or components.
	Unknown = "Unknown"
)

func binaryPrefix(n uint64, unit string) string {
	switch {
	case n != 0 && n&(1<<30-1) == 0:
		return fmt.Sprintf("%vGi%v (%v%v)", n>>30, unit, n, unit)
	case n != 0 && n&(1<<20-1) == 0:
		return fmt.Sprintf("%vMi%v (%v%v)", n>>20, unit, n, unit)
	case n != 0 && n&(1<<10-1) == 0:
		return fmt.Sprintf("%vki%v (%v%v)", n>>10, unit, n, unit)
	default:
		return fmt.Sprintf("%v%v", n, unit)
	}
}

// ZTEXSize represents the number of bytes in a ZTEX descriptor.
type ZTEXSize uint8

// ZTEXVersion represents the version of a ZTEX descriptor.
type ZTEXVersion uint8

// ZTEXMagic indicates the presence of a ZTEX descriptor.
type ZTEXMagic [4]uint8

// String returns a human-readable description of the ZTEX magic bytes.
func (z ZTEXMagic) String() string { return string(z.Bytes()) }

// Bytes returns a raw representation of the ZTEX magic bytes.
func (z ZTEXMagic) Bytes() []byte { return []byte{z[0], z[1], z[2], z[3]} }

// ZTEXProduct represents a ZTEX product ID.
type ZTEXProduct [4]uint8

// String returns a human-readable description of the ZTEX product ID.
func (z ZTEXProduct) String() string {
	p := Unknown
	switch {
	case z[0] == 0 && z[1] == 0 && z[2] == 0 && z[3] == 0:
		p = "Default"
	case z[0] == 1:
		p = "Experimental"
	case z[0] == 10 && z[1] == 0 && z[2] == 1 && z[3] == 1:
		p = "ZTEX BTCMiner"
	case z[0] == 10 && z[1] == 11:
		p = "ZTEX USB-FPGA Module 1.2"
	case z[0] == 10 && z[1] == 12 && z[2] == 2 && (1 <= z[3] && z[3] <= 4):
		p = "NIT"
	case z[0] == 10 && z[1] == 12:
		p = "ZTEX USB-FPGA Module 1.11"
	case z[0] == 10 && z[1] == 13:
		p = "ZTEX USB-FPGA Module 1.15"
	case z[0] == 10 && z[1] == 14:
		p = "ZTEX USB-FPGA Module 1.15x"
	case z[0] == 10 && z[1] == 15:
		p = "ZTEX USB-FPGA Module 1.15y"
	case z[0] == 10 && z[1] == 16:
		p = "ZTEX USB-FPGA Module 2.16"
	case z[0] == 10 && z[1] == 17:
		p = "ZTEX USB-FPGA Module 2.13"
	case z[0] == 10 && z[1] == 18:
		p = "ZTEX USB-FPGA Module 2.01"
	case z[0] == 10 && z[1] == 19:
		p = "ZTEX USB-FPGA Module 2.04"
	case z[0] == 10 && z[1] == 20:
		p = "ZTEX USB Module 1.0"
	case z[0] == 10 && z[1] == 30:
		p = "ZTEX USB-XMEGA Module 1.0"
	case z[0] == 10 && z[1] == 40:
		p = "ZTEX USB-FPGA Module 2.02"
	case z[0] == 10 && z[1] == 41:
		p = "ZTEX USB-FPGA Module 2.14"
	case z[0] == 10:
		p = "ZTEX"
	}
	return fmt.Sprintf("%v.%v.%v.%v (%v)", z[0], z[1], z[2], z[3], p)
}

// Bytes returns a raw representation of the ZTEX product ID.
func (z ZTEXProduct) Bytes() []byte { return []byte{z[0], z[1], z[2], z[3]} }

// ZTEXFirmware indicates the version of the ZTEX firmware.
type ZTEXFirmware uint8

// ZTEXInterface indicates the version of the ZTEX interface.
type ZTEXInterface uint8

// ZTEXCapability indicates the capabilities supported by the ZTEX device.
type ZTEXCapability [6]uint8

// String returns a human-readable description of the ZTEX capabilities
// supported by the device.
func (z ZTEXCapability) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("EEPROM: %v", z.EEPROM()))
	x = append(x, fmt.Sprintf("FPGA Configuration: %v", z.FPGAConfiguration()))
	x = append(x, fmt.Sprintf("Flash Memory: %v", z.FlashMemory()))
	x = append(x, fmt.Sprintf("Debug Helper: %v", z.DebugHelper()))
	x = append(x, fmt.Sprintf("XMEGA: %v", z.XMEGA()))
	x = append(x, fmt.Sprintf("High Speed FPGA Configuration: %v", z.HighSpeedFPGAConfiguration()))
	x = append(x, fmt.Sprintf("MAC EEPROM: %v", z.MACEEPROM()))
	x = append(x, fmt.Sprintf("MultiFPGA: %v", z.MultiFPGA()))
	x = append(x, fmt.Sprintf("Temperature Sensor: %v", z.TemperatureSensor()))
	x = append(x, fmt.Sprintf("Flash Memory 2: %v", z.FlashMemory2()))
	x = append(x, fmt.Sprintf("FX3 Firmware: %v", z.FX3Firmware()))
	x = append(x, fmt.Sprintf("Debug Helper 2: %v", z.DebugHelper2()))
	x = append(x, fmt.Sprintf("Default Firmware Interface: %v", z.DefaultFirmwareInterface()))
	return strings.Join(x, ", ")
}

// Function cap returns true if and only if ZTEX capability i.j is
// supported by the device.
func (z ZTEXCapability) cap(i, j uint) bool { return z[i]&(1<<j) != 0 }

// EEPROM returns true if and only if the device has EEPROM support.
func (z ZTEXCapability) EEPROM() bool { return z.cap(0, 0) }

// FPGAConfiguration returns true if and only if the device has basic
// FPGA configuration support.
func (z ZTEXCapability) FPGAConfiguration() bool { return z.cap(0, 1) }

// FlashMemory returns true if and only if the device has flash memory.
func (z ZTEXCapability) FlashMemory() bool { return z.cap(0, 2) }

// DebugHelper returns true if and only if the device has basic debug
// helper support.
func (z ZTEXCapability) DebugHelper() bool { return z.cap(0, 3) }

// XMEGA returns true if and only if the device has XMEGA support.
func (z ZTEXCapability) XMEGA() bool { return z.cap(0, 4) }

// HighSpeedFPGAConfiguration returns true if and only if the device
// supports high-speed FPGA configuration.
func (z ZTEXCapability) HighSpeedFPGAConfiguration() bool { return z.cap(0, 5) }

// MACEEPROM returns true if and only if the device has MAC EEPROM support.
func (z ZTEXCapability) MACEEPROM() bool { return z.cap(0, 6) }

// MultiFPGA returns true if and only if the device has multi-FPGA support.
func (z ZTEXCapability) MultiFPGA() bool { return z.cap(0, 7) }

// TemperatureSensor returns true if and only if the device has
// temperature sensor support.
func (z ZTEXCapability) TemperatureSensor() bool { return z.cap(1, 0) }

// FlashMemory2 returns true if and only if the device has advanced
// flash memory support.
func (z ZTEXCapability) FlashMemory2() bool { return z.cap(1, 1) }

// FX3Firmware returns true if and only if the device has FX3 firmware
// support.
func (z ZTEXCapability) FX3Firmware() bool { return z.cap(1, 2) }

// DebugHelper2 returns true if and only if the device has advanced debug
// helper support.
func (z ZTEXCapability) DebugHelper2() bool { return z.cap(1, 3) }

// DefaultFirmwareInterface returns true if and only if the device
// supports the default firmware interface.
func (z ZTEXCapability) DefaultFirmwareInterface() bool { return z.cap(1, 4) }

// ZTEXModule represents product specific configuration.
type ZTEXModule [12]uint8

// ZTEXSerial represents the device serial number.
type ZTEXSerial [10]uint8

// ZTEXConfig represents the ZTEX device descriptor.
type ZTEXConfig struct {
	ZTEXSize
	ZTEXVersion
	ZTEXMagic
	ZTEXProduct
	ZTEXFirmware
	ZTEXInterface
	ZTEXCapability
	ZTEXModule
	ZTEXSerial
}

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
		return "Xilinx Spartan-6 XC6SLX150 (x4)"
	case 14:
		return "Xilinx Artix-7 XC7A15T"
	default:
		return Unknown
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
		return Unknown
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
	x = append(x, fmt.Sprintf("Type %v", f.FPGAType))
	x = append(x, fmt.Sprintf("Package %v", f.FPGAPackage))
	x = append(x, fmt.Sprintf("Grade %v", f.FPGAGrade))
	return strings.Join(x, ", ")
}

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
		return Unknown
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
	x = append(x, fmt.Sprintf("Size %v", r.RAMSize))
	x = append(x, fmt.Sprintf("Type %v", r.RAMType))
	return strings.Join(x, ", ")
}

// BitstreamSize indicates the actual size of the FPGA bitstream in
// 4 kiB sectors.
type BitstreamSize [2]byte

// String returns a human-readable representation of the bitstream size.
func (b BitstreamSize) String() string {
	return binaryPrefix(uint64(b.Number())<<12, "B")
}

// Number returns a raw numeric representation of the bitstream size.
func (b BitstreamSize) Number() uint16 {
	return (uint16(b[0]) << 0) | (uint16(b[1]) << 8)
}

// BitstreamCapacity indicates the maximum size of the FPGA bitstream in
// 4 kiB sectors.
type BitstreamCapacity [2]byte

// String returns a human-readable representation of the bitstream size.
func (b BitstreamCapacity) String() string {
	return binaryPrefix(uint64(b.Number())<<12, "B")
}

// Number returns a raw numeric representation of the bitstream size.
func (b BitstreamCapacity) Number() uint16 {
	return (uint16(b[0]) << 0) | (uint16(b[1]) << 8)
}

// BitstreamStart indicates the start of the bitstream.
type BitstreamStart [2]byte

// String returns a human-readable representation of the bitstream size.
func (b BitstreamStart) String() string {
	return binaryPrefix(uint64(b.Number())<<12, "B")
}

// Number returns a raw numeric representation of the bitstream size.
func (b BitstreamStart) Number() uint16 {
	return (uint16(b[0]) << 0) | (uint16(b[1]) << 8)
}

// BitstreamConfig indicates the configuration of the bitstream in flash.
type BitstreamConfig struct {
	BitstreamSize
	BitstreamCapacity
	BitstreamStart
}

// String returns a human-readable representation of the bitstream
// configuration.
func (b BitstreamConfig) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Size %v", b.BitstreamSize))
	x = append(x, fmt.Sprintf("Capacity %v", b.BitstreamCapacity))
	x = append(x, fmt.Sprintf("Start %v", b.BitstreamStart))
	return strings.Join(x, ", ")
}

// DeviceConfig indicates the configuration of the device.
type DeviceConfig struct {
	BoardConfig
	FPGAConfig
	RAMConfig
	BitstreamConfig
}

// String returns a human-readable representation of the device configuration.
func (d DeviceConfig) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Board %v", d.BoardConfig))
	x = append(x, fmt.Sprintf("FPGA %v", d.FPGAConfig))
	x = append(x, fmt.Sprintf("RAM %v", d.RAMConfig))
	x = append(x, fmt.Sprintf("Bitstream %v", d.BitstreamConfig))
	return strings.Join(x, ", ")
}

// Device represents a ZTEX USB-FPGA module.
type Device struct {
	*gousb.Device

	ZTEXConfig
	DeviceConfig
}

// DeviceOption represents a device option.
type DeviceOption func(*Device) error

// ControlTimeout sets the timeout for control commands for the device.
func ControlTimeout(timeout time.Duration) DeviceOption {
	return func(d *Device) error {
		d.ControlTimeout = timeout
		return nil
	}
}

// OpenDevice opens a ZTEX USB-FPGA module and returns its device handle.
// If there are multiple modules present, then one is chosen arbitrarily.
func OpenDevice(ctx *gousb.Context, opt ...DeviceOption) (*Device, error) {
	d := &Device{}
	if dev, err := ctx.OpenDeviceWithVIDPID(VendorID, ProductID); err != nil {
		return nil, err
	} else if dev == nil {
		return nil, fmt.Errorf("(*gousb.Context).OpenDeviceWithVIDPID: got nil device, want non-nil device")
	} else {
		d.Device = dev
	}

	b := make([]byte, 128)

	// VR 0x22: ZTEX descriptor: read ZTEX descriptor
	if nbr, err := d.Control(0xc0, 0x22, 0, 0, b); err != nil {
		return nil, err
	} else if nbr != 40 {
		return nil, fmt.Errorf("(*gousb.Device).Control: ZTEX descriptor: read ZTEX descriptor: got %v bytes, want %v bytes", nbr, 40)
	} else if b[0] != 40 {
		return nil, fmt.Errorf("(*gousb.Device).Control: ZTEX descriptor: read ZTEX descriptor: got size %v, want size %v", b[0], 40)
	} else if b[1] != 1 {
		return nil, fmt.Errorf("(*gousb.Device).Control: ZTEX descriptor: read ZTEX descriptor: got version %v, want version %v", b[0], 1)
	}
	d.ZTEXConfig = ZTEXConfig{
		ZTEXSize(b[0]),
		ZTEXVersion(b[1]),
		ZTEXMagic([4]uint8{b[2], b[3], b[4], b[5]}),
		ZTEXProduct([4]uint8{b[6], b[7], b[8], b[9]}),
		ZTEXFirmware(b[10]),
		ZTEXInterface(b[11]),
		ZTEXCapability([6]uint8{b[12], b[13], b[14], b[15], b[16], b[17]}),
		ZTEXModule([12]uint8{b[18], b[19], b[20], b[21], b[22], b[23], b[24], b[25], b[26], b[27], b[28], b[29]}),
		ZTEXSerial([10]uint8{b[30], b[31], b[32], b[33], b[34], b[35], b[36], b[37], b[38], b[39]}),
	}

	// VR 0x3b: MAC EEPROM support: read from MAC EEPROM
	if nbr, err := d.Control(0xc0, 0x3b, 0, 0, b); err != nil {
		return nil, err
	} else if nbr != 128 {
		return nil, fmt.Errorf("(*gousb.Device).Control: MAC EEPROM support: read from MAC EEPROM: got %v bytes, want %v bytes", nbr, 128)
	} else if b[0] != 'C' || b[1] != 'D' || b[2] != '0' {
		return nil, fmt.Errorf("(*gousb.Device).Control: MAC EEPROM support: read from MAC EEPROM: got signature %v, want signature %v", b[:3], []byte{'C', 'D', '0'})
	}
	d.DeviceConfig = DeviceConfig{
		BoardConfig{
			BoardType(b[3]),
			BoardVersion{
				BoardSeries(b[4]),
				BoardNumber(b[5]),
				BoardVariant([2]byte{b[6], b[7]}),
			},
		},
		FPGAConfig{
			FPGAType([2]byte{b[8], b[9]}),
			FPGAPackage(b[10]),
			FPGAGrade([3]byte{b[11], b[12], b[13]}),
		},
		RAMConfig{
			RAMSize(b[14]),
			RAMType(b[15]),
		},
		BitstreamConfig{
			BitstreamSize([2]byte{b[26], b[27]}),
			BitstreamCapacity([2]byte{b[28], b[29]}),
			BitstreamStart([2]byte{b[30], b[31]}),
		},
	}

	for _, o := range opt {
		if err := o(d); err != nil {
			return nil, err
		}
	}
	return d, nil
}

// FPGAConfigured indicates whether or not the FPGA is configured.
type FPGAConfigured uint8

// String returns a human-readable description of the FPGA configuration
// indicator.
func (f FPGAConfigured) String() string {
	switch f {
	case 0:
		return "Unconfigured"
	case 1:
		return "Configured"
	default:
		return Unknown
	}
}

// Number returns the raw numeric representation of the FPGA configuration
// indicator.
func (f FPGAConfigured) Number() uint8 { return uint8(f) }

// Bool returns true if and only if the FPGA is configured.
func (f FPGAConfigured) Bool() bool { return f == 1 }

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
func (f FPGATransferred) Number() uint32 {
	return (uint32(f[0]) << 0) | (uint32(f[1]) << 8) | (uint32(f[2]) << 16) | (uint32(f[3]) << 24)
}

// FPGAInit represents the number of INIT_B states.
type FPGAInit uint8

// FPGAResult represents the result of previous FPGA configuration.
type FPGAResult uint8

// String returns a human-readable description of the FPGA configuration
// result.
func (f FPGAResult) String() string {
	switch f {
	case 0:
		return "Successful"
	default:
		return Unknown
	}
}

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
		return Unknown
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

// FPGAStatus retrieves the current status of the FPGA on the device.
func (d *Device) FPGAStatus() (*FPGAStatus, error) {
	b := make([]byte, 9)

	// VR 0x30: FPGA configuration: get FPGA state
	if nbr, err := d.Control(0xc0, 0x30, 0, 0, b); err != nil {
		return nil, err
	} else if nbr != 9 {
		return nil, fmt.Errorf("(*gousb.Device).Control: FPGA configuration: get FPGA state: got %v bytes, want %v bytes", nbr, 9)
	}

	return &FPGAStatus{
		FPGAConfigured(b[0]),
		FPGAChecksum(b[1]),
		FPGATransferred([4]uint8{b[2], b[3], b[4], b[5]}),
		FPGAInit(b[6]),
		FPGAResult(b[7]),
		FPGASwapped(b[8]),
	}, nil
}

// ResetFPGA resets the FPGA on the device.
func (d *Device) ResetFPGA() error {
	// VC 0x31: FPGA configuration: reset FPGA
	if nbr, err := d.Control(0x40, 0x31, 0, 0, nil); err != nil {
		return err
	} else if nbr != 0 {
		return fmt.Errorf("(*gousb.Device).Control: FPGA configuration: reset FPGA: got %v bytes, want %v bytes", nbr, 0)
	}

	return nil
}
