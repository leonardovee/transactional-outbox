package main

import (
	"context"
	"log"

	"github.com/jackc/pgx/v5"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/order"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/outbox"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/service"
)

func main() {
	ctx := context.Background()

	conn, err := pgx.Connect(ctx, "postgres://postgres:postgres@127.0.0.1/postgres")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close(ctx)

	err = conn.Ping(ctx)
	if err != nil {
		log.Fatal(err)
	}

	orderQueries := order.New(conn)
	outboxQueries := outbox.New(conn)
	service := service.New(conn, orderQueries, outboxQueries)

	for {
		service.ProcessOutbox(ctx)
	}
}
