package ztex

import (
	"fmt"
	"strings"
)

// DescriptorSize represents the number of bytes in a ZTEX descriptor.
type DescriptorSize uint8

// DescriptorVersion represents the version of a ZTEX descriptor.
type DescriptorVersion uint8

// DescriptorMagic indicates the presence of a ZTEX descriptor.
type DescriptorMagic [4]uint8

// String returns a human-readable description of the ZTEX magic bytes.
func (d DescriptorMagic) String() string { return string(d.Bytes()) }

// Bytes returns a raw representation of the ZTEX magic bytes.
func (d DescriptorMagic) Bytes() []byte { return []byte{d[0], d[1], d[2], d[3]} }

// DescriptorProduct represents a ZTEX product ID.
type DescriptorProduct [4]uint8

// String returns a human-readable description of the ZTEX product ID.
func (d DescriptorProduct) String() string {
	p := "Unknown"
	switch {
	case d[0] == 0 && d[1] == 0 && d[2] == 0 && d[3] == 0:
		p = "Default"
	case d[0] == 1:
		p = "Experimental"
	case d[0] == 10 && d[1] == 0 && d[2] == 1 && d[3] == 1:
		p = "ZTEX BTCMiner"
	case d[0] == 10 && d[1] == 11:
		p = "ZTEX USB-FPGA Module 1.2"
	case d[0] == 10 && d[1] == 12 && d[2] == 2 && (1 <= d[3] && d[3] <= 4):
		p = "NIT"
	case d[0] == 10 && d[1] == 12:
		p = "ZTEX USB-FPGA Module 1.11"
	case d[0] == 10 && d[1] == 13:
		p = "ZTEX USB-FPGA Module 1.15"
	case d[0] == 10 && d[1] == 14:
		p = "ZTEX USB-FPGA Module 1.15x"
	case d[0] == 10 && d[1] == 15:
		p = "ZTEX USB-FPGA Module 1.15y"
	case d[0] == 10 && d[1] == 16:
		p = "ZTEX USB-FPGA Module 2.16"
	case d[0] == 10 && d[1] == 17:
		p = "ZTEX USB-FPGA Module 2.13"
	case d[0] == 10 && d[1] == 18:
		p = "ZTEX USB-FPGA Module 2.01"
	case d[0] == 10 && d[1] == 19:
		p = "ZTEX USB-FPGA Module 2.04"
	case d[0] == 10 && d[1] == 20:
		p = "ZTEX USB Module 1.0"
	case d[0] == 10 && d[1] == 30:
		p = "ZTEX USB-XMEGA Module 1.0"
	case d[0] == 10 && d[1] == 40:
		p = "ZTEX USB-FPGA Module 2.02"
	case d[0] == 10 && d[1] == 41:
		p = "ZTEX USB-FPGA Module 2.14"
	case d[0] == 10 && d[1] == 42:
		p = "ZTEX USB3-FPGA Module 2.18"
	case d[0] == 10:
		p = "ZTEX"
	}
	return fmt.Sprintf("%v.%v.%v.%v [%v]", d[0], d[1], d[2], d[3], p)
}

// Bytes returns a raw representation of the ZTEX product ID.
func (d DescriptorProduct) Bytes() []byte { return []byte{d[0], d[1], d[2], d[3]} }

// DescriptorFirmware indicates the version of the ZTEX firmware.
type DescriptorFirmware uint8

// DescriptorInterface indicates the version of the ZTEX interface.
type DescriptorInterface uint8

// DescriptorCapability indicates the capabilities supported by the ZTEX device.
type DescriptorCapability [6]uint8

// String returns a human-readable description of the ZTEX capabilities
// supported by the device.
func (d DescriptorCapability) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("EEPROM(%v)", d.EEPROM()))
	x = append(x, fmt.Sprintf("FPGA Configuration(%v)", d.FPGAConfiguration()))
	x = append(x, fmt.Sprintf("Flash Memory(%v)", d.FlashMemory()))
	x = append(x, fmt.Sprintf("Debug Helper(%v)", d.DebugHelper()))
	x = append(x, fmt.Sprintf("XMEGA(%v)", d.XMEGA()))
	x = append(x, fmt.Sprintf("High Speed FPGA Configuration(%v)", d.HighSpeedFPGAConfiguration()))
	x = append(x, fmt.Sprintf("MAC EEPROM(%v)", d.MACEEPROM()))
	x = append(x, fmt.Sprintf("MultiFPGA(%v)", d.MultiFPGA()))
	x = append(x, fmt.Sprintf("Temperature Sensor(%v)", d.TemperatureSensor()))
	x = append(x, fmt.Sprintf("Flash Memory 2(%v)", d.FlashMemory2()))
	x = append(x, fmt.Sprintf("FX3 Firmware(%v)", d.FX3Firmware()))
	x = append(x, fmt.Sprintf("Debug Helper 2(%v)", d.DebugHelper2()))
	x = append(x, fmt.Sprintf("Default Firmware(%v)", d.DefaultFirmware()))
	return strings.Join(x, ", ")
}

