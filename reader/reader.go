package reader

import (
	"context"
	"encoding/json"
	"github.com/knexguy101/supabase-realtime-go/client"
	supabase_realtime_go "github.com/knexguy101/supabase-realtime-go/types"
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

		OnInsert: func(tableName string, message interface{}) {},
		OnDelete: func(oldRecords map[string]interface{}) {},
	}
}

func (x *Reader) Start() {
	x.disposed = false

	x.ctx, x.cancel = context.WithCancel(context.Background())

	go x.listen()
}

func (x *Reader) Dispose() {

	if x.ctx != nil {
		x.cancel()
	}

	x.disposed = true
}

func (x *Reader) listen() {
	for !x.disposed {
		var updateData supabase_realtime_go.PostgresMessage[interface{}]
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

					temp := reflect.New(schema).Interface()

					err = json.Unmarshal(b, &temp)
					if err != nil {
						continue
					}

					x.OnInsert(updateData.Payload.Data.Table, reflect.ValueOf(temp).Elem().Interface())
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
