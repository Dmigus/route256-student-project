package add

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/modifier"
	"strconv"
)

const (
	UserIdSegment = "userId"
	SkuIdSegment  = "skuId"
)

var errIncorrectUserId = fmt.Errorf("userId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)
var errIncorrectSkuId = fmt.Errorf("skuId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)
var errIncorrectCount = fmt.Errorf("request body must contain Count in range [%d, %d]", 1, math.MaxUint16)
var errIO = errors.New("reading bytes from request body failed")

type Add struct {
	cartService *modifier.CartModifierService
}

func New(cartService *modifier.CartModifierService) *Add {
	return &Add{
		cartService: cartService,
	}
}

func (h *Add) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := addItemReqFromR(r)
	if err != nil {
		if err == errIO {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
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
	userId, err1 := parseUserId(r)
	skuId, err2 := parseSkuId(r)
	bodyData, err3 := io.ReadAll(r.Body)
	if err3 != nil {
		return nil, errIO
	}
	count, err3 := parseCount(bodyData)
	if allErrs := errors.Join(err1, err2, err3); allErrs != nil {
		return nil, allErrs
	}
	return &addItemReq{
		userId: userId,
		skuId:  skuId,
		count:  count,
	}, nil
}

func parseUserId(r *http.Request) (service.User, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectUserId
	}
	return userId, nil
}

func parseSkuId(r *http.Request) (service.SkuId, error) {
	skuIdStr := r.PathValue(SkuIdSegment)
	skuId, err := strconv.ParseInt(skuIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectSkuId
	}
	return skuId, nil
}

func parseCount(data []byte) (service.ItemCount, error) {
	var reqBody addRequest
	err := json.Unmarshal(data, &reqBody)
	if err != nil || reqBody.Count == 0 {
		return 0, errIncorrectCount
	} else {
		return reqBody.Count, nil
	}
}
