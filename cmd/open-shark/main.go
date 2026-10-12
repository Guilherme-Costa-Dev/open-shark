// Command open-shark configures an Attack Shark X11 mouse from a TOML file
// and reads its battery level. It replaces the Windows-only GUI.
package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"time"

	"github.com/Guilherme-Costa-Dev/open-shark/internal/config"
	"github.com/Guilherme-Costa-Dev/open-shark/internal/hid"
	"github.com/Guilherme-Costa-Dev/open-shark/internal/protocol"
)

const (
	// pacingDelay is the pause between feature reports. The firmware hangs if
	// commands are sent too quickly.
	pacingDelay = 500 * time.Millisecond
	// batteryWait is how long to wait for a battery event.
	batteryWait = 5 * time.Second
)

func main() {
	if len(os.Args) < 2 {
		usage()
		os.Exit(2)
	}

	var err error
	switch os.Args[1] {
	case "apply":
		err = cmdApply(os.Args[2:])
	case "battery":
		err = cmdBattery(os.Args[2:])
	case "info":
		err = cmdInfo(os.Args[2:])
	case "help", "-h", "--help":
		usage()
		return
	default:
		fmt.Fprintf(os.Stderr, "unknown command %q\n\n", os.Args[1])
		usage()
		os.Exit(2)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

// cmdApply reads the config file and writes DPI and polling rate to the mouse.
func cmdApply(args []string) error {
	fs := flag.NewFlagSet("apply", flag.ExitOnError)
	cfgPath := fs.String("c", config.DefaultPath(), "path to the TOML configuration file")
	fs.Parse(args)

	cfg, err := config.Load(*cfgPath)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return fmt.Errorf("%w\ncopy config.example.toml to %s first", err, *cfgPath)
		}
		return err
	}

	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	fmt.Printf("device: %s (%s)\n", dev.Path, dev.Mode)

	dpiPkt, err := protocol.BuildDPI(protocol.DPISettings{
		Stages: [6]int(cfg.DPI.Stages),
		Active: cfg.DPI.Active,
	})
	if err != nil {
		return err
	}
	if err := dev.SendFeature(dpiPkt); err != nil {
		return fmt.Errorf("apply DPI: %w", err)
	}
	fmt.Printf("dpi: stages=%v active=%d\n", cfg.DPI.Stages, cfg.DPI.Active)
	time.Sleep(pacingDelay)

	pollPkt, err := protocol.BuildPolling(cfg.Polling.Rate)
	if err != nil {
		return err
	}
	if err := dev.SendFeature(pollPkt); err != nil {
		return fmt.Errorf("apply polling rate: %w", err)
	}
	fmt.Printf("polling rate: %d Hz\n", cfg.Polling.Rate)

	return nil
}

// cmdBattery prints the current battery level.
func cmdBattery(args []string) error {
	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	b, err := dev.ReadBattery(batteryWait)
	if err != nil {
		return err
	}
	fmt.Printf("battery: %d%% (%s)\n", b.Percent, b.Status)
	return nil
}

// cmdInfo prints device details and the battery level when available.
func cmdInfo(args []string) error {
	dev, err := hid.Open()
	if err != nil {
		return err
	}
	defer dev.Close()

	fmt.Printf("device: %s\nmode:   %s\n", dev.Path, dev.Mode)

	b, err := dev.ReadBattery(batteryWait)
	if err != nil {
		fmt.Printf("battery: unavailable (%v)\n", err)
		return nil
	}
	fmt.Printf("battery: %d%% (%s)\n", b.Percent, b.Status)
	return nil
}

func usage() {
	fmt.Fprint(os.Stderr, `open-shark - configure an Attack Shark X11 mouse from the command line

usage:
  open-shark apply [-c config.toml]   apply DPI and polling rate from the config file
  open-shark battery                  read the battery level
  open-shark info                     show the device and battery

The config file defaults to `+config.DefaultPath()+`.
`)
}
