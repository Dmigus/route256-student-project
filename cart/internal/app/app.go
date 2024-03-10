package app

import (
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"net/url"
	"route256.ozon.ru/project/cart/internal/client"
	addPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/add"
	clearPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/clear"
	deletePkg "route256.ozon.ru/project/cart/internal/controllers/handlers/delete"
	listPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/list"
	"route256.ozon.ru/project/cart/internal/controllers/middleware"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/itempresencechecker"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
	"route256.ozon.ru/project/cart/internal/providers/repository"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"route256.ozon.ru/project/cart/internal/usecases/modifier"
)

type App struct {
	config Config
}

func NewApp(config Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) Run() {
	prodServConfig := a.config.ProductService
	cartRepo := repository.New()
	baseUrl, err := url.Parse(prodServConfig.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	retryPolicy := client.NewRetryOnStatusCodes(prodServConfig.RetryPolicy.RetryStatusCodes, prodServConfig.RetryPolicy.MaxRetries)
	clientForProductService := client.NewRetryableClient(retryPolicy)
	rcPerformer := productservice.NewRCPerformer(clientForProductService, baseUrl, prodServConfig.AccessToken)
	itPresChecker := itempresencechecker.NewItemPresenceChecker(rcPerformer)
	prodInfoGetter := productinfogetter.NewProductInfoGetter(rcPerformer)
	cartModifierService := modifier.New(cartRepo, itPresChecker)
	cartListerService := lister.New(cartRepo, prodInfoGetter)

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

	serverConfig := a.config.Server
	hostAddr, err := netip.ParseAddr(serverConfig.Host)
	if err != nil {
		log.Fatal(err)
	}
	port := serverConfig.Port
	fullAddr := netip.AddrPortFrom(hostAddr, port)
	if err = http.ListenAndServe(fullAddr.String(), loggedReqsHandler); err != nil {
		log.Fatal(err)
	}
}
