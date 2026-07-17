package main

import (
	"log"
	"net/http"
	"os"
	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"github.com/dayzer0-hq/be-project-apple-launch-queue-template/internal/handlers"
	"github.com/dayzer0-hq/be-project-apple-launch-queue-template/internal/queue"
)
func main() {
	_ = godotenv.Load()
	rdb := redis.NewClient(&redis.Options{Addr: os.Getenv("REDIS_URL")})
	h := &handlers.Handler{Store: queue.NewRedisQueueStore(rdb)}
	port := os.Getenv("PORT"); if port == "" { port = "8080" }
	log.Printf("server on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, handlers.NewRouter(h)))
}
