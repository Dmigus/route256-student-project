// Package loms содержит приложение, в котором функционирует сервис loms
package loms

import (
	"context"
	"errors"
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	"route256.ozon.ru/project/loms/internal/providers/multipostgres"
	"route256.ozon.ru/project/loms/internal/services/loms/allordersgetter"
	"sync"
	"time"

	grpcContoller "route256.ozon.ru/project/loms/internal/controllers/grpc"
	mwGRPC "route256.ozon.ru/project/loms/internal/controllers/grpc/mw"
	httpContoller "route256.ozon.ru/project/loms/internal/controllers/http"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	eventsToModify "route256.ozon.ru/project/loms/internal/providers/multipostgres/modifiers/events"
	ordersToModify "route256.ozon.ru/project/loms/internal/providers/multipostgres/modifiers/orders"
	stocksToModify "route256.ozon.ru/project/loms/internal/providers/multipostgres/modifiers/stocks"
	ordersToRead "route256.ozon.ru/project/loms/internal/providers/multipostgres/readers/orders"
	stocksToRead "route256.ozon.ru/project/loms/internal/providers/multipostgres/readers/stocks"
	singleStocksToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/stocks"
	"route256.ozon.ru/project/loms/internal/services/loms"
	"route256.ozon.ru/project/loms/internal/services/loms/orderscanceller"
	"route256.ozon.ru/project/loms/internal/services/loms/orderscreator"
	"route256.ozon.ru/project/loms/internal/services/loms/ordersgetter"
	"route256.ozon.ru/project/loms/internal/services/loms/orderspayer"
	"route256.ozon.ru/project/loms/internal/services/loms/stocksinfogetter"
)

var bucketsForRequestDuration = []float64{0.001, 0.005, 0.01, 0.02, 0.05, 0.1, 0.2, 0.5, 1}

type App struct {
	config           Config
	grpcController   *grpcContoller.Server
	grpcInterceptors []grpc.UnaryServerInterceptor
	httpController   *httpContoller.Server
}

// NewApp возвращает иннициализарованный App, готовый к запуску
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
	service, err := a.initServiceWithPostgres()
	if err != nil {
		return err
	}
	a.grpcController = grpcContoller.NewServer(service)
	err = a.initInterceptors()
	if err != nil {
		return err
	}
	return nil
}

func (a *App) initServiceWithPostgres() (*loms.LOMService, error) {
	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "sql_request_duration_seconds",
		Help:      "Response time distribution made to PostgreSQL",
		Buckets:   bucketsForRequestDuration,
	},
		[]string{sqlmetrics.TableLabel, sqlmetrics.CategoryLabel, sqlmetrics.ErrLabel},
	)
	err := a.config.MetricsRegisterer.Register(responseTime)
	if err != nil {
		return nil, err
	}
	sqlDurationRecorder := sqlmetrics.NewSQLRequestDuration(responseTime)

	shardManager, err := newShardManager(a.config.Storages)
	if err != nil {
		return nil, err
	}

	defaultShard := shardManager.GetDefaultShard()
	// заполнение стоков начальными данными
	if err = fillStocksFromStockData(context.Background(), singleStocksToModify.NewStocks(defaultShard.Master(), sqlDurationRecorder)); err != nil {
		return nil, err
	}

	canceller := orderscanceller.NewOrderCanceller(multipostgres.NewTxManager3(
		func(tc multipostgres.TransactionCreator) orderscanceller.OrderRepo {
			return ordersToModify.NewOrders(tc, *shardManager, sqlDurationRecorder)
		}, func(tc multipostgres.TransactionCreator) orderscanceller.StockRepo {
			return stocksToModify.NewStocks(tc, *shardManager, sqlDurationRecorder)
		}, func(tc multipostgres.TransactionCreator) orderscanceller.EventSender {
			return eventsToModify.NewEventsToInsert(tc, *shardManager, sqlDurationRecorder)
		}))
	creator := orderscreator.NewOrdersCreator(multipostgres.NewTxManager3(
		func(tc multipostgres.TransactionCreator) orderscreator.OrderRepo {
			return ordersToModify.NewOrders(tc, *shardManager, sqlDurationRecorder)
		}, func(tc multipostgres.TransactionCreator) orderscreator.StockRepo {
			return stocksToModify.NewStocks(tc, *shardManager, sqlDurationRecorder)
		}, func(tc multipostgres.TransactionCreator) orderscreator.EventSender {
			return eventsToModify.NewEventsToInsert(tc, *shardManager, sqlDurationRecorder)
		}))
	getter := ordersgetter.NewOrdersGetter(multipostgres.NewTxManager1(
		func(tc multipostgres.TransactionCreator) ordersgetter.OrderRepo {
			return ordersToRead.NewOrders(tc, *shardManager, sqlDurationRecorder)
		}))
	allOrdersGetter := allordersgetter.NewOrdersGetter(multipostgres.NewTxManager1(
		func(tc multipostgres.TransactionCreator) allordersgetter.OrderRepo {
			return ordersToRead.NewOrders(tc, *shardManager, sqlDurationRecorder)
		}))
	payer := orderspayer.NewOrdersPayer(multipostgres.NewTxManager3(
		func(tc multipostgres.TransactionCreator) orderspayer.OrderRepo {
			return ordersToModify.NewOrders(tc, *shardManager, sqlDurationRecorder)
		}, func(tc multipostgres.TransactionCreator) orderspayer.StockRepo {
			return stocksToModify.NewStocks(tc, *shardManager, sqlDurationRecorder)
		}, func(tc multipostgres.TransactionCreator) orderspayer.EventSender {
			return eventsToModify.NewEventsToInsert(tc, *shardManager, sqlDurationRecorder)
		}))
	stocksInfoGetter := stocksinfogetter.NewGetter(multipostgres.NewTxManager1(
		func(tc multipostgres.TransactionCreator) stocksinfogetter.StockRepo {
			return stocksToRead.NewStocks(tc, *shardManager, sqlDurationRecorder)
		}))
	return loms.NewLOMService(
		creator,
		payer,
		stocksInfoGetter,
		getter,
		canceller,
		allOrdersGetter,
	), nil
}

