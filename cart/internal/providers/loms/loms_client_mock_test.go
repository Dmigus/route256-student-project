// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package loms

//go:generate minimock -i route256.ozon.ru/project/cart/internal/providers/loms.lomsClient -o loms_client_mock_test.go -n LomsClientMock -p loms

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/cart/internal/models"
)

// LomsClientMock implements lomsClient
type LomsClientMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetNumberOfItemInStocks          func(ctx context.Context, skuId int64) (u1 uint64, err error)
	inspectFuncGetNumberOfItemInStocks   func(ctx context.Context, skuId int64)
	afterGetNumberOfItemInStocksCounter  uint64
	beforeGetNumberOfItemInStocksCounter uint64
	GetNumberOfItemInStocksMock          mLomsClientMockGetNumberOfItemInStocks

	funcOrderCreate          func(ctx context.Context, userId int64, items []models.CartItem) (i1 int64, err error)
	inspectFuncOrderCreate   func(ctx context.Context, userId int64, items []models.CartItem)
	afterOrderCreateCounter  uint64
	beforeOrderCreateCounter uint64
	OrderCreateMock          mLomsClientMockOrderCreate
}

// NewLomsClientMock returns a mock for lomsClient
func NewLomsClientMock(t minimock.Tester) *LomsClientMock {
	m := &LomsClientMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetNumberOfItemInStocksMock = mLomsClientMockGetNumberOfItemInStocks{mock: m}
	m.GetNumberOfItemInStocksMock.callArgs = []*LomsClientMockGetNumberOfItemInStocksParams{}

	m.OrderCreateMock = mLomsClientMockOrderCreate{mock: m}
	m.OrderCreateMock.callArgs = []*LomsClientMockOrderCreateParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mLomsClientMockGetNumberOfItemInStocks struct {
	mock               *LomsClientMock
	defaultExpectation *LomsClientMockGetNumberOfItemInStocksExpectation
	expectations       []*LomsClientMockGetNumberOfItemInStocksExpectation

	callArgs []*LomsClientMockGetNumberOfItemInStocksParams
	mutex    sync.RWMutex
}

// LomsClientMockGetNumberOfItemInStocksExpectation specifies expectation struct of the lomsClient.GetNumberOfItemInStocks
type LomsClientMockGetNumberOfItemInStocksExpectation struct {
	mock    *LomsClientMock
	params  *LomsClientMockGetNumberOfItemInStocksParams
	results *LomsClientMockGetNumberOfItemInStocksResults
	Counter uint64
}

// LomsClientMockGetNumberOfItemInStocksParams contains parameters of the lomsClient.GetNumberOfItemInStocks
type LomsClientMockGetNumberOfItemInStocksParams struct {
	ctx   context.Context
	skuId int64
}

// LomsClientMockGetNumberOfItemInStocksResults contains results of the lomsClient.GetNumberOfItemInStocks
type LomsClientMockGetNumberOfItemInStocksResults struct {
	u1  uint64
	err error
}

// Expect sets up expected params for lomsClient.GetNumberOfItemInStocks
func (mmGetNumberOfItemInStocks *mLomsClientMockGetNumberOfItemInStocks) Expect(ctx context.Context, skuId int64) *mLomsClientMockGetNumberOfItemInStocks {
	if mmGetNumberOfItemInStocks.mock.funcGetNumberOfItemInStocks != nil {
		mmGetNumberOfItemInStocks.mock.t.Fatalf("LomsClientMock.GetNumberOfItemInStocks mock is already set by Set")
	}

	if mmGetNumberOfItemInStocks.defaultExpectation == nil {
		mmGetNumberOfItemInStocks.defaultExpectation = &LomsClientMockGetNumberOfItemInStocksExpectation{}
	}

	mmGetNumberOfItemInStocks.defaultExpectation.params = &LomsClientMockGetNumberOfItemInStocksParams{ctx, skuId}
	for _, e := range mmGetNumberOfItemInStocks.expectations {
		if minimock.Equal(e.params, mmGetNumberOfItemInStocks.defaultExpectation.params) {
			mmGetNumberOfItemInStocks.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetNumberOfItemInStocks.defaultExpectation.params)
		}
	}

	return mmGetNumberOfItemInStocks
}

