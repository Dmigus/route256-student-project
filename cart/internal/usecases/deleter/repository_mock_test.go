// Code generated by http://github.com/gojuno/minimock (dev). DO NOT EDIT.

package deleter

//go:generate minimock -i route256.ozon.ru/project/cart/internal/usecases/deleter.repository -o repository_mock_test.go -n RepositoryMock -p deleter

import (
	"context"
	"sync"
	mm_atomic "sync/atomic"
	mm_time "time"

	"github.com/gojuno/minimock/v3"
	"route256.ozon.ru/project/cart/internal/models"
)

// RepositoryMock implements repository
type RepositoryMock struct {
	t          minimock.Tester
	finishOnce sync.Once

	funcGetCart          func(ctx context.Context, user int64) (cp1 *models.Cart, err error)
	inspectFuncGetCart   func(ctx context.Context, user int64)
	afterGetCartCounter  uint64
	beforeGetCartCounter uint64
	GetCartMock          mRepositoryMockGetCart

	funcSaveCart          func(ctx context.Context, user int64, cart *models.Cart) (err error)
	inspectFuncSaveCart   func(ctx context.Context, user int64, cart *models.Cart)
	afterSaveCartCounter  uint64
	beforeSaveCartCounter uint64
	SaveCartMock          mRepositoryMockSaveCart
}

// NewRepositoryMock returns a mock for repository
func NewRepositoryMock(t minimock.Tester) *RepositoryMock {
	m := &RepositoryMock{t: t}

	if controller, ok := t.(minimock.MockController); ok {
		controller.RegisterMocker(m)
	}

	m.GetCartMock = mRepositoryMockGetCart{mock: m}
	m.GetCartMock.callArgs = []*RepositoryMockGetCartParams{}

	m.SaveCartMock = mRepositoryMockSaveCart{mock: m}
	m.SaveCartMock.callArgs = []*RepositoryMockSaveCartParams{}

	t.Cleanup(m.MinimockFinish)

	return m
}

type mRepositoryMockGetCart struct {
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockGetCartExpectation
	expectations       []*RepositoryMockGetCartExpectation

	callArgs []*RepositoryMockGetCartParams
	mutex    sync.RWMutex
}

// RepositoryMockGetCartExpectation specifies expectation struct of the repository.GetCart
type RepositoryMockGetCartExpectation struct {
	mock    *RepositoryMock
	params  *RepositoryMockGetCartParams
	results *RepositoryMockGetCartResults
	Counter uint64
}

// RepositoryMockGetCartParams contains parameters of the repository.GetCart
type RepositoryMockGetCartParams struct {
	ctx  context.Context
	user int64
}

// RepositoryMockGetCartResults contains results of the repository.GetCart
type RepositoryMockGetCartResults struct {
	cp1 *models.Cart
	err error
}

// Expect sets up expected params for repository.GetCart
func (mmGetCart *mRepositoryMockGetCart) Expect(ctx context.Context, user int64) *mRepositoryMockGetCart {
	if mmGetCart.mock.funcGetCart != nil {
		mmGetCart.mock.t.Fatalf("RepositoryMock.GetCart mock is already set by Set")
	}

	if mmGetCart.defaultExpectation == nil {
		mmGetCart.defaultExpectation = &RepositoryMockGetCartExpectation{}
	}

	mmGetCart.defaultExpectation.params = &RepositoryMockGetCartParams{ctx, user}
	for _, e := range mmGetCart.expectations {
		if minimock.Equal(e.params, mmGetCart.defaultExpectation.params) {
			mmGetCart.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmGetCart.defaultExpectation.params)
		}
	}

	return mmGetCart
}

// Inspect accepts an inspector function that has same arguments as the repository.GetCart
func (mmGetCart *mRepositoryMockGetCart) Inspect(f func(ctx context.Context, user int64)) *mRepositoryMockGetCart {
	if mmGetCart.mock.inspectFuncGetCart != nil {
		mmGetCart.mock.t.Fatalf("Inspect function is already set for RepositoryMock.GetCart")
	}

	mmGetCart.mock.inspectFuncGetCart = f

	return mmGetCart
}

