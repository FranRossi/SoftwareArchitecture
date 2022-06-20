module pipes_and_filters

go 1.18

require (
	gopkg.in/yaml.v2 v2.4.0
	own_logger v0.0.0-00010101000000-000000000000
)

replace own_logger => ./../own_logger
