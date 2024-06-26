// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package orderspayer

//go:generate minimock -i route256.ozon.ru/project/loms/internal/usecases/orderspayer.StockRepo -o stock_repo_mock_test.go -n StockRepoMock -p orderspayer

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/loms/internal/models"
)

// StockRepoMock implements StockRepo
type StockRepoMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcRemoveReserved          func(ctx context.Context, oa1 []models.OrderItem) (err error)
	inspectFuncRemoveReserved   func(ctx context.Context, oa1 []models.OrderItem)
	afterRemoveReservedCounter  uint64
	beforeRemoveReservedCounter uint64
	RemoveReservedMock          mStockRepoMockRemoveReserved
}

// NewStockRepoMock returns a mock for StockRepo
func NewStockRepoMock(t minimock.Tester) *StockRepoMock {
	m := &StockRepoMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.RemoveReservedMock = mStockRepoMockRemoveReserved{mock: m}
	m.RemoveReservedMock.callArgs = []*StockRepoMockRemoveReservedParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mStockRepoMockRemoveReserved struct {
	mock               *StockRepoMock
	defaultExpectation *StockRepoMockRemoveReservedExpectation
	expectations       []*StockRepoMockRemoveReservedExpectation

	callArgs []*StockRepoMockRemoveReservedParams
	mutex    sync.RWMutex
}

// StockRepoMockRemoveReservedExpectation specifies expectation struct of the StockRepo.RemoveReserved
type StockRepoMockRemoveReservedExpectation struct {
	mock    *StockRepoMock
	params  *StockRepoMockRemoveReservedParams
	results *StockRepoMockRemoveReservedResults
	Counter uint64
}

// StockRepoMockRemoveReservedParams contains parameters of the StockRepo.RemoveReserved
type StockRepoMockRemoveReservedParams struct {
	ctx context.Context
	oa1 []models.OrderItem
}

// StockRepoMockRemoveReservedResults contains results of the StockRepo.RemoveReserved
type StockRepoMockRemoveReservedResults struct {
	err error
}

// Expect sets up expected params for StockRepo.RemoveReserved
func (mmRemoveReserved *mStockRepoMockRemoveReserved) Expect(ctx context.Context, oa1 []models.OrderItem) *mStockRepoMockRemoveReserved {
	if mmRemoveReserved.mock.funcRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("StockRepoMock.RemoveReserved mock is already set by Set")
	}

	if mmRemoveReserved.defaultExpectation == nil {
		mmRemoveReserved.defaultExpectation = &StockRepoMockRemoveReservedExpectation{}
	}

	mmRemoveReserved.defaultExpectation.params = &StockRepoMockRemoveReservedParams{ctx, oa1}
	for _, e := range mmRemoveReserved.expectations {
		if minimock.Equal(e.params, mmRemoveReserved.defaultExpectation.params) {
			mmRemoveReserved.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmRemoveReserved.defaultExpectation.params)
		}
	}

	return mmRemoveReserved
}

// Inspect accepts an inspector function that has same arguments as the StockRepo.RemoveReserved
func (mmRemoveReserved *mStockRepoMockRemoveReserved) Inspect(f func(ctx context.Context, oa1 []models.OrderItem)) *mStockRepoMockRemoveReserved {
	if mmRemoveReserved.mock.inspectFuncRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("Inspect function is already set for StockRepoMock.RemoveReserved")
	}

	mmRemoveReserved.mock.inspectFuncRemoveReserved = f

	return mmRemoveReserved
}

// Return sets up results that will be returned by StockRepo.RemoveReserved
func (mmRemoveReserved *mStockRepoMockRemoveReserved) Return(err error) *StockRepoMock {
	if mmRemoveReserved.mock.funcRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("StockRepoMock.RemoveReserved mock is already set by Set")
	}

	if mmRemoveReserved.defaultExpectation == nil {
		mmRemoveReserved.defaultExpectation = &StockRepoMockRemoveReservedExpectation{mock: mmRemoveReserved.mock}
	}
	mmRemoveReserved.defaultExpectation.results = &StockRepoMockRemoveReservedResults{err}
	return mmRemoveReserved.mock
}

// Set uses given function f to mock the StockRepo.RemoveReserved method
func (mmRemoveReserved *mStockRepoMockRemoveReserved) Set(f func(ctx context.Context, oa1 []models.OrderItem) (err error)) *StockRepoMock {
	if mmRemoveReserved.defaultExpectation != nil {
		mmRemoveReserved.mock.t.Fatalf("Default expectation is already set for the StockRepo.RemoveReserved method")
	}

	if len(mmRemoveReserved.expectations) > 0 {
		mmRemoveReserved.mock.t.Fatalf("Some expectations are already set for the StockRepo.RemoveReserved method")
	}

	mmRemoveReserved.mock.funcRemoveReserved = f
	return mmRemoveReserved.mock
}

