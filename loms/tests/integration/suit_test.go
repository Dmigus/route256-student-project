//go:build integration

package integration

import (
	"context"
	"database/sql"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/pressly/goose/v3"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/modules/postgres"
	"github.com/testcontainers/testcontainers-go/wait"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"route256.ozon.ru/project/loms/internal/apps"
	"route256.ozon.ru/project/loms/internal/apps/loms"
	"route256.ozon.ru/project/loms/internal/pkg/api/loms/v1"
	"strconv"
	"time"
)

const (
	migrationsDir = "../../migrations"
	configDir     = "./config/testconfig.json"
)

type Suite struct {
	suite.Suite
	Pg       *postgres.PostgresContainer
	app      *loms.App
	ConnToDB *sql.DB
	client   v1.LOMServiceClient
}

func (s *Suite) SetupSuite() {
	config, err := apps.NewConfig[loms.Config](configDir)
	s.Require().NoError(err)
	ctx := context.Background()
	postgresContainer, err := postgres.RunContainer(ctx,
		testcontainers.WithImage("docker.io/postgres:16.2-bullseye"),
		postgres.WithDatabase(config.Storage.Master.Database),
		postgres.WithUsername(config.Storage.Master.User),
		postgres.WithPassword(config.Storage.Master.Password),
		testcontainers.WithWaitStrategy(
			wait.ForLog("database system is ready to accept connections").
				WithOccurrence(2).
				WithStartupTimeout(5*time.Second)),
	)
	s.Require().NoError(err)

	dsn, err := postgresContainer.ConnectionString(ctx, "sslmode=disable")
	s.Require().NoError(err)
	s.ConnToDB, err = sql.Open("pgx", dsn)
	s.Require().NoError(err)
	err = migrateUp(ctx, s.ConnToDB)
	s.Require().NoError(err)

	// app
	s.Pg = postgresContainer
	host, err := postgresContainer.Host(ctx)
	config.Storage.Master.Host = host
	config.Storage.Replica.Host = host
	exposedTcpPort, err := postgresContainer.MappedPort(ctx, "5432")
	s.Require().NoError(err)
	port := exposedTcpPort.Int()
	config.Storage.Master.Port = uint16(port)
	config.Storage.Replica.Port = uint16(port)
	s.app = loms.NewApp(config)
	go s.app.Run()
	conn, err := grpc.Dial(":"+strconv.Itoa(int(config.GRPCServer.Port)), grpc.WithTransportCredentials(insecure.NewCredentials()))
	s.Require().NoError(err)
	conn.Connect()
	s.client = v1.NewLOMServiceClient(conn)
}

func migrateUp(ctx context.Context, conn *sql.DB) error {
	return goose.UpContext(ctx, conn, migrationsDir)
}

func migrateDown(ctx context.Context, conn *sql.DB) error {
	return goose.DownContext(ctx, conn, migrationsDir)
}

func (s *Suite) TearDownSuite() {
	ctx := context.Background()
	err := migrateDown(ctx, s.ConnToDB)
	s.Assert().NoError(err)
	err = s.ConnToDB.Close()
	s.Assert().NoError(err)
	s.app.Stop()
	err = s.Pg.Terminate(ctx)
	s.Assert().NoError(err)
}

func (s *Suite) SetupTest() {
	_, err := s.ConnToDB.Exec("TRUNCATE item_unit;")
	s.Assert().NoError(err)
	_, err = s.ConnToDB.Exec("TRUNCATE \"order\";")
	s.Assert().NoError(err)
	_, err = s.ConnToDB.Exec("TRUNCATE order_item;")
	s.Assert().NoError(err)
}

