package protocol

import "fmt"

// pollingRates maps a polling frequency in Hz to its firmware code.
var pollingRates = map[int]byte{
	125:  0x08, // power saving
	250:  0x04, // office
	500:  0x02, // gaming
	1000: 0x01, // eSports
}

// BuildPolling returns the 9-byte feature report that sets the polling rate.
func BuildPolling(rate int) ([]byte, error) {
	code, ok := pollingRates[rate]
	if !ok {
		return nil, fmt.Errorf("unsupported polling rate %d Hz (use 125, 250, 500 or 1000)", rate)
	}

	buf := make([]byte, pollingPacketLen)
	buf[0] = pollingReportID
	buf[1] = pollingPacketLen
	buf[2] = ProfileID
	buf[3] = code
	buf[4] = code ^ 0xff // checksum
	return buf, nil
}
