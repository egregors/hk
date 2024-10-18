package sensors

import "math/rand"

type Dummy struct{}

func (d *Dummy) CurrentTemperature() (float64, error) {
	return 30 + 10*rand.Float64(), nil
}

func (d *Dummy) CurrentHumidity() (float64, error) {
	return 50 + 10*rand.Float64(), nil
}
