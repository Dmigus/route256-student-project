package app

import (
	"fmt"
	"github.com/gookit/config/v2"
	"github.com/gookit/config/v2/yaml"
	"log"
	"net/http"
	"net/netip"
	"net/url"
	addPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/add"
	clearPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/clear"
	deletePkg "route256.ozon.ru/project/cart/internal/controllers/handlers/delete"
	listPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/list"
	"route256.ozon.ru/project/cart/internal/controllers/middleware"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/retryableclient"
	"route256.ozon.ru/project/cart/internal/providers/repository"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"route256.ozon.ru/project/cart/internal/usecases/modifier"
)

func Run() {
	config.WithOptions(config.ParseEnv)
	config.AddDriver(yaml.Driver)
	err := config.LoadFiles("configs/product_service.yml", "configs/server.yml")
	if err != nil {
		panic(err)
	}
	cartRepo := repository.New()
	baseUrl, err := url.Parse(config.String("baseURL"))
	if err != nil {
		log.Fatal(err)
	}
	clientForProductService := retryableclient.NewRetryableClient(
		config.Int("maxRetries"),
		productservice.RetryCondition)
	prodService := productservice.New(clientForProductService, baseUrl, config.String("accessToken"))
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

	hostAddr, err := netip.ParseAddr(config.String("host"))
	if err != nil {
		log.Fatal(err)
	}
	port := uint16(config.Uint("port"))
	fullAddr := netip.AddrPortFrom(hostAddr, port)
	if err = http.ListenAndServe(fullAddr.String(), loggedReqsHandler); err != nil {
		log.Fatal(err)
	}
}
