// Package ztex manages ZTEX USB-FPGA modules.
package ztex

import (
	"errors"
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
		return "Unknown"
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
		return "Unknown"
	}
}

// Number returns the raw representation of a board series.
func (b BoardSeries) Number() uint8 { return uint8(b) }

// BoardNumber indicates a board in a series.
type BoardNumber uint8

// String returns a human-readable description of a board number.
func (b BoardNumber) String() string {
	if b == 255 {
		return "Unknown"
	}
	return fmt.Sprintf("%d", uint8(b))
}

// Number returns the raw representation of a board number.
func (b BoardNumber) Number() uint8 { return uint8(b) }

// BoardVariant indicates a variation on a board series and number.
type BoardVariant [2]byte

// String returns a human-readable description of a board variant.
func (b BoardVariant) String() string { return string(b.Bytes()) }

// Bytes returns the raw representation of a board variant.
func (b BoardVariant) Bytes() []byte {
	c := make([]byte, 0, 2)
	if b[0] == 0 {
		return c
	}
	c = append(c, b[0])
	if b[1] == 0 {
		return c
	}
	c = append(c, b[1])
	return c
}

// BoardVersion indicates the type, series, number, and variant of a ZTEX
// USB-FPGA module.  For example, a ZTEX USB3-FPGA module 2.18b device
// would be represented by
//
//   BoardVersion{
//     BoardType: BoardType(3),
//     BoardSeries: BoardSeries(2),
//     BoardNumber: BoardNumber(18),
//     BoardVariant: BoardVariant([2]byte{0x62, 0x00}]),
//   }
//
// as a BoardVersion structure.
type BoardVersion struct {
	BoardType
	BoardSeries
	BoardNumber
	BoardVariant
}

// String returns a human-readable representation of a board version.
func (b BoardVersion) String() string {
	return fmt.Sprintf("%v %v.%v%v", b.BoardType, b.BoardSeries, b.BoardNumber, b.BoardVariant)
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
		return "Unknown"
	}
}

// Bytes returns a raw representation of an FPGA type.
func (f FPGAType) Bytes() []byte { return []byte{f[0], f[1]} }

// Number returns a numeric representation of an FPGA type.
func (f FPGAType) Number() uint16 { return (uint16(f[0]) << 8) | (uint16(f[1]) << 0) }

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

// FPGAVersion indicates the type, package, speed grade, etc. of the FPGA
// present in a device.
type FPGAVersion struct {
	FPGAType
	FPGAPackage
}

func (f FPGAVersion) String() string {
	return fmt.Sprintf("Type: %v, Package: %v", f.FPGAType, f.FPGAPackage)
}

// Device represents a ZTEX USB-FPGA module.
type Device struct {
	*gousb.Device

	BoardVersion
	FPGAVersion

	Bytes []byte
}

// String returns a human-readable representation of the device.
func (d *Device) String() string {
	mfr, _ := d.Manufacturer()
	prd, _ := d.Product()
	snr, _ := d.SerialNumber()

	lines := []string{}
	lines = append(lines, fmt.Sprintf("Manufacturer: %v", mfr))
	lines = append(lines, fmt.Sprintf("Product: %v", prd))
	lines = append(lines, fmt.Sprintf("Serial Number: %v", snr))
	lines = append(lines, fmt.Sprintf("Board: %v", d.BoardVersion))
	lines = append(lines, fmt.Sprintf("FPGA: %v", d.FPGAVersion))

	return strings.Join(lines, "\n")
}

// DeviceOption represents a functional option for devices.
type DeviceOption func(*Device) error

// ControlTimeout is a device option that sets the timeout for control
// commands.
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
		return nil, errors.New("no device")
	} else {
		d.Device = dev
	}

	buf := make([]byte, 128)

	// VR 0x3b: MAC EEPROM support: Read from MAC EEPROM
	if n, err := d.Control(0xc0, 0x3b, 0, 0, buf); err != nil {
		return nil, err
	} else if n != 128 {
		return nil, fmt.Errorf("read from MAC EEPROM: got %v bytes, want %v bytes", n, 128)
	} else if buf[0] != 'C' || buf[1] != 'D' || buf[2] != '0' {
		return nil, fmt.Errorf("read from MAC EEPROM: got %v, want %v", buf[:3], []byte{'C', 'D', '0'})
	}

	d.BoardVersion = BoardVersion{
		BoardType(buf[3]),
		BoardSeries(buf[4]),
		BoardNumber(buf[5]),
		BoardVariant([2]byte{buf[6], buf[7]}),
	}

	d.FPGAVersion = FPGAVersion{
		FPGAType([2]byte{buf[8], buf[9]}),
		FPGAPackage(buf[10]),
	}

	for _, o := range opt {
		if err := o(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}
