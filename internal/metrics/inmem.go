package metrics

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
	"sort"
	"time"

	"github.com/egregors/hk/log"
)

const (
	cleanerWorkerSleep = 30 * time.Second
)

type DumpFn func() error

type Option func(m *InMem)

func WithRetention(dur time.Duration) Option {
	return func(m *InMem) {
		m.retentionDuration = dur
	}
}

func WithBackup() Option {
	return func(m *InMem) {
		m.backup = true
	}
}

type Value struct {
	T time.Time
	V float64
}

type valueChanMsg struct {
	key string
	m   Value
}

type InMem struct {
	GaugeTimeLine map[string][]Value
	gaugeTLch     chan valueChanMsg

	backup            bool
	retentionDuration time.Duration
}

func New(opts ...Option) (m *InMem, commitDump DumpFn) {
	m = &InMem{
		GaugeTimeLine: make(map[string][]Value),
		gaugeTLch:     make(chan valueChanMsg),
	}

	for _, opt := range opts {
		opt(m)
	}

	go m.collector()
	go m.cleaner()

	if m.backup {
		log.Info.Println("try to restore from dump")
		err := m.Restore()
		if err != nil {
			log.Erro.Printf("not this time: %s", err.Error())
		} else {
			log.Info.Println("got from dump:")
			for k, v := range m.GaugeTimeLine {
				log.Info.Printf("-- %s: %d", k, len(v))
			}
		}
	}

	commitDump = func() error {
		log.Debg.Println("close gaugeTL channel")
		close(m.gaugeTLch)

		if m.backup {
			return m.Dump()
		}

		return nil
	}

	return m, commitDump
}

func (m *InMem) Gauge(key string, val float64) {
	log.Debg.Printf("send: gauge %s: %v", key, val)
	go func() {
		m.gaugeTLch <- valueChanMsg{
			key: key,
			m:   Value{T: time.Now(), V: val},
		}
	}()
}

func (m *InMem) Avg(key string, dur time.Duration) []Value {
	// FIXME: mu?
	data, ok := m.GaugeTimeLine[key]
	if !ok {
		return nil
	}

	// get data for duration
	start, end := time.Now().Add(-dur), time.Now()
	var durData []Value
	for _, val := range data {
		t := val.T
		if t.After(start) && t.Before(end) {
			durData = append(durData, val)
		}
	}

	// calc hourly avg
	sort.Slice(durData, func(i, j int) bool {
		return durData[i].T.Before(durData[j].T)
	})

	hAvg := make(map[time.Time][]float64)
	for _, v := range durData {
		t := v.T.Truncate(time.Hour)
		hAvg[t] = append(hAvg[t], v.V)
	}

	avg := make([]Value, 0, len(hAvg))
	for k, v := range hAvg {
		sum := 0.0
		for _, vv := range v {
			sum += vv
		}

		avg = append(avg, Value{T: k, V: sum / (float64(len(v)))})
	}

	return avg
}

func (m *InMem) collector() {
	log.Debg.Println("collector started")
	for msg := range m.gaugeTLch {
		log.Debg.Printf("got: gauge %s: %v at %v", msg.key, msg.m.V, msg.m.T)
		m.GaugeTimeLine[msg.key] = append(m.GaugeTimeLine[msg.key], msg.m)
	}
}

func (m *InMem) cleaner() {
	if m.retentionDuration == 0 {
		log.Info.Println("retention isn't setted up")

		return
	}

	for {
		<-time.After(cleanerWorkerSleep)
		log.Debg.Printf("cleanup. retention period: %v\n", m.retentionDuration)
		log.Debg.Println("current size:")
		for k, v := range m.GaugeTimeLine {
			log.Debg.Printf("-- %s: %d", k, len(v))
		}

		cutoff := time.Now().Add(-m.retentionDuration)
		var totalVs, totalNewVs int
		for k, v := range m.GaugeTimeLine {
			var newV []Value
			for _, vv := range v {
				if vv.T.After(cutoff) {
					newV = append(newV, vv)
				}
			}
			totalVs += len(v)
			totalNewVs += len(newV)
			m.GaugeTimeLine[k] = newV
		}

		diff := totalVs - totalNewVs
		if diff != 0 {
			log.Debg.Printf("cleaner removed %d gauges by retention policy\n", diff)
		}
	}
}

func (m *InMem) Dump() error {
	buf := new(bytes.Buffer)
	encoder := gob.NewEncoder(buf)

	err := encoder.Encode(m.GaugeTimeLine)
	if err != nil {
		return fmt.Errorf("can't encode items: %w", err)
	}
	err = os.WriteFile("hk-dump.gob", buf.Bytes(), 0o600)
	if err != nil {
		return fmt.Errorf("can't save dump: %w", err)
	}

	return nil
}

func (m *InMem) Restore() error {
	f, err := os.ReadFile("hk-dump.gob")
	if err != nil {
		return fmt.Errorf("can't read dump: %w", err)
	}

	buf := bytes.NewBuffer(f)
	decoder := gob.NewDecoder(buf)

	err = decoder.Decode(&m.GaugeTimeLine)
	if err != nil {
		return fmt.Errorf("can't decode items: %w", err)
	}

	return nil
}
