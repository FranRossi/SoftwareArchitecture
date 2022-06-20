module stats_service

go 1.18

replace own_logger => ./../own_logger

replace encrypt => ./../encrypt

replace pipes_and_filters => ./../pipes_and_filters

replace message_queue => ./../message_queue

require message_queue v0.0.0-00010101000000-000000000000

require gopkg.in/yaml.v2 v2.4.0 // indirect

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	own_logger v0.0.0-00010101000000-000000000000
	pipes_and_filters v0.0.0-00010101000000-000000000000
)
