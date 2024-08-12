package reader

import (
	"context"
	"encoding/json"
	supabase_realtime_go "github.com/knexguy101/supabase-realtime-go"
	"github.com/knexguy101/supabase-realtime-go/client"
	"nhooyr.io/websocket/wsjson"
	"reflect"
)

type Reader struct {
	client  *client.RealtimeClient
	schemas map[string]reflect.Type

	OnInsert func(tableName string, message interface{})
	OnDelete func(oldRecords map[string]interface{})

	disposed bool
	ctx      context.Context
	cancel   func()
}

func NewReader(c *client.RealtimeClient, s *SchemaBuilder) *Reader {
	return &Reader{
		client:  c,
		schemas: s.List,
	}
}

func (x *Reader) Start() {
	x.disposed = false

	x.ctx, x.cancel = context.WithCancel(context.Background())

	go x.listen()
}

func (x *Reader) Dispose() {
	x.disposed = true
}

func (x *Reader) listen() {
	for !x.disposed {
		var updateData supabase_realtime_go.PostgresMessage[any]
		err := wsjson.Read(x.ctx, x.client.GetWS(), &updateData)
		if err != nil {
			continue
		}

		schema, found := x.schemas[updateData.Payload.Data.Table]
		if !found {
			continue
		}

		if updateData.Event == "postgres_changes" {
			switch updateData.Payload.Data.Type {
			case "INSERT":
				{
					b, err := json.Marshal(updateData.Payload.Data.Record)
					if err != nil {
						continue
					}

					temp := reflect.Zero(schema).Interface()

					err = json.Unmarshal(b, &temp)
					if err != nil {
						continue
					}

					x.OnInsert(updateData.Payload.Data.Table, temp)
				}
			case "DELETE":
				{
					x.OnDelete(updateData.Payload.Data.OldRecord)
				}
			case "UPDATE":

			}
		}
	}
}
