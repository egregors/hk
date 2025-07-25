//go:build linux && amd64

package sensors

import "math/rand/v2"

// BME280 is a sensor MOCK for temperature and humidity
type BME280 struct{}

func NewBME280() (*BME280, error) {
	return &BME280{}, nil
}

func (b *BME280) CurrentTemperature() (float64, error) {
	//nolint:gosec // this is a mock
	return 20 + 10*rand.Float64(), nil
}

func (b *BME280) CurrentHumidity() (float64, error) {
	//nolint:gosec // this is a mock
	return 50 + 20*rand.Float64(), nil
}