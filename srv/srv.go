package srv

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/brutella/hap"
)

const sensorPullingSleep = 5

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

type Server struct {
	webSrv  *http.Server
	hkSrv   HapServer
	climate ClimateSensor
	store   Store

	mu           *sync.Mutex
	currT, currH float64
}

func New(store Store, climate ClimateSensor, hapSrv HapServer) *Server {
	//mux := http.NewServeMux()
	//mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	//	t, h, err := getTempAndHumi(*sensor)
	//	if err != nil {
	//		fmt.Fprintf(w, fmt.Sprintf("ERR: %s", err.Error()))
	//
	//		return
	//	}
	//
	//	msg := fmt.Sprintf("Temp %v *C\nHumi %0.2f percent", t, h)
	//	fmt.Println(msg)
	//
	//	fmt.Fprintf(w, msg)
	//})

	return &Server{
		// TODO: do i need it here?
		webSrv:  nil,
		hkSrv:   hapSrv,
		climate: climate,
		store:   store,
	}
}

func (s *Server) Run(ctx context.Context) error {
	// go sensor pulling
	go func() {
		for {
			s.pullSensorData()
			<-time.After(sensorPullingSleep * time.Second)
		}
	}()

	g, ctx := errgroup.WithContext(ctx)

	// go web server
	g.Go(func() error {
		fmt.Println("web server listen on http://localhost:80")
		return s.webSrv.ListenAndServe()
	})
	// go hap server
	g.Go(func() error {
		fmt.Println("HAP server is up")
		return s.hkSrv.ListenAndServe(ctx)
	})

	return g.Wait()
}

func (s *Server) pullSensorData() {
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
		// TODO: replace it with a logger
		fmt.Printf("get sensor err: %s", err.Error())

		return
	}

	s.currT = t
	s.currH = h
}

func (s *Server) runWebServer() error {
	if s.webSrv != nil {
		return errors.New("web server already exist")
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = fmt.Fprintf(w, fmt.Sprintf("Temp %v *C\nHumi %0.2f percent", s.currT, s.currH))
	})

	s.webSrv = &http.Server{
		Addr:    ":80",
		Handler: mux,
	}

	return s.webSrv.ListenAndServe()
}

func (s *Server) runHapServer(ctx context.Context) error {
	return s.hkSrv.ListenAndServe(ctx)
}
