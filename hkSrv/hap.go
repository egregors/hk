package hkSrv

import (
	"context"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
)

type HapSrvOpts struct {
	DB  hap.Store
	Pin string

	Bridge      *accessory.Bridge
	Thermometer *accessory.Thermometer
	Humidifier  *accessory.Humidifier
}

type HapSrv struct {
	srv         *hap.Server
	thermometer *accessory.Thermometer
	humidifier  *accessory.Humidifier
}

func New(hapSrvOpts *HapSrvOpts) (*HapSrv, error) {
	s, err := hap.NewServer(
		hapSrvOpts.DB,
		hapSrvOpts.Bridge.A,
		hapSrvOpts.Thermometer.A,
		hapSrvOpts.Humidifier.A,
	)
	if err != nil {
		return nil, err
	}

	if hapSrvOpts.Pin != "" {
		s.Pin = hapSrvOpts.Pin
	}

	return &HapSrv{
		srv:         s,
		thermometer: hapSrvOpts.Thermometer,
		humidifier:  hapSrvOpts.Humidifier,
	}, nil
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
