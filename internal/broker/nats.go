package broker

import (
	"time"

	"github.com/car12o/alarm-digest/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/pkg/errors"
)

func NewNats(url string, log logger.Service) (*nats.EncodedConn, error) {
	nc, err := nats.Connect(
		url,
		nats.RetryOnFailedConnect(true),
		nats.MaxReconnects(15),
		nats.ReconnectWait(time.Second*2),
		nats.DisconnectErrHandler(func(_ *nats.Conn, err error) {
			log.Error(errors.Wrap(err, "Disconnected from NATS server"))
		}),
		nats.ReconnectHandler(func(_ *nats.Conn) {
			log.Info("Successfully reconnected to NATS server")
		}),
	)
	if err != nil {
		return nil, err
	}

	if err := nc.Flush(); err != nil {
		return nil, err
	}

	ec, err := nats.NewEncodedConn(nc, nats.JSON_ENCODER)
	if err != nil {
		return nil, err
	}

	return ec, nil
}
