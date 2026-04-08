package constants

const (
	MainQueueName  = "notification_queue"
	DLQName        = "notification_dlq"
	DLXExchange    = "dlx_exchange"
	RetryQueueName = "notification_retry"
	DLQRoutingKey  = "dlq_routing_key"
)

const MaxRetryCount = 3