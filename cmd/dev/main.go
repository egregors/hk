package main

import (
	"context"
	"fmt"
	"github.com/egregors/hk/utils/bp"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/brutella/hap"
	"github.com/d2r2/go-logger"

	"github.com/egregors/hk/internal/homekit"
	"github.com/egregors/hk/internal/metrics"
	"github.com/egregors/hk/internal/sensors"
	"github.com/egregors/hk/log"
	"github.com/egregors/hk/srv"
)

const (
	metricsRetention = 2 * time.Minute
)

//nolint:unused // for debug
func main_() {
	dataT := []float64{22.03, 21.30, 20.70, 20.74, 21.07, 21.32, 21.17, 21.28, 20.68, 20.75, 20.93, 20.80, 20.85, 21.32, 21.34, 20.80, 20.86, 21.09, 21.05, 21.10, 20.90, 20.05, 21.48, 20.94, 21.37}
	dataH := []float64{55.15, 55.15, 55.28, 54.47, 54.51, 55.34, 53.96, 53.62, 52.89, 52.51, 52.40, 52.38, 52.21, 52.08, 51.94, 51.74, 51.36, 50.48, 48.94, 49.92, 53.06, 53.18, 53.60, 54.20, 53.40}

	fmt.Println("T:")
	fmt.Println(bp.SimplePlot(4, dataT))
	fmt.Println(bp.SimplePlot(4, dataH))
}

func main() {
	setupLogger()

	db := hap.NewFsStore("./db")
	m, dumpFn := makeMetrics()
	server := srv.New(db, makeClimate(), makeFakeHkSrv(), m)

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

	log.Info.Println("buy")

	os.Exit(0)
}

func makeMetrics() (m srv.Metrics, dump metrics.DumpFn) {
	return metrics.New(metrics.WithRetention(metricsRetention), metrics.WithBackup())
}

func makeClimate() srv.ClimateSensor {
	bme280, err := sensors.NewBME280()
	if err != nil {
		log.Erro.Printf("can't create BME280 sensor: %s", err.Error())
		os.Exit(1)
	}

	return bme280
}

func makeFakeHkSrv() *homekit.NoopHap {
	return &homekit.NoopHap{}
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
