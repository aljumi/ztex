package ztex

import (
	"fmt"
	"strings"
)

// BitstreamSize indicates the actual size of the FPGA bitstream in
// 4 kiB sectors.
type BitstreamSize [2]byte

// String returns a human-readable representation of the bitstream size.
func (b BitstreamSize) String() string {
	return binaryPrefix(uint64(b.Number())<<12, "B")
}

// Number returns a raw numeric representation of the bitstream size.
func (b BitstreamSize) Number() uint16 { return bytesToUint16(b) }

// BitstreamCapacity indicates the maximum size of the FPGA bitstream in
// 4 kiB sectors.
type BitstreamCapacity [2]byte

// String returns a human-readable representation of the bitstream size.
func (b BitstreamCapacity) String() string {
	return binaryPrefix(uint64(b.Number())<<12, "B")
}

// Number returns a raw numeric representation of the bitstream size.
func (b BitstreamCapacity) Number() uint16 { return bytesToUint16(b) }

// BitstreamStart indicates the start of the bitstream.
type BitstreamStart [2]byte

// String returns a human-readable representation of the bitstream size.
func (b BitstreamStart) String() string {
	return binaryPrefix(uint64(b.Number())<<12, "B")
}

// Number returns a raw numeric representation of the bitstream size.
func (b BitstreamStart) Number() uint16 { return bytesToUint16(b) }

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
	x = append(x, fmt.Sprintf("Size(%v)", b.BitstreamSize))
	x = append(x, fmt.Sprintf("Capacity(%v)", b.BitstreamCapacity))
	x = append(x, fmt.Sprintf("Start(%v)", b.BitstreamStart))
	return strings.Join(x, ", ")
}
