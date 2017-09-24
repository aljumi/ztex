package ztex

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gousb"
)

// Device represents a ZTEX USB device.
type Device struct {
	*gousb.Device

	DescriptorConfig
	BoardConfig
	FPGAConfig
	RAMConfig
	BitstreamConfig
}

// String returns a human-readable representation of the device.
func (d *Device) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Device(%v)", d.Device))
	x = append(x, fmt.Sprintf("Descriptor(%v)", d.DescriptorConfig))
	x = append(x, fmt.Sprintf("Board(%v)", d.BoardConfig))
	x = append(x, fmt.Sprintf("FPGA(%v)", d.FPGAConfig))
	x = append(x, fmt.Sprintf("RAM(%v)", d.RAMConfig))
	x = append(x, fmt.Sprintf("Bitstream(%v)", d.BitstreamConfig))
	return strings.Join(x, ", ")
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

	if err := d.readDescriptorConfig(); err != nil {
		return nil, err
	}

	if err := d.readDeviceConfig(); err != nil {
		return nil, err
	}

	for _, o := range opt {
		if err := o(d); err != nil {
			return nil, err
		}
	}

	return d, nil
}

func (d *Device) readDescriptorConfig() error {
	b := make([]byte, 40)

	// VR 0x22: ZTEX descriptor: read ZTEX descriptor
	if nbr, err := d.Control(0xc0, 0x22, 0, 0, b); err != nil {
		return fmt.Errorf("(*ztex.Device).Control: ZTEX descriptor: read ZTEX descriptor: %v", err)
	} else if nbr != 40 {
		return fmt.Errorf("(*ztex.Device).Control: ZTEX descriptor: read ZTEX descriptor: got %v bytes, want %v bytes", nbr, 40)
	} else if b[0] != 40 {
		return fmt.Errorf("(*ztex.Device).Control: ZTEX descriptor: read ZTEX descriptor: got size %v, want size %v", b[0], 40)
	} else if b[1] != 1 {
		return fmt.Errorf("(*ztex.Device).Control: ZTEX descriptor: read ZTEX descriptor: got version %v, want version %v", b[0], 1)
	}

	d.DescriptorConfig = DescriptorConfig{
		DescriptorSize(b[0]),
		DescriptorVersion(b[1]),
		DescriptorMagic([4]uint8{b[2], b[3], b[4], b[5]}),
		DescriptorProduct([4]uint8{b[6], b[7], b[8], b[9]}),
		DescriptorFirmware(b[10]),
		DescriptorInterface(b[11]),
		DescriptorCapability([6]uint8{b[12], b[13], b[14], b[15], b[16], b[17]}),
		DescriptorModule([12]uint8{b[18], b[19], b[20], b[21], b[22], b[23], b[24], b[25], b[26], b[27], b[28], b[29]}),
		DescriptorSerial([10]uint8{b[30], b[31], b[32], b[33], b[34], b[35], b[36], b[37], b[38], b[39]}),
	}

	return nil
}

func (d *Device) readDeviceConfig() error {
	b := make([]byte, 128)

	// VR 0x3b: MAC EEPROM support: read from MAC EEPROM
	if nbr, err := d.Control(0xc0, 0x3b, 0, 0, b); err != nil {
		return fmt.Errorf("(*ztex.Device).Control: MAC EEPROM support: read from MAC EEPROM: %v", err)
	} else if nbr != 128 {
		return fmt.Errorf("(*ztex.Device).Control: MAC EEPROM support: read from MAC EEPROM: got %v bytes, want %v bytes", nbr, 128)
	} else if b[0] != 'C' || b[1] != 'D' || b[2] != '0' {
		return fmt.Errorf("(*ztex.Device).Control: MAC EEPROM support: read from MAC EEPROM: got signature %v, want signature %v", b[:3], []byte{'C', 'D', '0'})
	}

	d.BoardConfig = BoardConfig{
		BoardType(b[3]),
		BoardVersion{
			BoardSeries(b[4]),
			BoardNumber(b[5]),
			BoardVariant([2]byte{b[6], b[7]}),
		},
	}
	d.FPGAConfig = FPGAConfig{
		FPGAType([2]byte{b[8], b[9]}),
		FPGAPackage(b[10]),
		FPGAGrade([3]byte{b[11], b[12], b[13]}),
	}
	d.RAMConfig = RAMConfig{
		RAMSize(b[14]),
		RAMType(b[15]),
	}
	d.BitstreamConfig = BitstreamConfig{
		BitstreamSize([2]byte{b[26], b[27]}),
		BitstreamCapacity([2]byte{b[28], b[29]}),
		BitstreamStart([2]byte{b[30], b[31]}),
	}

	return nil
}

