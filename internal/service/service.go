package service

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"leonardovee.dev/microservices-patterns/transactional-outbox/internal/order"
)

type Command struct {
	Command     string
	AggregateID *string
	Total       *int
}

type Service struct {
	queries *order.Queries
}

func New(queries *order.Queries) *Service {
	return &Service{queries: queries}
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
	id := uuid.New().String()
	o, err := s.queries.CreateOrder(ctx, order.CreateOrderParams{
		ID:          id,
		AggregateID: id,
		Status: order.NullOrderStatus(order.NullOrderStatus{
			OrderStatus: "created",
			Valid:       true,
		}),
		Total: total,
	})
	if err != nil {
		return nil, err
	}

	return &o, nil
}

func (s *Service) UpdateAggregate(ctx context.Context, aggregateID string, status order.OrderStatus, total int32) (*order.Order, error) {
	o, err := s.queries.CreateOrder(ctx, order.CreateOrderParams{
		ID:          uuid.New().String(),
		AggregateID: aggregateID,
		Status: order.NullOrderStatus(order.NullOrderStatus{
			OrderStatus: status,
			Valid:       true,
		}),
		Total: total,
	})
	if err != nil {
		return nil, err
	}

	return &o, nil
}
