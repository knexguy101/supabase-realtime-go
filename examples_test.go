package supabase_realtime_go

import (
	"fmt"
	"github.com/knexguy101/supabase-realtime-go/client"
	"github.com/knexguy101/supabase-realtime-go/reader"
	"github.com/knexguy101/supabase-realtime-go/types"
	"os"
	"testing"
	"time"
)

type testData struct {
	ID        int8      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

func TestClientAndReader(t *testing.T) {
	c := client.CreateRealtimeClient(os.Getenv("PROJECT_REF"), os.Getenv("API_KEY"), []types.PostgresChanges{
		{
			Event: "*",
			Table: "test",
		},
	})

	c.OnDial = func(client *client.RealtimeClient) {
		fmt.Println("client connected")
	}
	c.OnHeartbeatSent = func(client *client.RealtimeClient) {
		fmt.Println("heartbeat sent")
	}
	c.OnHeartbeatFailed = func(isConAlive bool, err error, client *client.RealtimeClient) {
		fmt.Printf("heartbeat failed, alive: %v\n", isConAlive)
	}

	err := c.Connect()
	if err != nil {
		t.Error(err)
	}
	defer c.Disconnect()

	s := reader.NewSchemaBuilder()
	s.Set("test", testData{})

	r := reader.NewReader(c, s)

	r.OnInsert = func(tableName string, message interface{}) {
		switch tableName {
		case "test":
			{
				m := message.(testData)
				fmt.Println("inserted item to test", m.CreatedAt.String())
			}
		}
	}
	r.OnDelete = func(oldRecords map[string]interface{}) {
		fmt.Println("item deleted", oldRecords)
	}

	r.Start()
	defer r.Dispose()

	ch := make(chan bool)
	<-ch
}
