package add

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"route256.ozon.ru/project/cart/internal/models"
	"route256.ozon.ru/project/cart/internal/usecases/adder"
	"strconv"
)

const (
	UserIdSegment = "userId"
	SkuIdSegment  = "skuId"
)

var (
	errIncorrectUserId = fmt.Errorf("userId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)
	errIncorrectSkuId  = fmt.Errorf("skuId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)
	errIncorrectCount  = fmt.Errorf("request body must contain Count in range [%d, %d]", 1, math.MaxUint16)
	errIO              = errors.New("reading bytes from request body failed")
)

type itemAdderService interface {
	AddItem(ctx context.Context, user int64, skuId int64, count uint16) error
}

type Add struct {
	adder itemAdderService
}

func New(cartService itemAdderService) *Add {
	return &Add{
		adder: cartService,
	}
}

func (h *Add) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	req, err := addItemReqFromR(r)
	if err != nil {
		if errors.Is(err, errIO) {
			w.WriteHeader(http.StatusInternalServerError)
		} else {
			http.Error(w, err.Error(), http.StatusBadRequest)
		}
		return
	}
	if err = h.adder.AddItem(r.Context(), req.userId, req.skuId, req.count); err != nil {
		if errors.Is(err, adder.ErrNotEnoughNumInStocks) {
			w.WriteHeader(http.StatusPreconditionFailed)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
		}
		return
	}
	w.WriteHeader(http.StatusOK)
}

type addItemReq struct {
	userId int64
	skuId  int64
	count  uint16
}

func addItemReqFromR(r *http.Request) (*addItemReq, error) {
	bodyData, err := io.ReadAll(r.Body)
	if err != nil {
		return nil, errIO
	}
	userId, errUserId := parseUserId(r)
	skuId, errSkuId := parseSkuId(r)
	count, errCount := parseCount(bodyData)
	if allErrs := errors.Join(errUserId, errSkuId, errCount); allErrs != nil {
		return nil, allErrs
	}
	return &addItemReq{
		userId: userId,
		skuId:  skuId,
		count:  count,
	}, nil
}

func parseUserId(r *http.Request) (int64, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectUserId
	}
	return userId, nil
}

func parseSkuId(r *http.Request) (int64, error) {
	skuIdStr := r.PathValue(SkuIdSegment)
	skuId, err := strconv.ParseInt(skuIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectSkuId
	}
	return skuId, nil
}

func parseCount(data []byte) (uint16, error) {
	var reqBody addRequest
	err := json.Unmarshal(data, &reqBody)
	if err != nil || reqBody.Count == nil || !models.IsNumberValidCount(*reqBody.Count) {
		return 0, errIncorrectCount
	}
	return *reqBody.Count, nil
}
