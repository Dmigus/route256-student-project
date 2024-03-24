package checkout

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"route256.ozon.ru/project/cart/internal/models"
)

var (
	errIncorrectUserId = fmt.Errorf("request body must contain user in range [%d, %d]", math.MinInt64, math.MaxInt64)
	errIO              = errors.New("reading bytes from request body failed")
)

type checkoutService interface {
	Checkout(ctx context.Context, userId int64) (int64, error)
}

type Checkout struct {
	checkouter checkoutService
}

func New(checkouter checkoutService) *Checkout {
	return &Checkout{
		checkouter: checkouter,
	}
}

func (c *Checkout) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := parseRequest(r)
	if err != nil {
		if errors.Is(err, errIO) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	orderId, err := c.checkouter.Checkout(r.Context(), userId)
	if err != nil {
		fillHeaderFromError(w, err)
		return
	}
	dto := orderIdToDTO(orderId)
	marshalled, err := json.Marshal(dto)
	if err != nil {
		fillHeaderFromError(w, err)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write(marshalled)
}

func fillHeaderFromError(w http.ResponseWriter, err error) {
	switch {
	case errors.Is(err, models.ErrFailedPrecondition):
		w.WriteHeader(http.StatusPreconditionFailed)
	case errors.Is(err, models.ErrNotFound):
		w.WriteHeader(http.StatusNotFound)
	case err != nil:
		w.WriteHeader(http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}

func parseRequest(r *http.Request) (int64, error) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return 0, errIO
	}
	var request checkoutRequest
	err = json.Unmarshal(bodyData, &request)
	if err != nil || request.User == nil {
		return 0, errIncorrectUserId
	}
	return *request.User, nil
}

func orderIdToDTO(orderId int64) checkoutResponse {
	return checkoutResponse{OrderId: orderId}
}
