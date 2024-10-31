package srv

import (
	"context"
	"errors"
	"fmt"
	"github.com/egregors/hk/internal/metrics"
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
	pullPushSleep = 1

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
	Avg(key string, dur time.Duration) []metrics.Value
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

		temp := s.metrics.Avg(temperatureKey, 24*time.Hour)
		humi := s.metrics.Avg(humidityKey, 24*time.Hour)

		_, _ = fmt.Fprintf(
			w,
			"Temp %v *C\nHumi %0.2f percent\n\n%s\n\n",
			s.currT, s.currH,
			renderHourlyAvgTable(temp, humi),
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

func renderHourlyAvgTable(hourlyAverageT, hourlyAverageH []metrics.Value) string {
	var builder strings.Builder
	builder.WriteString("+-------------------------------+----------------+----------------+\n")
	builder.WriteString("| Hour                          |        T       |        H       |\n")
	builder.WriteString("+-------------------------------+----------------+----------------+\n")

	merge := make(map[string][]float64)

	// collect temp
	for _, v := range hourlyAverageT {
		if _, ok := merge[v.T.String()]; !ok {
			merge[v.T.String()] = make([]float64, 2)
		}

		merge[v.T.String()][0] = v.V
	}

	// collect humi
	for _, v := range hourlyAverageH {
		if _, ok := merge[v.T.String()]; !ok {
			merge[v.T.String()] = make([]float64, 2)
		}

		merge[v.T.String()][1] = v.V
	}

	allKeys := make([]string, 0, len(merge))
	for k := range merge {
		allKeys = append(allKeys, k)
	}
	sort.Strings(allKeys)

	// HH: { tt.t hh.h }
	// 01: { 23.5 60.0 }
	up, down, same := "^", "v", "~"
	var prevT = merge[allKeys[0]][0]
	for _, hour := range allKeys {
		nowMark := ""
		progMark := same
		if hour == time.Now().Truncate(time.Hour).String() {
			nowMark = " <--"
		}
		val := merge[hour]
		if val[0] > prevT {
			progMark = up
		} else if val[0] < prevT {
			progMark = down
		} else {
			progMark = same
		}
		builder.WriteString(fmt.Sprintf("| %-4s | %s %12.2f | %14.2f |%s\n", hour, progMark, val[0], val[1], nowMark))
		prevT = val[0]
	}

	builder.WriteString("+-------------------------------+----------------+----------------+\n")

	return builder.String()
}
