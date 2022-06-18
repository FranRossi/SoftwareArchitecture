module notification_center

go 1.18

require (
	github.com/google/uuid v1.3.0
	github.com/streadway/amqp v1.0.0
	message_queue v0.0.0-00010101000000-000000000000
	own_logger v0.0.0-00010101000000-000000000000
)

replace message_queue => ./../message_queue

replace own_logger => ./../own_logger