// Return sets up results that will be returned by repository.GetCart
func (mmGetCart *mRepositoryMockGetCart) Return(cp1 *models.Cart, err error) *RepositoryMock {
	if mmGetCart.mock.funcGetCart != nil {
		mmGetCart.mock.t.Fatalf("RepositoryMock.GetCart mock is already set by Set")
	}

	if mmGetCart.defaultExpectation == nil {
		mmGetCart.defaultExpectation = &RepositoryMockGetCartExpectation{mock: mmGetCart.mock}
	}
	mmGetCart.defaultExpectation.results = &RepositoryMockGetCartResults{cp1, err}
	return mmGetCart.mock
}

// Set uses given function f to mock the repository.GetCart method
func (mmGetCart *mRepositoryMockGetCart) Set(f func(ctx context.Context, user int64) (cp1 *models.Cart, err error)) *RepositoryMock {
	if mmGetCart.defaultExpectation != nil {
		mmGetCart.mock.t.Fatalf("Default expectation is already set for the repository.GetCart method")
	}

	if len(mmGetCart.expectations) > 0 {
		mmGetCart.mock.t.Fatalf("Some expectations are already set for the repository.GetCart method")
	}

	mmGetCart.mock.funcGetCart = f
	return mmGetCart.mock
}

// When sets expectation for the repository.GetCart which will trigger the result defined by the following
// Then helper
func (mmGetCart *mRepositoryMockGetCart) When(ctx context.Context, user int64) *RepositoryMockGetCartExpectation {
	if mmGetCart.mock.funcGetCart != nil {
		mmGetCart.mock.t.Fatalf("RepositoryMock.GetCart mock is already set by Set")
	}

	expectation := &RepositoryMockGetCartExpectation{
		mock:   mmGetCart.mock,
		params: &RepositoryMockGetCartParams{ctx, user},
	}
	mmGetCart.expectations = append(mmGetCart.expectations, expectation)
	return expectation
}

// Then sets up repository.GetCart return parameters for the expectation previously defined by the When method
func (e *RepositoryMockGetCartExpectation) Then(cp1 *models.Cart, err error) *RepositoryMock {
	e.results = &RepositoryMockGetCartResults{cp1, err}
	return e.mock
}

// GetCart implements repository
func (mmGetCart *RepositoryMock) GetCart(ctx context.Context, user int64) (cp1 *models.Cart, err error) {
	mm_atomic.AddUint64(&mmGetCart.beforeGetCartCounter, 1)
	defer mm_atomic.AddUint64(&mmGetCart.afterGetCartCounter, 1)

	if mmGetCart.inspectFuncGetCart != nil {
		mmGetCart.inspectFuncGetCart(ctx, user)
	}

	mm_params := RepositoryMockGetCartParams{ctx, user}

	// Record call args
	mmGetCart.GetCartMock.mutex.Lock()
	mmGetCart.GetCartMock.callArgs = append(mmGetCart.GetCartMock.callArgs, &mm_params)
	mmGetCart.GetCartMock.mutex.Unlock()

	for _, e := range mmGetCart.GetCartMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.cp1, e.results.err
		}
	}

	if mmGetCart.GetCartMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmGetCart.GetCartMock.defaultExpectation.Counter, 1)
		mm_want := mmGetCart.GetCartMock.defaultExpectation.params
		mm_got := RepositoryMockGetCartParams{ctx, user}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmGetCart.t.Errorf("RepositoryMock.GetCart got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmGetCart.GetCartMock.defaultExpectation.results
		if mm_results == nil {
			mmGetCart.t.Fatal("No results are set for the RepositoryMock.GetCart")
		}
		return (*mm_results).cp1, (*mm_results).err
	}
	if mmGetCart.funcGetCart != nil {
		return mmGetCart.funcGetCart(ctx, user)
	}
	mmGetCart.t.Fatalf("Unexpected call to RepositoryMock.GetCart. %v %v", ctx, user)
	return
}

