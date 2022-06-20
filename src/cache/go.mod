module cache

go 1.18

require github.com/go-redis/redis/v8 v8.11.5

require (
	github.com/cespare/xxhash/v2 v2.1.2 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/joho/godotenv v1.4.0
	own_logger v0.0.0-00010101000000-000000000000
)

replace own_logger => ./../own_logger
