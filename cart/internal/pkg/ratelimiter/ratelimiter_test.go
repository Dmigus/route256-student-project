//go:build unit

package ratelimiter

import (
	"context"
	"github.com/gojuno/minimock/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"sync/atomic"
	"testing"
	"time"
)

func TestRateLimiter_ReturnWithCancel(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	tck := NewSystemTimeTicker(1)
	rl := NewRateLimiter(0, tck)
	err := rl.Acquire(ctx)
	assert.Error(t, err)
}

func TestRateLimiter_GetInitialResources(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 500*time.Millisecond)
	defer cancel()
	mc := minimock.NewController(t)
	tickerMock := NewTickerMock(mc)
	tickerMock.StartMock.Return()
	tickerMock.StopMock.Return()

	channelToTick := make(chan struct{}, 200)
	tickerMock.GetTickChMock.Return(channelToTick)
	rl := NewRateLimiter(100, tickerMock)
	for i := 0; i < 100; i++ {
		err := rl.Acquire(ctx)
		assert.NoError(t, err)
	}

	for i := 0; i < 200; i++ {
		channelToTick <- struct{}{}
	}

	for i := 0; i < 100; i++ {
		err := rl.Acquire(ctx)
		assert.NoError(t, err)
	}

	// рейт лимитер может не успеть вызвать stop. Minimock всё равно проверяет наличие хотя бы одного вызова, поэтому установим счётчик
	atomic.CompareAndSwapUint64(&tickerMock.StopMock.mock.afterStopCounter, 0, 1)
}

func TestRateLimiter_AcquireTwoTimes(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	tickerMock := NewTickerMock(mc)
	tickerMock.StartMock.Return()

	channelToTick := make(chan struct{}, 100)
	tickerMock.GetTickChMock.Return(channelToTick)

	maxTestDuration := 2 * time.Second
	rl := NewRateLimiter(1, tickerMock)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(maxTestDuration)
		cancel()
	}()

	err := rl.Acquire(ctx)
	require.NoError(t, err)

	// проверка, что проходит достаточно времени перед получением второй раз
	startWait := time.Now()
	go func() {
		time.Sleep(50 * time.Millisecond)
		channelToTick <- struct{}{}
	}()
	err = rl.Acquire(ctx)
	require.NoError(t, err)
	endWait := time.Now()
	assert.GreaterOrEqual(t, endWait.Sub(startWait), 50*time.Millisecond)
}

func TestRateLimiter_AcquireManyTimes(t *testing.T) {
	t.Parallel()
	tck := NewSystemTimeTicker(1000)
	rl := NewRateLimiter(0, tck)
	maxTestDuration := 2 * time.Second
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(maxTestDuration)
		cancel()
	}()
	startWait := time.Now()
	for i := 0; i < 1000; i++ {
		err := rl.Acquire(ctx)
		assert.NoError(t, err)
	}
	endWait := time.Now()
	assert.GreaterOrEqual(t, endWait.Sub(startWait), 50*time.Millisecond)
}

func TestRateLimiter_checkStartAndStopTicker(t *testing.T) {
	t.Parallel()
	mc := minimock.NewController(t)
	tickerMock := NewTickerMock(mc)
	tickerMock.StartMock.Return()
	tickerMock.StopMock.Return()
	maxTestDuration := 2 * time.Second
	rl := NewRateLimiter(0, tickerMock)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		time.Sleep(maxTestDuration)
		cancel()
	}()
	channelToTick := make(chan struct{}, 100)
	tickerMock.GetTickChMock.Return(channelToTick)
	channelToTick <- struct{}{}

	// тикер не стартовал в начале работы
	assert.Equal(t, uint64(0), tickerMock.StartBeforeCounter())
	assert.False(t, rl.tickerRunning)

	// делаем запрос и проверим, что тикер стартовал
	err := rl.Acquire(ctx)
	require.NoError(t, err)
	assert.Equal(t, uint64(1), tickerMock.StartBeforeCounter())

	channelToTick <- struct{}{}
	// проверим, что через некорое время он остановился
	time.Sleep(50 * time.Millisecond)
	assert.Equal(t, uint64(1), tickerMock.StopBeforeCounter())
	assert.False(t, rl.tickerRunning)

}
