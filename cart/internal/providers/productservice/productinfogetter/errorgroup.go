package productinfogetter

import (
	"context"
	"sync"
	"sync/atomic"
)

type errorGroupResult struct {
	mu        sync.RWMutex
	result    error
	isPresent atomic.Bool
}

// getFirstError сохраняет значение newErr первого вызова этого метода и возвращает его. Последующие вызовы этого
// метода возвращают первое сохранённое значение.
func (er *errorGroupResult) getFirstError(newErr error) error {
	er.mu.RLock()
	present := er.isPresent.Load()
	er.mu.RUnlock()
	if present {
		return er.result
	}
	er.mu.Lock()
	if er.isPresent.CompareAndSwap(false, true) {
		er.result = newErr
	}
	er.mu.Unlock()
	return er.result
}

// ErrorGroup это самописный аналог errorgroup.
type ErrorGroup struct {
	cancel context.CancelFunc
	wg     sync.WaitGroup
	errs   chan error
	ctx    context.Context
	result errorGroupResult
}

// NewErrorGroup возвращает новую ErrorGroup и производный от ctx контекст. Производный контекст будет отменён, когда первая функция, переданная в метод Go вернёт не nil ошибку, либо метод Wait вернёт какое-либо значение.
func NewErrorGroup(ctx context.Context) (*ErrorGroup, context.Context) {
	ctxToManage, cancelFunc := context.WithCancel(ctx)
	return &ErrorGroup{
		ctx:    ctxToManage,
		cancel: cancelFunc,
		errs:   make(chan error, 1),
	}, ctxToManage
}

// Go запускает функцию f в новой горутине. Первое не нулевое возвращаемое значение из f фиксируется и будет возращаться функцией Wait в дальнейшем
func (e *ErrorGroup) Go(f func() error) {
	e.wg.Add(1)
	go func() {
		defer e.wg.Done()
		err := f()
		if err == nil || e.isCancelled() {
			return
		}
		select {
		case e.errs <- err:
			e.cancel()
		default:
		}
	}()
}

// Wait дожидается, пока вернётся первое ненулевое значение функций f, переданных в метод Go и возращает его. Если все f отработали без ошибки, то будет возвращать nil
func (e *ErrorGroup) Wait() error {
	defer e.cancel()
	go func() {
		defer e.cancel()
		e.wg.Wait()
	}()
	var receivedErr error
	// ожидание завершения всех ошибок, либо первой ошибики
	select {
	case <-e.ctx.Done():
		select {
		case receivedErr = <-e.errs:
		default:
		}
	case receivedErr = <-e.errs:
	}
	return e.result.getFirstError(receivedErr)
}

func (e *ErrorGroup) isCancelled() bool {
	select {
	case <-e.ctx.Done():
		return true
	default:
		return false
	}
}
