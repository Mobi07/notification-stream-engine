package events

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}
