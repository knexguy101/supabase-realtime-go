# Supabase Realtime Golang
This is an unofficial wrapper for realtime functionality in golang. Please feel free to fork and customize.

## Auth
```md
ProjectRef = https://[PROJECT REF HERE].supabase.co
ApiKey = eyJ...
```

## Quickstart
```go
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
	panic(err)
}
defer c.Disconnect()
```

## Layout and Events
This package calls functions as events, making handling a breeze. It also allows for raw
websocket usage if desired.

#### Reader
```go
//postgres change events
OnInsert func(tableName string, message interface{})
OnDelete func(oldRecords map[string]interface{})
```
#### Client
```go
// connection events
OnDial       func(client *RealtimeClient)
OnDisconnect func(client *RealtimeClient)

// heartbeat events
OnHeartbeatFailed       func(isConAlive bool, err error, client *RealtimeClient)
OnHeartbeatRedialFailed func(err error, client *RealtimeClient)
OnHeartbeatSent         func(client *RealtimeClient)

// websocket connection
func (x *RealtimeClient) GetWS() *websocket.Conn
```

## Client Events
When creating a `Client`, you are able to specify which events from what tables you would like the `Client` to listen for.
You can specify 1 or more `PostgresChanges` when initializing the `Client`
```go
// this example reads all events from the table "test"
c := client.CreateRealtimeClient(os.Getenv("PROJECT_REF"), os.Getenv("API_KEY"), []types.PostgresChanges{
    {
        Event: "*",
        Table: "test",
    },
})
```

## Reader and Parser
Don't want to deal with reading and parsing the postgres events? Try using a `Reader`
<br />
A `Reader` takes in a `SchemaBuilder` and a `Client`, in return it will automatically dispatch events. Message data
is automatically parsed allowing you to cast the received message to the appropriate type.
```go
type testData struct {
	ID        int8      `json:"id"`
	CreatedAt time.Time `json:"created_at"`
}

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
```

## Schemas
Schemas are easy with the `SchemaBuilder`, any struct works.

``Note: All records come through as a JSON format.``
```go
type nameOfMySchema struct {
	ID int `json:"id"`
}

type nameOfSecondSchema struct {
    ID int `json:"id"`
}

s := reader.NewSchemaBuilder()
s.Set("TABLE NAME HERE", nameOfMySchema{})
s.Set("ANOTHER TABLE NAME", nameOfSecondSchema{})
```