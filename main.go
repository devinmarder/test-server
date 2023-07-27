package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

var addr = flag.String("port", ":8000", "the server port")

func main() {
	flag.Parse()
	log.Printf("starting test server on port: %s", *addr)
	s := &server{metrics: metrics{metrics: make(map[string][]metric)}}
	go s.writeMetrics()
	log.Fatal(http.ListenAndServe(*addr, s))
}

type server struct {
	metrics metrics
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
	s.metrics.Update(path)
	subs := strings.Split(path, "/")
	if len(subs) == 0 {
		return
	}
	if len(subs) == 1 {
		return
	}
	switch subs[1] {
	case "codes":
		code := subs[2]
		setHeaderStatus(w, code)
		b, _ := json.Marshal(response{Code: subs[2]})
		w.Write(b)
		return
	case "delay":
		delay := subs[2]
		d, err := strconv.Atoi(delay)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		time.Sleep(time.Duration(d) * time.Millisecond)
		b, _ := json.Marshal(response{Code: "200"})
		w.Write(b)
		return
	}
}

func setHeaderStatus(w http.ResponseWriter, c string) {
	if c == "" {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	code, err := strconv.Atoi(c)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	w.WriteHeader(code)
}

type response struct {
	Code string
}

type metrics struct {
	mx      sync.Mutex
	metrics map[string][]metric
}

func (m *metrics) Update(path string) {
	m.mx.Lock()
	defer m.mx.Unlock()
	if _, ok := m.metrics[path]; !ok {
		m.metrics[path] = []metric{&totalMetric{}, NewRateMetric(1 * time.Second)}
	}
	for _, v := range m.metrics[path] {
		v.inc()
	}
}

func (s *server) writeMetrics() {
	w := io.Writer(os.Stdout)
	for {
		lines := s.metrics.displayMetrics(w)
		time.Sleep(1 * time.Second)
		_, _ = fmt.Fprint(w, strings.Repeat(clear, lines))
	}
}

func (m *metrics) displayMetrics(w io.Writer) int {
	m.mx.Lock()
	defer m.mx.Unlock()
	lines := 0
	keys := make([]string, 0, len(m.metrics))
	for k, _ := range m.metrics {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		fmt.Fprintf(w, "%s\n", k)
		lines++
		for _, m := range m.metrics[k] {
			fmt.Fprintf(w, "\t-%s\n", m.string())
			lines++
		}
	}
	return lines
}

type metric interface {
	inc()
	string() string
}

type totalMetric struct {
	total atomic.Int64
}

func (t *totalMetric) inc() {
	t.total.Add(1)
}

func (t *totalMetric) string() string {
	return "total:\t" + strconv.FormatInt(t.total.Load(), 10) + " reqs"
}

type rateMetric struct {
	period time.Duration
	count  atomic.Int64
	rate   atomic.Int64
}

func NewRateMetric(period time.Duration) *rateMetric {
	r := &rateMetric{period: period}
	r.Run()
	return r
}

func (r *rateMetric) Run() {
	go func() {
		t := time.NewTicker(r.period)
		for range t.C {
			r.rate.Store(r.count.Load() / int64(r.period.Seconds()))
			r.count.Store(0)
		}
	}()
}

func (r *rateMetric) inc() {
	r.count.Add(1)
}

func (r *rateMetric) string() string {
	return "rate:\t" + strconv.FormatInt(r.rate.Load(), 10) + " reqs/s"
}

const ESC = 27

var clear = fmt.Sprintf("%c[%dA%c[2K", ESC, 1, ESC)
