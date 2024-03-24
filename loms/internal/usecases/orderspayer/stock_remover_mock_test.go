// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package orderspayer

//go:generate minimock -i route256.ozon.ru/project/loms/internal/usecases/orderspayer.stockRemover -o stock_remover_mock_test.go -n StockRemoverMock -p orderspayer

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/loms/internal/models"
)

// StockRemoverMock implements stockRemover
type StockRemoverMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcRemoveReserved          func(ctx context.Context, oa1 []models.OrderItem) (err error)
	inspectFuncRemoveReserved   func(ctx context.Context, oa1 []models.OrderItem)
	afterRemoveReservedCounter  uint64
	beforeRemoveReservedCounter uint64
	RemoveReservedMock          mStockRemoverMockRemoveReserved
}

// NewStockRemoverMock returns a mock for stockRemover
func NewStockRemoverMock(t minimock.Tester) *StockRemoverMock {
	m := &StockRemoverMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveReservedMock = mStockRemoverMockRemoveReserved{mock: m}
	m.RemoveReservedMock.callArgs = []*StockRemoverMockRemoveReservedParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mStockRemoverMockRemoveReserved struct {
	mock               *StockRemoverMock
	defaultExpectation *StockRemoverMockRemoveReservedExpectation
	expectations       []*StockRemoverMockRemoveReservedExpectation

	callArgs []*StockRemoverMockRemoveReservedParams
	mutex    sync.RWMutex
}

// StockRemoverMockRemoveReservedExpectation specifies expectation struct of the stockRemover.RemoveReserved
type StockRemoverMockRemoveReservedExpectation struct {
	mock    *StockRemoverMock
	params  *StockRemoverMockRemoveReservedParams
	results *StockRemoverMockRemoveReservedResults
	Counter uint64
}

// StockRemoverMockRemoveReservedParams contains parameters of the stockRemover.RemoveReserved
type StockRemoverMockRemoveReservedParams struct {
	ctx context.Context
	oa1 []models.OrderItem
}

// StockRemoverMockRemoveReservedResults contains results of the stockRemover.RemoveReserved
type StockRemoverMockRemoveReservedResults struct {
	err error
}

// Expect sets up expected params for stockRemover.RemoveReserved
func (mmRemoveReserved *mStockRemoverMockRemoveReserved) Expect(ctx context.Context, oa1 []models.OrderItem) *mStockRemoverMockRemoveReserved {
	if mmRemoveReserved.mock.funcRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("StockRemoverMock.RemoveReserved mock is already set by Set")
	}

	if mmRemoveReserved.defaultExpectation == nil {
		mmRemoveReserved.defaultExpectation = &StockRemoverMockRemoveReservedExpectation{}
	}

	mmRemoveReserved.defaultExpectation.params = &StockRemoverMockRemoveReservedParams{ctx, oa1}
	for _, e := range mmRemoveReserved.expectations {
		if minimock.Equal(e.params, mmRemoveReserved.defaultExpectation.params) {
			mmRemoveReserved.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRemoveReserved.defaultExpectation.params)
		}
	}

	return mmRemoveReserved
}

// Inspect accepts an inspector function that has same arguments as the stockRemover.RemoveReserved
func (mmRemoveReserved *mStockRemoverMockRemoveReserved) Inspect(f func(ctx context.Context, oa1 []models.OrderItem)) *mStockRemoverMockRemoveReserved {
	if mmRemoveReserved.mock.inspectFuncRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("Inspect function is already set for StockRemoverMock.RemoveReserved")
	}

	mmRemoveReserved.mock.inspectFuncRemoveReserved = f

	return mmRemoveReserved
}

// Return sets up results that will be returned by stockRemover.RemoveReserved
func (mmRemoveReserved *mStockRemoverMockRemoveReserved) Return(err error) *StockRemoverMock {
	if mmRemoveReserved.mock.funcRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("StockRemoverMock.RemoveReserved mock is already set by Set")
	}

	if mmRemoveReserved.defaultExpectation == nil {
		mmRemoveReserved.defaultExpectation = &StockRemoverMockRemoveReservedExpectation{mock: mmRemoveReserved.mock}
	}
	mmRemoveReserved.defaultExpectation.results = &StockRemoverMockRemoveReservedResults{err}
	return mmRemoveReserved.mock
}

// Set uses given function f to mock the stockRemover.RemoveReserved method
func (mmRemoveReserved *mStockRemoverMockRemoveReserved) Set(f func(ctx context.Context, oa1 []models.OrderItem) (err error)) *StockRemoverMock {
	if mmRemoveReserved.defaultExpectation != nil {
		mmRemoveReserved.mock.t.Fatalf("Default expectation is already set for the stockRemover.RemoveReserved method")
	}

	if len(mmRemoveReserved.expectations) > 0 {
		mmRemoveReserved.mock.t.Fatalf("Some expectations are already set for the stockRemover.RemoveReserved method")
	}

	mmRemoveReserved.mock.funcRemoveReserved = f
	return mmRemoveReserved.mock
}

