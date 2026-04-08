package events

type UserRegistration struct {
	UserID string `json:"user_id" mapstructure:"user_id"`
	Email  string `json:"email" mapstructure:"email"`
}

type TransactionCompleted struct {
	TransactionID string  `json:"transaction_id" mapstructure:"transaction_id"`
	UserID        string  `json:"user_id"  mapstructure:"user_id"`
	Amount        float64 `json:"amount" mapstructure:"amount"`
}

type PasswordResetPayload struct {
	UserID string `json:"user_id"  mapstructure:"user_id"`
	Email  string `json:"email" mapstructure:"email"`
}

type Event struct {
	ID        string      `json:"id"`
	Type      string      `json:"type"`
	Timestamp int64       `json:"timestamp"`
	Payload   interface{} `json:"payload"`
}
