package stocks

import (
	"context"
	"errors"
	"fmt"
	"route256.ozon.ru/project/loms/internal/models"
	"sync"
)

var ErrItemIsNotExists = errors.New("item is not exists")

type InMemoryStockStorage struct {
	mu   sync.RWMutex
	data map[int64]*ItemUnits
}

func NewInMemoryStockStorage() *InMemoryStockStorage {
	return &InMemoryStockStorage{
		data: make(map[int64]*ItemUnits),
	}
}

func (i *InMemoryStockStorage) SetItemUnits(skuId int64, units *ItemUnits) {
	i.mu.Lock()
	defer i.mu.Unlock()
	i.data[skuId] = units
}

// Reserve резервирует товары. Если хотя бы один товар зарезервировать не удалось, то результат такой же, как будто операции не было. Возвращает первую произошедшую ошибку.
func (i *InMemoryStockStorage) Reserve(ctx context.Context, items []models.OrderItem) error {
	succeededItems := make([]models.OrderItem, len(items))
	failed := false
	var err error
	for _, it := range items {
		err = i.reserveOne(it)
		if err != nil {
			failed = true
			err = fmt.Errorf("error reserving %d units if item with skuId = %d: %w", it.Count, it.SkuId, err)
			break
		}
	}
	if !failed {
		return nil
	}
	_ = i.CancelReserved(ctx, succeededItems)
	return err
}

// AddItems добавляет незарезервированные позиции к итемам. Если произошла ошибка хотя бы с одним из них(не найден товар), то результат такой же, как будто операции не было.
func (i *InMemoryStockStorage) AddItems(ctx context.Context, items []models.OrderItem) error {
	addedReserved := make([]models.OrderItem, 0)
	failed := false
	var err error
	for _, it := range items {
		err = i.addReserved(it)
		if err != nil {
			failed = true
			err = fmt.Errorf("error returning %d units if item with skuId = %d: %w", it.Count, it.SkuId, err)
			break
		}
		addedReserved = append(addedReserved, it)
	}
	if !failed {
		_ = i.CancelReserved(ctx, addedReserved)
		return nil
	}
	_ = i.RemoveReserved(ctx, addedReserved)
	return err
}

func (i *InMemoryStockStorage) addReserved(it models.OrderItem) error {
	itemUni, err := i.getItemOrErr(it.SkuId)
	if err != nil {
		return err
	}
	itemUni.addReserved(it.Count)
	return nil
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
		return nil, fmt.Errorf("error getting item with skuid %d: %w", skuId, ErrItemIsNotExists)
	}
	return itemUni, nil
}
