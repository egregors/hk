//go:build linux && arm64

package light

import (
	"os/exec"

	"github.com/egregors/hk/utils/cli"
)

type UsbGarland struct {
	cmdOn, cmdOff *exec.Cmd
}

func NewUsbGarland() (usbGarland *UsbGarland, err error) {
	err = cli.CheckCommandExists("uhubctl")
	if err != nil {
		return nil, err
	}

	// sudo uhubctl -l 1 -a on
	usbGarland = new(UsbGarland)
	usbGarland.cmdOn = exec.Command("sudo", "uhubctl", "-l", "1", "-a", "on")
	usbGarland.cmdOff = exec.Command("sudo", "uhubctl", "-l", "1", "-a", "off")

	return
}

func (u *UsbGarland) On() error {
	return u.cmdOn.Run()
}

func (u *UsbGarland) Off() error {
	return u.cmdOff.Run()
}
