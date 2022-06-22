module pipes_and_filters_example

go 1.18

replace pipes_and_filters => ./../../pipes_and_filters

replace own_logger => ./../../own_logger

require pipes_and_filters v0.0.0-00010101000000-000000000000

require (
	gopkg.in/yaml.v2 v2.4.0 // indirect
	own_logger v0.0.0-00010101000000-000000000000 // indirect
)
