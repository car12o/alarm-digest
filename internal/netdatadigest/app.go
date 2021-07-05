package netdatadigest

import (
	"github.com/car12o/netdata-digest/internal/alarms"
	"github.com/car12o/netdata-digest/internal/broker"
	"github.com/car12o/netdata-digest/internal/messenger"
	"github.com/car12o/netdata-digest/internal/storage"
	"github.com/car12o/netdata-digest/pkg/logger"
	"github.com/go-redis/redis/v8"
	"github.com/nats-io/nats.go"
)

type App struct {
	Logger    logger.Service
	RedisConn *redis.Client
	NatsConn  *nats.EncodedConn
}

func NewApp(redisAddr string, natsUrl string) (*App, error) {
	log := logger.New()

	log.Info("Connectiong to Redis server")
	rc, err := storage.NewRedis(redisAddr)
	if err != nil {
		return nil, err
	}
	log.Info("Successfully connected to Redis server")

	log.Info("Connectiong to NATS server")
	nc, err := broker.NewNats(natsUrl, log)
	if err != nil {
		return nil, err
	}
	log.Info("Successfully connected to NATS server")

	if err := alarms.NewService(
		// "",
		messenger.NewService(nc),
	).SubscribeTopics(); err != nil {
		return nil, err
	}

	return &App{log, rc, nc}, nil
}

func (app *App) Close() {
	app.RedisConn.Close()
	app.NatsConn.Close()
}
