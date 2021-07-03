package netdatadigest

import (
	"fmt"
	"net/http"

	"github.com/car12o/netdata-digest/pkg/logger"
)

type App struct {
	Logger logger.Service
}

func NewApp() *App {
	return &App{
		Logger: logger.New(),
	}
}

func (app *App) Serve(host string, port uint64) error {
	addr := fmt.Sprintf("%s:%d", host, port)
	app.Logger.Info(fmt.Sprintf("Serving netdata-digest at %s", addr))
	return http.ListenAndServe(addr, nil)
}
