// Package service содержит сервис нотификатор
package service

import (
	"context"

	"route256.ozon.ru/project/notifier/internal/models"
)

type (
	// EventHandler это обработчик событий
	EventHandler interface {
		Handle(context.Context, *models.OrderStatusChangedEvent) error
	}
	eventHandlingRunner interface {
		Run(context.Context, EventHandler) error
	}
	// Notifier это сервис, который предназначен для организации обработки поступающих сообщений.
	Notifier struct {
		handler        EventHandler
		handlingRunner eventHandlingRunner
	}
)

// NewNotifier новый Notifier, который может запустить обработку сообщений handler'ом
func NewNotifier(handler EventHandler, handlingRunner eventHandlingRunner) *Notifier {
	return &Notifier{handler: handler, handlingRunner: handlingRunner}
}

// Run запускает обработку сообщений
func (n *Notifier) Run(ctx context.Context) error {
	return n.handlingRunner.Run(ctx, n.handler)
}
