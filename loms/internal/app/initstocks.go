package app

import (
	"bytes"
	"context"
	"encoding/json"
)

type stocksRepoToInit interface {
	SetItemUnits(ctx context.Context, skuID int64, total, reserved uint64) error
}

func fillStocksFromStockData(ctx context.Context, stocksRepo stocksRepoToInit) error {
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
		err := stocksRepo.SetItemUnits(ctx, it.Sku, it.TotalCount, it.Reserved)
		if err != nil {
			return err
		}
	}
	return nil
}
