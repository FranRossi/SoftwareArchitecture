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
	own_logger v0.0.0-00010101000000-000000000000
	pipes_and_filters v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/go-redis/redis/v8 v8.11.5 // indirect
)

require (
	cache v0.0.0-00010101000000-000000000000
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/gofiber/fiber/v2 v2.34.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/google/uuid v1.3.0 // indirect
	github.com/joho/godotenv v1.4.0
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.37.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/net v0.0.0-20220225172249-27dd8689420f // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20200526211855-cb27e3aa2013 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
)

replace auth => ./../auth/

replace cache => ./../cache

replace pipes_and_filters => ./../pipes_and_filters

replace electoral_service => ./../electoral_service/

replace message_queue => ./../message_queue

replace encrypt => ./../encrypt

replace own_logger => ./../own_logger

