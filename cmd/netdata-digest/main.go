package main

import (
	"flag"

	netdatadigest "github.com/car12o/netdata-digest/internal/netdata-digest"
)

type config struct {
	host string
	port uint64
}

var cfg config

func init() {
	flag.StringVar(&cfg.host, "host", "0.0.0.0", "Server host")
	flag.Uint64Var(&cfg.port, "port", 3000, "Server port")
	flag.Parse()
}

func main() {
	if err := netdatadigest.NewApp().Serve(
		cfg.host, cfg.port,
	); err != nil {
		panic(err)
	}
}
