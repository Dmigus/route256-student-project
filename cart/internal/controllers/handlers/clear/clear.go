package clear

import (
	"context"
	"fmt"
	"math"
	"net/http"
	"strconv"
)

const UserIdSegment = "userId"

var errIncorrectUserId = fmt.Errorf("userId must be number in range [%d, %d]", math.MinInt64, math.MaxInt64)

type clearCartService interface {
	ClearCart(ctx context.Context, user int64) error
}

type Clear struct {
	cartCleaner clearCartService
}

func New(cartService clearCartService) *Clear {
	return &Clear{
		cartCleaner: cartService,
	}
}

func (h *Clear) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userId, err := parseUserId(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.cartCleaner.ClearCart(r.Context(), userId); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func parseUserId(r *http.Request) (int64, error) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err := strconv.ParseInt(userIdStr, 10, 64)
	if err != nil {
		return 0, errIncorrectUserId
	}
	return userId, nil
}