func (s *Suite) TestOrderCreate() {
	ctx := context.Background()
	_, err := s.ConnToDB.ExecContext(ctx, "INSERT INTO item_unit(sku_id, total, reserved) VALUES ($1, $2, $3)", 773297411, 150, 10)
	s.Require().NoError(err)

	items := []*v1.Item{{Sku: 773297411, Count: 50}}
	req := &v1.OrderCreateRequest{User: 123, Items: items}
	ordId, err := s.client.OrderCreate(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(int64(1), ordId.OrderID)

	orderRow := s.ConnToDB.QueryRowContext(ctx, "SELECT status, are_items_reserved FROM \"order\";")
	var orderStatus string
	var reserved bool
	err = orderRow.Scan(&orderStatus, &reserved)
	s.Require().NoError(err)
	s.Assert().Equal("AwaitingPayment", orderStatus)
	s.Assert().True(reserved)

	ordItRow := s.ConnToDB.QueryRowContext(ctx, "SELECT sku_id, count FROM order_item WHERE order_id = $1;", ordId.OrderID)
	var savedSku int
	var itCnt int
	err = ordItRow.Scan(&savedSku, &itCnt)
	s.Require().NoError(err)
	s.Assert().Equal(773297411, savedSku)
	s.Assert().Equal(50, itCnt)
}

func (s *Suite) TestOrderGet() {
	ctx := context.Background()
	_, err := s.ConnToDB.ExecContext(ctx, "INSERT INTO \"order\"(id, user_id, status, are_items_reserved) VALUES ($1, $2, $3, $4);", 123, 456, "Payed", false)
	s.Require().NoError(err)
	_, err = s.ConnToDB.ExecContext(ctx, "INSERT INTO order_item(order_id, sku_id, count) VALUES ($1, $2, $3)", 123, 10, 20)
	s.Require().NoError(err)

	req := &v1.OrderId{OrderID: 123}
	order, err := s.client.OrderInfo(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(int64(456), order.User)
	s.Require().Len(order.Items, 1)
	s.Assert().Equal(uint32(10), order.Items[0].Sku)
}

func (s *Suite) TestOrderPay() {
	ctx := context.Background()
	_, err := s.ConnToDB.ExecContext(ctx, "INSERT INTO \"order\"(id, user_id, status, are_items_reserved) VALUES ($1, $2, $3, $4);", 123, 456, "AwaitingPayment", true)
	s.Require().NoError(err)

	_, err = s.ConnToDB.ExecContext(ctx, "INSERT INTO order_item(order_id, sku_id, count) VALUES ($1, $2, $3)", 123, 1005, 50)
	s.Require().NoError(err)

	_, err = s.ConnToDB.ExecContext(ctx, "INSERT INTO item_unit(sku_id, total, reserved) VALUES ($1, $2, $3)", 1005, 350, 50)
	s.Require().NoError(err)

	req := &v1.OrderId{OrderID: 123}
	_, err = s.client.OrderPay(ctx, req)
	s.Require().NoError(err)

	orderRow := s.ConnToDB.QueryRowContext(ctx, "SELECT status FROM \"order\" WHERE id = $1;", 123)
	var newStatus string
	err = orderRow.Scan(&newStatus)
	s.Require().NoError(err)
	s.Assert().Equal("Payed", newStatus)

	orderRow = s.ConnToDB.QueryRowContext(ctx, "SELECT total, reserved FROM item_unit WHERE sku_id = $1;", 1005)
	var total, reserved int
	err = orderRow.Scan(&total, &reserved)
	s.Require().NoError(err)
	s.Assert().Equal(300, total)
	s.Assert().Equal(0, reserved)
}

func (s *Suite) TestOrderCancel() {
	ctx := context.Background()
	_, err := s.ConnToDB.ExecContext(ctx, "INSERT INTO \"order\"(id, user_id, status, are_items_reserved) VALUES ($1, $2, $3, $4);", 123, 456, "Payed", false)
	s.Require().NoError(err)

	_, err = s.ConnToDB.ExecContext(ctx, "INSERT INTO order_item(order_id, sku_id, count) VALUES ($1, $2, $3)", 123, 1005, 50)
	s.Require().NoError(err)

	_, err = s.ConnToDB.ExecContext(ctx, "INSERT INTO item_unit(sku_id, total, reserved) VALUES ($1, $2, $3)", 1005, 0, 0)
	s.Require().NoError(err)

	req := &v1.OrderId{OrderID: 123}
	_, err = s.client.OrderCancel(ctx, req)
	s.Require().NoError(err)

	orderRow := s.ConnToDB.QueryRowContext(ctx, "SELECT status FROM \"order\" WHERE id = $1;", 123)
	var newStatus string
	err = orderRow.Scan(&newStatus)
	s.Require().NoError(err)
	s.Assert().Equal("Cancelled", newStatus)

	orderRow = s.ConnToDB.QueryRowContext(ctx, "SELECT total FROM item_unit WHERE sku_id = $1;", 1005)
	var total int
	err = orderRow.Scan(&total)
	s.Require().NoError(err)
	s.Assert().Equal(50, total)
}

func (s *Suite) TestStockGet() {
	ctx := context.Background()
	_, err := s.ConnToDB.ExecContext(ctx, "INSERT INTO item_unit(sku_id, total, reserved) VALUES ($1, $2, $3)", 1005, 40, 15)
	s.Require().NoError(err)

	req := &v1.StocksInfoRequest{Sku: 1005}
	resp, err := s.client.StocksInfo(ctx, req)
	s.Require().NoError(err)
	s.Assert().Equal(uint64(25), resp.Count)
}
