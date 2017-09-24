package ztex

import (
	"fmt"

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
