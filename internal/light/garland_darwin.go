//go:build darwin

package light

import "github.com/egregors/hk/log"

type UsbGarland struct{}

func NewUsbGarland() (*UsbGarland, error) {
	return &UsbGarland{}, nil
}

func (u *UsbGarland) On() error {
	log.Debg.Printf("call garland ON")

	return nil
}

func (u *UsbGarland) Off() error {
	log.Debg.Printf("call garland OFF")

	return nil
}
