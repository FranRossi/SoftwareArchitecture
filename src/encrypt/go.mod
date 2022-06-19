module encrypt

go 1.18

replace own_logger => ./../own_logger

replace auth => ./../auth/

replace electoral_service => ./../electoral_service/

replace voter_api => ./../voter_api

replace pipes_and_filters => ./../pipes_and_filters

replace message_queue => ./../message_queue

require electoral_service v0.0.0-00010101000000-000000000000
