package client

import (
	"context"
	"fmt"
	"github.com/knexguy101/supabase-realtime-go/types"
	"nhooyr.io/websocket/wsjson"
	"time"
)

// Start sending heartbeats to the server to maintain connection
func (x *RealtimeClient) startHeartbeats() {
	ticker := time.NewTicker(x.heartbeatDuration)
	defer ticker.Stop()

	done := false
	for !done {
		select {
		case <-ticker.C:
			{
				err := x.sendHeartbeat()
				if err != nil {
					alive := x.isConnectionAlive(err)
					x.OnHeartbeatFailed(alive, err, x)
					if alive {
						continue
					}

					err = x.dial()
					if err != nil {
						x.OnHeartbeatRedialFailed(err, x)
					}
				} else {
					x.OnHeartbeatSent(x)
				}
			}
		case <-x.statusCtx.Done():
			{
				done = true
			}
		}
	}
}

func (x *RealtimeClient) sendHeartbeat() error {
	msg := types.HearbeatMsg{
		TemplateMsg: types.TemplateMsg{
			Event: HEARTBEAT_EVENT,
			Topic: "phoenix",
			Ref:   "",
		},
		Payload: struct{}{},
	}

	ctx, cancel := context.WithTimeout(context.Background(), x.heartbeatDuration)
	defer cancel()

	err := wsjson.Write(ctx, x.conn, msg)
	if err != nil {
		return fmt.Errorf("failed to send hearbeat in %f seconds: %w", x.heartbeatDuration.Seconds(), err)
	}

	return nil
}
