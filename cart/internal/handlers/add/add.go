package add

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/modifier"
	"strconv"
)

const (
	UserIdSegment = "userId"
	SkuIdSegment  = "skuId"
)

type Add struct {
	cartService modifier.CartModifierService
}

func New(cartService modifier.CartModifierService) *Add {
	return &Add{
		cartService: cartService,
	}
}

func (h *Add) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := addItemReqFromR(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.cartService.AddItem(r.Context(), req.userId, req.skuId, req.count); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

type addItemReq struct {
	userId service.User
	skuId  service.SkuId
	count  service.ItemCount
}

func addItemReqFromR(r *http.Request) (*addItemReq, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err1 := strconv.Atoi(userIdStr)
	skuIdStr := r.PathValue(SkuIdSegment)
	skuId, err2 := strconv.Atoi(skuIdStr)
	bodyData, err3 := io.ReadAll(r.Body)
	if err3 != nil {
		return nil, errors.Join(err1, err2, err3)
	}
	var reqBody addRequest
	err3 = json.Unmarshal(bodyData, &reqBody)
	allErrs := errors.Join(err1, err2, err3)
	if allErrs != nil {
		return nil, allErrs
	}
	return &addItemReq{
		userId: service.User(userId),
		skuId:  service.SkuId(skuId),
		count:  reqBody.Count,
	}, nil
}
