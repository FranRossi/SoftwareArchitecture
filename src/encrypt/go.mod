module encrypt

go 1.18

replace message_queue => ./../message_queue

replace electoral_service => ./../electoral_service/

replace pipes_and_filters => ./../pipes_and_filters

replace own_logger => ./../own_logger


require electoral_service v0.0.0-00010101000000-000000000000
