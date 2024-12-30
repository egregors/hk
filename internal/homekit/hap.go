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
	Light       *accessory.Lightbulb
}

type HapSrv struct {
	srv         *hap.Server
	thermometer *accessory.Thermometer
	humidifier  *accessory.Humidifier
	light       *accessory.Lightbulb
}

func NewHapSrv(hapSrvOpts *HapSrvOpts) (*HapSrv, error) {
	log.Info.Println("make HapSrv")

	// see: https://github.com/brutella/hap/pull/53
	hapSrvOpts.Bridge.A.Id = 1
	hapSrvOpts.Thermometer.A.Id = 2
	hapSrvOpts.Humidifier.A.Id = 3
	hapSrvOpts.Light.A.Id = 4

	s, err := hap.NewServer(
		hapSrvOpts.DB,
		hapSrvOpts.Bridge.A,
		hapSrvOpts.Thermometer.A,
		hapSrvOpts.Humidifier.A,
		hapSrvOpts.Light.A,
	)
	if err != nil {
		return nil, err
	}

	if hapSrvOpts.Pin != "" {
		log.Info.Printf("set custom PIN")
		s.Pin = hapSrvOpts.Pin
	}

	return &HapSrv{
		srv:         s,
		thermometer: hapSrvOpts.Thermometer,
		humidifier:  hapSrvOpts.Humidifier,
		light:       hapSrvOpts.Light,
	}, nil
}

func (s *HapSrv) LightEventsCh() chan bool {
	ch := make(chan bool)
	s.light.Lightbulb.On.OnValueUpdate(func(old, new bool, req *http.Request) {
		log.Debg.Printf("got light %v -> %v\n", old, new)
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
