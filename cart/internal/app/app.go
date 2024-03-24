package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jasongerard/healthz"
	"log"
	"net/http"
	"net/netip"
	"net/url"
	lomsClientPkg "route256.ozon.ru/project/cart/internal/clients/loms"
	"route256.ozon.ru/project/cart/internal/clients/retriablehttp"
	"route256.ozon.ru/project/cart/internal/clients/retriablehttp/policies"
	addPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/add"
	checkoutPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/checkout"
	clearPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/clear"
	deletePkg "route256.ozon.ru/project/cart/internal/controllers/handlers/delete"
	listPkg "route256.ozon.ru/project/cart/internal/controllers/handlers/list"
	"route256.ozon.ru/project/cart/internal/controllers/middleware"
	lomsProviderPkg "route256.ozon.ru/project/cart/internal/providers/loms"
	"route256.ozon.ru/project/cart/internal/providers/productservice"
	"route256.ozon.ru/project/cart/internal/providers/productservice/itempresencechecker"
	"route256.ozon.ru/project/cart/internal/providers/productservice/productinfogetter"
	"route256.ozon.ru/project/cart/internal/providers/repository"
	"route256.ozon.ru/project/cart/internal/usecases"
	"route256.ozon.ru/project/cart/internal/usecases/adder"
	"route256.ozon.ru/project/cart/internal/usecases/checkouter"
	"route256.ozon.ru/project/cart/internal/usecases/clearer"
	"route256.ozon.ru/project/cart/internal/usecases/deleter"
	"route256.ozon.ru/project/cart/internal/usecases/lister"
	"sync/atomic"
	"time"
)

type App struct {
	config         Config
	httpController http.Handler
	server         atomic.Pointer[http.Server]
}

// NewApp возращает App, готовый к запуску
func NewApp(config Config) *App {
	app := &App{
		config: config,
	}
	app.init()
	return app
}

func (a *App) init() {
	cartRepo := repository.New()
	prodServConfig := a.config.ProductService
	baseUrl, err := url.Parse(prodServConfig.BaseURL)
	if err != nil {
		log.Fatal(err)
	}
	retryPolicy := policies.NewRetryOnStatusCodes(prodServConfig.RetryPolicy.RetryStatusCodes, prodServConfig.RetryPolicy.MaxRetries)
	clientForProductService := retriablehttp.NewRetryableClient(retryPolicy)
	rcPerformer := productservice.NewRCPerformer(clientForProductService, baseUrl, prodServConfig.AccessToken)
	itPresChecker := itempresencechecker.NewItemPresenceChecker(rcPerformer)
	prodInfoGetter := productinfogetter.NewProductInfoGetter(rcPerformer)

	lomsConfig := a.config.LOMS
	lomsClient, err := lomsClientPkg.NewClient(lomsConfig.Address)
	if err != nil {
		log.Fatal(err)
	}
	loms := lomsProviderPkg.NewLOMSProvider(lomsClient)
	cartAdder := adder.New(cartRepo, itPresChecker, loms)
	cartDeleter := deleter.NewCartDeleter(cartRepo)
	cartClearer := clearer.NewCartClearer(cartRepo)
	cartLister := lister.New(cartRepo, prodInfoGetter)
	checkouterUsecase := checkouter.NewCheckouter(cartRepo, loms)
	wholeCartService := usecases.NewCartService(cartAdder, cartDeleter, cartClearer, cartLister, checkouterUsecase)
	mux := http.NewServeMux()
	addHandler := addPkg.New(wholeCartService)
	mux.Handle(fmt.Sprintf("POST /user/{%s}/cart/{%s}", addPkg.UserIdSegment, addPkg.SkuIdSegment), addHandler)
	clearHandler := clearPkg.New(wholeCartService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart", clearPkg.UserIdSegment), clearHandler)
	deleteHandler := deletePkg.New(wholeCartService)
	mux.Handle(fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", deletePkg.UserIdSegment, deletePkg.SkuIdSegment), deleteHandler)
	listHandler := listPkg.New(wholeCartService)
	mux.Handle(fmt.Sprintf("GET /user/{%s}/cart", listPkg.UserIdSegment), listHandler)
	checkoutHandler := checkoutPkg.New(wholeCartService)
	mux.Handle("POST /cart/checkout", checkoutHandler)
	probesMux := healthz.CreateMux()
	mux.Handle("GET /healthz/alive", probesMux)
	a.httpController = middleware.NewLogger(mux)
}

func (a *App) GetAddrFromConfig() (string, error) {
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
	if a.isRunning() {
		return
	}
	addr, err := a.GetAddrFromConfig()
	if err != nil {
		log.Fatal(err)
	}
	newServer := &http.Server{
		Addr:    addr,
		Handler: a.httpController,
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
	srvToShutdown := a.server.Load()
	if srvToShutdown == nil {
		return
	}
	timeout := time.Duration(a.config.Server.ShutdownTimoutSeconds) * time.Second
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()
	if err := srvToShutdown.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("HTTP shutdown error: %v", err)
	}
	a.server.Store(nil)
}

func (a *App) isRunning() bool {
	return a.server.Load() != nil
}
