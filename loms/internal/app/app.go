package app

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5/pgxpool"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	grpcContoller "route256.ozon.ru/project/loms/internal/controllers/grpc"
	mwGRPC "route256.ozon.ru/project/loms/internal/controllers/grpc/mw"
	"route256.ozon.ru/project/loms/internal/controllers/grpc/protoc/v1"
	httpContoller "route256.ozon.ru/project/loms/internal/controllers/http"
	"route256.ozon.ru/project/loms/internal/providers/inmemory"
	"route256.ozon.ru/project/loms/internal/providers/inmemory/orders"
	"route256.ozon.ru/project/loms/internal/providers/inmemory/orders/orderidgenerator"
	"route256.ozon.ru/project/loms/internal/providers/inmemory/stocks"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/modifier"
	"route256.ozon.ru/project/loms/internal/providers/singlepostgres/reader"
	"route256.ozon.ru/project/loms/internal/usecases"
	"route256.ozon.ru/project/loms/internal/usecases/orderscanceller"
	"route256.ozon.ru/project/loms/internal/usecases/orderscreator"
	"route256.ozon.ru/project/loms/internal/usecases/ordersgetter"
	"route256.ozon.ru/project/loms/internal/usecases/orderspayer"
	"route256.ozon.ru/project/loms/internal/usecases/stocksinfogetter"
	"sync/atomic"
	"time"
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
	var service *usecases.LOMService
	if a.config.Storage == nil {
		service = a.initServiceWithInMemory()
	} else {
		service = a.initServiceWithPostgres()
	}
	a.grpcController = grpcContoller.NewServer(service)
}

func (a *App) initServiceWithInMemory() *usecases.LOMService {
	idGenerator := orderidgenerator.NewSequentialGenerator(1)
	ordersRepo := orders.NewInMemoryOrdersStorage(idGenerator)
	stocksRepo := stocks.NewInMemoryStockStorage()
	err := fillStocksFromStockData(context.Background(), stocksRepo)
	if err != nil {
		log.Fatal(err)
	}
	canceller := orderscanceller.NewOrderCanceller(inmemory.NewTxManager[orderscanceller.OrderRepo, orderscanceller.StockRepo](ordersRepo, stocksRepo))
	creator := orderscreator.NewOrdersCreator(inmemory.NewTxManager[orderscreator.OrderRepo, orderscreator.StockRepo](ordersRepo, stocksRepo))
	getter := ordersgetter.NewOrdersGetter(inmemory.NewTxManager[ordersgetter.OrderRepo, any](ordersRepo, stocksRepo))
	payer := orderspayer.NewOrdersPayer(inmemory.NewTxManager[orderspayer.OrderRepo, orderspayer.StockRepo](ordersRepo, stocksRepo))
	stocksInfoGetter := stocksinfogetter.NewGetter(inmemory.NewTxManager[any, stocksinfogetter.StockRepo](ordersRepo, stocksRepo))
	return usecases.NewLOMService(
		creator,
		payer,
		stocksInfoGetter,
		getter,
		canceller,
	)
}

func (a *App) initServiceWithPostgres() *usecases.LOMService {
	connStr := a.config.getPostgresDSN()
	conn, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		log.Fatal(err)
	}
	err = conn.Ping(context.Background())
	if err != nil {
		log.Fatal(err)
	}

	stocksRepo := modifier.NewStocks(conn)
	ctxToInitStocks := context.Background()
	err = fillStocksFromStockData(ctxToInitStocks, stocksRepo)
	if err != nil {
		log.Fatal(err)
	}

	canceller := orderscanceller.NewOrderCanceller(singlepostgres.NewTxManager(conn,
		func(singlepostgres.TxBeginner) orderscanceller.OrderRepo {
			return modifier.NewOrders(conn)
		}, func(singlepostgres.TxBeginner) orderscanceller.StockRepo {
			return modifier.NewStocks(conn)
		}))
	creator := orderscreator.NewOrdersCreator(singlepostgres.NewTxManager(conn,
		func(singlepostgres.TxBeginner) orderscreator.OrderRepo {
			return modifier.NewOrders(conn)
		}, func(singlepostgres.TxBeginner) orderscreator.StockRepo {
			return modifier.NewStocks(conn)
		}))
	getter := ordersgetter.NewOrdersGetter(singlepostgres.NewTxManager(conn,
		func(singlepostgres.TxBeginner) ordersgetter.OrderRepo {
			return reader.NewOrders(conn)
		}, func(singlepostgres.TxBeginner) any {
			return nil
		}))
	payer := orderspayer.NewOrdersPayer(singlepostgres.NewTxManager(conn,
		func(singlepostgres.TxBeginner) orderspayer.OrderRepo {
			return modifier.NewOrders(conn)
		}, func(singlepostgres.TxBeginner) orderspayer.StockRepo {
			return modifier.NewStocks(conn)
		}))
	stocksInfoGetter := stocksinfogetter.NewGetter(singlepostgres.NewTxManager(conn,
		func(singlepostgres.TxBeginner) any {
			return nil
		}, func(singlepostgres.TxBeginner) stocksinfogetter.StockRepo {
			return reader.NewStocks(conn)
		}))
	return usecases.NewLOMService(
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