// When sets expectation for the StockRepo.RemoveReserved which will trigger the result defined by the following
// Then helper
func (mmRemoveReserved *mStockRepoMockRemoveReserved) When(ctx context.Context, oa1 []models.OrderItem) *StockRepoMockRemoveReservedExpectation {
	if mmRemoveReserved.mock.funcRemoveReserved != nil {
		mmRemoveReserved.mock.t.Fatalf("StockRepoMock.RemoveReserved mock is already set by Set")
	}

	expectation := &StockRepoMockRemoveReservedExpectation{
		mock:   mmRemoveReserved.mock,
		params: &StockRepoMockRemoveReservedParams{ctx, oa1},
	}
	mmRemoveReserved.expectations = append(mmRemoveReserved.expectations, expectation)
	return expectation
}

// Then sets up StockRepo.RemoveReserved return parameters for the expectation previously defined by the When method
func (e *StockRepoMockRemoveReservedExpectation) Then(err error) *StockRepoMock {
	e.results = &StockRepoMockRemoveReservedResults{err}
	return e.mock
}

// RemoveReserved implements StockRepo
func (mmRemoveReserved *StockRepoMock) RemoveReserved(ctx context.Context, oa1 []models.OrderItem) (err error) {
	mm_atomic.AddUint64(&mmRemoveReserved.beforeRemoveReservedCounter, 1)
	defer mm_atomic.AddUint64(&mmRemoveReserved.afterRemoveReservedCounter, 1)

	if mmRemoveReserved.inspectFuncRemoveReserved != nil {
		mmRemoveReserved.inspectFuncRemoveReserved(ctx, oa1)
	}

	mm_params := StockRepoMockRemoveReservedParams{ctx, oa1}

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
		mm_got := StockRepoMockRemoveReservedParams{ctx, oa1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmRemoveReserved.t.Errorf("StockRepoMock.RemoveReserved got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmRemoveReserved.RemoveReservedMock.defaultExpectation.results
		if mm_results == nil {
			mmRemoveReserved.t.Fatal("No results are set for the StockRepoMock.RemoveReserved")
		}
		return (*mm_results).err
	}
	if mmRemoveReserved.funcRemoveReserved != nil {
		return mmRemoveReserved.funcRemoveReserved(ctx, oa1)
	}
	mmRemoveReserved.t.Fatalf("Unexpected call to StockRepoMock.RemoveReserved. %v %v", ctx, oa1)
	return
}

// RemoveReservedAfterCounter returns a count of finished StockRepoMock.RemoveReserved invocations
func (mmRemoveReserved *StockRepoMock) RemoveReservedAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemoveReserved.afterRemoveReservedCounter)
}

// RemoveReservedBeforeCounter returns a count of StockRepoMock.RemoveReserved invocations
func (mmRemoveReserved *StockRepoMock) RemoveReservedBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmRemoveReserved.beforeRemoveReservedCounter)
}

// Calls returns a list of arguments used in each call to StockRepoMock.RemoveReserved.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmRemoveReserved *mStockRepoMockRemoveReserved) Calls() []*StockRepoMockRemoveReservedParams {
	mmRemoveReserved.mutex.RLock()

	argCopy := make([]*StockRepoMockRemoveReservedParams, len(mmRemoveReserved.callArgs))
	copy(argCopy, mmRemoveReserved.callArgs)

	mmRemoveReserved.mutex.RUnlock()

	return argCopy
}

// MinimockRemoveReservedDone returns true if the count of the RemoveReserved invocations corresponds
// the number of defined expectations
func (m *StockRepoMock) MinimockRemoveReservedDone() bool {
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
func (m *StockRepoMock) MinimockRemoveReservedInspect() {
	for _, e := range m.RemoveReservedMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to StockRepoMock.RemoveReserved with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.RemoveReservedMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterRemoveReservedCounter) < 1 {
		if m.RemoveReservedMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to StockRepoMock.RemoveReserved")
		} else {
			m.t.Errorf("Expected call to StockRepoMock.RemoveReserved with params: %#v", *m.RemoveReservedMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcRemoveReserved != nil && mm_atomic.LoadUint64(&m.afterRemoveReservedCounter) < 1 {
		m.t.Error("Expected call to StockRepoMock.RemoveReserved")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *StockRepoMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockRemoveReservedInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *StockRepoMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *StockRepoMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockRemoveReservedDone()
}
