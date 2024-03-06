package app

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"route256.ozon.ru/project/cart/internal/cartrepository"
	"route256.ozon.ru/project/cart/internal/handlers/add"
	clear2 "route256.ozon.ru/project/cart/internal/handlers/clear"
	delete2 "route256.ozon.ru/project/cart/internal/handlers/delete"
	"route256.ozon.ru/project/cart/internal/handlers/list"
	"route256.ozon.ru/project/cart/internal/inmemorycart"
	"route256.ozon.ru/project/cart/internal/middleware"
	"route256.ozon.ru/project/cart/internal/productservice"
	"route256.ozon.ru/project/cart/internal/retryableclient"
	"route256.ozon.ru/project/cart/internal/service/lister"
	"route256.ozon.ru/project/cart/internal/service/modifier"
)

func Run() {
	inMemoryCartFabric := &inmemorycart.Fabric{}
	cartRepo := cartrepository.New(inMemoryCartFabric)
	baseUrl, err := url.Parse("http://route256.pavl.uk:8080/")
	if err != nil {
		log.Fatal(err)
	}
	clientForProductService := retryableclient.NewRetryableClient(3, productservice.RetryCondition)
	prodService := productservice.New(clientForProductService, baseUrl, "testtoken")
	cartModifierService := modifier.New(cartRepo, prodService)
	cartListerService := lister.New(cartRepo, prodService)

	mux := http.NewServeMux()
	addHandler := add.New(cartModifierService)
	mux.Handle(fmt.Sprintf("POST /user/{%s}/cart/{%s}", add.UserIdSegment, add.SkuIdSegment), addHandler)

	clearHandler := clear2.New(cartModifierService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart", clear2.UserIdSegment), clearHandler)

	deleteHandler := delete2.New(cartModifierService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", delete2.UserIdSegment, delete2.SkuIdSegment), deleteHandler)

	listHandler := list.New(cartListerService)
	mux.Handle(fmt.Sprintf("GET /user/{%s}/cart", list.UserIdSegment), listHandler)

	loggedReqsHandler := middleware.NewLogger(mux)
	if err = http.ListenAndServe("0.0.0.0:8080", loggedReqsHandler); err != nil {
		log.Fatal(err)
	}
}
