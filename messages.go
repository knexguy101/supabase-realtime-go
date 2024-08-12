package supabase_realtime_go

type TemplateMsg struct {
	Event string `json:"event"`
	Topic string `json:"topic"`
	Ref   string `json:"ref"`
}

type ConnectionMsg struct {
	TemplateMsg

	Payload struct {
		Data struct {
			Schema     string            `json:"schema"`
			Table      string            `json:"table"`
			CommitTime string            `json:"commit_timestamp"`
			EventType  string            `json:"eventType"`
			New        map[string]string `json:"new"`
			Old        map[string]string `json:"old"`
			Errors     string            `json:"errors"`
		} `json:"data"`
	} `json:"payload"`
}

type PostgresMessage[RecordObj interface{}] struct {
	Event   string `json:"event"`
	Topic   string `json:"topic"`
	Payload struct {
		Data struct {
			Table     string                 `json:"table"`
			Type      string                 `json:"type"`
			Record    RecordObj              `json:"record"`
			OldRecord map[string]interface{} `json:"old_record"`
		} `json:"data"`
	} `json:"payload"`
}

type SubscribeMsg struct {
	Event   string           `json:"event"`
	Topic   string           `json:"topic"`
	Ref     string           `json:"ref"`
	Payload SubscribePayload `json:"payload"`
}

type SubscribePayload struct {
	Config SubscribeConfig `json:"config"`
}

type SubscribeConfig struct {
	PostgresChanges []PostgresChanges `json:"postgres_changes"`
}

type PostgresChanges struct {
	Event  string `json:"event"`
	Schema string `json:"schema,omitempty"`
	Table  string `json:"table"`
	Filter string `json:"filter,omitempty"`
}

type HearbeatMsg struct {
	TemplateMsg

	Payload struct {
	} `json:"payload"`
}
