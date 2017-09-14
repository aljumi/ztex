// Package ztex manages ZTEX USB FPGA modules.
package ztex

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/gousb"
)

const (
	// VendorID is the ZTEX USB vendor ID (VID).
	VendorID = gousb.ID(0x221A)
	// ProductID is the standard ZTEX USB product ID (PID)
	ProductID = gousb.ID(0x0100)

	// ControlTimeout is the default timeout for control transfers.
	ControlTimeout = 1000 * time.Millisecond
)

// FirmwareVersion indicates which version of the ZTEX firmware is
// present.  Currently only FX1 and FX2 are supported.
type FirmwareVersion uint8

// String returns a human-readable description of a firmware version.
func (f FirmwareVersion) String() string {
	switch f {
	case 2:
		return "FX2"
	case 3:
		return "FX3"
	case 255:
		fallthrough
	default:
		return "Unknown"
	}
}

// Number returns the raw representation of a firmware version.
func (f FirmwareVersion) Number() uint8 { return uint8(f) }

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
type BoardVariant []byte

// String returns a human-readable description of a board variant.
func (b BoardVariant) String() string { return string(b) }

// Number returns the raw representation of a board variant.
func (b BoardVariant) Bytes() []byte { return []byte(b) }

// BoardVersion comprises the series, number, and variant of a ZTEX USB
// FPGA module.  For example, a ZTEX USB FPGA module 2.14b device would be
// represented by
//
//   BoardVersion{
//     Series: BoardSeries(2),
//     Number: BoardNumber(14),
//     Variant: BoardVariant("b"),
//   }
//
// as a BoardVersion structure.
type BoardVersion struct {
	Series  BoardSeries
	Number  BoardNumber
	Variant BoardVariant
}

// String returns a human-readable representation of a board version.
func (b BoardVersion) String() string {
	return b.Series.String() + "." + b.Number.String() + b.Variant.String()
}

// Device represents a ZTEX USB FPGA module.
type Device struct {
	*gousb.Device

	Firmware FirmwareVersion
	Board    BoardVersion

	// Endpoint for fast FPGA configuration.  (0 = Unsupported)
	FastConfigurationEndpoint byte
	// Interface for fast FPGA configuration.  (0 = Unsupported)
	FastConfigurationInterface byte
	// Default interface major version number.  (0 = Not Available)
	DefaultInterfaceMajorVersion byte
	// Default interface minor version number.  (0 = Not Available)
	DefaultInterfaceMinorVersion byte
	// Output endpoint of default interface.  (255 = Not Available)
	DefaultInterfaceOutputEndpoint byte
	// Input endpoint of default interface.  (255 = Not Available)
	DefaultInterfaceInputEndpoint byte
}

// OpenDevice opens the
func OpenDevice(ctx *gousb.Context) (*Device, error) {
	d := &Device{}
	if dev, err := ctx.OpenDeviceWithVIDPID(VendorID, ProductID); err != nil {
		return nil, err
	} else if dev == nil {
		return nil, errors.New("no device")
	} else {
		d.Device = dev
		d.Device.ControlTimeout = ControlTimeout
	}

	buf := make([]byte, 128)

	// VR 0x3b: MAC EEPROM support: Read from MAC EEPROM
	if nbt, err := d.Control(0xc0, 0x3b, 0, 0, buf); err != nil {
		return nil, err
	} else if nbt == 128 {
		d.Firmware = FirmwareVersion(buf[3])
		d.Board = BoardVersion{
			Series:  BoardSeries(buf[4]),
			Number:  BoardNumber(buf[5]),
			Variant: BoardVariant(buf[6:8]),
		}
	} else {
		d.Firmware = FirmwareVersion(255)
		d.Board = BoardVersion{
			Series:  BoardSeries(255),
			Number:  BoardNumber(255),
			Variant: BoardVariant([]byte{255, 255}),
		}
	}

	// VR 0x33: High speed FPGA configuration support: Read Endpoint settings
	if nbt, err := d.Control(0xc0, 0x33, 0, 0, buf); err != nil {
		return nil, err
	} else if nbt == 2 {
		d.FastConfigurationEndpoint = buf[0]
		d.FastConfigurationInterface = buf[1]
	} else {
		return nil, errors.New("internal error")
	}

	// VR 0x64: Default firmware interface: Return Default Interface information
	if nbt, err := d.Control(0xc0, 0x64, 0, 0, buf); err != nil {
		return nil, err
	} else if nbt == 4 {
		d.DefaultInterfaceMajorVersion = buf[0]
		d.DefaultInterfaceMinorVersion = buf[3]
		d.DefaultInterfaceOutputEndpoint = buf[1] & 127
		d.DefaultInterfaceInputEndpoint = buf[2] | 128
	} else if nbt == 3 {
		d.DefaultInterfaceMajorVersion = buf[0]
		d.DefaultInterfaceOutputEndpoint = buf[1] & 127
		d.DefaultInterfaceInputEndpoint = buf[2] | 128
	} else {
		d.DefaultInterfaceOutputEndpoint = 255
		d.DefaultInterfaceInputEndpoint = 255
	}

	return d, nil
}
