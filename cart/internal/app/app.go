package app

import (
	"context"
	"errors"
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
	"sync/atomic"
	"time"
)

type App struct {
	config              Config
	cartRepo            *repository.InMemoryCartRepository
	itPresChecker       *itempresencechecker.ItemPresenceChecker
	prodInfoGetter      *productinfogetter.ProductInfoGetter
	cartModifierService *modifier.CartModifierService
	cartListerService   *lister.CartListerService
	HttpController      http.Handler
	server              atomic.Pointer[http.Server]
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

func (a *App) InitController() {
	if a.HttpController != nil {
		return
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
}

func (a *App) GetAddr() (string, error) {
	serverConfig := a.config.Server
	hostAddr, err := netip.ParseAddr(serverConfig.Host)
	if err != nil {
		return "", err
	}
	port := serverConfig.Port
	fullAddr := netip.AddrPortFrom(hostAddr, port)
	return fullAddr.String(), nil
}

// Run представляет из себя блокирующий вызов, который запускает новый сервер, согласно текущей конфигурации.
// Если он уже запущен, то функция ничего не делает. Если не удалось запустить, вся программа завершается с ошибкой
func (a *App) Run() {
	if a.server.Load() != nil {
		return
	}
	a.InitController()
	addr, err := a.GetAddr()
	if err != nil {
		log.Fatal(err)
	}
	newServer := &http.Server{
		Addr:    addr,
		Handler: a.HttpController,
	}
	if !a.server.CompareAndSwap(nil, newServer) {
		return
	}
	if err = newServer.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}
}

// Stop останавливает запущенный сервер в течение ShutdownTimoutSeconds секунд. Если не был запущен, функция ничего не делает. Если не удалось
// остановить в течение таймаута, вся программа завершается с ошибкой. Возврат из функции произойдёт, когда shutdown завершится.
func (a *App) Stop() {
	srvToShutdown := a.server.Swap(nil)
	if srvToShutdown == nil {
		return
	}
	timeout := time.Duration(a.config.Server.ShutdownTimoutSeconds) * time.Second
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()
	if err := srvToShutdown.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
}