// ResetFX3 resets the Cypress CYUSB3033 EZ-USB FX3S controller on the
// device, if one is present.
func (d *Device) ResetFX3() error {
	if !d.DescriptorCapability.FX3Firmware() {
		return fmt.Errorf("operation not supported")
	}

	// VC 0xa1: FX3 support: reset FX3 controller
	if nbr, err := d.Control(0x40, 0xa1, 1, 0, nil); err != nil {
		return fmt.Errorf("(*gousb.Device).Control: FX3 firmware: reset and boot from flash: %v", err)
	} else if nbr != 0 {
		return fmt.Errorf("(*gousb.Device).Control: FX3 firmware: reset and boot from flash: got %v bytes, want %v bytes", nbr, 0)
	}

	return nil
}

// FPGAStatus retrieves the current FPGA status.
func (d *Device) FPGAStatus() (*FPGAStatus, error) {
	if !d.DescriptorCapability.FPGAConfiguration() {
		return nil, fmt.Errorf("operation not supported")
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

// ResetFPGA resets the FPGA on the device.
func (d *Device) ResetFPGA() error {
	if !d.DescriptorCapability.FPGAConfiguration() {
		return fmt.Errorf("operation not supported")
	}

	// VC 0x31: FPGA configuration: reset FPGA
	if nbr, err := d.Control(0x40, 0x31, 0, 0, nil); err != nil {
		return fmt.Errorf("(*gousb.Device).Control: FPGA configuration: reset FPGA: %v", err)
	} else if nbr != 0 {
		return fmt.Errorf("(*gousb.Device).Control: FPGA configuration: reset FPGA: got %v bytes, want %v bytes", nbr, 0)
	}

	return nil
}

// FlashStatus retrieves the current flash memory status.
func (d *Device) FlashStatus() (*FlashStatus, error) {
	if !d.DescriptorCapability.FlashMemory() {
		return nil, fmt.Errorf("operation not supported")
	}

	b := make([]byte, 8)

	// VR 0x40: flash memory support: get flash state
	if nbr, err := d.Control(0xc0, 0x40, 0, 0, b); err != nil {
		return nil, fmt.Errorf("(*gousb.Device).Control: flash memory support: get flash state: %v", err)
	} else if nbr != 8 {
		return nil, fmt.Errorf("(*gousb.Device).Control: flash memory support: get flash state: got %v bytes, want %v bytes", nbr, 8)
	}

	return &FlashStatus{
		FlashEnabled(b[0]),
		FlashSector([2]uint8{b[1], b[2]}),
		FlashCount([4]uint8{b[3], b[4], b[5], b[6]}),
		FlashError(b[7]),
	}, nil
}

// ResetDefaultFirmware resets the default firmware, if it is present.
func (d *Device) ResetDefaultFirmware() error {
	if !d.DescriptorCapability.DefaultFirmware() {
		return fmt.Errorf("operation not supported")
	}

	// VC 0x60: default firmware interface: reset
	if nbr, err := d.Control(0x40, 0x60, 0, 0, nil); err != nil {
		return fmt.Errorf("(*gousb.Device).Control: default firmware interface: reset: %v", err)
	} else if nbr != 0 {
		return fmt.Errorf("(*gousb.Device).Control: default firmware interface: reset: got %v bytes, want %v bytes", nbr, 0)
	}

	return nil
}
