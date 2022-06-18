module voting_simulator

go 1.18



require (
	encrypt v0.0.0-00010101000000-000000000000
	google.golang.org/grpc v1.47.0
	google.golang.org/protobuf v1.28.0
)

require (
	electoral_service v0.0.0-00010101000000-000000000000 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	golang.org/x/net v0.0.0-20201021035429-f5854403a974 // indirect
	golang.org/x/sys v0.0.0-20210119212857-b64e53b001e4 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
)


replace voter_api => ./../voter_api
replace encrypt => ./../encrypt
replace electoral_service => ./../electoral_service // TODO capaz pasar a un paquete /models