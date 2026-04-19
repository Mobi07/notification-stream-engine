package rabbitmq

import (
	"github.com/streadway/amqp"
)

func GetRetryCount(msg amqp.Delivery) int {
	if val, ok := msg.Headers["x-retry-count"]; ok {
		if count, ok := val.(int32); ok {
			return int(count)
		}
	}

	return 0
}

func PublishToDLQ(ch *amqp.Channel, dlxExchange, dlqName string, msg amqp.Delivery) error {
	return ch.Publish(
		dlxExchange,
		dlqName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         msg.Body,
			Headers:      msg.Headers,
			DeliveryMode: amqp.Persistent, // survive restart
		},
	)
}

func PublishToRetryQueue(ch *amqp.Channel, retryQueueName string, msg amqp.Delivery, retryCount int) error {
	headers := msg.Headers
	if headers == nil {
		headers = make(amqp.Table)
	}
	headers["x-retry-count"] = int32(retryCount)

	return ch.Publish(
		"",
		retryQueueName,
		false,
		false,
		amqp.Publishing{
			ContentType:  "application/json",
			Body:         msg.Body,
			Headers:      headers,
			DeliveryMode: amqp.Persistent, // survive restart
		},
	)
}
