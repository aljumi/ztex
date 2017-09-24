// Package ztex manages ZTEX USB-FPGA modules.
package ztex

import (
	"fmt"
	"strings"
)

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
	case z[0] == 10 && z[1] == 42:
		p = "ZTEX USB3-FPGA Module 2.18"
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
