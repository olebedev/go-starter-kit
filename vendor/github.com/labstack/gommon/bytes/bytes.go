package bytes

import (
	"fmt"
	"math"
	"strconv"
)

var (
	global = New()
)

type (
	Bytes struct {
		iec bool
	}
)

// New creates a Bytes instance.
func New() *Bytes {
	return &Bytes{}
}

// SetBinaryPrefix sets binary prefix format.
func (g *Bytes) SetBinaryPrefix(on bool) {
	g.iec = on
}

// Format formats bytes to string. For example, 1323 bytes will return 1.32 KB.
// If binary prefix is set, it will return 1.29 KiB.
func (g *Bytes) Format(b uint64) string {
	unit := uint64(1000)
	if g.iec {
		unit = 1024
	}
	if b < unit {
		return strconv.FormatUint(b, 10) + " B"
	}
	bb := float64(b)
	uunit := float64(unit)
	x := math.Floor(math.Log(bb) / math.Log(uunit))
	pre := make([]byte, 1, 2)
	pre[0] = "KMGTPE"[uint8(x)-1]
	if g.iec {
		pre = pre[:2]
		pre[1] = 'i'
	}
	// TODO: Improve performance?
	return fmt.Sprintf("%.02f %sB", bb/math.Pow(uunit, x), pre)

}

// BinaryPrefix wraps global Bytes's BinaryPrefix function.
func BinaryPrefix(on bool) {
	global.SetBinaryPrefix(on)
}

// Format wraps global Bytes's Format function.
func Format(b uint64) string {
	return global.Format(b)
}