// GetCartAfterCounter returns a count of finished RepositoryMock.GetCart invocations
func (mmGetCart *RepositoryMock) GetCartAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCart.afterGetCartCounter)
}

// GetCartBeforeCounter returns a count of RepositoryMock.GetCart invocations
func (mmGetCart *RepositoryMock) GetCartBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmGetCart.beforeGetCartCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.GetCart.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmGetCart *mRepositoryMockGetCart) Calls() []*RepositoryMockGetCartParams {
	mmGetCart.mutex.RLock()

	argCopy := make([]*RepositoryMockGetCartParams, len(mmGetCart.callArgs))
	copy(argCopy, mmGetCart.callArgs)

	mmGetCart.mutex.RUnlock()

	return argCopy
}

// MinimockGetCartDone returns true if the count of the GetCart invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockGetCartDone() bool {
	for _, e := range m.GetCartMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCartMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCartCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCart != nil && mm_atomic.LoadUint64(&m.afterGetCartCounter) < 1 {
		return false
	}
	return true
}

// MinimockGetCartInspect logs each unmet expectation
func (m *RepositoryMock) MinimockGetCartInspect() {
	for _, e := range m.GetCartMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.GetCart with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.GetCartMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterGetCartCounter) < 1 {
		if m.GetCartMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.GetCart")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.GetCart with params: %#v", *m.GetCartMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcGetCart != nil && mm_atomic.LoadUint64(&m.afterGetCartCounter) < 1 {
		m.t.Error("Expected call to RepositoryMock.GetCart")
	}
}

type mRepositoryMockSaveCart struct {
	mock               *RepositoryMock
	defaultExpectation *RepositoryMockSaveCartExpectation
	expectations       []*RepositoryMockSaveCartExpectation

	callArgs []*RepositoryMockSaveCartParams
	mutex    sync.RWMutex
}

// RepositoryMockSaveCartExpectation specifies expectation struct of the repository.SaveCart
type RepositoryMockSaveCartExpectation struct {
	mock    *RepositoryMock
	params  *RepositoryMockSaveCartParams
	results *RepositoryMockSaveCartResults
	Counter uint64
}

// RepositoryMockSaveCartParams contains parameters of the repository.SaveCart
type RepositoryMockSaveCartParams struct {
	ctx  context.Context
	user int64
	cart *models.Cart
}

// RepositoryMockSaveCartResults contains results of the repository.SaveCart
type RepositoryMockSaveCartResults struct {
	err error
}

// Expect sets up expected params for repository.SaveCart
func (mmSaveCart *mRepositoryMockSaveCart) Expect(ctx context.Context, user int64, cart *models.Cart) *mRepositoryMockSaveCart {
	if mmSaveCart.mock.funcSaveCart != nil {
		mmSaveCart.mock.t.Fatalf("RepositoryMock.SaveCart mock is already set by Set")
	}

	if mmSaveCart.defaultExpectation == nil {
		mmSaveCart.defaultExpectation = &RepositoryMockSaveCartExpectation{}
	}

	mmSaveCart.defaultExpectation.params = &RepositoryMockSaveCartParams{ctx, user, cart}
	for _, e := range mmSaveCart.expectations {
		if minimock.Equal(e.params, mmSaveCart.defaultExpectation.params) {
			mmSaveCart.mock.t.Fatalf("Expectation set by When has same params: %#v", *mmSaveCart.defaultExpectation.params)
		}
	}

	return mmSaveCart
}

// Inspect accepts an inspector function that has same arguments as the repository.SaveCart
func (mmSaveCart *mRepositoryMockSaveCart) Inspect(f func(ctx context.Context, user int64, cart *models.Cart)) *mRepositoryMockSaveCart {
	if mmSaveCart.mock.inspectFuncSaveCart != nil {
		mmSaveCart.mock.t.Fatalf("Inspect function is already set for RepositoryMock.SaveCart")
	}

	mmSaveCart.mock.inspectFuncSaveCart = f

	return mmSaveCart
}

