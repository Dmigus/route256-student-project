// Package outboxsender содержит реализацию работы transactional outbox
package outboxsender

import (
	"context"
	"fmt"
	"time"

	"route256.ozon.ru/project/loms/internal/models"
)

type (
	// Outbox представляет из себя хранилище событий, из которого можно получать набор событий не более
	Outbox interface {
		PullNextEvents(ctx context.Context, batchSize int32) ([]models.OrderStatusChangedEvent, error)
	}
	eventsPusher interface {
		PushOrderStatusChangedEvents(ctx context.Context, events []models.OrderStatusChangedEvent) error
	}
	txManager interface {
		WithinTransaction(context.Context, func(ctx context.Context, outbox Outbox) bool) error
	}
	// Service представляет сервис, осуществляющий регулярную выгрузку ивентов из хранилища событий в брокер
	Service struct {
		tx           txManager
		broker       eventsPusher
		iterInterval time.Duration
		batchSize    int32
	}
)

// NewService создаёт новый Service.
func NewService(tx txManager, broker eventsPusher, iterInterval time.Duration, batchSize int32) *Service {
	return &Service{
		tx:           tx,
		broker:       broker,
		iterInterval: iterInterval,
		batchSize:    batchSize,
	}
}

// Run запускает обработку событий. Блокирующая. Остановиться, когда переданный контекст будет отменён
func (s *Service) Run(ctx context.Context) {
	t := time.NewTicker(s.iterInterval)
	defer t.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-t.C:
			// пока непонятно, куда логировать ошибку
			s.iteration(ctx)
		}
	}
}

func (s *Service) iteration(serviceLiveCtx context.Context) error {
	iterCtx, cancel := context.WithCancel(serviceLiveCtx)
	defer cancel()
	var serviceErr error
	trErr := s.tx.WithinTransaction(iterCtx, func(ctx context.Context, outbox Outbox) bool {
		events, err := outbox.PullNextEvents(ctx, s.batchSize)
		if err != nil {
			serviceErr = fmt.Errorf("could not pull events from outbox: %w", err)
			return false
		}
		err = s.broker.PushOrderStatusChangedEvents(ctx, events)
		if err != nil {
			serviceErr = fmt.Errorf("could not push events to broker: %w", err)
			return false
		}
		return true
	})
	if serviceErr != nil {
		return serviceErr
	}
	return trErr
}
