package service

import "context"
import "route256.ozon.ru/project/notifier/internal/models"

type (
	EventHandler interface {
		Handle(context.Context, *models.OrderStatusChangedEvent) error
	}
	eventHandlingRunner interface {
		Run(context.Context, EventHandler) error
	}
	Notifier struct {
		handler        EventHandler
		handlingRunner eventHandlingRunner
	}
)

func NewNotifier(handler EventHandler, handlingRunner eventHandlingRunner) *Notifier {
	return &Notifier{handler: handler, handlingRunner: handlingRunner}
}

func (n *Notifier) Run(ctx context.Context) error {
	return n.handlingRunner.Run(ctx, n.handler)
}