// Return sets up results that will be returned by repository.SaveCart
func (mmSaveCart *mRepositoryMockSaveCart) Return(err error) *RepositoryMock {
	if mmSaveCart.mock.funcSaveCart != nil {
		mmSaveCart.mock.t.Fatalf("RepositoryMock.SaveCart mock is already set by Set")
	}

	if mmSaveCart.defaultExpectation == nil {
		mmSaveCart.defaultExpectation = &RepositoryMockSaveCartExpectation{mock: mmSaveCart.mock}
	}
	mmSaveCart.defaultExpectation.results = &RepositoryMockSaveCartResults{err}
	return mmSaveCart.mock
}

// Set uses given function f to mock the repository.SaveCart method
func (mmSaveCart *mRepositoryMockSaveCart) Set(f func(ctx context.Context, user int64, cart *models.Cart) (err error)) *RepositoryMock {
	if mmSaveCart.defaultExpectation != nil {
		mmSaveCart.mock.t.Fatalf("Default expectation is already set for the repository.SaveCart method")
	}

	if len(mmSaveCart.expectations) > 0 {
		mmSaveCart.mock.t.Fatalf("Some expectations are already set for the repository.SaveCart method")
	}

	mmSaveCart.mock.funcSaveCart = f
	return mmSaveCart.mock
}

// When sets expectation for the repository.SaveCart which will trigger the result defined by the following
// Then helper
func (mmSaveCart *mRepositoryMockSaveCart) When(ctx context.Context, user int64, cart *models.Cart) *RepositoryMockSaveCartExpectation {
	if mmSaveCart.mock.funcSaveCart != nil {
		mmSaveCart.mock.t.Fatalf("RepositoryMock.SaveCart mock is already set by Set")
	}

	expectation := &RepositoryMockSaveCartExpectation{
		mock:   mmSaveCart.mock,
		params: &RepositoryMockSaveCartParams{ctx, user, cart},
	}
	mmSaveCart.expectations = append(mmSaveCart.expectations, expectation)
	return expectation
}

// Then sets up repository.SaveCart return parameters for the expectation previously defined by the When method
func (e *RepositoryMockSaveCartExpectation) Then(err error) *RepositoryMock {
	e.results = &RepositoryMockSaveCartResults{err}
	return e.mock
}

// SaveCart implements repository
func (mmSaveCart *RepositoryMock) SaveCart(ctx context.Context, user int64, cart *models.Cart) (err error) {
	mm_atomic.AddUint64(&mmSaveCart.beforeSaveCartCounter, 1)
	defer mm_atomic.AddUint64(&mmSaveCart.afterSaveCartCounter, 1)

	if mmSaveCart.inspectFuncSaveCart != nil {
		mmSaveCart.inspectFuncSaveCart(ctx, user, cart)
	}

	mm_params := RepositoryMockSaveCartParams{ctx, user, cart}

	// Record call args
	mmSaveCart.SaveCartMock.mutex.Lock()
	mmSaveCart.SaveCartMock.callArgs = append(mmSaveCart.SaveCartMock.callArgs, &mm_params)
	mmSaveCart.SaveCartMock.mutex.Unlock()

	for _, e := range mmSaveCart.SaveCartMock.expectations {
		if minimock.Equal(*e.params, mm_params) {
			mm_atomic.AddUint64(&e.Counter, 1)
			return e.results.err
		}
	}

	if mmSaveCart.SaveCartMock.defaultExpectation != nil {
		mm_atomic.AddUint64(&mmSaveCart.SaveCartMock.defaultExpectation.Counter, 1)
		mm_want := mmSaveCart.SaveCartMock.defaultExpectation.params
		mm_got := RepositoryMockSaveCartParams{ctx, user, cart}
		if mm_want != nil && !minimock.Equal(*mm_want, mm_got) {
			mmSaveCart.t.Errorf("RepositoryMock.SaveCart got unexpected parameters, want: %#v, got: %#v%s\n", *mm_want, mm_got, minimock.Diff(*mm_want, mm_got))
		}

		mm_results := mmSaveCart.SaveCartMock.defaultExpectation.results
		if mm_results == nil {
			mmSaveCart.t.Fatal("No results are set for the RepositoryMock.SaveCart")
		}
		return (*mm_results).err
	}
	if mmSaveCart.funcSaveCart != nil {
		return mmSaveCart.funcSaveCart(ctx, user, cart)
	}
	mmSaveCart.t.Fatalf("Unexpected call to RepositoryMock.SaveCart. %v %v %v", ctx, user, cart)
	return
}