// Inspect accepts an inspector function that has same arguments as the lomsClient.GetNumberOfItemInStocks
func (mmGetNumberOfItemInStocks *mLomsClientMockGetNumberOfItemInStocks) Inspect(f func(ctx context.Context, skuId int64)) *mLomsClientMockGetNumberOfItemInStocks {
	if mmGetNumberOfItemInStocks.mock.inspectFuncGetNumberOfItemInStocks != nil {
		mmGetNumberOfItemInStocks.mock.t.Fatalf("Inspect function is already set for LomsClientMock.GetNumberOfItemInStocks")
	}

	mmGetNumberOfItemInStocks.mock.inspectFuncGetNumberOfItemInStocks = f

	return mmGetNumberOfItemInStocks
}

// Return sets up results that will be returned by lomsClient.GetNumberOfItemInStocks
func (mmGetNumberOfItemInStocks *mLomsClientMockGetNumberOfItemInStocks) Return(u1 uint64, err error) *LomsClientMock {
	if mmGetNumberOfItemInStocks.mock.funcGetNumberOfItemInStocks != nil {
		mmGetNumberOfItemInStocks.mock.t.Fatalf("LomsClientMock.GetNumberOfItemInStocks mock is already set by Set")
	}

	if mmGetNumberOfItemInStocks.defaultExpectation == nil {
		mmGetNumberOfItemInStocks.defaultExpectation = &LomsClientMockGetNumberOfItemInStocksExpectation{mock: mmGetNumberOfItemInStocks.mock}
	}
	mmGetNumberOfItemInStocks.defaultExpectation.results = &LomsClientMockGetNumberOfItemInStocksResults{u1, err}
	return mmGetNumberOfItemInStocks.mock
}

// Set uses given function f to mock the lomsClient.GetNumberOfItemInStocks method
func (mmGetNumberOfItemInStocks *mLomsClientMockGetNumberOfItemInStocks) Set(f func(ctx context.Context, skuId int64) (u1 uint64, err error)) *LomsClientMock {
	if mmGetNumberOfItemInStocks.defaultExpectation != nil {
		mmGetNumberOfItemInStocks.mock.t.Fatalf("Default expectation is already set for the lomsClient.GetNumberOfItemInStocks method")
	}

	if len(mmGetNumberOfItemInStocks.expectations) > 0 {
		mmGetNumberOfItemInStocks.mock.t.Fatalf("Some expectations are already set for the lomsClient.GetNumberOfItemInStocks method")
	}

	mmGetNumberOfItemInStocks.mock.funcGetNumberOfItemInStocks = f
	return mmGetNumberOfItemInStocks.mock
}

// When sets expectation for the lomsClient.GetNumberOfItemInStocks which will trigger the result defined by the following
// Then helper
func (mmGetNumberOfItemInStocks *mLomsClientMockGetNumberOfItemInStocks) When(ctx context.Context, skuId int64) *LomsClientMockGetNumberOfItemInStocksExpectation {
	if mmGetNumberOfItemInStocks.mock.funcGetNumberOfItemInStocks != nil {
		mmGetNumberOfItemInStocks.mock.t.Fatalf("LomsClientMock.GetNumberOfItemInStocks mock is already set by Set")
	}

	expectation := &LomsClientMockGetNumberOfItemInStocksExpectation{
		mock:   mmGetNumberOfItemInStocks.mock,
		params: &LomsClientMockGetNumberOfItemInStocksParams{ctx, skuId},
	}
	mmGetNumberOfItemInStocks.expectations = append(mmGetNumberOfItemInStocks.expectations, expectation)
	return expectation
}

// Then sets up lomsClient.GetNumberOfItemInStocks return parameters for the expectation previously defined by the When method
func (e *LomsClientMockGetNumberOfItemInStocksExpectation) Then(u1 uint64, err error) *LomsClientMock {
	e.results = &LomsClientMockGetNumberOfItemInStocksResults{u1, err}
	return e.mock
}

