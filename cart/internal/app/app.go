package app

import (
	"context"
	"errors"
	"fmt"
	"github.com/jasongerard/healthz"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"log"
	"net/http"
	"net/netip"
	"net/url"
	"route256.ozon.ru/project/cart/internal/clients/durationobserverhttp"
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
	"route256.ozon.ru/project/cart/internal/pkg/ratelimiter"
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

	psResponseTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "http_product_service_duration_seconds",
		Help:      "Response time distribution made to Product service by cart service",
	},
		[]string{durationobserverhttp.MethodNameLabel, durationobserverhttp.CodeLabel, durationobserverhttp.UrlLabel},
	)
	observerTripper := durationobserverhttp.NewDurationObserverTripper(psResponseTime, http.DefaultTransport)

	retryPolicy := policies.NewRetryOnStatusCodes(prodServConfig.RetryPolicy.RetryStatusCodes, prodServConfig.RetryPolicy.MaxRetries)
	ticker := ratelimiter.NewSystemTimeTicker(int64(prodServConfig.RPS))
	rateLimiter := ratelimiter.NewRateLimiter(prodServConfig.RPS, ticker)
	rateLimitedTripper := ratelimiterhttp.NewRateLimitedTripper(rateLimiter, observerTripper)

	clientForProductService := &http.Client{Transport: retriablehttp.NewRetryRoundTripper(rateLimitedTripper, retryPolicy)}
	rcPerformer := productservice.NewRCPerformer(clientForProductService, baseUrl, prodServConfig.AccessToken)
	itPresChecker := itempresencechecker.NewItemPresenceChecker(rcPerformer)
	prodInfoGetter := productinfogetter.NewProductInfoGetter(rcPerformer)

	lomsConfig := a.config.LOMS
	lomsResponseTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "grpc_loms_duration_seconds",
		Help:      "Response time distribution made to loms by cart service",
	},
		[]string{lomsClientPkg.MethodNameLabel, lomsClientPkg.CodeLabel},
	)
	lomsClient, err := lomsClientPkg.NewClient(lomsConfig.Address, lomsResponseTime)
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
	responseTime := promauto.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "cart",
		Name:      "http_duration_seconds",
		Help:      "Response time distribution made to cart",
	},
		[]string{middleware.MethodNameLabel, middleware.CodeLabel, middleware.UrlLabel},
	)
	addHandler := middleware.NewDurationObserverMW(addPkg.New(wholeCartService), responseTime, "/user/<user_id>/cart/<cart_id>")
	addPattern := fmt.Sprintf("POST /user/{%s}/cart/{%s}", addPkg.UserIdSegment, addPkg.SkuIdSegment)
	mux.Handle(addPattern, addHandler)
	clearHandler := middleware.NewDurationObserverMW(clearPkg.New(wholeCartService), responseTime, "/user/<user_id>/cart")
	clearPattern := fmt.Sprintf("DELETE /user/{%s}/cart", clearPkg.UserIdSegment)
	mux.Handle(clearPattern, clearHandler)
	deleteHandler := middleware.NewDurationObserverMW(deletePkg.New(wholeCartService), responseTime, "/user/<user_id>/cart/<cart_id>")
	deletePattern := fmt.Sprintf("DELETE /user/{%s}/cart/{%s}", deletePkg.UserIdSegment, deletePkg.SkuIdSegment)
	mux.Handle(deletePattern, deleteHandler)
	listHandler := middleware.NewDurationObserverMW(listPkg.New(wholeCartService), responseTime, "/user/<user_id>/cart")
	listPattern := fmt.Sprintf("GET /user/{%s}/cart", listPkg.UserIdSegment)
	mux.Handle(listPattern, listHandler)
	checkoutHandler := middleware.NewDurationObserverMW(checkoutPkg.New(wholeCartService), responseTime, "/cart/checkout")
	checkoutPattern := "POST /cart/checkout"
	mux.Handle(checkoutPattern, checkoutHandler)

	probesMux := middleware.NewDurationObserverMW(healthz.CreateMux(), responseTime, "/healthz/alive")
	probesPattern := "GET /healthz/alive"
	mux.Handle(probesPattern, probesMux)
	metricsPattern := "/metrics"
	metricsHandler := middleware.NewDurationObserverMW(promhttp.Handler(), responseTime, "/metrics")
	mux.Handle(metricsPattern, metricsHandler)
	a.httpController = otelhttp.NewHandler(middleware.NewLogger(mux), "processing request by cart")
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
	addr, err := a.GetAddrFromConfig()
	if err != nil {
		log.Fatal(err)
	}
	server := &http.Server{
		Addr:    addr,
		Handler: a.httpController,
	}
	errStopping := make(chan error, 1)
	appRunCtx, cancel := context.WithCancel(ctx)
	go func() {
		<-appRunCtx.Done()
		errStopping <- a.stop(server)
		close(errStopping)
	}()
	errServing := server.ListenAndServe()
	cancel()
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
