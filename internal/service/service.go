package service

import (
	"context"
	"encoding/json"
	"errors"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/order"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/outbox"
)

type Command struct {
	Command     string
	AggregateID *string
	Total       *int
}

type Service struct {
	db            *pgx.Conn
	orderQueries  *order.Queries
	outboxQueries *outbox.Queries
}

func New(db *pgx.Conn, orderQueries *order.Queries, outboxQueries *outbox.Queries) *Service {
	return &Service{
		db,
		orderQueries,
		outboxQueries,
	}
}

func (s *Service) Execute(ctx context.Context, command *Command) (*order.Order, error) {
	switch command.Command {
	case "create":
		return s.CreateOrder(ctx, int32(*command.Total))
	case "approve":
		return s.UpdateAggregate(ctx, *command.AggregateID, "approved", int32(*command.Total))
	case "ready":
		return s.UpdateAggregate(ctx, *command.AggregateID, "ready", int32(*command.Total))
	case "ship":
		return s.UpdateAggregate(ctx, *command.AggregateID, "shipped", int32(*command.Total))
	case "arrive":
		return s.UpdateAggregate(ctx, *command.AggregateID, "arrived", int32(*command.Total))
	default:
		return nil, errors.New("unknown command")
	}
}

func (s *Service) CreateOrder(ctx context.Context, total int32) (*order.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	qtx := s.orderQueries.WithTx(tx)

	id := uuid.New().String()
	o, err := qtx.CreateOrder(ctx, order.CreateOrderParams{
		ID:          id,
		AggregateID: id,
		Status: order.NullOrderStatus(order.NullOrderStatus{
			OrderStatus: "created",
			Valid:       true,
		}),
		Total: total,
	})
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	payload, err := json.Marshal(o)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	_, err = s.outboxQueries.CreateOutbox(ctx, outbox.CreateOutboxParams{
		ID:            uuid.New().String(),
		AggregateID:   id,
		AggregateType: "order",
		Type:          "create",
		Payload:       payload,
	})
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *Service) UpdateAggregate(ctx context.Context, aggregateID string, status order.OrderStatus, total int32) (*order.Order, error) {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return nil, err
	}

	qtx := s.orderQueries.WithTx(tx)

	o, err := qtx.CreateOrder(ctx, order.CreateOrderParams{
		ID:          uuid.New().String(),
		AggregateID: aggregateID,
		Status: order.NullOrderStatus(order.NullOrderStatus{
			OrderStatus: status,
			Valid:       true,
		}),
		Total: total,
	})
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	payload, err := json.Marshal(o)
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	_, err = s.outboxQueries.CreateOutbox(ctx, outbox.CreateOutboxParams{
		ID:            uuid.New().String(),
		AggregateID:   aggregateID,
		AggregateType: "order",
		Type:          string(status),
		Payload:       payload,
	})
	if err != nil {
		tx.Rollback(ctx)
		return nil, err
	}

	err = tx.Commit(ctx)
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *Service) ProcessOutbox(ctx context.Context) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}

	qtx := s.outboxQueries.WithTx(tx)

	outboxes, err := qtx.ListOutbox(ctx, 10)
	if err != nil {
		tx.Rollback(ctx)
		return err
	}

	for _, o := range outboxes {
		err := s.processOutboxItem(ctx, o)
		if err != nil {
			tx.Rollback(ctx)
			return err
		}
	}

	err = tx.Commit(ctx)
	if err != nil {
		return err
	}

	return nil
}

func (s *Service) processOutboxItem(ctx context.Context, o outbox.Outbox) error {
	log.Printf("Processing outbox item: %v", o)

	_, err := s.outboxQueries.DeleteOutbox(ctx, o.ID)
	if err != nil {
		return err
	}

	return nil
}