// GetNumberOfItemInStocks implements lomsClient
func (mmGetNumberOfItemInStocks *LomsClientMock) GetNumberOfItemInStocks(ctx context.Context, skuId int64) (u1 uint64, err error) {
	mm_atomic.AddUint64(&mmGetNumberOfItemInStocks.beforeGetNumberOfItemInStocksCounter, 1)
	defer mm_atomic.AddUint64(&mmGetNumberOfItemInStocks.afterGetNumberOfItemInStocksCounter, 1)

	if mmGetNumberOfItemInStocks.inspectFuncGetNumberOfItemInStocks != nil {
		mmGetNumberOfItemInStocks.inspectFuncGetNumberOfItemInStocks(ctx, skuId)
	}

	mm_params := LomsClientMockGetNumberOfItemInStocksParams{ctx, skuId}

	// Record call args
	mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.mutex.Lock()
	mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.callArgs = append(mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.callArgs, &mm_params)
	mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.mutex.Unlock()

	for _, e := range mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.u1, e.results.err
		}
	}

	if mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.defaultExpectation.Counter, 1)
		mm_want := mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.defaultExpectation.params
		mm_got := LomsClientMockGetNumberOfItemInStocksParams{ctx, skuId}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetNumberOfItemInStocks.t.Errorf("LomsClientMock.GetNumberOfItemInStocks got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetNumberOfItemInStocks.GetNumberOfItemInStocksMock.defaultExpectation.results
		if mm_results == nil {
			mmGetNumberOfItemInStocks.t.Fatal("No results are set for the LomsClientMock.GetNumberOfItemInStocks")
		}
		return (*mm_results).u1, (*mm_results).err
	}
	if mmGetNumberOfItemInStocks.funcGetNumberOfItemInStocks != nil {
		return mmGetNumberOfItemInStocks.funcGetNumberOfItemInStocks(ctx, skuId)
	}
	mmGetNumberOfItemInStocks.t.Fatalf("Unexpected call to LomsClientMock.GetNumberOfItemInStocks. %v %v", ctx, skuId)
	return
}

// GetNumberOfItemInStocksAfterCounter returns a count of finished LomsClientMock.GetNumberOfItemInStocks invocations
func (mmGetNumberOfItemInStocks *LomsClientMock) GetNumberOfItemInStocksAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetNumberOfItemInStocks.afterGetNumberOfItemInStocksCounter)
}

// GetNumberOfItemInStocksBeforeCounter returns a count of LomsClientMock.GetNumberOfItemInStocks invocations
func (mmGetNumberOfItemInStocks *LomsClientMock) GetNumberOfItemInStocksBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetNumberOfItemInStocks.beforeGetNumberOfItemInStocksCounter)
}

// Calls returns a list of arguments used in each call to LomsClientMock.GetNumberOfItemInStocks.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetNumberOfItemInStocks *mLomsClientMockGetNumberOfItemInStocks) Calls() []*LomsClientMockGetNumberOfItemInStocksParams {
	mmGetNumberOfItemInStocks.mutex.RLock()

	argCopy := make([]*LomsClientMockGetNumberOfItemInStocksParams, len(mmGetNumberOfItemInStocks.callArgs))
	copy(argCopy, mmGetNumberOfItemInStocks.callArgs)

	mmGetNumberOfItemInStocks.mutex.RUnlock()

	return argCopy
}

// MinimockGetNumberOfItemInStocksDone returns true if the count of the GetNumberOfItemInStocks invocations corresponds
// the number of defined expectations
func (m *LomsClientMock) MinimockGetNumberOfItemInStocksDone() bool {
	for _, e := range m.GetNumberOfItemInStocksMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetNumberOfItemInStocksMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetNumberOfItemInStocksCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetNumberOfItemInStocks != nil && mm_atomic.LoadUint64(&m.afterGetNumberOfItemInStocksCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetNumberOfItemInStocksInspect logs each unmet expectation
func (m *LomsClientMock) MinimockGetNumberOfItemInStocksInspect() {
	for _, e := range m.GetNumberOfItemInStocksMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LomsClientMock.GetNumberOfItemInStocks with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetNumberOfItemInStocksMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetNumberOfItemInStocksCounter) < 1 {
		if m.GetNumberOfItemInStocksMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LomsClientMock.GetNumberOfItemInStocks")
		} else {
			m.t.Errorf("Expected call to LomsClientMock.GetNumberOfItemInStocks with params: %#v", *m.GetNumberOfItemInStocksMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetNumberOfItemInStocks != nil && mm_atomic.LoadUint64(&m.afterGetNumberOfItemInStocksCounter) < 1 {
		m.t.Error("Expected call to LomsClientMock.GetNumberOfItemInStocks")
	}
}

type mLomsClientMockOrderCreate struct {
	mock               *LomsClientMock
	defaultExpectation *LomsClientMockOrderCreateExpectation
	expectations       []*LomsClientMockOrderCreateExpectation

	callArgs []*LomsClientMockOrderCreateParams
	mutex    sync.RWMutex
}

// LomsClientMockOrderCreateExpectation specifies expectation struct of the lomsClient.OrderCreate
type LomsClientMockOrderCreateExpectation struct {
	mock    *LomsClientMock
	params  *LomsClientMockOrderCreateParams
	results *LomsClientMockOrderCreateResults
	Counter uint64
}

// LomsClientMockOrderCreateParams contains parameters of the lomsClient.OrderCreate
type LomsClientMockOrderCreateParams struct {
	ctx    context.Context
	userId int64
	items  []models.CartItem
}

// LomsClientMockOrderCreateResults contains results of the lomsClient.OrderCreate
type LomsClientMockOrderCreateResults struct {
	i1  int64
	err error
}

// Expect sets up expected params for lomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) Expect(ctx context.Context, userId int64, items []models.CartItem) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{}
	}

	mmOrderCreate.defaultExpectation.params = &LomsClientMockOrderCreateParams{ctx, userId, items}
	for _, e := range mmOrderCreate.expectations {
		if minimock.Equal(e.params, mmOrderCreate.defaultExpectation.params) {
			mmOrderCreate.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmOrderCreate.defaultExpectation.params)
		}
	}

	return mmOrderCreate
}

