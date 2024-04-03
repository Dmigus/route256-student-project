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
	"route256.ozon.ru/project/cart/internal/clients/ratelimiterhttp"
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
	"time"
)

type App struct {
	config         Config
	httpController http.Handler
	rateLimiter    *ratelimiterhttp.RateLimiter
}

// NewApp возращает App, готовый к запуску, либо ошибку
func NewApp(config Config) (*App, error) {
	app := &App{
		config: config,
	}
	if err := app.init(); err != nil {
		return nil, err
	}
	return app, nil
}

func (a *App) init() error {
	cartRepo := repository.New()
	prodServConfig := a.config.ProductService
	baseUrl, err := url.Parse(prodServConfig.BaseURL)
	if err != nil {
		return err
	}
	retryPolicy := policies.NewRetryOnStatusCodes(prodServConfig.RetryPolicy.RetryStatusCodes, prodServConfig.RetryPolicy.MaxRetries)
	rateLimiter := ratelimiterhttp.NewRateLimiter(prodServConfig.RPS)
	a.rateLimiter = rateLimiter
	rateLimitedTripper := ratelimiterhttp.NewRateLimitedTripper(rateLimiter, http.DefaultTransport)
	clientForProductService := &http.Client{Transport: retriablehttp.NewRetryRoundTripper(rateLimitedTripper, retryPolicy)}
	rcPerformer := productservice.NewRCPerformer(clientForProductService, baseUrl, prodServConfig.AccessToken)
	itPresChecker := itempresencechecker.NewItemPresenceChecker(rcPerformer)
	prodInfoGetter := productinfogetter.NewProductInfoGetter(rcPerformer)

	lomsConfig := a.config.LOMS
	lomsClient, err := lomsClientPkg.NewClient(lomsConfig.Address)
	if err != nil {
		return err
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
	return nil
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
// Возвращает критические ошибки, произошедшие при работе http сервера и его остановке.
// Сервер начнёт прекращение своей работы, когда переданный контекст ctx будет отменён
func (a *App) Run(ctx context.Context) error {
	a.rateLimiter.Run(ctx)
	addr, err := a.GetAddrFromConfig()
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr:    addr,
		Handler: a.httpController,
	}
	errStopping := make(chan error, 1)
	go func() {
		<-ctx.Done()
		errStopping <- a.stop(server)
		close(errStopping)
	}()
	errServing := server.ListenAndServe()
	if errors.Is(errServing, http.ErrServerClosed) {
		errServing = nil
	}
	return errors.Join(errServing, <-errStopping)
}

// stop пытается произвести graceful shotdown server в течение ShutdownTimoutSeconds секунд. Если не удалось
// остановить в течение таймаута, то сервер заверщается немедленнло.
func (a *App) stop(server *http.Server) error {
	timeout := time.Duration(a.config.Server.ShutdownTimoutSeconds) * time.Second
	shutdownCtx, shutdownRelease := context.WithTimeout(context.Background(), timeout)
	defer shutdownRelease()
	err := server.Shutdown(shutdownCtx)
	if err != nil {
		errClosing := server.Close()
		err = errors.Join(err, errClosing)
	}
	return err
}
