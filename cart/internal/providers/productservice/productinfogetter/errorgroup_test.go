//go:build unit

package productinfogetter

import (
	"context"
	"fmt"
	"github.com/stretchr/testify/assert"
	"sync"
	"testing"
)

func Test_errorGroupResult_getFirstError(t *testing.T) {
	t.Parallel()
	result := errorGroupResult{}
	numOfCalls := 100000
	returnedErrs := make([]error, numOfCalls)

	wg := sync.WaitGroup{}
	for i := 0; i < numOfCalls; i++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			err := result.getFirstError(fmt.Errorf("error num %d", num))
			returnedErrs[num] = err
		}(i)
	}
	wg.Wait()
	// проверим, что все полученные ошибки одинаковы
	set := make(map[error]struct{})
	for _, val := range returnedErrs {
		set[val] = struct{}{}
	}
	assert.Len(t, set, 1)
}

func TestErrorGroup_Wait(t *testing.T) {
	t.Parallel()
	errGr, ctx := NewErrorGroup(context.Background())
	errorToThrow := fmt.Errorf("oops error")
	errGr.Go(func() error {
		return nil
	})
	errGr.Go(func() error {
		return errorToThrow
	})
	errFromWait := errGr.Wait()
	assert.ErrorIs(t, errFromWait, errorToThrow)
	// проверка, что контекст завершён
	assert.Error(t, ctx.Err())
}

func TestErrorGroup_WaitNil(t *testing.T) {
	t.Parallel()
	errGr, _ := NewErrorGroup(context.Background())
	errGr.Go(func() error {
		return nil
	})
	errFromWait := errGr.Wait()
	assert.Nil(t, errFromWait)

	errorToThrow := fmt.Errorf("oops error")
	errGr.Go(func() error {
		return errorToThrow
	})
	errFromWait = errGr.Wait()
	assert.Nil(t, errFromWait)
}