// Inspect accepts an inspector function that has same arguments as the lomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) Inspect(f func(ctx context.Context, userId int64, items []models.CartItem)) *mLomsClientMockOrderCreate {
	if mmOrderCreate.mock.inspectFuncOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("Inspect function is already set for LomsClientMock.OrderCreate")
	}

	mmOrderCreate.mock.inspectFuncOrderCreate = f

	return mmOrderCreate
}

// Return sets up results that will be returned by lomsClient.OrderCreate
func (mmOrderCreate *mLomsClientMockOrderCreate) Return(i1 int64, err error) *LomsClientMock {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	if mmOrderCreate.defaultExpectation == nil {
		mmOrderCreate.defaultExpectation = &LomsClientMockOrderCreateExpectation{mock: mmOrderCreate.mock}
	}
	mmOrderCreate.defaultExpectation.results = &LomsClientMockOrderCreateResults{i1, err}
	return mmOrderCreate.mock
}

// Set uses given function f to mock the lomsClient.OrderCreate method
func (mmOrderCreate *mLomsClientMockOrderCreate) Set(f func(ctx context.Context, userId int64, items []models.CartItem) (i1 int64, err error)) *LomsClientMock {
	if mmOrderCreate.defaultExpectation != nil {
		mmOrderCreate.mock.t.Fatalf("Default expectation is already set for the lomsClient.OrderCreate method")
	}

	if len(mmOrderCreate.expectations) > 0 {
		mmOrderCreate.mock.t.Fatalf("Some expectations are already set for the lomsClient.OrderCreate method")
	}

	mmOrderCreate.mock.funcOrderCreate = f
	return mmOrderCreate.mock
}

// When sets expectation for the lomsClient.OrderCreate which will trigger the result defined by the following
// Then helper
func (mmOrderCreate *mLomsClientMockOrderCreate) When(ctx context.Context, userId int64, items []models.CartItem) *LomsClientMockOrderCreateExpectation {
	if mmOrderCreate.mock.funcOrderCreate != nil {
		mmOrderCreate.mock.t.Fatalf("LomsClientMock.OrderCreate mock is already set by Set")
	}

	expectation := &LomsClientMockOrderCreateExpectation{
		mock:   mmOrderCreate.mock,
		params: &LomsClientMockOrderCreateParams{ctx, userId, items},
	}
	mmOrderCreate.expectations = append(mmOrderCreate.expectations, expectation)
	return expectation
}

// Then sets up lomsClient.OrderCreate return parameters for the expectation previously defined by the When method
func (e *LomsClientMockOrderCreateExpectation) Then(i1 int64, err error) *LomsClientMock {
	e.results = &LomsClientMockOrderCreateResults{i1, err}
	return e.mock
}

