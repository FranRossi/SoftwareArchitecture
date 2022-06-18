module voter_api

go 1.18

require (
	auth v0.0.0-00010101000000-000000000000
	electoral_service v0.0.0-00010101000000-000000000000
	encrypt v0.0.0-00010101000000-000000000000
	go.mongodb.org/mongo-driver v1.9.1
	google.golang.org/grpc v1.46.2
	google.golang.org/protobuf v1.28.0
	message_queue v0.0.0-00010101000000-000000000000
	pipes_and_filters v0.0.0-00010101000000-000000000000
)

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20211112202133-69e39bad7dc2 // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20210615035016-665e8c7367d1 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace auth => ./../auth/

replace pipes_and_filters => ./../pipes_and_filters

replace electoral_service => ./../electoral_service/

replace message_queue => ./../message_queue

replace encrypt => ./../encrypt
replace own_logger => ./../own_logger
