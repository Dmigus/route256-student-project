package clear

import (
	"net/http"
	"route256.ozon.ru/project/cart/internal/service"
	"route256.ozon.ru/project/cart/internal/service/modifier"
	"strconv"
)

const UserIdSegment = "userId"

type Clear struct {
	cartService *modifier.CartModifierService
}

func New(cartService *modifier.CartModifierService) *Clear {
	return &Clear{
		cartService: cartService,
	}
}

func (h *Clear) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	userIdStr := r.PathValue(UserIdSegment)
	userId, err := strconv.Atoi(userIdStr)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	if err = h.cartService.ClearCart(r.Context(), service.User(userId)); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}
