// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package ordersgetter

//go:generate minimock -i route256.ozon.ru/project/loms/internal/usecases/ordersgetter.txManager -o tx_manager_mock_test.go -n TxManagerMock -p ordersgetter

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
)

// TxManagerMock implements txManager
type TxManagerMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcWithinTransaction          func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool) (err error)
	inspectFuncWithinTransaction   func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool)
	afterWithinTransactionCounter  uint64
	beforeWithinTransactionCounter uint64
	WithinTransactionMock          mTxManagerMockWithinTransaction
}

// NewTxManagerMock returns a mock for txManager
func NewTxManagerMock(t minimock.Tester) *TxManagerMock {
	m := &TxManagerMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.WithinTransactionMock = mTxManagerMockWithinTransaction{mock: m}
	m.WithinTransactionMock.callArgs = []*TxManagerMockWithinTransactionParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mTxManagerMockWithinTransaction struct {
	mock               *TxManagerMock
	defaultExpectation *TxManagerMockWithinTransactionExpectation
	expectations       []*TxManagerMockWithinTransactionExpectation

	callArgs []*TxManagerMockWithinTransactionParams
	mutex    sync.RWMutex
}

// TxManagerMockWithinTransactionExpectation specifies expectation struct of the txManager.WithinTransaction
type TxManagerMockWithinTransactionExpectation struct {
	mock    *TxManagerMock
	params  *TxManagerMockWithinTransactionParams
	results *TxManagerMockWithinTransactionResults
	Counter uint64
}

// TxManagerMockWithinTransactionParams contains parameters of the txManager.WithinTransaction
type TxManagerMockWithinTransactionParams struct {
	ctx context.Context
	f1  func(ctx context.Context, orders OrderRepo) bool
}

// TxManagerMockWithinTransactionResults contains results of the txManager.WithinTransaction
type TxManagerMockWithinTransactionResults struct {
	err error
}

// Expect sets up expected params for txManager.WithinTransaction
func (mmWithinTransaction *mTxManagerMockWithinTransaction) Expect(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool) *mTxManagerMockWithinTransaction {
	if mmWithinTransaction.mock.funcWithinTransaction != nil {
		mmWithinTransaction.mock.t.Fatalf("TxManagerMock.WithinTransaction mock is already set by Set")
	}

	if mmWithinTransaction.defaultExpectation == nil {
		mmWithinTransaction.defaultExpectation = &TxManagerMockWithinTransactionExpectation{}
	}

	mmWithinTransaction.defaultExpectation.params = &TxManagerMockWithinTransactionParams{ctx, f1}
	for _, e := range mmWithinTransaction.expectations {
		if minimock.Equal(e.params, mmWithinTransaction.defaultExpectation.params) {
			mmWithinTransaction.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmWithinTransaction.defaultExpectation.params)
		}
	}

	return mmWithinTransaction
}

// Inspect accepts an inspector function that has same arguments as the txManager.WithinTransaction
func (mmWithinTransaction *mTxManagerMockWithinTransaction) Inspect(f func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool)) *mTxManagerMockWithinTransaction {
	if mmWithinTransaction.mock.inspectFuncWithinTransaction != nil {
		mmWithinTransaction.mock.t.Fatalf("Inspect function is already set for TxManagerMock.WithinTransaction")
	}

	mmWithinTransaction.mock.inspectFuncWithinTransaction = f

	return mmWithinTransaction
}

// Return sets up results that will be returned by txManager.WithinTransaction
func (mmWithinTransaction *mTxManagerMockWithinTransaction) Return(err error) *TxManagerMock {
	if mmWithinTransaction.mock.funcWithinTransaction != nil {
		mmWithinTransaction.mock.t.Fatalf("TxManagerMock.WithinTransaction mock is already set by Set")
	}

	if mmWithinTransaction.defaultExpectation == nil {
		mmWithinTransaction.defaultExpectation = &TxManagerMockWithinTransactionExpectation{mock: mmWithinTransaction.mock}
	}
	mmWithinTransaction.defaultExpectation.results = &TxManagerMockWithinTransactionResults{err}
	return mmWithinTransaction.mock
}

// Set uses given function f to mock the txManager.WithinTransaction method
func (mmWithinTransaction *mTxManagerMockWithinTransaction) Set(f func(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool) (err error)) *TxManagerMock {
	if mmWithinTransaction.defaultExpectation != nil {
		mmWithinTransaction.mock.t.Fatalf("Default expectation is already set for the txManager.WithinTransaction method")
	}

	if len(mmWithinTransaction.expectations) > 0 {
		mmWithinTransaction.mock.t.Fatalf("Some expectations are already set for the txManager.WithinTransaction method")
	}

	mmWithinTransaction.mock.funcWithinTransaction = f
	return mmWithinTransaction.mock
}

