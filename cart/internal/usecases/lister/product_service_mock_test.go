// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package lister

//go:generate minimock -i route256.ozon.ru/project/cart/internal/usecases/lister.productService -o product_service_mock_test.go -n ProductServiceMock -p lister

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/cart/internal/models"
)

// ProductServiceMock implements productService
type ProductServiceMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetProductsInfo          func(ctx context.Context, skuIds []int64) (pa1 []models.ProductInfo, err error)
	inspectFuncGetProductsInfo   func(ctx context.Context, skuIds []int64)
	afterGetProductsInfoCounter  uint64
	beforeGetProductsInfoCounter uint64
	GetProductsInfoMock          mProductServiceMockGetProductsInfo
}

// NewProductServiceMock returns a mock for productService
func NewProductServiceMock(t minimock.Tester) *ProductServiceMock {
	m := &ProductServiceMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetProductsInfoMock = mProductServiceMockGetProductsInfo{mock: m}
	m.GetProductsInfoMock.callArgs = []*ProductServiceMockGetProductsInfoParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mProductServiceMockGetProductsInfo struct {
	mock               *ProductServiceMock
	defaultExpectation *ProductServiceMockGetProductsInfoExpectation
	expectations       []*ProductServiceMockGetProductsInfoExpectation

	callArgs []*ProductServiceMockGetProductsInfoParams
	mutex    sync.RWMutex
}

// ProductServiceMockGetProductsInfoExpectation specifies expectation struct of the productService.GetProductsInfo
type ProductServiceMockGetProductsInfoExpectation struct {
	mock    *ProductServiceMock
	params  *ProductServiceMockGetProductsInfoParams
	results *ProductServiceMockGetProductsInfoResults
	Counter uint64
}

// ProductServiceMockGetProductsInfoParams contains parameters of the productService.GetProductsInfo
type ProductServiceMockGetProductsInfoParams struct {
	ctx    context.Context
	skuIds []int64
}

// ProductServiceMockGetProductsInfoResults contains results of the productService.GetProductsInfo
type ProductServiceMockGetProductsInfoResults struct {
	pa1 []models.ProductInfo
	err error
}

// Expect sets up expected params for productService.GetProductsInfo
func (mmGetProductsInfo *mProductServiceMockGetProductsInfo) Expect(ctx context.Context, skuIds []int64) *mProductServiceMockGetProductsInfo {
	if mmGetProductsInfo.mock.funcGetProductsInfo != nil {
		mmGetProductsInfo.mock.t.Fatalf("ProductServiceMock.GetProductsInfo mock is already set by Set")
	}

	if mmGetProductsInfo.defaultExpectation == nil {
		mmGetProductsInfo.defaultExpectation = &ProductServiceMockGetProductsInfoExpectation{}
	}

	mmGetProductsInfo.defaultExpectation.params = &ProductServiceMockGetProductsInfoParams{ctx, skuIds}
	for _, e := range mmGetProductsInfo.expectations {
		if minimock.Equal(e.params, mmGetProductsInfo.defaultExpectation.params) {
			mmGetProductsInfo.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetProductsInfo.defaultExpectation.params)
		}
	}

	return mmGetProductsInfo
}

// Inspect accepts an inspector function that has same arguments as the productService.GetProductsInfo
func (mmGetProductsInfo *mProductServiceMockGetProductsInfo) Inspect(f func(ctx context.Context, skuIds []int64)) *mProductServiceMockGetProductsInfo {
	if mmGetProductsInfo.mock.inspectFuncGetProductsInfo != nil {
		mmGetProductsInfo.mock.t.Fatalf("Inspect function is already set for ProductServiceMock.GetProductsInfo")
	}

	mmGetProductsInfo.mock.inspectFuncGetProductsInfo = f

	return mmGetProductsInfo
}

// Return sets up results that will be returned by productService.GetProductsInfo
func (mmGetProductsInfo *mProductServiceMockGetProductsInfo) Return(pa1 []models.ProductInfo, err error) *ProductServiceMock {
	if mmGetProductsInfo.mock.funcGetProductsInfo != nil {
		mmGetProductsInfo.mock.t.Fatalf("ProductServiceMock.GetProductsInfo mock is already set by Set")
	}

	if mmGetProductsInfo.defaultExpectation == nil {
		mmGetProductsInfo.defaultExpectation = &ProductServiceMockGetProductsInfoExpectation{mock: mmGetProductsInfo.mock}
	}
	mmGetProductsInfo.defaultExpectation.results = &ProductServiceMockGetProductsInfoResults{pa1, err}
	return mmGetProductsInfo.mock
}

// Set uses given function f to mock the productService.GetProductsInfo method
func (mmGetProductsInfo *mProductServiceMockGetProductsInfo) Set(f func(ctx context.Context, skuIds []int64) (pa1 []models.ProductInfo, err error)) *ProductServiceMock {
	if mmGetProductsInfo.defaultExpectation != nil {
		mmGetProductsInfo.mock.t.Fatalf("Default expectation is already set for the productService.GetProductsInfo method")
	}

	if len(mmGetProductsInfo.expectations) > 0 {
		mmGetProductsInfo.mock.t.Fatalf("Some expectations are already set for the productService.GetProductsInfo method")
	}

	mmGetProductsInfo.mock.funcGetProductsInfo = f
	return mmGetProductsInfo.mock
}

