module encrypt

go 1.18

replace electoral_service => ./../electoral_service/
replace own_logger => ./../own_logger
replace message_queue => ./../message_queue
replace pipes_and_filters => ./../pipes_and_filters


require (
	electoral_service v0.0.0-00010101000000-000000000000
	own_logger v0.0.0-00010101000000-000000000000
)
