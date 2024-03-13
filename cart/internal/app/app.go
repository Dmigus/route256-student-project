package app

import (
	"fmt"
	"log"
	"net/http"
	"net/netip"
	"net/url"
	"route256.ozon.ru/project/cart/internal/client"
	"route256.ozon.ru/project/cart/internal/client/policies"
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
	config              Config
	cartRepo            *repository.InMemoryCartRepository
	itPresChecker       *itempresencechecker.ItemPresenceChecker
	prodInfoGetter      *productinfogetter.ProductInfoGetter
	cartModifierService *modifier.CartModifierService
	cartListerService   *lister.CartListerService
	HttpController      http.Handler
}

func NewApp(config Config) *App {
	return &App{
		config: config,
	}
}

func (a *App) initRepo() {
	a.cartRepo = repository.New()
}

func (a *App) initProductService() {
	prodServConfig := a.config.ProductService
	baseUrl, err := url.Parse(prodServConfig.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	retryPolicy := policies.NewRetryOnStatusCodes(prodServConfig.RetryPolicy.RetryStatusCodes, prodServConfig.RetryPolicy.MaxRetries)
	clientForProductService := client.NewRetryableClient(retryPolicy)
	rcPerformer := productservice.NewRCPerformer(clientForProductService, baseUrl, prodServConfig.AccessToken)
	a.itPresChecker = itempresencechecker.NewItemPresenceChecker(rcPerformer)
	a.prodInfoGetter = productinfogetter.NewProductInfoGetter(rcPerformer)
}

func (a *App) InitCartServices() {
	if a.cartRepo == nil {
		a.initRepo()
	}
	if a.itPresChecker == nil && a.prodInfoGetter == nil {
		a.initProductService()
	}
	a.cartModifierService = modifier.New(a.cartRepo, a.itPresChecker)
	a.cartListerService = lister.New(a.cartRepo, a.prodInfoGetter)
}

func (a *App) InitController() http.Handler {
	if a.HttpController != nil {
		return a.HttpController
	}
	if a.cartModifierService == nil && a.cartListerService == nil {
		a.InitCartServices()
	}
	mux := http.NewServeMux()
	addHandler := addPkg.New(a.cartModifierService)
	mux.Handle(fmt.Sprintf("POST /user/{%s}/cart/{%s}", addPkg.UserIdSegment, addPkg.SkuIdSegment), addHandler)
	clearHandler := clearPkg.New(a.cartModifierService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart", clearPkg.UserIdSegment), clearHandler)
	deleteHandler := deletePkg.New(a.cartModifierService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", deletePkg.UserIdSegment, deletePkg.SkuIdSegment), deleteHandler)
	listHandler := listPkg.New(a.cartListerService)
	mux.Handle(fmt.Sprintf("GET /user/{%s}/cart", listPkg.UserIdSegment), listHandler)
	a.HttpController = middleware.NewLogger(mux)
	return a.HttpController
}

func (a *App) Run() {
	a.InitController()
	serverConfig := a.config.Server
	hostAddr, err := netip.ParseAddr(serverConfig.Host)
	if err != nil {
		log.Fatal(err)
	}
	port := serverConfig.Port
	fullAddr := netip.AddrPortFrom(hostAddr, port)
	if err = http.ListenAndServe(fullAddr.String(), a.HttpController); err != nil {
		log.Fatal(err)
	}
}
