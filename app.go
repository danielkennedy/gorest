package main

import (
	"code.google.com/p/log4go"
	"fmt"
	"net/http"
	"os"
)

const (
	HostVar = "VCAP_APP_HOST"
	PortVar = "VCAP_APP_PORT"
)

func main() {
	log := make(log4go.Logger)
	log.AddFilter("stdout", log4go.DEBUG, log4go.NewConsoleLogWriter())

	http.HandleFunc("/", handler)
	var port string
	if port = os.Getenv(PortVar); port == "" {
		port = "8080"
	}
	log.Debug("Listening at port %v\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		panic(err)
	}
}

func handler(res http.ResponseWriter, req *http.Request) {
	// Dump ENV
	fmt.Fprint(res, "ENV:\n")
	env := os.Environ()
	for _, e := range env {
		fmt.Fprintln(res, e)
	}
  fmt.Fprintf(res, "Serving request for %s", req.URL.Path[1:])
}

