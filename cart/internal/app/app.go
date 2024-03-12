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
	config              Config
	cartRepo            *repository.InMemoryCartRepository
	itPresChecker       *itempresencechecker.ItemPresenceChecker
	prodInfoGetter      *productinfogetter.ProductInfoGetter
	cartModifierService *modifier.CartModifierService
	cartListerService   *lister.CartListerService
	AddHandler          *addPkg.Add
	ClearHandler        *clearPkg.Clear
	DeleteHandler       *deletePkg.Delete
	ListHandler         *listPkg.List
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
	retryPolicy := client.NewRetryOnStatusCodes(prodServConfig.RetryPolicy.RetryStatusCodes, prodServConfig.RetryPolicy.MaxRetries)
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

func (a *App) InitHandlers() {
	if a.cartModifierService == nil && a.cartListerService == nil {
		a.InitCartServices()
	}
	a.AddHandler = addPkg.New(a.cartModifierService)
	a.ClearHandler = clearPkg.New(a.cartModifierService)
	a.DeleteHandler = deletePkg.New(a.cartModifierService)
	a.ListHandler = listPkg.New(a.cartListerService)
}

func (a *App) Run() {
	a.InitHandlers()
	mux := http.NewServeMux()
	mux.Handle(fmt.Sprintf("POST /user/{%s}/cart/{%s}", addPkg.UserIdSegment, addPkg.SkuIdSegment), a.AddHandler)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart", clearPkg.UserIdSegment), a.ClearHandler)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", deletePkg.UserIdSegment, deletePkg.SkuIdSegment), a.DeleteHandler)
	mux.Handle(fmt.Sprintf("GET /user/{%s}/cart", listPkg.UserIdSegment), a.ListHandler)

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