// When sets expectation for the stockRemover.RemoveReserved which will trigger the result defined by the following
// Then helper
func (mmRemoveReserved *mStockRemoverMockRemoveReserved) When(ctx context.Context, oa1 []models.OrderItem) *StockRemoverMockRemoveReservedExpectation {
	if mmRemoveReserved.mock.funcRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("StockRemoverMock.RemoveReserved mock is already set by Set")
	}

	expectation := &StockRemoverMockRemoveReservedExpectation{
		mock:   mmRemoveReserved.mock,
		params: &StockRemoverMockRemoveReservedParams{ctx, oa1},
	}
	mmRemoveReserved.expectations = append(mmRemoveReserved.expectations, expectation)
	return expectation
}

// Then sets up stockRemover.RemoveReserved return parameters for the expectation previously defined by the When method
func (e *StockRemoverMockRemoveReservedExpectation) Then(err error) *StockRemoverMock {
	e.results = &StockRemoverMockRemoveReservedResults{err}
	return e.mock
}

// RemoveReserved implements stockRemover
func (mmRemoveReserved *StockRemoverMock) RemoveReserved(ctx context.Context, oa1 []models.OrderItem) (err error) {
	mm_atomic.AddUint64(&mmRemoveReserved.beforeRemoveReservedCounter, 1)
	defer mm_atomic.AddUint64(&mmRemoveReserved.afterRemoveReservedCounter, 1)

	if mmRemoveReserved.inspectFuncRemoveReserved != nil {
		mmRemoveReserved.inspectFuncRemoveReserved(ctx, oa1)
	}

	mm_params := StockRemoverMockRemoveReservedParams{ctx, oa1}

	// Record call args
	mmRemoveReserved.RemoveReservedMock.mutex.Lock()
	mmRemoveReserved.RemoveReservedMock.callArgs = append(mmRemoveReserved.RemoveReservedMock.callArgs, &mm_params)
	mmRemoveReserved.RemoveReservedMock.mutex.Unlock()

	for _, e := range mmRemoveReserved.RemoveReservedMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmRemoveReserved.RemoveReservedMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmRemoveReserved.RemoveReservedMock.defaultExpectation.Counter, 1)
		mm_want := mmRemoveReserved.RemoveReservedMock.defaultExpectation.params
		mm_got := StockRemoverMockRemoveReservedParams{ctx, oa1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRemoveReserved.t.Errorf("StockRemoverMock.RemoveReserved got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRemoveReserved.RemoveReservedMock.defaultExpectation.results
		if mm_results == nil {
			mmRemoveReserved.t.Fatal("No results are set for the StockRemoverMock.RemoveReserved")
		}
		return (*mm_results).err
	}
	if mmRemoveReserved.funcRemoveReserved != nil {
		return mmRemoveReserved.funcRemoveReserved(ctx, oa1)
	}
	mmRemoveReserved.t.Fatalf("Unexpected call to StockRemoverMock.RemoveReserved. %v %v", ctx, oa1)
	return
}

// RemoveReservedAfterCounter returns a count of finished StockRemoverMock.RemoveReserved invocations
func (mmRemoveReserved *StockRemoverMock) RemoveReservedAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemoveReserved.afterRemoveReservedCounter)
}

// RemoveReservedBeforeCounter returns a count of StockRemoverMock.RemoveReserved invocations
func (mmRemoveReserved *StockRemoverMock) RemoveReservedBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemoveReserved.beforeRemoveReservedCounter)
}

// Calls returns a list of arguments used in each call to StockRemoverMock.RemoveReserved.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRemoveReserved *mStockRemoverMockRemoveReserved) Calls() []*StockRemoverMockRemoveReservedParams {
	mmRemoveReserved.mutex.RLock()

	argCopy := make([]*StockRemoverMockRemoveReservedParams, len(mmRemoveReserved.callArgs))
	copy(argCopy, mmRemoveReserved.callArgs)

	mmRemoveReserved.mutex.RUnlock()

	return argCopy
}

// MinimockRemoveReservedDone returns true if the count of the RemoveReserved invocations corresponds
// the number of defined expectations
func (m *StockRemoverMock) MinimockRemoveReservedDone() bool {
	for _, e := range m.RemoveReservedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveReservedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveReservedCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoveReserved != nil && mm_atomic.LoadUint64(&m.afterRemoveReservedCounter) < 1 {
		return false
	}
	return true
}

// MinimockRemoveReservedInspect logs each unmet expectation
func (m *StockRemoverMock) MinimockRemoveReservedInspect() {
	for _, e := range m.RemoveReservedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to StockRemoverMock.RemoveReserved with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveReservedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveReservedCounter) < 1 {
		if m.RemoveReservedMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to StockRemoverMock.RemoveReserved")
		} else {
			m.t.Errorf("Expected call to StockRemoverMock.RemoveReserved with params: %#v", *m.RemoveReservedMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoveReserved != nil && mm_atomic.LoadUint64(&m.afterRemoveReservedCounter) < 1 {
		m.t.Error("Expected call to StockRemoverMock.RemoveReserved")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *StockRemoverMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockRemoveReservedInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *StockRemoverMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *StockRemoverMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRemoveReservedDone()
}