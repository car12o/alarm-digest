package test

import (
	"flag"
	"math/rand"
	"os"
	"testing"
	"time"

	"github.com/car12o/netdata-digest/internal/alarm"
	"github.com/car12o/netdata-digest/internal/broker"
	"github.com/car12o/netdata-digest/internal/messenger"
	"github.com/car12o/netdata-digest/pkg/logger"
	"github.com/nats-io/nats.go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type config struct {
	natsUrl string
	nc      *nats.EncodedConn
}

var cfg config

func init() {
	flag.StringVar(&cfg.natsUrl, "nats-url", nats.DefaultURL, "Nats URL")
}

func setup() error {
	var err error
	flag.Parse()
	cfg.nc, err = broker.NewNats(cfg.natsUrl, &logger.MockService{})
	return err
}

func shutdown() {
	cfg.nc.Close()
}

func TestMain(m *testing.M) {
	if err := setup(); err != nil {
		panic(err)
	}
	code := m.Run()
	shutdown()
	os.Exit(code)
}

func TestOrderDelivery(t *testing.T) {
	ch := make(chan *messenger.AlarmDigest)
	ns, err := cfg.nc.BindRecvChan("AlarmDigest", ch)
	require.Nil(t, err)
	defer ns.Unsubscribe()

	now := time.Now()
	alarms := []*messenger.AlarmStatusChanged{
		{
			AlarmID:   "alarmID-1",
			UserID:    "userID-1",
			Status:    alarm.StatusWarning,
			ChangedAt: now,
		},
		{
			AlarmID:   "alarmID-2",
			UserID:    "userID-1",
			Status:    alarm.StatusWarning,
			ChangedAt: now.Add(time.Millisecond * 15),
		},
		{
			AlarmID:   "alarmID-2",
			UserID:    "userID-1",
			Status:    alarm.StatusCritical,
			ChangedAt: now.Add(time.Millisecond * 20),
		},
		{
			AlarmID:   "alarmID-1",
			UserID:    "userID-2",
			Status:    alarm.StatusCritical,
			ChangedAt: now.Add(time.Millisecond * 22),
		},
		{
			AlarmID:   "alarmID-1",
			UserID:    "userID-1",
			Status:    alarm.StatusCritical,
			ChangedAt: now.Add(time.Millisecond * 25),
		},
		{
			AlarmID:   "alarmID-1",
			UserID:    "userID-1",
			Status:    alarm.StatusWarning,
			ChangedAt: now.Add(time.Millisecond * 30),
		},
		{
			AlarmID:   "alarmID-1",
			UserID:    "userID-2",
			Status:    alarm.StatusWarning,
			ChangedAt: now.Add(time.Millisecond * 32),
		},
		{
			AlarmID:   "alarmID-2",
			UserID:    "userID-1",
			Status:    alarm.StatusCleared,
			ChangedAt: now.Add(time.Millisecond * 35),
		},
	}

	expected := []messenger.AlarmStatusChanged{
		*alarms[2],
		*alarms[5],
		*alarms[6],
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(alarms), func(i, j int) { alarms[i], alarms[j] = alarms[j], alarms[i] })
	for _, alarm := range alarms {
		err = cfg.nc.Publish("AlarmStatusChanged", alarm)
		require.Nil(t, err)
		time.Sleep(time.Millisecond * 1)
	}

	err = cfg.nc.Publish("SendAlarmDigest", &messenger.SendAlarmDigest{UserID: "userID-1"})
	require.Nil(t, err)

	msg := <-ch
	assert.Equal(t, msg.UserID, "userID-1")
	assert.Len(t, msg.ActiveAlarms, 2)
	assert.Equal(t, expected[0].AlarmID, msg.ActiveAlarms[0].AlarmID)
	assert.Equal(t, expected[0].Status, msg.ActiveAlarms[0].Status)
	assert.True(t, expected[0].ChangedAt.Equal(msg.ActiveAlarms[0].LatestChangedAt))
	assert.Equal(t, expected[1].AlarmID, msg.ActiveAlarms[1].AlarmID)
	assert.Equal(t, expected[1].Status, msg.ActiveAlarms[1].Status)
	assert.True(t, expected[1].ChangedAt.Equal(msg.ActiveAlarms[1].LatestChangedAt))

	err = cfg.nc.Publish("SendAlarmDigest", &messenger.SendAlarmDigest{UserID: "userID-2"})
	require.Nil(t, err)

	msg = <-ch
	assert.Equal(t, msg.UserID, "userID-2")
	assert.Len(t, msg.ActiveAlarms, 1)
	assert.Equal(t, expected[2].AlarmID, msg.ActiveAlarms[0].AlarmID)
	assert.Equal(t, expected[2].Status, msg.ActiveAlarms[0].Status)
	assert.True(t, expected[2].ChangedAt.Equal(msg.ActiveAlarms[0].LatestChangedAt))
}

func TestNoAlarmsToDigest(t *testing.T) {
	ch := make(chan *messenger.AlarmDigest)
	ns, err := cfg.nc.BindRecvChan("AlarmDigest", ch)
	require.Nil(t, err)
	defer ns.Unsubscribe()

	err = cfg.nc.Publish("SendAlarmDigest", &messenger.SendAlarmDigest{UserID: "userID"})
	require.Nil(t, err)

	select {
	case msg := <-ch:
		assert.Fail(t, "channel should be empty", msg)
	default:
	}
}
