package main

import (
	"context"
	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/egregors/hk/hkSrv"
	"github.com/egregors/hk/sensors"
	"github.com/egregors/hk/srv"
	"os"
	"os/signal"
	"syscall"
)

func makeClimate() srv.ClimateSensor {
	climate, err := sensors.NewBME280()
	if err != nil {
		panic(err)
	}

	return climate
}

func makeHkSrv(db hap.Store) *hkSrv.HapSrv {
	hk, err := hkSrv.New(&hkSrv.HapSrvOpts{
		DB:  db,
		Pin: "11112222", // TODO: do not forget change pin!
		Bridge: accessory.NewBridge(accessory.Info{
			Name:         "Raspberry Pi5",
			SerialNumber: "-",
			Manufacturer: "Raspberry Pi",
			Model:        "Model 5",
			Firmware:     "-",
		}),
		Thermometer: accessory.NewTemperatureSensor(accessory.Info{
			Name:         "Temperature",
			SerialNumber: "-",
			Manufacturer: "bosch",
			Model:        "BME280",
			Firmware:     "-",
		}),
		Humidifier: accessory.NewHumidifier(accessory.Info{
			Name:         "Humidity",
			SerialNumber: "-",
			Manufacturer: "bosch",
			Model:        "BME280",
			Firmware:     "-",
		}),
	})
	if err != nil {
		panic(err)
	}

	return hk
}

func main() {
	// TODO:
	// 	- [ ] don't panic
	// 	- [ ] add proper logger
	//  - [ ] collect metrics (t, h)

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-c
		signal.Stop(c)
		cancel()
	}()

	db := hap.NewFsStore("./db")
	s := srv.New(db, makeClimate(), makeHkSrv(db))

	if err := s.Run(ctx); err != nil {
		panic(err)
	}
}
