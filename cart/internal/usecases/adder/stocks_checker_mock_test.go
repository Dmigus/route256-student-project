// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package adder

//go:generate minimock -i route256.ozon.ru/project/cart/internal/usecases/modifier.stocksChecker -o stocks_checker_mock_test.go -n StocksCheckerMock -p modifier

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// StocksCheckerMock implements stocksChecker
type StocksCheckerMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcIsItemAvailable          func(ctx context.Context, skuId int64, count uint16) (b1 bool, err error)
	inspectFuncIsItemAvailable   func(ctx context.Context, skuId int64, count uint16)
	afterIsItemAvailableCounter  uint64
	beforeIsItemAvailableCounter uint64
	IsItemAvailableMock          mStocksCheckerMockIsItemAvailable
}

// NewStocksCheckerMock returns a mock for stocksChecker
func NewStocksCheckerMock(t minimock.Tester) *StocksCheckerMock {
	m := &StocksCheckerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.IsItemAvailableMock = mStocksCheckerMockIsItemAvailable{mock: m}
	m.IsItemAvailableMock.callArgs = []*StocksCheckerMockIsItemAvailableParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mStocksCheckerMockIsItemAvailable struct {
	mock               *StocksCheckerMock
	defaultExpectation *StocksCheckerMockIsItemAvailableExpectation
	expectations       []*StocksCheckerMockIsItemAvailableExpectation

	callArgs []*StocksCheckerMockIsItemAvailableParams
	mutex    sync.RWMutex
}

// StocksCheckerMockIsItemAvailableExpectation specifies expectation struct of the stocksChecker.IsItemAvailable
type StocksCheckerMockIsItemAvailableExpectation struct {
	mock    *StocksCheckerMock
	params  *StocksCheckerMockIsItemAvailableParams
	results *StocksCheckerMockIsItemAvailableResults
	Counter uint64
}

// StocksCheckerMockIsItemAvailableParams contains parameters of the stocksChecker.IsItemAvailable
type StocksCheckerMockIsItemAvailableParams struct {
	ctx   context.Context
	skuId int64
	count uint16
}

// StocksCheckerMockIsItemAvailableResults contains results of the stocksChecker.IsItemAvailable
type StocksCheckerMockIsItemAvailableResults struct {
	b1  bool
	err error
}

// Expect sets up expected params for stocksChecker.IsItemAvailable
func (mmIsItemAvailable *mStocksCheckerMockIsItemAvailable) Expect(ctx context.Context, skuId int64, count uint16) *mStocksCheckerMockIsItemAvailable {
	if mmIsItemAvailable.mock.funcIsItemAvailable != nil {
		mmIsItemAvailable.mock.t.Fatalf("StocksCheckerMock.IsItemAvailable mock is already set by Set")
	}

	if mmIsItemAvailable.defaultExpectation == nil {
		mmIsItemAvailable.defaultExpectation = &StocksCheckerMockIsItemAvailableExpectation{}
	}

	mmIsItemAvailable.defaultExpectation.params = &StocksCheckerMockIsItemAvailableParams{ctx, skuId, count}
	for _, e := range mmIsItemAvailable.expectations {
		if minimock.Equal(e.params, mmIsItemAvailable.defaultExpectation.params) {
			mmIsItemAvailable.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmIsItemAvailable.defaultExpectation.params)
		}
	}

	return mmIsItemAvailable
}

// Inspect accepts an inspector function that has same arguments as the stocksChecker.IsItemAvailable
func (mmIsItemAvailable *mStocksCheckerMockIsItemAvailable) Inspect(f func(ctx context.Context, skuId int64, count uint16)) *mStocksCheckerMockIsItemAvailable {
	if mmIsItemAvailable.mock.inspectFuncIsItemAvailable != nil {
		mmIsItemAvailable.mock.t.Fatalf("Inspect function is already set for StocksCheckerMock.IsItemAvailable")
	}

	mmIsItemAvailable.mock.inspectFuncIsItemAvailable = f

	return mmIsItemAvailable
}

// Return sets up results that will be returned by stocksChecker.IsItemAvailable
func (mmIsItemAvailable *mStocksCheckerMockIsItemAvailable) Return(b1 bool, err error) *StocksCheckerMock {
	if mmIsItemAvailable.mock.funcIsItemAvailable != nil {
		mmIsItemAvailable.mock.t.Fatalf("StocksCheckerMock.IsItemAvailable mock is already set by Set")
	}

	if mmIsItemAvailable.defaultExpectation == nil {
		mmIsItemAvailable.defaultExpectation = &StocksCheckerMockIsItemAvailableExpectation{mock: mmIsItemAvailable.mock}
	}
	mmIsItemAvailable.defaultExpectation.results = &StocksCheckerMockIsItemAvailableResults{b1, err}
	return mmIsItemAvailable.mock
}

// Set uses given function f to mock the stocksChecker.IsItemAvailable method
func (mmIsItemAvailable *mStocksCheckerMockIsItemAvailable) Set(f func(ctx context.Context, skuId int64, count uint16) (b1 bool, err error)) *StocksCheckerMock {
	if mmIsItemAvailable.defaultExpectation != nil {
		mmIsItemAvailable.mock.t.Fatalf("Default expectation is already set for the stocksChecker.IsItemAvailable method")
	}

	if len(mmIsItemAvailable.expectations) > 0 {
		mmIsItemAvailable.mock.t.Fatalf("Some expectations are already set for the stocksChecker.IsItemAvailable method")
	}

	mmIsItemAvailable.mock.funcIsItemAvailable = f
	return mmIsItemAvailable.mock
}

