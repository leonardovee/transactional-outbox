package main

import (
	"context"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/handler"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/order"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/outbox"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/service"
)

func main() {
	ctx := context.Background()

	conn, err := pgxpool.New(ctx, "postgres://postgres:postgres@127.0.0.1/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	orderQueries := order.New(conn)
	outboxQueries := outbox.New(conn)
	service := service.New(conn, orderQueries, outboxQueries)
	handler := handler.New(service)

	go startWorker(ctx, service)
	startServer(handler)
}

func startServer(handler *handler.Handler) {
	r := gin.Default()

	r.POST("/orders/command", handler.Command)

	r.Run()
}

func startWorker(ctx context.Context, service *service.Service) {
	for {
		time.Sleep(200 * time.Millisecond)

		err := service.ProcessOutbox(ctx)
		if err != nil {
			log.Println(err)
		}
	}
}
