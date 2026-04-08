# Notification Stream Engine

##  Overview

Notification Stream Engine is a scalable, event-driven notification system built in Go, designed to process and deliver notifications reliably using asynchronous messaging patterns.

The system decouples event producers from notification delivery using a message broker, enabling high throughput, fault tolerance, and extensibility for multiple event types and delivery channels.

---

##  Problem Statement

In distributed systems, sending notifications synchronously leads to increased latency, tight coupling, and poor fault tolerance.

This project solves that by:

* Offloading notification processing to asynchronous workers
* Ensuring reliable delivery with retries and Dead Letter Queues (DLQ)
* Preventing duplicate notifications using idempotency mechanisms

---

##  Architecture

```
+-------------------+       +-------------------+       +----------------------+
|  Producer Service | ----> | Message Broker    | ----> | Notification Worker  |
|  (Events)         |       | (Queue / Stream)  |       | (Consumers)          |
+-------------------+       +-------------------+       +----------+-----------+
                                                              |
                                                              v
                                                   +----------------------+
                                                   | Event Handlers       |
                                                   | (Modular Routing)    |
                                                   +----------+-----------+
                                                              |
                                                              v
                                                   +----------------------+
                                                   | Notification Channel |
                                                   | (Email, etc.)        |
                                                   +----------------------+
```

---

##  Features

* **Event-Driven Architecture**

  * Decoupled producers and consumers using message queues

* **Reliable Message Processing**

  * Retry mechanism with delayed retries (TTL-based queues)
  * Dead Letter Queue (DLQ) for failed messages

* **Idempotency Handling**

  * Redis-based deduplication to ensure at-least-once delivery without duplicates

* **Modular Event Handlers**

  * Pluggable handler architecture for different event types

* **Scalable Design**

  * Horizontal scaling of consumers supported

---

##  Tech Stack

* **Language:** Golang
* **Framework:** Gin (for APIs, if applicable)
* **Messaging:** RabbitMQ / Kafka (based on implementation)
* **Cache:** Redis (for idempotency)
* **Architecture:** Event-Driven Microservices

---

##  Event Flow

1. Producer emits an event (e.g., `UserRegistered`)
2. Event is pushed to the message broker
3. Consumer picks up the message
4. Event is routed to the appropriate handler
5. Notification is sent (e.g., email)
6. On failure:

   * Retry with delay
   * Move to DLQ after max retries

---

##  Supported Event Types

* User Registration
* KYC Approval
* Transaction Events
* (Extensible for more)

---

##  Reliability Mechanisms

| Mechanism       | Purpose                             |
| --------------- | ----------------------------------- |
| Retry Queue     | Handles transient failures          |
| DLQ             | Stores failed messages for analysis |
| Idempotency Key | Prevents duplicate processing       |

---

##  Getting Started

### Prerequisites

* Go installed
* RabbitMQ / Kafka running
* Redis running

### Run the Service

```bash
git clone https://github.com/Mobi07/notification-stream-engine.git
cd notification-stream-engine

go mod tidy
go run main.go
```

---

##  Future Improvements

* Multi-channel notifications (SMS, Push, Slack)
* Observability (Prometheus + Grafana)
* Circuit breaker for external services
* Rate limiting & prioritization
* Dashboard for monitoring DLQ events

---

##  Use Cases

* Fintech transaction alerts
* User onboarding notifications
* KYC / compliance updates
* System event alerts

---

##  Author

Developed as a backend-focused distributed systems project to demonstrate scalable notification infrastructure using event-driven design principles.