// SaveCartAfterCounter returns a count of finished RepositoryMock.SaveCart invocations
func (mmSaveCart *RepositoryMock) SaveCartAfterCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSaveCart.afterSaveCartCounter)
}

// SaveCartBeforeCounter returns a count of RepositoryMock.SaveCart invocations
func (mmSaveCart *RepositoryMock) SaveCartBeforeCounter() uint64 {
	return mm_atomic.LoadUint64(&mmSaveCart.beforeSaveCartCounter)
}

// Calls returns a list of arguments used in each call to RepositoryMock.SaveCart.
// The list is in the same order as the calls were made (i.e. recent calls have a higher index)
func (mmSaveCart *mRepositoryMockSaveCart) Calls() []*RepositoryMockSaveCartParams {
	mmSaveCart.mutex.RLock()

	argCopy := make([]*RepositoryMockSaveCartParams, len(mmSaveCart.callArgs))
	copy(argCopy, mmSaveCart.callArgs)

	mmSaveCart.mutex.RUnlock()

	return argCopy
}

// MinimockSaveCartDone returns true if the count of the SaveCart invocations corresponds
// the number of defined expectations
func (m *RepositoryMock) MinimockSaveCartDone() bool {
	for _, e := range m.SaveCartMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			return false
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SaveCartMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSaveCartCounter) < 1 {
		return false
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSaveCart != nil && mm_atomic.LoadUint64(&m.afterSaveCartCounter) < 1 {
		return false
	}
	return true
}

// MinimockSaveCartInspect logs each unmet expectation
func (m *RepositoryMock) MinimockSaveCartInspect() {
	for _, e := range m.SaveCartMock.expectations {
		if mm_atomic.LoadUint64(&e.Counter) < 1 {
			m.t.Errorf("Expected call to RepositoryMock.SaveCart with params: %#v", *e.params)
		}
	}

	// if default expectation was set then invocations count should be greater than zero
	if m.SaveCartMock.defaultExpectation != nil && mm_atomic.LoadUint64(&m.afterSaveCartCounter) < 1 {
		if m.SaveCartMock.defaultExpectation.params == nil {
			m.t.Error("Expected call to RepositoryMock.SaveCart")
		} else {
			m.t.Errorf("Expected call to RepositoryMock.SaveCart with params: %#v", *m.SaveCartMock.defaultExpectation.params)
		}
	}
	// if func was set then invocations count should be greater than zero
	if m.funcSaveCart != nil && mm_atomic.LoadUint64(&m.afterSaveCartCounter) < 1 {
		m.t.Error("Expected call to RepositoryMock.SaveCart")
	}
}

// MinimockFinish checks that all mocked methods have been called the expected number of times
func (m *RepositoryMock) MinimockFinish() {
	m.finishOnce.Do(func() {
		if !m.minimockDone() {
			m.MinimockGetCartInspect()

			m.MinimockSaveCartInspect()
			m.t.FailNow()
		}
	})
}

// MinimockWait waits for all mocked methods to be called the expected number of times
func (m *RepositoryMock) MinimockWait(timeout mm_time.Duration) {
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

func (m *RepositoryMock) minimockDone() bool {
	done := true
	return done &&
		m.MinimockGetCartDone() &&
		m.MinimockSaveCartDone()
}
