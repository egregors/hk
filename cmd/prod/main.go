package main

import (
	"context"
	"github.com/egregors/hk/internal/light"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brutella/hap"
	"github.com/brutella/hap/accessory"
	"github.com/d2r2/go-logger"

	"github.com/egregors/hk/internal/homekit"
	"github.com/egregors/hk/internal/metrics"
	"github.com/egregors/hk/internal/sensors"
	"github.com/egregors/hk/log"
	"github.com/egregors/hk/srv"
)

const (
	metricsRetention = 30 * 24 * time.Hour
	hapPIN           = "11112222" // TODO: use secure pin (not this one)
)

var revision = "HEAD"

func main() {
	setupLogger()
	log.Info.Printf("🇭🇰 revision: %s", revision)

	db := hap.NewFsStore("./db")
	m, dumpFn := makeMetrics()
	server := srv.New(db, makeClimate(), makeLight(), makeHkSrv(db), m)

	ctx, cancel := context.WithCancel(context.Background())
	go graceful(cancel, dumpFn)

	if err := server.Run(ctx); err != nil {
		log.Erro.Printf("can't run server: %s", err.Error())
		os.Exit(1)
	}
}

func graceful(cancel context.CancelFunc, dumpFn metrics.DumpFn) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	<-c
	log.Info.Println("server shutdown...")

	signal.Stop(c)

	log.Info.Println("ctx cancel")
	cancel()

	log.Info.Println("try make a dump to restore it next time...")
	if err := dumpFn(); err != nil {
		log.Erro.Printf("can't make a metrics dump: %s", err.Error())
	} else {
		log.Info.Println("done")
	}

	log.Info.Println("bye")

	os.Exit(0)
}

func makeMetrics() (m srv.Metrics, dump metrics.DumpFn) {
	return metrics.New(
		metrics.WithRetention(metricsRetention),
		metrics.WithBackup(),
		metrics.WithAutosave(60*time.Minute),
	)
}

func makeClimate() srv.ClimateSensor {
	bme280, err := sensors.NewBME280()
	if err != nil {
		log.Erro.Printf("can't create BME280 sensor: %s", err.Error())
		os.Exit(1)
	}

	return bme280
}

func makeLight() srv.LightCtrl {
	// TODO: make two different external devices: required and options,
	//  in case of fail of optional device setup just skip it.
	garland, err := light.NewUsbGarland()
	if err != nil {
		log.Erro.Printf("can't create USB garland: %s", err.Error())
		os.Exit(1)
	}

	return garland
}

func makeHkSrv(db hap.Store) *homekit.HapSrv {
	hk, err := homekit.NewHapSrv(&homekit.HapSrvOpts{
		DB:  db,
		Pin: hapPIN,
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
	log.Debg.Off()

	err := logger.ChangePackageLogLevel("i2c", logger.InfoLevel)
	if err != nil {
		log.Erro.Printf("can't setup i2c logger to INTO: %s", err.Error())
	}

	err = logger.ChangePackageLogLevel("bsbmp", logger.InfoLevel)
	if err != nil {
		log.Erro.Printf("can't setup bsbmp logger to INTO: %s", err.Error())
	}
}
