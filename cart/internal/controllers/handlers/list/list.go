package list

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"net/http"
	"route256.ozon.ru/project/cart/internal/usecases"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"sort"
	"strconv"
)

const UserIdSegment = "userId"

var errIncorrectUserId = fmt.Errorf("userId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)

type cartListerService interface {
	ListCartContent(ctx context.Context, user usecases.User) (*lister.CartContent, error)
}

type List struct {
	cartLister cartListerService
}

func New(cartService cartListerService) *List {
	return &List{
		cartLister: cartService,
	}
}

func (h *List) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := parseUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	list, err := h.cartLister.ListCartContent(r.Context(), userId)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	dto := listToDTO(list)
	if isCartEmpty(dto) {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	marshalled, err := json.Marshal(dto)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	_, err = w.Write(marshalled)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func parseUserId(r *http.Request) (usecases.User, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectUserId
	}
	return userId, nil
}

func listToDTO(content *lister.CartContent) listResponse {
	lr := listResponse{
		TotalPrice: content.TotalPrice,
		Items:      make([]listResponseItem, 0, len(content.Items)),
	}
	for _, item := range content.Items {
		lr.Items = append(lr.Items, listResponseItem{
			SkuId: item.SkuId,
			Name:  item.Name,
			Count: item.Count,
			Price: item.Price,
		})
	}
	sort.Slice(lr.Items, func(i, j int) bool {
		return lr.Items[i].SkuId < lr.Items[j].SkuId
	})
	return lr
}

func isCartEmpty(list listResponse) bool {
	return len(list.Items) == 0
}