func (a *App) initInterceptors() error {
	a.grpcInterceptors = append(a.grpcInterceptors, mwGRPC.SetUpErrorCode)
	mwLogger := mwGRPC.NewLoggerMW(a.config.Logger)
	a.grpcInterceptors = append(a.grpcInterceptors, mwLogger.LogReqAndResp)
	a.grpcInterceptors = append(a.grpcInterceptors, mwGRPC.RecoverPanic)
	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "grpc_request_duration_seconds",
		Help:      "Response time distribution made to loms. Example: Median of all queries duration histogram_quantile(0.5, loms_grpc_request_duration_seconds_bucket)",
		Buckets:   bucketsForRequestDuration,
	},
		[]string{mwGRPC.MethodNameLabel, mwGRPC.CodeLabel},
	)
	err := a.config.MetricsRegisterer.Register(responseTime)
	if err != nil {
		return err
	}
	a.grpcInterceptors = append(a.grpcInterceptors, mwGRPC.NewRequestDurationInterceptor(responseTime).RecordDuration)
	a.grpcInterceptors = append(a.grpcInterceptors, mwGRPC.Validate)
	return nil
}

// Run представляет из себя блокирующий вызов, который запускает новый grpc и http контроллеры
func (a *App) Run(ctx context.Context) error {
	grpcServer := grpc.NewServer(
		grpc.StatsHandler(otelgrpc.NewServerHandler()),
		grpc.ChainUnaryInterceptor(a.grpcInterceptors...),
	)
	reflection.Register(grpcServer)
	v1.RegisterLOMServiceServer(grpcServer, a.grpcController)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPCServer.Port))
	if err != nil {
		return err
	}
	wg := &sync.WaitGroup{}
	wg.Add(4)
	serverRunCtx, cancel := context.WithCancel(ctx)
	defer cancel()
	go func() {
		defer wg.Done()
		<-serverRunCtx.Done()
		a.stop(grpcServer)
	}()
	var errGRPC error
	go func() {
		defer wg.Done()
		errGRPC = grpcServer.Serve(lis)
		cancel()
	}()
	a.httpController, err = httpContoller.NewServer(
		fmt.Sprintf(":%d", a.config.GRPCServer.Port),
		fmt.Sprintf(":%d", a.config.HTTPGateway.Port),
		a.config.Swagger.Path,
		a.config.MetricsHandler)
	if err != nil {
		return err
	}
	var errHTTP, errHTTPClose error
	go func() {
		defer wg.Done()
		<-serverRunCtx.Done()
		errHTTPClose = a.stopGateway(a.httpController)
	}()
	go func() {
		defer wg.Done()
		errHTTP = a.httpController.Serve()
		cancel()
	}()
	wg.Wait()
	return errors.Join(errGRPC, errHTTP, errHTTPClose)
}

// stop останавливает запущенный grpc сервер в течение ShutdownTimoutSeconds секунд.
func (a *App) stop(server *grpc.Server) {
	stopped := make(chan struct{})
	go func() {
		server.GracefulStop()
		close(stopped)
	}()
	timerToForceStop := time.NewTimer(time.Duration(a.config.GRPCServer.ShutdownTimoutSeconds) * time.Second)
	select {
	case <-timerToForceStop.C:
		server.Stop()
	case <-stopped:
		timerToForceStop.Stop()
	}
}

func (a *App) stopGateway(gwToStop *httpContoller.Server) error {
	return gwToStop.Stop(time.Duration(a.config.HTTPGateway.ShutdownTimoutSeconds))
}
