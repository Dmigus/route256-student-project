package singlepostres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"route256.ozon.ru/project/loms/internal/providers/orderidgenerator"
)

const selectLastOrderId = `SELECT COALESCE(max(id), 0) from "order"`

func CreateSequentialGenerator(ctx context.Context) (*orderidgenerator.SequentialGenerator, error) {
	tx := ctx.Value(trKey).(pgx.Tx)
	var maxOrderIdInRepo int64
	err := tx.QueryRow(ctx, selectLastOrderId).Scan(&maxOrderIdInRepo)
	if err != nil {
		return nil, fmt.Errorf("could not define last order id in repo: %w", err)
	}
	return orderidgenerator.NewSequentialGenerator(maxOrderIdInRepo + 1), nil
}
