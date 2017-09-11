// Package ztex manages ZTEX USB FPGA modules.
package ztex

import (
	"time"

	"github.com/google/gousb"
)

var (
	// VendorID is the ZTEX USB vendor ID (VID).
	VendorID = gousb.ID(0x221A)

	// ProductID is the standard ZTEX USB product ID (PID)
	ProductID = gousb.ID(0x0100)

	// ControlTimeout is the timeout for control transfers.
	ControlTimeout = 1000 * time.Millisecond
)

// Device represents a ZTEX USB FPGA module.
type Device struct {
	ProductName  [128]byte
	SerialNumber [64]byte
	FXVersion    byte    // 2 = FX2, 3 = FX3.
	BoardSeries  byte    // 1 = Series 1, 2 = Series 2.  (255 = Unknown)
	BoardNumber  byte    // E.g., 2.14b -> 2.  (255 = Unknown)
	BoardVariant [2]byte // E.g., 2.14b -> [2]byte{'b', 0}.

	// Endpoint for fast FPGA configuration.  (0 = Unsupported)
	FastConfigurationEndpoint byte
	// Interface for fast FPGA configuration.  (0 = Unsupported)
	FastConfigurationInterface byte
	/// Default interface major version number.  (0 = Not Available)
	DefaultInterfaceMajorVersion byte
	/// Default interface minor version number.  (0 = Not Available)
	DefaultInterfaceMinorVersion byte
	/// Output endpoint of default interface.  (255 = Not Available)
	DefaultInterfaceOutputEndpoint byte
	/// Input endpoint of default interface.  (255 = Not Available)
	DefaultInterfaceInputEndpoint byte

	// Device is the USB device handle for the module.
	Device *gousb.Device
}

func OpenDevice(ctx *gousb.Context) (*Device, error) {
	d := &Device{}

	if dev, err := ctx.OpenDeviceWithVIDPID(VendorID, ProductID); err != nil {
		return nil, err
	} else {
		dev.ControlTimeout = ControlTimeout
		d.Device = dev
	}

	buf := make([]byte, 128)

	// VR 0x33: High speed FPGA configuration support: Read Endpoint settings
	if nbt, err := d.Device.Control(0xc0, 0x33, 0, 0, buf); err != nil {
		return nil, err
	} else if nbt == 2 {
		d.FastConfigurationEndpoint = buf[0]
		d.FastConfigurationInterface = buf[1]
	}

	// VR 0x3b: MAC EEPROM support: Read from MAC EEPROM
	if nbt, err := d.Device.Control(0xc0, 0x3b, 0, 0, buf); err != nil {
		return nil, err
	} else if nbt == 128 {
		d.FXVersion = buf[3]
		d.BoardSeries = buf[4]
		d.BoardNumber = buf[5]
		d.BoardVariant[0] = buf[6]
		d.BoardVariant[1] = buf[7]
	} else {
		d.BoardSeries = 255
		d.BoardNumber = 255
	}

	// VR 0x64: Default firmware interface: Return Default Interface information
	if nbt, err := d.Device.Control(0xc0, 0x64, 0, 0, buf); err != nil {
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

// Close releases all external resources associated with the device.
func (d *Device) Close() error { return d.Device.Close() }
