package sensors

import (
	"github.com/d2r2/go-bsbmp"
	"github.com/d2r2/go-i2c"
	"github.com/d2r2/go-logger"
)

type BME280 struct {
	sensor *bsbmp.BMP
}

func NewBME280() (*BME280, error) {
	_ = logger.ChangePackageLogLevel("conn", logger.InfoLevel)
	_ = logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)

	conn, err := i2c.NewI2C(0x77, 1)
	if err != nil {
		return nil, err
	}

	defer func() { _ = conn.Close() }()

	sensor, err := bsbmp.NewBMP(bsbmp.BME280, conn)
	if err != nil {
		return nil, err
	}

	return &BME280{sensor: sensor}, nil
}

func (b *BME280) CurrentTemperature() (float64, error) {
	t, err := b.sensor.ReadTemperatureC(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return 0, err
	}

	return float64(t), nil
}

func (b *BME280) CurrentHumidity() (float64, error) {
	_, h, err := b.sensor.ReadHumidityRH(bsbmp.ACCURACY_STANDARD)
	if err != nil {
		return 0, err
	}

	return float64(h), nil
}
