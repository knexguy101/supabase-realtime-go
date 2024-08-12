package client

import (
	"context"
	"fmt"
	"github.com/knexguy101/supabase-realtime-go/types"
	"nhooyr.io/websocket"
	"time"
)

type RealtimeClient struct {
	realtimeUrl string
	apiKey      string
	events      []types.PostgresChanges

	conn            *websocket.Conn
	statusCtx       context.Context
	statusCtxCancel func()

	dialTimeout       time.Duration
	joinDuration      time.Duration
	heartbeatDuration time.Duration

	OnDial       func(client *RealtimeClient)
	OnDisconnect func(client *RealtimeClient)

	OnHeartbeatFailed       func(isConAlive bool, err error, client *RealtimeClient)
	OnHeartbeatRedialFailed func(err error, client *RealtimeClient)
	OnHeartbeatSent         func(client *RealtimeClient)
}

func CreateRealtimeClient(projectRef string, apiKey string, events []types.PostgresChanges) *RealtimeClient {
	realtimeUrl := fmt.Sprintf(
		"wss://%s.supabase.co/realtime/v1/websocket?apikey=%s&log_level=info&vsn=1.0.0",
		projectRef,
		apiKey,
	)

	ctx, cancel := context.WithCancel(context.Background())

	return &RealtimeClient{
		realtimeUrl:       realtimeUrl,
		apiKey:            apiKey,
		events:            events,
		statusCtx:         ctx,
		statusCtxCancel:   cancel,
		dialTimeout:       10 * time.Second,
		heartbeatDuration: 5 * time.Second,
		joinDuration:      5 * time.Second,

		OnDial:       func(client *RealtimeClient) {},
		OnDisconnect: func(client *RealtimeClient) {},

		OnHeartbeatFailed:       func(isConAlive bool, err error, client *RealtimeClient) {},
		OnHeartbeatRedialFailed: func(err error, client *RealtimeClient) {},
		OnHeartbeatSent:         func(client *RealtimeClient) {},
	}
}

func (x *RealtimeClient) GetWS() *websocket.Conn {
	return x.conn
}

func (x *RealtimeClient) SetDelays(dialTimeout, heartbeatDuration, joinDuration time.Duration) {
	x.dialTimeout = dialTimeout
	x.heartbeatDuration = heartbeatDuration
	x.joinDuration = joinDuration
}
