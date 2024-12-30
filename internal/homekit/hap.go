package homekit

import (
	"context"
	"net/http"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"

	"github.com/egregors/hk/log"
)

type HapSrvOpts struct {
	DB  hap.Store
	Pin string

	Bridge      *accessory.Bridge
	Thermometer *accessory.Thermometer
	Humidifier  *accessory.Humidifier
	USB2Power   *accessory.Switch
}

type HapSrv struct {
	srv         *hap.Server
	thermometer *accessory.Thermometer
	humidifier  *accessory.Humidifier
	usb2power   *accessory.Switch
}

func NewHapSrv(hapSrvOpts *HapSrvOpts) (*HapSrv, error) {
	log.Info.Println("make HapSrv")

	// see: https://github.com/brutella/hap/pull/53
	hapSrvOpts.Bridge.A.Id = 1
	hapSrvOpts.Thermometer.A.Id = 2
	hapSrvOpts.Humidifier.A.Id = 3
	hapSrvOpts.USB2Power.A.Id = 4

	s, err := hap.NewServer(
		hapSrvOpts.DB,
		hapSrvOpts.Bridge.A,
		hapSrvOpts.Thermometer.A,
		hapSrvOpts.Humidifier.A,
		hapSrvOpts.USB2Power.A,
	)
	if err != nil {
		return nil, err
	}

	if hapSrvOpts.Pin != "" {
		log.Info.Printf("set custom PIN")
		s.Pin = hapSrvOpts.Pin
	}

	hapSrvOpts.USB2Power.Switch.On.SetValue(true)

	return &HapSrv{
		srv:         s,
		thermometer: hapSrvOpts.Thermometer,
		humidifier:  hapSrvOpts.Humidifier,
		usb2power:   hapSrvOpts.USB2Power,
	}, nil
}

func (s *HapSrv) USB2PowerChan() chan bool {
	log.Debg.Printf("usb2power is %v after the start", s.usb2power.Switch.On.Value())

	ch := make(chan bool)
	s.usb2power.Switch.On.OnValueUpdate(func(old, new bool, req *http.Request) {
		log.Debg.Printf("got usb2power %v -> %v\n", old, new)
		ch <- new
	})

	// FIXME: nobody close it
	return ch
}

func (s *HapSrv) SetCurrentTemperature(t float64) {
	s.thermometer.TempSensor.CurrentTemperature.SetValue(t)
}

func (s *HapSrv) SetCurrentHumidity(h float64) {
	s.humidifier.Humidifier.CurrentRelativeHumidity.SetValue(h)
}

func (s *HapSrv) ListenAndServe(ctx context.Context) error {
	return s.srv.ListenAndServe(ctx)
}
