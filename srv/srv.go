package srv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/brutella/hap"
	"golang.org/x/sync/errgroup"

	"github.com/egregors/hk/log"
)

const (
	pullPushSleep = 5

	temperatureKey = "current_temperature"
	humidityKey    = "current_humidity"
)

type HapServer interface {
	SetCurrentTemperature(t float64)
	SetCurrentHumidity(h float64)

	ListenAndServe(ctx context.Context) error
}

type ClimateSensor interface {
	CurrentTemperature() (float64, error)
	CurrentHumidity() (float64, error)
}

type Store interface {
	hap.Store
}

type Metrics interface {
	Gauge(key string, val float64)
	GetForPeriodByH(key string, dur time.Duration) map[string]float64
}

type Server struct {
	webSrv  *http.Server
	hkSrv   HapServer
	climate ClimateSensor
	store   Store
	metrics Metrics

	mu           *sync.RWMutex
	currT, currH float64
}

func New(store Store, climate ClimateSensor, hapSrv HapServer, metrics Metrics) *Server {
	return &Server{
		webSrv:  nil,
		hkSrv:   hapSrv,
		climate: climate,
		store:   store,
		metrics: metrics,
		mu:      &sync.RWMutex{},
	}
}

func (s *Server) Run(ctx context.Context) error {
	go func() {
		log.Info.Printf("start syncing sensor data with %d seconds sleep", pullPushSleep)
		for {
			s.pullDataFromSensor()
			s.pushDataToHK()
			<-time.After(pullPushSleep * time.Second)
		}
	}()

	g, ctx := errgroup.WithContext(ctx)

	// go web server
	g.Go(func() error {
		log.Info.Println("start web server on http://localhost:80")
		return s.runWebServer()
	})
	// go hap server
	g.Go(func() error {
		log.Info.Println("start HAP server")
		return s.runHapServer(ctx)
	})

	return g.Wait()
}

func (s *Server) pullDataFromSensor() {
	s.mu.Lock()
	defer s.mu.Unlock()

	var (
		err  error
		t, h float64
	)

	g := new(errgroup.Group)
	g.Go(func() error {
		t, err = s.climate.CurrentTemperature()
		return err
	})
	g.Go(func() error {
		h, err = s.climate.CurrentHumidity()
		return err
	})
	if err = g.Wait(); err != nil {
		log.Erro.Printf("can't get sensor data: %s", err.Error())

		return
	}

	s.currT, s.currH = t, h

	s.metrics.Gauge(temperatureKey, s.currT)
	s.metrics.Gauge(humidityKey, s.currH)
}

func (s *Server) pushDataToHK() {
	s.mu.RLock()
	defer s.mu.RUnlock()

	s.hkSrv.SetCurrentTemperature(s.currT)
	s.hkSrv.SetCurrentHumidity(s.currH)
}

func (s *Server) runWebServer() error {
	if s.webSrv != nil {
		return errors.New("web server already exist")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, _ *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()

		_, _ = fmt.Fprintf(
			w,
			"Temp %v *C\nHumi %0.2f percent\n\n%s",
			s.currT, s.currH, renderHourlyAverageTable(
				s.metrics.GetForPeriodByH(temperatureKey, 24*time.Hour),
				s.metrics.GetForPeriodByH(humidityKey, 24*time.Hour),
			),
		)
	})

	s.webSrv = &http.Server{
		Addr:              ":80",
		Handler:           mux,
		ReadHeaderTimeout: 1 * time.Second,
	}

	return s.webSrv.ListenAndServe()
}

func (s *Server) runHapServer(ctx context.Context) error {
	return s.hkSrv.ListenAndServe(ctx)
}

func renderHourlyAverageTable(hourlyAverageT, hourlyAverageH map[string]float64) string {
	var builder strings.Builder

	builder.WriteString("+------+----------------+----------------+\n")
	builder.WriteString("| Hour |        T       |        H       |\n")
	builder.WriteString("+------+----------------+----------------+\n")

	// HH: { tt.t hh.h }
	// 01: { 23.5 60.0 }
	allData := make(map[string][]float64)
	for k, v := range hourlyAverageT {
		allData[k] = []float64{v, hourlyAverageH[k]}
	}
	allKeys := make([]string, 0, len(allData))
	for k := range allData {
		allKeys = append(allKeys, k)
	}

	sort.Strings(allKeys)

	for _, hour := range allKeys {
		val := allData[hour]
		builder.WriteString(fmt.Sprintf("| %-4s | %14.2f | %14.2f |\n", hour, val[0], val[1]))
	}

	builder.WriteString("+------+----------------+----------------+\n")

	return builder.String()
}
