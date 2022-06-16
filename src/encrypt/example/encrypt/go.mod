module encrypt_example

go 1.18

replace encrypt => ./../../../encrypt

replace electoral_service => ./../../../electoral_service/

require (
	electoral_service v0.0.0-00010101000000-000000000000
	encrypt v0.0.0-00010101000000-000000000000
)
