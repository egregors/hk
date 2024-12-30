//go:build linux && arm64

package light

import (
	"os/exec"

	"github.com/egregors/hk/utils/cli"
)

type UsbGarland struct{}

func NewUsbGarland() (*UsbGarland, error) {
	// sudo uhubctl -l 1 -a on
	err := cli.CheckCommandExists("uhubctl")
	if err != nil {
		return nil, err
	}

	return &UsbGarland{}, nil
}

func (u *UsbGarland) On() error {
	cmdOn := exec.Command("sudo", "uhubctl", "-l", "1", "-a", "on")

	return cmdOn.Run()
}

func (u *UsbGarland) Off() error {
	cmdOff := exec.Command("sudo", "uhubctl", "-l", "1", "-a", "off")

	return cmdOff.Run()
}
