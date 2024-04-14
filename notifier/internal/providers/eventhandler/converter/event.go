package converter

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	"route256.ozon.ru/project/notifier/internal/models"
	"route256.ozon.ru/project/notifier/internal/pkg/api/loms/v1"
)

func MessageToChangeOrderStatusEvent(message *models.EventMessage) (*v1.ChangeOrderStatusEvent, error) {
	evMessage := v1.ChangeOrderStatusEvent{}
	if err := proto.Unmarshal(message.Payload, &evMessage); err != nil {
		return nil, errors.Wrap(err, "failed to parse message")
	}
	return &evMessage, nil
}
