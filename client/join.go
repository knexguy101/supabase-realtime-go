package client

import (
	"context"
	"fmt"
	"github.com/knexguy101/supabase-realtime-go/types"
	"nhooyr.io/websocket/wsjson"
)

func (x *RealtimeClient) sendJoin() error {
	msg := types.SubscribeMsg{
		Event: JOIN_EVENT,
		Topic: "realtime:postgres",
		Ref:   "",
		Payload: types.SubscribePayload{
			Config: types.SubscribeConfig{
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
