package agent

import (
	"testing"
	"time"
)

func TestDetectorTriggersOnPPS(t *testing.T) {
	cfg := Config{
		AlertPPS:       1000,
		AlertBPS:       1 << 30,
		AlertSYNRatio:  0.99,
		AlertUniqueSrc: 100000,
	}
	d := newDetector(3 * time.Second)

	event, hit := d.evaluate(time.Now(), PacketStats{
		PacketsPerSecond: 2500,
		BitsPerSecond:    1000000,
		UniqueSources:    10,
	}, cfg, []string{"eth0"}, "test-vm")

	if !hit {
		t.Fatal("expected detector to trigger")
	}
	if event.HostID != "test-vm" || event.Severity == "" {
		t.Fatalf("unexpected event: %+v", event)
	}
}

func TestDetectorDoesNotTriggerBelowThresholds(t *testing.T) {
	cfg := Config{
		AlertPPS:       10000,
		AlertBPS:       1000000000,
		AlertSYNRatio:  0.99,
		AlertUniqueSrc: 10000,
	}
	d := newDetector(3 * time.Second)

	_, hit := d.evaluate(time.Now(), PacketStats{
		PacketsPerSecond: 100,
		BitsPerSecond:    1000000,
		UniqueSources:    5,
	}, cfg, []string{"eth0"}, "test-vm")

	if hit {
		t.Fatal("did not expect detector to trigger")
	}
}

func TestDetectorReset(t *testing.T) {
	d := newDetector(3 * time.Second)
	d.evaluate(time.Now(), PacketStats{PacketsPerSecond: 1}, Config{AlertPPS: 1}, nil, "test")
	d.reset()

	if len(d.pps) != 0 || len(d.bps) != 0 || len(d.synRatio) != 0 || len(d.unique) != 0 {
		t.Fatal("detector was not reset")
	}
}