// Function cap returns true if and only if ZTEX capability i.j is
// supported by the device.
func (d DescriptorCapability) cap(i, j uint) bool { return d[i]&(1<<j) != 0 }

// EEPROM returns true if and only if the device has EEPROM support.
func (d DescriptorCapability) EEPROM() bool { return d.cap(0, 0) }

// FPGAConfiguration returns true if and only if the device has basic
// FPGA configuration support.
func (d DescriptorCapability) FPGAConfiguration() bool { return d.cap(0, 1) }

// FlashMemory returns true if and only if the device has flash memory.
func (d DescriptorCapability) FlashMemory() bool { return d.cap(0, 2) }

// DebugHelper returns true if and only if the device has basic debug
// helper support.
func (d DescriptorCapability) DebugHelper() bool { return d.cap(0, 3) }

// XMEGA returns true if and only if the device has XMEGA support.
func (d DescriptorCapability) XMEGA() bool { return d.cap(0, 4) }

// HighSpeedFPGAConfiguration returns true if and only if the device
// supports high-speed FPGA configuration.
func (d DescriptorCapability) HighSpeedFPGAConfiguration() bool { return d.cap(0, 5) }

// MACEEPROM returns true if and only if the device has MAC EEPROM support.
func (d DescriptorCapability) MACEEPROM() bool { return d.cap(0, 6) }

// MultiFPGA returns true if and only if the device has multi-FPGA support.
func (d DescriptorCapability) MultiFPGA() bool { return d.cap(0, 7) }

// TemperatureSensor returns true if and only if the device has
// temperature sensor support.
func (d DescriptorCapability) TemperatureSensor() bool { return d.cap(1, 0) }

// FlashMemory2 returns true if and only if the device has advanced
// flash memory support.
func (d DescriptorCapability) FlashMemory2() bool { return d.cap(1, 1) }

// FX3Firmware returns true if and only if the device has FX3 firmware
// support.
func (d DescriptorCapability) FX3Firmware() bool { return d.cap(1, 2) }

// DebugHelper2 returns true if and only if the device has advanced debug
// helper support.
func (d DescriptorCapability) DebugHelper2() bool { return d.cap(1, 3) }

// DefaultFirmware returns true if and only if the device supports the
// default firmware interface.
func (d DescriptorCapability) DefaultFirmware() bool { return d.cap(1, 4) }

// DescriptorModule represents product specific configuration.
type DescriptorModule [12]uint8

// DescriptorSerial represents the device serial number.
type DescriptorSerial [10]uint8

// String returns a human-readable description of the device serial number.
func (d DescriptorSerial) String() string { return string(d.Bytes()) }

// Bytes returns a raw representation of the device serial number.
func (d DescriptorSerial) Bytes() []byte {
	return []byte{d[0], d[1], d[2], d[3], d[4], d[5], d[6], d[7], d[8], d[9]}
}

// DescriptorConfig represents the ZTEX device descriptor.
type DescriptorConfig struct {
	DescriptorSize
	DescriptorVersion
	DescriptorMagic
	DescriptorProduct
	DescriptorFirmware
	DescriptorInterface
	DescriptorCapability
	DescriptorModule
	DescriptorSerial
}

// String returns a human-readable description of a ZTEX device descriptor.
func (d DescriptorConfig) String() string {
	x := []string{}
	x = append(x, fmt.Sprintf("Size(%v)", d.DescriptorSize))
	x = append(x, fmt.Sprintf("Version(%v)", d.DescriptorVersion))
	x = append(x, fmt.Sprintf("Magic(%v)", d.DescriptorMagic))
	x = append(x, fmt.Sprintf("Product(%v)", d.DescriptorProduct))
	x = append(x, fmt.Sprintf("Firmware(%v)", d.DescriptorFirmware))
	x = append(x, fmt.Sprintf("Interface(%v)", d.DescriptorInterface))
	x = append(x, fmt.Sprintf("Capability(%v)", d.DescriptorCapability))
	x = append(x, fmt.Sprintf("Module(%v)", d.DescriptorModule))
	x = append(x, fmt.Sprintf("Serial(%v)", d.DescriptorSerial))
	return strings.Join(x, ", ")
}
