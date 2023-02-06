package main

import (
	"flag"
	"log"
	"net/http"
	"strconv"
	"strings"
)

var addr = flag.String("addr", ":8000", "the server port")

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
		setHeaderStatus(w, subs[2:])

		return
	}
}

func setHeaderStatus(w http.ResponseWriter, subs []string) {
	if len(subs) != 1 {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	code, err := strconv.Atoi(subs[0])
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)

		return
	}
	w.WriteHeader(code)
}
