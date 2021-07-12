package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/car12o/alarm-digest/internal/alarmdigest"
	"github.com/nats-io/nats.go"
)

type config struct {
	host      string
	port      uint64
	redisAddr string
	natsUrl   string
}

var cfg config

func init() {
	flag.StringVar(&cfg.host, "host", "0.0.0.0", "Server host")
	flag.Uint64Var(&cfg.port, "port", 3000, "Server port")
	flag.StringVar(&cfg.redisAddr, "redis-addr", "0.0.0.0:6379", "Redis address")
	flag.StringVar(&cfg.natsUrl, "nats-url", nats.DefaultURL, "Nats URL")
	flag.Parse()
}

func main() {
	app, err := alarmdigest.NewApp(cfg.redisAddr, cfg.natsUrl)
	if err != nil {
		panic(err)
	}
	defer app.Close()

	addr := fmt.Sprintf("%s:%d", cfg.host, cfg.port)
	app.Log.Info(fmt.Sprintf("Serving alarm-digest at %s", addr))
	if err := http.ListenAndServe(addr, nil); err != nil {
		panic(err)
	}
}
