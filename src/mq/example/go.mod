module mq-example

go 1.18

replace mq => ./../../mq

require mq v0.0.0-00010101000000-000000000000

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
)
