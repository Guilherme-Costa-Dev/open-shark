// Package protocol implements the HID packet formats used by the
// Attack Shark X11 mouse.
//
// The layouts are ported from the community reverse-engineering effort in
// HarukaYamamoto0/attack-shark-x11-driver and cross-checked against the
// device's HID report descriptor.
package protocol

const (
	// VendorID is shared by the whole Attack Shark X11 family.
	VendorID = 0x1d57

	// Product IDs for the supported connection modes.
	ProductWireless = 0xfa60 // 2.4 GHz dongle
	ProductWired    = 0xfa55 // X11 connected by cable
	ProductR1Wired  = 0xfa61 // R1 connected by cable

	// ConfigInterface is the HID interface that carries the configuration
	// feature reports and the battery event stream.
	ConfigInterface = 2

	// ProfileID is the on-board profile that settings are written to.
	ProfileID = 0x01

	// DPI limits for the PAW3311 sensor in this mouse.
	MinDPI = 50
	MaxDPI = 22000

	// Report IDs.
	dpiReportID     = 0x04
	pollingReportID = 0x06

	// Report packet lengths (excluding the leading report ID).
	dpiPacketLen     = 0x38 // 56 bytes
	pollingPacketLen = 0x09 // 9 bytes
)
