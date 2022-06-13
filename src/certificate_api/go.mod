module certificate_api

go 1.18

require (
	github.com/gofiber/fiber/v2 v2.34.0
	go.mongodb.org/mongo-driver v1.9.1
)

require (
	github.com/google/uuid v1.3.0 // indirect
	github.com/streadway/amqp v1.0.0 // indirect
)

require (
	electoral_service v0.0.0-00010101000000-000000000000
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/go-stack/stack v1.8.1 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.15.6 // indirect
	github.com/pkg/errors v0.9.1 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.37.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.1 // indirect
	github.com/xdg-go/stringprep v1.0.3 // indirect
	github.com/youmark/pkcs8 v0.0.0-20201027041543-1326539a0a0a // indirect
	golang.org/x/crypto v0.0.0-20220525230936-793ad666bf5e // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220227234510-4e6760a101f9 // indirect
	golang.org/x/text v0.3.7 // indirect
	message_queue v0.0.0-00010101000000-000000000000
)

replace message_queue => ./../message_queue

replace electoral_service => ./../electoral_service/
