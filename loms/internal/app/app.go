package app

import (
	"bytes"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	grpcContoller "route256.ozon.ru/project/loms/internal/controllers/grpc"
	mwGRPC "route256.ozon.ru/project/loms/internal/controllers/grpc/mw"
	"route256.ozon.ru/project/loms/internal/controllers/grpc/protoc/v1"
	httpContoller "route256.ozon.ru/project/loms/internal/controllers/http"
	"route256.ozon.ru/project/loms/internal/providers/orderidgenerator"
	"route256.ozon.ru/project/loms/internal/providers/orders"
	"route256.ozon.ru/project/loms/internal/providers/stocks"
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
	ordersRepo := orders.NewInMemoryOrdersStorage()
	stocksRepo := stocks.NewInMemoryStockStorage()
	err := fillStocksFromStockData(stocksRepo)
	if err != nil {
		log.Fatal(err)
	}

	canceller := orderscanceller.NewOrderCanceller(ordersRepo, stocksRepo)
	idGenerator := orderidgenerator.NewSequentialGenerator()
	creator := orderscreator.NewOrdersCreator(idGenerator, ordersRepo, stocksRepo)
	getter := ordersgetter.NewOrdersGetter(ordersRepo)
	payer := orderspayer.NewOrdersPayer(ordersRepo, stocksRepo)
	stocksInfoGetter := stocksinfogetter.NewGetter(stocksRepo)
	wholeService := usecases.NewLOMService(
		creator,
		payer,
		stocksInfoGetter,
		getter,
		canceller,
	)
	a.grpcController = grpcContoller.NewServer(wholeService)
}

func fillStocksFromStockData(stocksRepo *stocks.InMemoryStockStorage) error {
	reader := bytes.NewReader(stockdata)
	jsonParser := json.NewDecoder(reader)
	var items []struct {
		Sku        int64  `json:"sku"`
		TotalCount uint64 `json:"total_count"`
		Reserved   uint64 `json:"reserved"`
	}
	if err := jsonParser.Decode(&items); err != nil {
		return err
	}
	for _, it := range items {
		itemUnits := stocks.NewItemUnits(it.TotalCount, it.Reserved)
		stocksRepo.SetItemUnits(it.Sku, itemUnits)
	}
	return nil
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
