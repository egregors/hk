package main

import (
	"context"
	"flag"
	"github.com/egregors/hk/internal/light"
	"github.com/egregors/hk/internal/ntfy"
	"os"
	"os/signal"
	"strconv"
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
	metricsRetention = 3600 * time.Hour
)

var revision string = "HEAD"

func main() {
	// Parse command line flags
	ntfyURL := flag.String("ntfy-url", "", "ntfy.sh URL with topic for error notifications")
	ntfyEnabled := flag.Bool("ntfy-enabled", false, "enable ntfy notifications")
	flag.Parse()

	// Check environment variables if flags not provided
	if *ntfyURL == "" {
		*ntfyURL = os.Getenv("NTFY_URL")
	}
	if !*ntfyEnabled {
		if enabled := os.Getenv("NTFY_ENABLED"); enabled != "" {
			if parsed, err := strconv.ParseBool(enabled); err == nil {
				*ntfyEnabled = parsed
			}
		}
	}

	setupLogger()
	log.Info.Printf("🇭🇰 revision: %s", revision)

	db := hap.NewFsStore("./db")
	m, dumpFn := makeMetrics()
	notifier := makeNotifier(*ntfyURL, *ntfyEnabled)
	server := srv.New(db, makeClimate(), makeLight(), makeFakeHkSrv(), m, notifier)

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

func makeLight() srv.USB2PowerCtrl {
	// TODO: make two different external devices: required and options,
	//  in case of fail of optional device setup just skip it.
	garland, err := light.NewUsbGarland()
	if err != nil {
		log.Erro.Printf("can't create USB garland: %s", err.Error())
		os.Exit(1)
	}

	return garland
}

func makeFakeHkSrv() *homekit.NoopHap {
	return &homekit.NoopHap{}
}

func makeNotifier(url string, enabled bool) srv.Notifier {
	if !enabled || url == "" {
		log.Info.Println("ntfy notifications disabled")
		return &ntfy.NoopNotifier{}
	}

	log.Info.Printf("ntfy notifications enabled with URL: %s", url)
	return ntfy.New(ntfy.Config{
		URL:     url,
		Enabled: enabled,
	})
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
