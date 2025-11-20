//go:build linux && arm64

package light

import (
	"bytes"
	"fmt"
	"os/exec"
	"regexp"
	"strings"

	"github.com/egregors/hk/log"
	"github.com/egregors/hk/utils/cli"
)

type UsbGarland struct {
	hubLocation string
}

func NewUsbGarland() (*UsbGarland, error) {
	err := cli.CheckCommandExists("uhubctl")
	if err != nil {
		return nil, err
	}

	hubLocation, err := detectHubLocation()
	if err != nil {
		return nil, fmt.Errorf("failed to detect USB hub location: %w", err)
	}

	log.Info.Printf("detected USB hub location: %s", hubLocation)
	return &UsbGarland{hubLocation: hubLocation}, nil
}

// detectHubLocation runs uhubctl to list all hubs and finds a controllable one
func detectHubLocation() (string, error) {
	cmd := exec.Command("sudo", "uhubctl")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return "", fmt.Errorf("failed to run uhubctl: %w (output: %s)", err, out.String())
	}

	output := out.String()
	log.Debg.Printf("uhubctl output:\n%s", output)

	// Parse uhubctl output to find hub locations
	// Example output format:
	// Current status for hub 1-1 [05e3:0610 USB2.0 Hub, USB 2.10, 4 ports, ppps]
	//   Port 1: 0100 power
	// or
	// Current status for hub 1 [1d6b:0002 Linux 6.1.21-v8+ dwc_otg_hcd DWC OTG Controller, USB 2.00, 1 ports, ppps]
	locationRegex := regexp.MustCompile(`Current status for hub ([0-9-]+)`)
	matches := locationRegex.FindAllStringSubmatch(output, -1)

	if len(matches) == 0 {
		return "", fmt.Errorf("no USB hubs found in uhubctl output")
	}

	// Try to find a hub that supports per-port power switching (ppps)
	// This is indicated by "ppps" in the hub description
	for _, match := range matches {
		if len(match) > 1 {
			location := match[1]
			// Check if this hub line contains "ppps" (per-port power switching)
			lineStart := strings.Index(output, match[0])
			lineEnd := strings.Index(output[lineStart:], "\n")
			if lineEnd == -1 {
				lineEnd = len(output[lineStart:])
			}
			line := output[lineStart : lineStart+lineEnd]

			if strings.Contains(line, "ppps") {
				log.Info.Printf("found controllable USB hub at location: %s", location)
				return location, nil
			}
		}
	}

	// If no hub with ppps found, use the first hub as fallback
	if len(matches) > 0 && len(matches[0]) > 1 {
		location := matches[0][1]
		log.Info.Printf("using first available USB hub at location: %s (fallback)", location)
		return location, nil
	}

	return "", fmt.Errorf("no suitable USB hub found")
}

func (u *UsbGarland) On() error {
	cmd := exec.Command("sudo", "uhubctl", "-l", u.hubLocation, "-a", "on")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to turn on USB power (hub %s): %w (output: %s)", u.hubLocation, err, out.String())
	}

	log.Debg.Printf("USB power ON (hub %s): %s", u.hubLocation, out.String())
	return nil
}

func (u *UsbGarland) Off() error {
	cmd := exec.Command("sudo", "uhubctl", "-l", u.hubLocation, "-a", "off")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	err := cmd.Run()
	if err != nil {
		return fmt.Errorf("failed to turn off USB power (hub %s): %w (output: %s)", u.hubLocation, err, out.String())
	}

	log.Debg.Printf("USB power OFF (hub %s): %s", u.hubLocation, out.String())
	return nil
}
