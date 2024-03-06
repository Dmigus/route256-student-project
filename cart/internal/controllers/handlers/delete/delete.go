package delete

import (
	"context"
	"errors"
	"fmt"
	"math"
	"net/http"
	"route256.ozon.ru/project/cart/internal/usecases"
	"strconv"
)

const (
	UserIdSegment = "userId"
	SkuIdSegment  = "skuId"
)

var (
	errIncorrectUserId = fmt.Errorf("userId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)
	errIncorrectSkuId  = fmt.Errorf("skuId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)
)

type itemDeleterService interface {
	DeleteItem(ctx context.Context, user usecases.User, skuId usecases.SkuId) error
}

type Delete struct {
	deleter itemDeleterService
}

func New(cartService itemDeleterService) *Delete {
	return &Delete{
		deleter: cartService,
	}
}

func (h *Delete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := deleteItemReqFromR(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.deleter.DeleteItem(r.Context(), req.userId, req.skuId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type deleteItemReq struct {
	userId usecases.User
	skuId  usecases.SkuId
}

func deleteItemReqFromR(r *http.Request) (*deleteItemReq, error) {
	userId, errUserId := parseUserId(r)
	skuId, errSkuId := parseSkuId(r)
	allErrs := errors.Join(errUserId, errSkuId)
	if allErrs != nil {
		return nil, allErrs
	}
	return &deleteItemReq{
		userId: userId,
		skuId:  skuId,
	}, nil
}

func parseUserId(r *http.Request) (usecases.User, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectUserId
	}
	return userId, nil
}

func parseSkuId(r *http.Request) (usecases.SkuId, error) {
	skuIdStr := r.PathValue(SkuIdSegment)
	skuId, err := strconv.ParseInt(skuIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectSkuId
	}
	return skuId, nil
}