// When sets expectation for the productService.GetProductsInfo which will trigger the result defined by the following
// Then helper
func (mmGetProductsInfo *mProductServiceMockGetProductsInfo) When(ctx context.Context, skuIds []int64) *ProductServiceMockGetProductsInfoExpectation {
	if mmGetProductsInfo.mock.funcGetProductsInfo != nil {
		mmGetProductsInfo.mock.t.Fatalf("ProductServiceMock.GetProductsInfo mock is already set by Set")
	}

	expectation := &ProductServiceMockGetProductsInfoExpectation{
		mock:   mmGetProductsInfo.mock,
		params: &ProductServiceMockGetProductsInfoParams{ctx, skuIds},
	}
	mmGetProductsInfo.expectations = append(mmGetProductsInfo.expectations, expectation)
	return expectation
}

// Then sets up productService.GetProductsInfo return parameters for the expectation previously defined by the When method
func (e *ProductServiceMockGetProductsInfoExpectation) Then(pa1 []models.ProductInfo, err error) *ProductServiceMock {
	e.results = &ProductServiceMockGetProductsInfoResults{pa1, err}
	return e.mock
}

// GetProductsInfo implements productService
func (mmGetProductsInfo *ProductServiceMock) GetProductsInfo(ctx context.Context, skuIds []int64) (pa1 []models.ProductInfo, err error) {
	mm_atomic.AddUint64(&mmGetProductsInfo.beforeGetProductsInfoCounter, 1)
	defer mm_atomic.AddUint64(&mmGetProductsInfo.afterGetProductsInfoCounter, 1)

	if mmGetProductsInfo.inspectFuncGetProductsInfo != nil {
		mmGetProductsInfo.inspectFuncGetProductsInfo(ctx, skuIds)
	}

	mm_params := ProductServiceMockGetProductsInfoParams{ctx, skuIds}

	// Record call args
	mmGetProductsInfo.GetProductsInfoMock.mutex.Lock()
	mmGetProductsInfo.GetProductsInfoMock.callArgs = append(mmGetProductsInfo.GetProductsInfoMock.callArgs, &mm_params)
	mmGetProductsInfo.GetProductsInfoMock.mutex.Unlock()

	for _, e := range mmGetProductsInfo.GetProductsInfoMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.pa1, e.results.err
		}
	}

	if mmGetProductsInfo.GetProductsInfoMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetProductsInfo.GetProductsInfoMock.defaultExpectation.Counter, 1)
		mm_want := mmGetProductsInfo.GetProductsInfoMock.defaultExpectation.params
		mm_got := ProductServiceMockGetProductsInfoParams{ctx, skuIds}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetProductsInfo.t.Errorf("ProductServiceMock.GetProductsInfo got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetProductsInfo.GetProductsInfoMock.defaultExpectation.results
		if mm_results == nil {
			mmGetProductsInfo.t.Fatal("No results are set for the ProductServiceMock.GetProductsInfo")
		}
		return (*mm_results).pa1, (*mm_results).err
	}
	if mmGetProductsInfo.funcGetProductsInfo != nil {
		return mmGetProductsInfo.funcGetProductsInfo(ctx, skuIds)
	}
	mmGetProductsInfo.t.Fatalf("Unexpected call to ProductServiceMock.GetProductsInfo. %v %v", ctx, skuIds)
	return
}

// GetProductsInfoAfterCounter returns a count of finished ProductServiceMock.GetProductsInfo invocations
func (mmGetProductsInfo *ProductServiceMock) GetProductsInfoAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProductsInfo.afterGetProductsInfoCounter)
}

// GetProductsInfoBeforeCounter returns a count of ProductServiceMock.GetProductsInfo invocations
func (mmGetProductsInfo *ProductServiceMock) GetProductsInfoBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetProductsInfo.beforeGetProductsInfoCounter)
}

// Calls returns a list of arguments used in each call to ProductServiceMock.GetProductsInfo.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetProductsInfo *mProductServiceMockGetProductsInfo) Calls() []*ProductServiceMockGetProductsInfoParams {
	mmGetProductsInfo.mutex.RLock()

	argCopy := make([]*ProductServiceMockGetProductsInfoParams, len(mmGetProductsInfo.callArgs))
	copy(argCopy, mmGetProductsInfo.callArgs)

	mmGetProductsInfo.mutex.RUnlock()

	return argCopy
}

// MinimockGetProductsInfoDone returns true if the count of the GetProductsInfo invocations corresponds
// the number of defined expectations
func (m *ProductServiceMock) MinimockGetProductsInfoDone() bool {
	for _, e := range m.GetProductsInfoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductsInfoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetProductsInfoCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProductsInfo != nil && mm_atomic.LoadUint64(&m.afterGetProductsInfoCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetProductsInfoInspect logs each unmet expectation
func (m *ProductServiceMock) MinimockGetProductsInfoInspect() {
	for _, e := range m.GetProductsInfoMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to ProductServiceMock.GetProductsInfo with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetProductsInfoMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetProductsInfoCounter) < 1 {
		if m.GetProductsInfoMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to ProductServiceMock.GetProductsInfo")
		} else {
			m.t.Errorf("Expected call to ProductServiceMock.GetProductsInfo with params: %#v", *m.GetProductsInfoMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetProductsInfo != nil && mm_atomic.LoadUint64(&m.afterGetProductsInfoCounter) < 1 {
		m.t.Error("Expected call to ProductServiceMock.GetProductsInfo")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *ProductServiceMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetProductsInfoInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *ProductServiceMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *ProductServiceMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetProductsInfoDone()
}
