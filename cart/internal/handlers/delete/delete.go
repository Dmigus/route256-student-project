package delete

import (
	"errors"
	"net/http"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/modifier"
	"strconv"
)

const (
	UserIdSegment = "userId"
	SkuIdSegment  = "skuId"
)

type Delete struct {
	cartService *modifier.CartModifierService
}

func New(cartService *modifier.CartModifierService) *Delete {
	return &Delete{
		cartService: cartService,
	}
}

func (h *Delete) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := deleteItemReqFromR(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.cartService.DeleteItem(r.Context(), req.userId, req.skuId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

type deleteItemReq struct {
	userId service.User
	skuId  service.SkuId
}

func deleteItemReqFromR(r *http.Request) (*deleteItemReq, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err1 := strconv.Atoi(userIdStr)
	skuIdStr := r.PathValue(SkuIdSegment)
	skuId, err2 := strconv.Atoi(skuIdStr)
	allErrs := errors.Join(err1, err2)
	if allErrs != nil {
		return nil, allErrs
	}
	return &deleteItemReq{
		userId: service.User(userId),
		skuId:  service.SkuId(skuId),
	}, nil
}
