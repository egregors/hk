//go:build linux && amd64

package light

import "github.com/egregors/hk/log"

type UsbGarland struct{}

func NewUsbGarland() (*UsbGarland, error) {
	return &UsbGarland{}, nil
}

func (u *UsbGarland) On() error {
	log.Info.Println("USB garland ON (mock)")
	return nil
}

func (u *UsbGarland) Off() error {
	log.Info.Println("USB garland OFF (mock)")
	return nil
}