// When sets expectation for the stocksChecker.IsItemAvailable which will trigger the result defined by the following
// Then helper
func (mmIsItemAvailable *mStocksCheckerMockIsItemAvailable) When(ctx context.Context, skuId int64, count uint16) *StocksCheckerMockIsItemAvailableExpectation {
	if mmIsItemAvailable.mock.funcIsItemAvailable != nil {
		mmIsItemAvailable.mock.t.Fatalf("StocksCheckerMock.IsItemAvailable mock is already set by Set")
	}

	expectation := &StocksCheckerMockIsItemAvailableExpectation{
		mock:   mmIsItemAvailable.mock,
		params: &StocksCheckerMockIsItemAvailableParams{ctx, skuId, count},
	}
	mmIsItemAvailable.expectations = append(mmIsItemAvailable.expectations, expectation)
	return expectation
}

// Then sets up stocksChecker.IsItemAvailable return parameters for the expectation previously defined by the When method
func (e *StocksCheckerMockIsItemAvailableExpectation) Then(b1 bool, err error) *StocksCheckerMock {
	e.results = &StocksCheckerMockIsItemAvailableResults{b1, err}
	return e.mock
}

// IsItemAvailable implements stocksChecker
func (mmIsItemAvailable *StocksCheckerMock) IsItemAvailable(ctx context.Context, skuId int64, count uint16) (b1 bool, err error) {
	mm_atomic.AddUint64(&mmIsItemAvailable.beforeIsItemAvailableCounter, 1)
	defer mm_atomic.AddUint64(&mmIsItemAvailable.afterIsItemAvailableCounter, 1)

	if mmIsItemAvailable.inspectFuncIsItemAvailable != nil {
		mmIsItemAvailable.inspectFuncIsItemAvailable(ctx, skuId, count)
	}

	mm_params := StocksCheckerMockIsItemAvailableParams{ctx, skuId, count}

	// Record call args
	mmIsItemAvailable.IsItemAvailableMock.mutex.Lock()
	mmIsItemAvailable.IsItemAvailableMock.callArgs = append(mmIsItemAvailable.IsItemAvailableMock.callArgs, &mm_params)
	mmIsItemAvailable.IsItemAvailableMock.mutex.Unlock()

	for _, e := range mmIsItemAvailable.IsItemAvailableMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.b1, e.results.err
		}
	}

	if mmIsItemAvailable.IsItemAvailableMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmIsItemAvailable.IsItemAvailableMock.defaultExpectation.Counter, 1)
		mm_want := mmIsItemAvailable.IsItemAvailableMock.defaultExpectation.params
		mm_got := StocksCheckerMockIsItemAvailableParams{ctx, skuId, count}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmIsItemAvailable.t.Errorf("StocksCheckerMock.IsItemAvailable got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmIsItemAvailable.IsItemAvailableMock.defaultExpectation.results
		if mm_results == nil {
			mmIsItemAvailable.t.Fatal("No results are set for the StocksCheckerMock.IsItemAvailable")
		}
		return (*mm_results).b1, (*mm_results).err
	}
	if mmIsItemAvailable.funcIsItemAvailable != nil {
		return mmIsItemAvailable.funcIsItemAvailable(ctx, skuId, count)
	}
	mmIsItemAvailable.t.Fatalf("Unexpected call to StocksCheckerMock.IsItemAvailable. %v %v %v", ctx, skuId, count)
	return
}

// IsItemAvailableAfterCounter returns a count of finished StocksCheckerMock.IsItemAvailable invocations
func (mmIsItemAvailable *StocksCheckerMock) IsItemAvailableAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmIsItemAvailable.afterIsItemAvailableCounter)
}

// IsItemAvailableBeforeCounter returns a count of StocksCheckerMock.IsItemAvailable invocations
func (mmIsItemAvailable *StocksCheckerMock) IsItemAvailableBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmIsItemAvailable.beforeIsItemAvailableCounter)
}

// Calls returns a list of arguments used in each call to StocksCheckerMock.IsItemAvailable.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmIsItemAvailable *mStocksCheckerMockIsItemAvailable) Calls() []*StocksCheckerMockIsItemAvailableParams {
	mmIsItemAvailable.mutex.RLock()

	argCopy := make([]*StocksCheckerMockIsItemAvailableParams, len(mmIsItemAvailable.callArgs))
	copy(argCopy, mmIsItemAvailable.callArgs)

	mmIsItemAvailable.mutex.RUnlock()

	return argCopy
}

// MinimockIsItemAvailableDone returns true if the count of the IsItemAvailable invocations corresponds
// the number of defined expectations
func (m *StocksCheckerMock) MinimockIsItemAvailableDone() bool {
	for _, e := range m.IsItemAvailableMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.IsItemAvailableMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterIsItemAvailableCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcIsItemAvailable != nil && mm_atomic.LoadUint64(&m.afterIsItemAvailableCounter) < 1 {
		return false
	}
	return true
}

// MinimockIsItemAvailableInspect logs each unmet expectation
func (m *StocksCheckerMock) MinimockIsItemAvailableInspect() {
	for _, e := range m.IsItemAvailableMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to StocksCheckerMock.IsItemAvailable with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.IsItemAvailableMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterIsItemAvailableCounter) < 1 {
		if m.IsItemAvailableMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to StocksCheckerMock.IsItemAvailable")
		} else {
			m.t.Errorf("Expected call to StocksCheckerMock.IsItemAvailable with params: %#v", *m.IsItemAvailableMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcIsItemAvailable != nil && mm_atomic.LoadUint64(&m.afterIsItemAvailableCounter) < 1 {
		m.t.Error("Expected call to StocksCheckerMock.IsItemAvailable")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *StocksCheckerMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockIsItemAvailableInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *StocksCheckerMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *StocksCheckerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockIsItemAvailableDone()
}
