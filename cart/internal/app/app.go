package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	addPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/add"
	clearPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/clear"
	deletePkg "route256.ozon.ru/project/cart/internal/controllers/handlers/delete"
	listPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/list"
	"route256.ozon.ru/project/cart/internal/controllers/middleware"
	productservice2 "route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/retryableclient"
	"route256.ozon.ru/project/cart/internal/providers/repository"
	"route256.ozon.ru/project/cart/internal/providers/repository/inmemorycart"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"route256.ozon.ru/project/cart/internal/usecases/modifier"
)

func Run() {
	inMemoryCartCreator := &inmemorycart.CartCreator{}
	cartRepo := repository.New(inMemoryCartCreator)
	baseUrl, err := url.Parse("http://route256.pavl.uk:8080/")
	if err != nil {
		log.Fatal(err)
	}
	clientForProductService := retryableclient.NewRetryableClient(3, productservice2.RetryCondition)
	prodService := productservice2.New(clientForProductService, baseUrl, "testtoken")
	cartModifierService := modifier.New(cartRepo, prodService)
	cartListerService := lister.New(cartRepo, prodService)

	mux := http.NewServeMux()
	addHandler := addPkg.New(cartModifierService)
	mux.Handle(fmt.Sprintf("POST /user/{%s}/cart/{%s}", addPkg.UserIdSegment, addPkg.SkuIdSegment), addHandler)

	clearHandler := clearPkg.New(cartModifierService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart", clearPkg.UserIdSegment), clearHandler)

	deleteHandler := deletePkg.New(cartModifierService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", deletePkg.UserIdSegment, deletePkg.SkuIdSegment), deleteHandler)

	listHandler := listPkg.New(cartListerService)
	mux.Handle(fmt.Sprintf("GET /user/{%s}/cart", listPkg.UserIdSegment), listHandler)

	loggedReqsHandler := middleware.NewLogger(mux)
	if err = http.ListenAndServe("0.0.0.0:8080", loggedReqsHandler); err != nil {
		log.Fatal(err)
	}
}
