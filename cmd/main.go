package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/d2r2/go-logger"

	"github.com/egregors/hk/hkSrv"
	"github.com/egregors/hk/log"
	"github.com/egregors/hk/sensors"
	"github.com/egregors/hk/srv"
)

func main() {
	setupLogger()

	ctx, cancel := context.WithCancel(context.Background())
	db := hap.NewFsStore("./db")
	server := srv.New(db, makeClimate(), makeHkSrv(db))

	go graceful(cancel)

	if err := server.Run(ctx); err != nil {
		log.Erro.Printf("can't run server: %s", err.Error())
		os.Exit(1)
	}
}

func graceful(cancel context.CancelFunc) {
	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info.Println("server shutdown...")

	signal.Stop(c)
	cancel()

	os.Exit(0)
}

func makeClimate() srv.ClimateSensor {
	climate, err := sensors.NewBME280()
	if err != nil {
		log.Erro.Printf("can't create BME280 sensor: %s", err.Error())
		os.Exit(1)
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
		log.Erro.Printf("can't create HAP server: %s", err.Error())
		os.Exit(1)
	}

	return hk
}

func setupLogger() {
	err := logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	if err != nil {
		log.Erro.Printf("can't setup i2c logger to INTO: %s", err.Error())
	}

	err = logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)
	if err != nil {
		log.Erro.Printf("can't setup bsbmp logger to INTO: %s", err.Error())
	}
}
