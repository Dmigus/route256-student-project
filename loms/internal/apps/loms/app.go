// Package loms содержит приложение, в котором функционирует сервис loms
package loms

import (
	"context"
	"errors"
	"fmt"
	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/prometheus/client_golang/prometheus"
	"go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
	grpcContoller "route256.ozon.ru/project/loms/internal/controllers/grpc"
	mwGRPC "route256.ozon.ru/project/loms/internal/controllers/grpc/mw"
	httpContoller "route256.ozon.ru/project/loms/internal/controllers/http"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"route256.ozon.ru/project/loms/internal/pkg/sqlmetrics"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres"
	eventsToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/events"
	ordersToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/orders"
	stocksToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/stocks"
	ordersToRead "route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/orders"
	stocksToRead "route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/stocks"
	"sync"

	"time"

	"route256.ozon.ru/project/loms/internal/services/loms"
	"route256.ozon.ru/project/loms/internal/services/loms/orderscanceller"
	"route256.ozon.ru/project/loms/internal/services/loms/orderscreator"
	"route256.ozon.ru/project/loms/internal/services/loms/ordersgetter"
	"route256.ozon.ru/project/loms/internal/services/loms/orderspayer"
	"route256.ozon.ru/project/loms/internal/services/loms/stocksinfogetter"
)

type App struct {
	config           Config
	grpcController   *grpcContoller.Server
	grpcInterceptors []grpc.UnaryServerInterceptor
	httpController   *httpContoller.Server
}

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

func createConnToPostgres(dsn string) (*pgxpool.Pool, error) {
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("create connection pool: %w", err)
	}
	cfg.ConnConfig.Tracer = otelpgx.NewTracer(otelpgx.WithTrimSQLInSpanName())
	conn, err := pgxpool.NewWithConfig(context.Background(), cfg)
	if err != nil {
		return nil, fmt.Errorf("connect to database: %w", err)
	}
	return conn, nil
}

func (a *App) initServiceWithPostgres() (*loms.LOMService, error) {
	responseTime := prometheus.NewHistogramVec(prometheus.HistogramOpts{
		Namespace: "loms",
		Name:      "sql_request_duration_seconds",
		Help:      "Response time distribution made to PostgreSQL",
	},
		[]string{sqlmetrics.TableLabel, sqlmetrics.CategoryLabel, sqlmetrics.ErrLabel},
	)
	err := a.config.MetricsRegisterer.Register(responseTime)
	if err != nil {
		return nil, err
	}
	sqlDurationRecorder := sqlmetrics.NewSQLRequestDuration(responseTime)

	connMaster, err := createConnToPostgres(a.config.Storage.Master.GetPostgresDSN())
	if err != nil {
		return nil, err
	}
	// заполнение стоков начальными данными
	if err = fillStocksFromStockData(context.Background(), stocksToModify.NewStocks(connMaster, sqlDurationRecorder)); err != nil {
		return nil, err
	}

	connReplica, err := createConnToPostgres(a.config.Storage.Replica.GetPostgresDSN())
	if err != nil {
		return nil, err
	}
	canceller := orderscanceller.NewOrderCanceller(singlepostgres.NewTxManagerThree(connMaster,
		func(tx pgx.Tx) orderscanceller.OrderRepo {
			return ordersToModify.NewOrders(tx, sqlDurationRecorder)
		}, func(tx pgx.Tx) orderscanceller.StockRepo {
			return stocksToModify.NewStocks(tx, sqlDurationRecorder)
		}, func(tx pgx.Tx) orderscanceller.EventSender {
			return eventsToModify.NewEvents(tx, sqlDurationRecorder)
		}))
	creator := orderscreator.NewOrdersCreator(singlepostgres.NewTxManagerThree(connMaster,
		func(tx pgx.Tx) orderscreator.OrderRepo {
			return ordersToModify.NewOrders(tx, sqlDurationRecorder)
		}, func(tx pgx.Tx) orderscreator.StockRepo {
			return stocksToModify.NewStocks(tx, sqlDurationRecorder)
		}, func(tx pgx.Tx) orderscreator.EventSender {
			return eventsToModify.NewEvents(tx, sqlDurationRecorder)
		}))
	getter := ordersgetter.NewOrdersGetter(singlepostgres.NewTxManagerOne(connReplica,
		func(tx pgx.Tx) ordersgetter.OrderRepo {
			return ordersToRead.NewOrders(tx, sqlDurationRecorder)
		}))
	payer := orderspayer.NewOrdersPayer(singlepostgres.NewTxManagerThree(connMaster,
		func(tx pgx.Tx) orderspayer.OrderRepo {
			return ordersToModify.NewOrders(tx, sqlDurationRecorder)
		}, func(tx pgx.Tx) orderspayer.StockRepo {
			return stocksToModify.NewStocks(tx, sqlDurationRecorder)
		}, func(tx pgx.Tx) orderspayer.EventSender {
			return eventsToModify.NewEvents(tx, sqlDurationRecorder)
		}))
	stocksInfoGetter := stocksinfogetter.NewGetter(singlepostgres.NewTxManagerOne(connReplica,
		func(tx pgx.Tx) stocksinfogetter.StockRepo {
			return stocksToRead.NewStocks(tx, sqlDurationRecorder)
		}))
	return loms.NewLOMService(
		creator,
		payer,
		stocksInfoGetter,
		getter,
		canceller,
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
		Help:      "Response time distribution made to loms",
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
	go func() {
		defer wg.Done()
		err = grpcServer.Serve(lis)
		cancel()
	}()
	a.httpController = httpContoller.NewServer(
		fmt.Sprintf(":%d", a.config.GRPCServer.Port),
		fmt.Sprintf(":%d", a.config.HTTPGateway.Port),
		a.config.Swagger.Path,
		a.config.MetricsHandler)
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
	return errors.Join(err, errHTTP, errHTTPClose)
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
