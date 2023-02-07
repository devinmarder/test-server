package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var addr = flag.String("port", ":8000", "the server port")

func main() {
	flag.Parse()
	log.Printf("starting test server on port: %s", *addr)
	log.Fatal(http.ListenAndServe(*addr, &server{}))
}

type server struct {
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := r.URL.Path
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
		log.Printf("received request for code: %s", code)
		setHeaderStatus(w, code)
		b, _ := json.Marshal(response{Code: subs[2]})
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
