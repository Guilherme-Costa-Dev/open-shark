// Package hid locates the Attack Shark configuration interface and exchanges
// feature reports and event notifications with it.
package hid

import (
	"errors"
	"fmt"
	"time"

	"github.com/Guilherme-Costa-Dev/open-shark/internal/protocol"
	gohid "github.com/sstallion/go-hid"
)

// ErrNotFound is returned when no Attack Shark mouse is connected.
var ErrNotFound = errors.New("no Attack Shark X11 configuration interface found (is the mouse connected?)")

// Device is an open handle to the mouse's configuration interface.
type Device struct {
	dev  *gohid.Device
	Path string
	Mode string
}

// Open initializes hidapi and opens the configuration interface of the first
// supported mouse it finds.
func Open() (*Device, error) {
	if err := gohid.Init(); err != nil {
		return nil, fmt.Errorf("hidapi init: %w", err)
	}

	targets := []struct {
		pid  uint16
		mode string
	}{
		{protocol.ProductWireless, "wireless (2.4 GHz)"},
		{protocol.ProductWired, "wired"},
		{protocol.ProductR1Wired, "R1 (wired)"},
	}

	var path, mode string
	for _, t := range targets {
		var found string
		err := gohid.Enumerate(protocol.VendorID, t.pid, func(info *gohid.DeviceInfo) error {
			if info.InterfaceNbr == protocol.ConfigInterface {
				found = info.Path
			}
			return nil
		})
		if err != nil {
			continue
		}
		if found != "" {
			path, mode = found, t.mode
			break
		}
	}

	if path == "" {
		_ = gohid.Exit()
		return nil, ErrNotFound
	}

	dev, err := gohid.OpenPath(path)
	if err != nil {
		_ = gohid.Exit()
		return nil, fmt.Errorf("open %s (permission denied? run with sudo): %w", path, err)
	}

	return &Device{dev: dev, Path: path, Mode: mode}, nil
}

// Close releases the device and the hidapi context.
func (d *Device) Close() error {
	if d.dev != nil {
		_ = d.dev.Close()
	}
	return gohid.Exit()
}

// SendFeature transmits a feature report. The first byte must be the report ID.
func (d *Device) SendFeature(payload []byte) error {
	n, err := d.dev.SendFeatureReport(payload)
	if err != nil {
		return fmt.Errorf("send feature report: %w", err)
	}
	if n != len(payload) {
		return fmt.Errorf("send feature report: wrote %d of %d bytes", n, len(payload))
	}
	return nil
}

// ReadBattery waits for the next battery update pushed by the mouse.
func (d *Device) ReadBattery(timeout time.Duration) (protocol.Battery, error) {
	const noUpdate = "timed out waiting for a battery update (the battery is only reported in 2.4 GHz wireless mode)"

	deadline := time.Now().Add(timeout)
	buf := make([]byte, 64)
	for {
		remaining := time.Until(deadline)
		if remaining <= 0 {
			return protocol.Battery{}, errors.New(noUpdate)
		}

		n, err := d.dev.ReadWithTimeout(buf, remaining)
		if err != nil {
			if errors.Is(err, gohid.ErrTimeout) {
				return protocol.Battery{}, errors.New(noUpdate)
			}
			return protocol.Battery{}, fmt.Errorf("read event: %w", err)
		}

		if b, ok := protocol.ParseBattery(buf[:n]); ok {
			return b, nil
		}
	}
}
