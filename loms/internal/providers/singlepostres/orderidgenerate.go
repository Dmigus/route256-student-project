package singlepostres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/providers/orderidgenerator"
)

const selectLastOrderID = `SELECT COALESCE(max(id), 0) from "order"`

// CreateSequentialGenerator создаёт SequentialGenerator, учитывая максимальное значение id заказа из уже имеющихся в БД.
func CreateSequentialGenerator(ctx context.Context) (*orderidgenerator.SequentialGenerator, error) {
	tx := ctx.Value(trKey).(pgx.Tx)
	var maxOrderIDInRepo int64
	err := tx.QueryRow(ctx, selectLastOrderID).Scan(&maxOrderIDInRepo)
	if err != nil {
		return nil, fmt.Errorf("could not define last order id in repo: %w", err)
	}
	return orderidgenerator.NewSequentialGenerator(maxOrderIDInRepo + 1), nil
}
