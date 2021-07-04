package netdatadigest

import (
	"github.com/car12o/netdata-digest/internal/alarms"
	"github.com/car12o/netdata-digest/internal/messagebroker"
	"github.com/car12o/netdata-digest/pkg/logger"
)

type App struct {
	Logger        logger.Service
	Messagebroker messagebroker.Service
}

func NewApp(natsUrl string) (*App, error) {
	log := logger.New()

	log.Info("Connectiong to NATS server")
	mb, err := messagebroker.NewService(natsUrl, log)
	if err != nil {
		return nil, err
	}
	log.Info("Successfully connected to NATS server")

	if err := alarms.NewService(mb).SubscribeTopics(); err != nil {
		return nil, err
	}

	return &App{log, mb}, nil
}

func (app *App) Close() {
	app.Messagebroker.Close()
}
