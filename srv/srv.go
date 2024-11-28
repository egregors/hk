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
	"github.com/egregors/hk/internal/metrics"
	"golang.org/x/sync/errgroup"

	"github.com/egregors/hk/log"
	"github.com/egregors/hk/utils/bp"
)

const (
	pullPushSleep = 60 * time.Second

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
		log.Info.Printf("start syncing sensor data with %s sleep", pullPushSleep)
		for {
			s.pullDataFromSensor()
			s.pushDataToHK()
			<-time.After(pullPushSleep)
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

		temp := s.metrics.Avg(temperatureKey, 24*time.Hour*3)
		humi := s.metrics.Avg(humidityKey, 24*time.Hour*3)

		_, _ = fmt.Fprintf(
			w,
			"Temp %0.2f °C\nHumi %0.2f %%\n\n%s\n\n%s\n\n",
			s.currT, s.currH,
			renderHourlyAvgVisualisation(temp, humi),
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
	builder.WriteString("+----------------+---------+--------+\n")
	builder.WriteString("|    Datetime    |    T    |    H   |\n")
	builder.WriteString("+----------------+---------+--------+\n")

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
	if len(allKeys) == 0 {
		// show "nothing to show"
		builder.WriteString("|      -         |    -    |    -   |\n")

		return builder.String()
	}

	sort.Slice(allKeys, func(i, j int) bool {
		return allKeys[i] > allKeys[j]
	})

	// HH: { tt.t hh.h }
	// 01: { 23.5 60.0 }
	for i := 0; i < len(allKeys); i++ {
		hour := allKeys[i]

		val := merge[hour]

		// hour: 2024-11-06 15:00:00 +0100 CET
		split := strings.Split(hour, " ")
		timeMark := strings.Join([]string{
			split[0],
			fmt.Sprint(strings.Split(split[1], ":")[0] + "h"),
		}, " ")
		// TODO: put tempProgressionMark back, instead of ""
		builder.WriteString(fmt.Sprintf("| %-14s | %1s%6.2f | %6.2f |\n", timeMark, "", val[0], val[1]))
	}
	//                      | 2024-11-08 18h | ~ 34.93 |  54.58 |
	builder.WriteString("+----------------+---------+--------+\n")

	return builder.String()
}

func renderHourlyAvgVisualisation(hourlyAverageT, hourlyAverageH []metrics.Value) string {
	tData := make([]float64, 0, len(hourlyAverageT))
	for _, v := range hourlyAverageT {
		tData = append(tData, v.V)
	}
	hData := make([]float64, 0, len(hourlyAverageH))
	for _, v := range hourlyAverageH {
		hData = append(hData, v.V)
	}

	return bp.SimplePlot(6, tData) + "\n\n" + bp.SimplePlot(6, hData)
}
