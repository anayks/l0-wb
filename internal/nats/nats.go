package nats

import (
	"github.com/nats-io/stan.go"
)

const (
	clusterID = "test-cluster"
	clientID  = "event-store"
)

func NewNats() (*stan.Conn, error) {
	sc, err := stan.Connect(clusterID, clientID, stan.NatsURL("nats:4222"))
	if err != nil {
		return nil, err
	}

	return &sc, nil
}
