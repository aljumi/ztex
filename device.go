package ztex

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gousb"
)

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
		return nil, fmt.Errorf("(*gousb.Context).OpenDeviceWithVIDPID: %v", err)
	} else if dev == nil {
		return nil, fmt.Errorf("(*gousb.Context).OpenDeviceWithVIDPID: got nil device, want non-nil device")
	} else {
		d.Device = dev
	}

	b := make([]byte, 128)

	// VR 0x22: ZTEX descriptor: read ZTEX descriptor
	if nbr, err := d.Control(0xc0, 0x22, 0, 0, b); err != nil {
		return nil, fmt.Errorf("(*gousb.Device).Control: ZTEX descriptor: read ZTEX descriptor: %v", err)
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
		return nil, fmt.Errorf("(*gousb.Device).Control: MAC EEPROM support: read from MAC EEPROM: %v", err)
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

// ResetFX3 resets the Cypress CYUSB3033 EZ-USB FX3S controller on the
// device, if one is present.
func (d *Device) ResetFX3() error {
	if !d.ZTEXCapability.FX3Firmware() {
		return fmt.Errorf("(*ztex.Device).ResetFX3: operation not supported")
	}

	if nbr, err := d.Control(0x40, 0xa1, 1, 0, nil); err != nil {
		return fmt.Errorf("(*gousb.Device).Control: FX3 firmware: reset and boot from flash: %v", err)
	} else if nbr != 0 {
		return fmt.Errorf("(*gousb.Device).Control: FX3 firmware: reset and boot from flash: got %v bytes, want %v bytes", nbr, 0)
	}

	return nil
}

// FPGAStatus retrieves the current status of the FPGA on the device.
func (d *Device) FPGAStatus() (*FPGAStatus, error) {
	if !d.ZTEXCapability.FPGAConfiguration() {
		return nil, fmt.Errorf("(*ztex.Device).FPGAStatus: operation not supported")
	}

	b := make([]byte, 9)

	// VR 0x30: FPGA configuration: get FPGA state
	if nbr, err := d.Control(0xc0, 0x30, 0, 0, b); err != nil {
		return nil, fmt.Errorf("(*gousb.Device).Control: FPGA configuration: get FPGA state: %v", err)
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

// ResetFPGA resets the FPGA on the device, if one is present.
func (d *Device) ResetFPGA() error {
	if !d.ZTEXCapability.FPGAConfiguration() {
		return fmt.Errorf("(*ztex.Device).ResetFPGA: operation not supported")
	}

	// VC 0x31: FPGA configuration: reset FPGA
	if nbr, err := d.Control(0x40, 0x31, 0, 0, nil); err != nil {
		return fmt.Errorf("(*gousb.Device).Control: FPGA configuration: reset FPGA: %v", err)
	} else if nbr != 0 {
		return fmt.Errorf("(*gousb.Device).Control: FPGA configuration: reset FPGA: got %v bytes, want %v bytes", nbr, 0)
	}

	return nil
}
