package protocol

import (
	"fmt"
	"math"
)

// The sensor exposes six user-selectable stages, but the firmware packet has
// room for eight DPI slots. We always emit eight so the packet matches what
// the reference driver produces.
const (
	dpiStages      = 6 // user-configurable stages
	dpiSlots       = 8 // slots present in the packet
	activeStageMsk = 0x3f
)

// defaultStageColors keeps the per-stage LED colours intact so applying DPI
// does not blank the indicator LEDs. Ported from dpiBuilderDefaultOptions.
var defaultStageColors = [dpiSlots][3]byte{
	{0xff, 0x00, 0x00},
	{0x00, 0xff, 0x00},
	{0x00, 0x00, 0xff},
	{0xff, 0xff, 0x00},
	{0x00, 0xff, 0xff},
	{0xff, 0x00, 0xff},
	{0xff, 0x40, 0x00},
	{0xff, 0xff, 0xff},
}

// DPISettings describes the DPI values written to the mouse.
type DPISettings struct {
	Stages [dpiStages]int // DPI per stage, MinDPI..MaxDPI
	Active int            // active stage, 1..dpiStages
}

// BuildDPI returns the 56-byte feature report that applies DPI settings.
func BuildDPI(s DPISettings) ([]byte, error) {
	if s.Active < 1 || s.Active > dpiStages {
		return nil, fmt.Errorf("active stage %d out of range 1..%d", s.Active, dpiStages)
	}

	buf := make([]byte, dpiPacketLen)
	buf[0] = dpiReportID
	buf[1] = dpiPacketLen
	buf[2] = ProfileID
	buf[3] = 0x00 // lift-off distance | angle snap
	buf[4] = 0x01 // motion sync | ripple control (ripple enabled)
	buf[5] = activeStageMsk

	var doubleFlags, tripleFlags byte
	for i := 0; i < dpiSlots; i++ {
		dpi := 0 // stages 7 and 8 are unused
		if i < dpiStages {
			dpi = s.Stages[i]
			if dpi < MinDPI || dpi > MaxDPI {
				return nil, fmt.Errorf("stage %d: DPI %d out of range %d..%d", i+1, dpi, MinDPI, MaxDPI)
			}
		}

		x, y, isDouble, isTriple := encodeDPI(dpi)
		buf[8+i] = x
		buf[16+i] = y
		if isDouble {
			doubleFlags |= 1 << i
		}
		if isTriple {
			tripleFlags |= 1 << i
		}

		rgb := defaultStageColors[i]
		buf[25+i*3] = rgb[0]
		buf[25+i*3+1] = rgb[1]
		buf[25+i*3+2] = rgb[2]
	}
	buf[6] = doubleFlags
	buf[7] = tripleFlags
	buf[24] = byte(s.Active)
	buf[49] = 0x02 // DPI active indicator

	// 16-bit big-endian checksum of bytes 3..49.
	sum := 0
	for i := 3; i <= 49; i++ {
		sum += int(buf[i])
	}
	buf[50] = byte((sum >> 8) & 0xff)
	buf[51] = byte(sum & 0xff)

	return buf, nil
}

// encodeDPI converts a DPI value into the sensor's register bytes and the
// double/triple speed flags. Ported from utils/dpi.ts.
func encodeDPI(dpi int) (x, y byte, isDouble, isTriple bool) {
	if dpi < MinDPI {
		dpi = MinDPI
	}
	if dpi > 26000 {
		dpi = 26000
	}

	if dpi == 20100 {
		return 0xeb, 0x01, true, false
	}

	target := dpi
	if dpi > 10000 {
		target = int(math.Round(float64(dpi) / 2))
		isDouble = true
	}

	switch {
	case target > 10000:
		extra := (target - 10100) / 100
		combined := 199 + extra
		return byte(combined & 0xff), byte((combined >> 8) & 0xff), isDouble, false

	case target > 5000 && target%100 == 0:
		idx := target/100 - 1
		if idx >= 0 && idx < len(dpiMap3311) {
			return dpiMap3311[idx], 0x01, isDouble, false
		}
		return 0x00, 0x01, isDouble, false

	default:
		idx := (target - MinDPI) / 50
		if idx >= 0 && idx < len(dpiMap3311) {
			return dpiMap3311[idx], 0x00, isDouble, false
		}
		return 0x00, 0x00, isDouble, false
	}
}
