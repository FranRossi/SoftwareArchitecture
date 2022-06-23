module consulting_api

go 1.18

require (
	auth v0.0.0-00010101000000-000000000000
	github.com/gofiber/fiber/v2 v2.34.1
	go.mongodb.org/mongo-driver v1.9.1
	own_logger v0.0.0-00010101000000-000000000000
)

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.4.0 // indirect
)

require (
	cache v0.0.0-00010101000000-000000000000
	github.com/andybalholm/brotli v1.0.4 // indirect
	github.com/dgrijalva/jwt-go v3.2.0+incompatible // indirect
	github.com/go-redis/redis/v8 v8.11.5
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
	golang.org/x/crypto v0.0.0-20220622213112-05595931fe9d // indirect
	golang.org/x/sync v0.0.0-20220601150217-0de741cfad7f // indirect
	golang.org/x/sys v0.0.0-20220622161953-175b2fd9d664 // indirect
	golang.org/x/text v0.3.7 // indirect
)

replace own_logger => ./../own_logger

replace auth => ./../auth

replace cache => ./../cache