// OrderCreate implements lomsClient
func (mmOrderCreate *LomsClientMock) OrderCreate(ctx context.Context, userId int64, items []models.CartItem) (i1 int64, err error) {
	mm_atomic.AddUint64(&mmOrderCreate.beforeOrderCreateCounter, 1)
	defer mm_atomic.AddUint64(&mmOrderCreate.afterOrderCreateCounter, 1)

	if mmOrderCreate.inspectFuncOrderCreate != nil {
		mmOrderCreate.inspectFuncOrderCreate(ctx, userId, items)
	}

	mm_params := LomsClientMockOrderCreateParams{ctx, userId, items}

	// Record call args
	mmOrderCreate.OrderCreateMock.mutex.Lock()
	mmOrderCreate.OrderCreateMock.callArgs = append(mmOrderCreate.OrderCreateMock.callArgs, &mm_params)
	mmOrderCreate.OrderCreateMock.mutex.Unlock()

	for _, e := range mmOrderCreate.OrderCreateMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.i1, e.results.err
		}
	}

	if mmOrderCreate.OrderCreateMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmOrderCreate.OrderCreateMock.defaultExpectation.Counter, 1)
		mm_want := mmOrderCreate.OrderCreateMock.defaultExpectation.params
		mm_got := LomsClientMockOrderCreateParams{ctx, userId, items}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmOrderCreate.t.Errorf("LomsClientMock.OrderCreate got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmOrderCreate.OrderCreateMock.defaultExpectation.results
		if mm_results == nil {
			mmOrderCreate.t.Fatal("No results are set for the LomsClientMock.OrderCreate")
		}
		return (*mm_results).i1, (*mm_results).err
	}
	if mmOrderCreate.funcOrderCreate != nil {
		return mmOrderCreate.funcOrderCreate(ctx, userId, items)
	}
	mmOrderCreate.t.Fatalf("Unexpected call to LomsClientMock.OrderCreate. %v %v %v", ctx, userId, items)
	return
}

// OrderCreateAfterCounter returns a count of finished LomsClientMock.OrderCreate invocations
func (mmOrderCreate *LomsClientMock) OrderCreateAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOrderCreate.afterOrderCreateCounter)
}

// OrderCreateBeforeCounter returns a count of LomsClientMock.OrderCreate invocations
func (mmOrderCreate *LomsClientMock) OrderCreateBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmOrderCreate.beforeOrderCreateCounter)
}

// Calls returns a list of arguments used in each call to LomsClientMock.OrderCreate.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmOrderCreate *mLomsClientMockOrderCreate) Calls() []*LomsClientMockOrderCreateParams {
	mmOrderCreate.mutex.RLock()

	argCopy := make([]*LomsClientMockOrderCreateParams, len(mmOrderCreate.callArgs))
	copy(argCopy, mmOrderCreate.callArgs)

	mmOrderCreate.mutex.RUnlock()

	return argCopy
}

// MinimockOrderCreateDone returns true if the count of the OrderCreate invocations corresponds
// the number of defined expectations
func (m *LomsClientMock) MinimockOrderCreateDone() bool {
	for _, e := range m.OrderCreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.OrderCreateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterOrderCreateCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcOrderCreate != nil && mm_atomic.LoadUint64(&m.afterOrderCreateCounter) < 1 {
		return false
	}
	return true
}

// MinimockOrderCreateInspect logs each unmet expectation
func (m *LomsClientMock) MinimockOrderCreateInspect() {
	for _, e := range m.OrderCreateMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to LomsClientMock.OrderCreate with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.OrderCreateMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterOrderCreateCounter) < 1 {
		if m.OrderCreateMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to LomsClientMock.OrderCreate")
		} else {
			m.t.Errorf("Expected call to LomsClientMock.OrderCreate with params: %#v", *m.OrderCreateMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcOrderCreate != nil && mm_atomic.LoadUint64(&m.afterOrderCreateCounter) < 1 {
		m.t.Error("Expected call to LomsClientMock.OrderCreate")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *LomsClientMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetNumberOfItemInStocksInspect()

			m.MinimockOrderCreateInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *LomsClientMock) MinimockWait(timeout mm_time.Duration) {
	timeoutCh := mm_time.After(timeout)
	for {
		if m.minimockDone() {
			return
		}
		select {
		case <-timeoutCh:
			m.MinimockFinish()
			return
		case <-mm_time.After(10 * mm_time.Millisecond):
		}
	}
}

func (m *LomsClientMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetNumberOfItemInStocksDone() &&
		m.MinimockOrderCreateDone()
}
