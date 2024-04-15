// Package loms содержит приложение, в котором функционирует сервис loms
package loms

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	grpcContoller "route256.ozon.ru/project/loms/internal/controllers/grpc"
	mwGRPC "route256.ozon.ru/project/loms/internal/controllers/grpc/mw"
	httpContoller "route256.ozon.ru/project/loms/internal/controllers/http"
	v1 "route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres"
	eventsToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/events"
	ordersToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/orders"
	stocksToModify "route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifiers/stocks"
	ordersToRead "route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/orders"
	stocksToRead "route256.ozon.ru/project/loms/internal/providers/singlepostgres/readers/stocks"

	"sync/atomic"
	"time"

	"route256.ozon.ru/project/loms/internal/services/loms"
	"route256.ozon.ru/project/loms/internal/services/loms/orderscanceller"
	"route256.ozon.ru/project/loms/internal/services/loms/orderscreator"
	"route256.ozon.ru/project/loms/internal/services/loms/ordersgetter"
	"route256.ozon.ru/project/loms/internal/services/loms/orderspayer"
	"route256.ozon.ru/project/loms/internal/services/loms/stocksinfogetter"
)

type App struct {
	config         Config
	grpcController *grpcContoller.Server
	grpcServer     atomic.Pointer[grpc.Server]
	httpGateway    atomic.Pointer[httpContoller.Server]
}

func NewApp(config Config) *App {
	app := &App{
		config: config,
	}
	app.init()
	return app
}

func (a *App) init() {
	service := a.initServiceWithPostgres()
	a.grpcController = grpcContoller.NewServer(service)
}

func createConnToPostgres(dsn string) *pgxpool.Pool {
	conn, err := pgxpool.New(context.Background(), dsn)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}
	return conn
}

func (a *App) initServiceWithPostgres() *loms.LOMService {
	connMaster := createConnToPostgres(a.config.Storage.Master.GetPostgresDSN())
	if err := fillStocksFromStockData(context.Background(), stocksToModify.NewStocks(connMaster)); err != nil {
		log.Fatal(err)
	}

	connReplica := createConnToPostgres(a.config.Storage.Replica.GetPostgresDSN())
	canceller := orderscanceller.NewOrderCanceller(singlepostgres.NewTxManagerThree(connMaster,
		func(tx pgx.Tx) orderscanceller.OrderRepo {
			return ordersToModify.NewOrders(tx)
		}, func(tx pgx.Tx) orderscanceller.StockRepo {
			return stocksToModify.NewStocks(tx)
		}, func(tx pgx.Tx) orderscanceller.EventSender {
			return eventsToModify.NewEvents(tx)
		}))
	creator := orderscreator.NewOrdersCreator(singlepostgres.NewTxManagerThree(connMaster,
		func(tx pgx.Tx) orderscreator.OrderRepo {
			return ordersToModify.NewOrders(tx)
		}, func(tx pgx.Tx) orderscreator.StockRepo {
			return stocksToModify.NewStocks(tx)
		}, func(tx pgx.Tx) orderscreator.EventSender {
			return eventsToModify.NewEvents(tx)
		}))
	getter := ordersgetter.NewOrdersGetter(singlepostgres.NewTxManagerOne(connReplica,
		func(tx pgx.Tx) ordersgetter.OrderRepo {
			return ordersToRead.NewOrders(tx)
		}))
	payer := orderspayer.NewOrdersPayer(singlepostgres.NewTxManagerThree(connMaster,
		func(tx pgx.Tx) orderspayer.OrderRepo {
			return ordersToModify.NewOrders(tx)
		}, func(tx pgx.Tx) orderspayer.StockRepo {
			return stocksToModify.NewStocks(tx)
		}, func(tx pgx.Tx) orderspayer.EventSender {
			return eventsToModify.NewEvents(tx)
		}))
	stocksInfoGetter := stocksinfogetter.NewGetter(singlepostgres.NewTxManagerOne(connReplica,
		func(tx pgx.Tx) stocksinfogetter.StockRepo {
			return stocksToRead.NewStocks(tx)
		}))
	return loms.NewLOMService(
		creator,
		payer,
		stocksInfoGetter,
		getter,
		canceller,
	)
}

// Run представляет из себя блокирующий вызов, который запускает новый сервер, согласно текущей конфигурации.
// Если он уже запущен, то функция ничего не делает. Если не удалось запустить, вся программа завершается с ошибкой
func (a *App) Run() {
	if a.grpcServer.Load() != nil {
		return
	}
	grpcServer := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			mwGRPC.SetUpErrorCode,
			mwGRPC.LogReqAndResp,
			mwGRPC.RecoverPanic,
			mwGRPC.Validate,
		),
	)
	if !a.grpcServer.CompareAndSwap(nil, grpcServer) {
		return
	}
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.config.GRPCServer.Port))
	if err != nil {
		panic(err)
	}
	reflection.Register(grpcServer)
	v1.RegisterLOMServiceServer(grpcServer, a.grpcController)
	if err = grpcServer.Serve(lis); err != nil {
		log.Fatal(err)
	}
}

func (a *App) RunGateway() {
	if a.httpGateway.Load() != nil {
		return
	}
	newGW := httpContoller.NewServer(
		fmt.Sprintf(":%d", a.config.GRPCServer.Port),
		fmt.Sprintf(":%d", a.config.HTTPGateway.Port),
		a.config.Swagger.Path)
	if !a.httpGateway.CompareAndSwap(nil, newGW) {
		return
	}
	newGW.Serve()
}

// Stop останавливает запущенный сервер в течение ShutdownTimoutSeconds секунд. Если не был запущен, функция ничего не делает. Если не удалось
// остановить в течение таймаута, вся программа завершается с ошибкой. Возврат из функции произойдёт, когда shutdown завершится.
func (a *App) Stop() {
	srvToShutdown := a.grpcServer.Load()
	if srvToShutdown == nil {
		return
	}
	stopped := make(chan struct{})
	go func() {
		srvToShutdown.GracefulStop()
		close(stopped)
	}()
	timerToForceStop := time.NewTimer(time.Duration(a.config.GRPCServer.ShutdownTimoutSeconds) * time.Second)
	select {
	case <-timerToForceStop.C:
		srvToShutdown.Stop()
	case <-stopped:
		timerToForceStop.Stop()
	}
	a.grpcServer.Store(nil)
}

func (a *App) StopGateway() {
	gwToStop := a.httpGateway.Load()
	if gwToStop == nil {
		return
	}
	gwToStop.Stop(time.Duration(a.config.HTTPGateway.ShutdownTimoutSeconds))
	a.httpGateway.Store(nil)
}