// When sets expectation for the txManager.WithinTransaction which will trigger the result defined by the following
// Then helper
func (mmWithinTransaction *mTxManagerMockWithinTransaction) When(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool) *TxManagerMockWithinTransactionExpectation {
	if mmWithinTransaction.mock.funcWithinTransaction != nil {
		mmWithinTransaction.mock.t.Fatalf("TxManagerMock.WithinTransaction mock is already set by Set")
	}

	expectation := &TxManagerMockWithinTransactionExpectation{
		mock:   mmWithinTransaction.mock,
		params: &TxManagerMockWithinTransactionParams{ctx, f1},
	}
	mmWithinTransaction.expectations = append(mmWithinTransaction.expectations, expectation)
	return expectation
}

// Then sets up txManager.WithinTransaction return parameters for the expectation previously defined by the When method
func (e *TxManagerMockWithinTransactionExpectation) Then(err error) *TxManagerMock {
	e.results = &TxManagerMockWithinTransactionResults{err}
	return e.mock
}

// WithinTransaction implements txManager
func (mmWithinTransaction *TxManagerMock) WithinTransaction(ctx context.Context, f1 func(ctx context.Context, orders OrderRepo) bool) (err error) {
	mm_atomic.AddUint64(&mmWithinTransaction.beforeWithinTransactionCounter, 1)
	defer mm_atomic.AddUint64(&mmWithinTransaction.afterWithinTransactionCounter, 1)

	if mmWithinTransaction.inspectFuncWithinTransaction != nil {
		mmWithinTransaction.inspectFuncWithinTransaction(ctx, f1)
	}

	mm_params := TxManagerMockWithinTransactionParams{ctx, f1}

	// Record call args
	mmWithinTransaction.WithinTransactionMock.mutex.Lock()
	mmWithinTransaction.WithinTransactionMock.callArgs = append(mmWithinTransaction.WithinTransactionMock.callArgs, &mm_params)
	mmWithinTransaction.WithinTransactionMock.mutex.Unlock()

	for _, e := range mmWithinTransaction.WithinTransactionMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmWithinTransaction.WithinTransactionMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmWithinTransaction.WithinTransactionMock.defaultExpectation.Counter, 1)
		mm_want := mmWithinTransaction.WithinTransactionMock.defaultExpectation.params
		mm_got := TxManagerMockWithinTransactionParams{ctx, f1}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmWithinTransaction.t.Errorf("TxManagerMock.WithinTransaction got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmWithinTransaction.WithinTransactionMock.defaultExpectation.results
		if mm_results == nil {
			mmWithinTransaction.t.Fatal("No results are set for the TxManagerMock.WithinTransaction")
		}
		return (*mm_results).err
	}
	if mmWithinTransaction.funcWithinTransaction != nil {
		return mmWithinTransaction.funcWithinTransaction(ctx, f1)
	}
	mmWithinTransaction.t.Fatalf("Unexpected call to TxManagerMock.WithinTransaction. %v %v", ctx, f1)
	return
}

// WithinTransactionAfterCounter returns a count of finished TxManagerMock.WithinTransaction invocations
func (mmWithinTransaction *TxManagerMock) WithinTransactionAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWithinTransaction.afterWithinTransactionCounter)
}

// WithinTransactionBeforeCounter returns a count of TxManagerMock.WithinTransaction invocations
func (mmWithinTransaction *TxManagerMock) WithinTransactionBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmWithinTransaction.beforeWithinTransactionCounter)
}

// Calls returns a list of arguments used in each call to TxManagerMock.WithinTransaction.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmWithinTransaction *mTxManagerMockWithinTransaction) Calls() []*TxManagerMockWithinTransactionParams {
	mmWithinTransaction.mutex.RLock()

	argCopy := make([]*TxManagerMockWithinTransactionParams, len(mmWithinTransaction.callArgs))
	copy(argCopy, mmWithinTransaction.callArgs)

	mmWithinTransaction.mutex.RUnlock()

	return argCopy
}

// MinimockWithinTransactionDone returns true if the count of the WithinTransaction invocations corresponds
// the number of defined expectations
func (m *TxManagerMock) MinimockWithinTransactionDone() bool {
	for _, e := range m.WithinTransactionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WithinTransactionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterWithinTransactionCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWithinTransaction != nil && mm_atomic.LoadUint64(&m.afterWithinTransactionCounter) < 1 {
		return false
	}
	return true
}

// MinimockWithinTransactionInspect logs each unmet expectation
func (m *TxManagerMock) MinimockWithinTransactionInspect() {
	for _, e := range m.WithinTransactionMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to TxManagerMock.WithinTransaction with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.WithinTransactionMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterWithinTransactionCounter) < 1 {
		if m.WithinTransactionMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to TxManagerMock.WithinTransaction")
		} else {
			m.t.Errorf("Expected call to TxManagerMock.WithinTransaction with params: %#v", *m.WithinTransactionMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcWithinTransaction != nil && mm_atomic.LoadUint64(&m.afterWithinTransactionCounter) < 1 {
		m.t.Error("Expected call to TxManagerMock.WithinTransaction")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *TxManagerMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockWithinTransactionInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *TxManagerMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *TxManagerMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockWithinTransactionDone()
}
