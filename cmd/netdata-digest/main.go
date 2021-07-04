package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/car12o/netdata-digest/internal/netdatadigest"
	"github.com/nats-io/nats.go"
)

type config struct {
	host    string
	port    uint64
	natsUrl string
}

var cfg config

func init() {
	flag.StringVar(&cfg.host, "host", "0.0.0.0", "Server host")
	flag.Uint64Var(&cfg.port, "port", 3000, "Server port")
	flag.StringVar(&cfg.natsUrl, "nats-url", nats.DefaultURL, "Nets URL")
	flag.Parse()
}

func main() {
	app, err := netdatadigest.NewApp(cfg.natsUrl)
	if err != nil {
		panic(err)
	}
	defer app.Close()

	addr := fmt.Sprintf("%s:%d", cfg.host, cfg.port)
	app.Logger.Info(fmt.Sprintf("Serving netdata-digest at %s", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
