package protocol

import (
	"bytes"
	"testing"
)

func TestBuildPolling1000Hz(t *testing.T) {
	got, err := BuildPolling(1000)
	if err != nil {
		t.Fatal(err)
	}
	want := []byte{0x06, 0x09, 0x01, 0x01, 0xfe, 0x00, 0x00, 0x00, 0x00}
	if !bytes.Equal(got, want) {
		t.Fatalf("got % x, want % x", got, want)
	}
}

func TestBuildPollingRejectsUnknownRate(t *testing.T) {
	if _, err := BuildPolling(750); err == nil {
		t.Fatal("expected an error for an unsupported rate")
	}
}

func TestBuildDPILayout(t *testing.T) {
	got, err := BuildDPI(DPISettings{
		Stages: [6]int{800, 1600, 2400, 3200, 5000, 22000},
		Active: 2,
	})
	if err != nil {
		t.Fatal(err)
	}
	if len(got) != dpiPacketLen {
		t.Fatalf("packet length = %d, want %d", len(got), dpiPacketLen)
	}
	if got[0] != dpiReportID || got[1] != dpiPacketLen || got[2] != ProfileID {
		t.Fatalf("bad header: % x", got[:3])
	}
	if got[5] != activeStageMsk {
		t.Fatalf("active stage mask = %#x, want %#x", got[5], activeStageMsk)
	}

	// 800 DPI encodes to X=0x12; 22000 DPI uses the double-speed path.
	if got[8] != 0x12 || got[16] != 0x00 {
		t.Fatalf("stage 1 bytes = %#x/%#x, want 0x12/0x00", got[8], got[16])
	}
	if got[13] != 0xd0 || got[21] != 0x00 {
		t.Fatalf("stage 6 bytes = %#x/%#x, want 0xd0/0x00", got[13], got[21])
	}
	if got[6] != 0x20 {
		t.Fatalf("double flags = %#x, want 0x20", got[6])
	}
	if got[7] != 0x00 {
		t.Fatalf("triple flags = %#x, want 0x00", got[7])
	}
	if got[24] != 2 {
		t.Fatalf("active stage = %d, want 2", got[24])
	}

	sum := 0
	for i := 3; i <= 49; i++ {
		sum += int(got[i])
	}
	if got[50] != byte(sum>>8) || got[51] != byte(sum&0xff) {
		t.Fatalf("checksum = %#x %#x, want %#x %#x", got[50], got[51], byte(sum>>8), byte(sum&0xff))
	}
}

func TestBuildDPIRejectsOutOfRangeStage(t *testing.T) {
	_, err := BuildDPI(DPISettings{Stages: [6]int{10, 1600, 2400, 3200, 5000, 22000}, Active: 1})
	if err == nil {
		t.Fatal("expected an error for a DPI below the minimum")
	}
}

func TestParseBattery(t *testing.T) {
	b, ok := ParseBattery([]byte{0x03, 0x55, 0x40, 0x01, 0x50})
	if !ok || b.Percent != 80 || b.Status != BatteryNormal {
		t.Fatalf("got %+v ok=%v, want 80%%/normal", b, ok)
	}

	if _, ok := ParseBattery([]byte{0x02, 0x00, 0x00, 0x00, 0x00}); ok {
		t.Fatal("expected a non-battery report to be rejected")
	}
	if _, ok := ParseBattery([]byte{0x03, 0x55, 0x40, 0x09, 0x50}); ok {
		t.Fatal("expected an unknown status byte to be rejected")
	}
}
