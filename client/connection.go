package client

import (
	"context"
	"errors"
	"fmt"
	"io"
	"nhooyr.io/websocket"
	"strings"
)

func (x *RealtimeClient) Connect() error {
	// Attempt to dial the server
	err := x.dial()
	if err != nil {
		return fmt.Errorf("cannot connect to the server: %w", err)
	}

	go x.startHeartbeats()

	return nil
}

func (x *RealtimeClient) Disconnect() error {
	x.statusCtxCancel()

	x.OnDisconnect(x)

	if x.conn != nil {
		return x.conn.Close(websocket.StatusNormalClosure, "Closing the connection")
	}

	return nil
}

func (x *RealtimeClient) dial() error {
	ctx, cancel := context.WithTimeout(context.Background(), x.dialTimeout)
	defer cancel()

	conn, _, err := websocket.Dial(ctx, x.realtimeUrl, nil)
	if err != nil {
		return err
	}

	x.OnDial(x)

	x.conn = conn

	return nil
}

func (x *RealtimeClient) isConnectionAlive(err error) bool {
	return !errors.Is(err, io.EOF) && !strings.Contains(err.Error(), "use of closed network connection")
}
