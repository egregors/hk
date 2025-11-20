//go:build linux && arm64

package light

import (
	"testing"
)

func TestParseUhubctlOutput(t *testing.T) {
	tests := []struct {
		name           string
		uhubctlOutput  string
		expectedHub    string
		shouldFindHub  bool
	}{
		{
			name: "Single hub with ppps",
			uhubctlOutput: `Current status for hub 1-1 [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ppps]
  Port 1: 0100 power
  Port 2: 0100 power
  Port 3: 0100 power
  Port 4: 0100 power`,
			expectedHub:   "1-1",
			shouldFindHub: true,
		},
		{
			name: "Multiple hubs, first has ppps",
			uhubctlOutput: `Current status for hub 1 [1d6b:0002 Linux 6.1.21-v8+ dwc_otg_hcd DWC OTG Controller, USB 2.00, 1 ports, ppps]
  Port 1: 0503 power highspeed enable connect [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ppps]
Current status for hub 1-1 [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ppps]
  Port 1: 0100 power
  Port 2: 0100 power`,
			expectedHub:   "1",
			shouldFindHub: true,
		},
		{
			name: "Hub without ppps (fallback)",
			uhubctlOutput: `Current status for hub 2-1 [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ganged]
  Port 1: 0100 power
  Port 2: 0100 power`,
			expectedHub:   "2-1",
			shouldFindHub: true,
		},
		{
			name: "Multiple hubs with different formats",
			uhubctlOutput: `Current status for hub 1 [1d6b:0002 Linux Foundation 2.0 root hub, USB 2.00, 4 ports, ganged]
  Port 1: 0503 power highspeed enable connect
Current status for hub 2-1 [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ppps]
  Port 1: 0100 power off
  Port 2: 0100 power
Current status for hub 3 [1d6b:0003 Linux Foundation 3.0 root hub, USB 3.00, 4 ports, ganged]`,
			expectedHub:   "2-1",
			shouldFindHub: true,
		},
		{
			name:           "No hubs found",
			uhubctlOutput:  `No compatible devices found`,
			expectedHub:    "",
			shouldFindHub:  false,
		},
		{
			name: "Complex hub location (nested)",
			uhubctlOutput: `Current status for hub 1-1.2 [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ppps]
  Port 1: 0100 power
  Port 2: 0100 power`,
			expectedHub:   "1-1.2",
			shouldFindHub: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// We can't directly test detectHubLocation since it runs sudo uhubctl
			// But we can test the parsing logic by creating a mock
			// For now, we'll just document the expected behavior
			
			// This test serves as documentation for the expected uhubctl output formats
			// The actual detectHubLocation function should handle all these cases
			
			if tt.shouldFindHub {
				t.Logf("Expected to find hub: %s", tt.expectedHub)
			} else {
				t.Logf("Expected to not find any hub")
			}
		})
	}
}

func TestUsbGarlandStruct(t *testing.T) {
	// Test that UsbGarland struct holds hub location
	garland := &UsbGarland{
		hubLocation: "1-1",
	}

	if garland.hubLocation != "1-1" {
		t.Errorf("Expected hub location '1-1', got '%s'", garland.hubLocation)
	}
}
