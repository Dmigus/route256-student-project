package stocks

import (
	"context"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
	"sync"
)

type InMemoryStockStorage struct {
	mu   sync.RWMutex
	data map[int64]*ItemUnits
}

func (i *InMemoryStockStorage) SetItemUnits(skuId int64, units *ItemUnits) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[skuId] = units
}

// Reserve резервирует товары. Если хотя бы один товар зарезервировать не удалось, то эффект такой же, как будто операции не было. Возвращает первую произошедшую ошибку.
func (i *InMemoryStockStorage) Reserve(ctx context.Context, items []models.OrderItem) error {
	succeededItems := make([]models.OrderItem, len(items))
	failed := false
	var err error
	for _, it := range items {
		err = i.reserveOne(it)
		if err != nil {
			failed = true
			break
		}
	}
	if !failed {
		return nil
	}
	_ = i.CancelReserved(ctx, succeededItems)
	return err
}

func (i *InMemoryStockStorage) reserveOne(it models.OrderItem) error {
	itemUni, err := i.getItemOrErr(it.SkuId)
	if err != nil {
		return err
	}
	return itemUni.reserve(it.Count)
}

// RemoveReserved удаляет элементы из имеющихся и зарезервированных. Товары не найденные в базе, игнорируются. Реализация никогда не возвращает шибку.
func (i *InMemoryStockStorage) RemoveReserved(_ context.Context, items []models.OrderItem) error {
	for _, it := range items {
		itemUni, err := i.getItemOrErr(it.SkuId)
		if err != nil {
			continue
		}
		itemUni.removeReserved(it.Count)
	}
	return nil
}

// GetNumOfAvailable возвращает
func (i *InMemoryStockStorage) GetNumOfAvailable(_ context.Context, skuId int64) (uint64, error) {
	itemUni, err := i.getItemOrErr(skuId)
	if err != nil {
		return 0, err
	}
	return itemUni.getNumOfAvailable(), nil
}

// CancelReserved отменяет резервирование для итемов items. Товары не найденные в базе, игнорируются. Реализация никогда не возвращает шибку.
func (i *InMemoryStockStorage) CancelReserved(_ context.Context, items []models.OrderItem) error {
	for _, it := range items {
		itemUni, err := i.getItemOrErr(it.SkuId)
		if err != nil {
			continue
		}
		itemUni.cancelReserve(it.Count)
	}
	return nil
}

func (i *InMemoryStockStorage) getItemOrErr(skuId int64) (*ItemUnits, error) {
	i.mu.RLock()
	defer i.mu.RUnlock()
	itemUni, exists := i.data[skuId]
	if !exists {
		return nil, fmt.Errorf("item with skuid %d is not exists", skuId)
	}
	return itemUni, nil
}
