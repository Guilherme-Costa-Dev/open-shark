package protocol

// Battery event stream markers (pushed by the mouse in wireless mode).
const (
	batteryReportTag = 0x03
	batteryPrefix    = 0x55
	batteryMsgType   = 0x40
	batteryMsgType1  = 0x41
)

// BatteryStatus is the charging state reported alongside the percentage.
type BatteryStatus byte

const (
	BatteryNormal   BatteryStatus = 0x01
	BatteryFull     BatteryStatus = 0x02
	BatteryCharging BatteryStatus = 0x03
)

func (s BatteryStatus) String() string {
	switch s {
	case BatteryNormal:
		return "on battery"
	case BatteryFull:
		return "fully charged"
	case BatteryCharging:
		return "charging"
	default:
		return "unknown"
	}
}

// Battery is a decoded battery update.
type Battery struct {
	Status  BatteryStatus
	Percent int
}

// ParseBattery decodes a 5-byte event notification of the form
//
//	03 55 40 <status> <percent>
//
// It returns false for any report that is not a battery update.
func ParseBattery(report []byte) (Battery, bool) {
	if len(report) < 5 {
		return Battery{}, false
	}
	if report[0] != batteryReportTag || report[1] != batteryPrefix {
		return Battery{}, false
	}
	if report[2] != batteryMsgType && report[2] != batteryMsgType1 {
		return Battery{}, false
	}

	status := BatteryStatus(report[3])
	switch status {
	case BatteryNormal, BatteryFull, BatteryCharging:
	default:
		return Battery{}, false
	}

	percent := int(report[4])
	if percent > 100 {
		return Battery{}, false
	}

	return Battery{Status: status, Percent: percent}, true
}
