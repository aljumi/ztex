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
)

func binaryPrefix(n uint64, unit string) string {
	switch {
	case n != 0 && n&(1<<30-1) == 0:
		return fmt.Sprintf("%v%v [%vGi%v]", n, unit, n>>30, unit)
	case n != 0 && n&(1<<20-1) == 0:
		return fmt.Sprintf("%v%v [%vMi%v]", n, unit, n>>20, unit)
	case n != 0 && n&(1<<10-1) == 0:
		return fmt.Sprintf("%v%v [%vki%v]", n, unit, n>>10, unit)
	default:
		return fmt.Sprintf("%v%v", n, unit)
	}
}
