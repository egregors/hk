package homekit

import (
	"context"
)

type NoopHap struct{}

func (n NoopHap) SetCurrentTemperature(t float64) {}

func (n NoopHap) SetCurrentHumidity(h float64) {}

func (n NoopHap) ListenAndServe(ctx context.Context) error {
	<-ctx.Done()

	return nil
}
