package client

import (
	"context"
	"fmt"
	"github.com/knexguy101/supabase-realtime-go"
	"nhooyr.io/websocket/wsjson"
)

func (x *RealtimeClient) sendJoin() error {
	msg := supabase_realtime_go.SubscribeMsg{
		Event: JOIN_EVENT,
		Topic: "realtime:postgres",
		Ref:   "",
		Payload: supabase_realtime_go.SubscribePayload{
			Config: supabase_realtime_go.SubscribeConfig{
				PostgresChanges: x.events,
			},
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), x.joinDuration)
	defer cancel()

	err := wsjson.Write(ctx, x.conn, msg)
	if err != nil {
		return fmt.Errorf("failed to send hearbeat in %f seconds: %w", x.joinDuration.Seconds(), err)
	}

	return nil
}
