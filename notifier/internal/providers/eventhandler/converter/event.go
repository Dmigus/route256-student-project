// Package converter содержит функции конвертации для типов, сгенерированных из proto файла
package converter

import (
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"route256.ozon.ru/project/notifier/internal/models"
	v1 "route256.ozon.ru/project/notifier/internal/pkg/api/loms/v1"
)

// MessageToChangeOrderStatusEvent конвертирует EventMessage в событие изменения статуса заказа
func MessageToChangeOrderStatusEvent(message *models.EventMessage) (*v1.ChangeOrderStatusEvent, error) {
	evMessage := v1.ChangeOrderStatusEvent{}
	if err := proto.Unmarshal(message.Payload, &evMessage); err != nil {
		return nil, errors.Wrap(err, "failed to parse message")
	}
	return &evMessage, nil
}
