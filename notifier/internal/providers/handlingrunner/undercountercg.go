package handlingrunner

import (
	"github.com/IBM/sarama"
	"sync"
)

// underCounterCG consumer group, которая позволяет защитить cg от преждевременного закрытия
type underCounterCG struct {
	wg sync.WaitGroup
	cg sarama.ConsumerGroup
}

func newUnderCounterCG(cg sarama.ConsumerGroup) *underCounterCG {
	return &underCounterCG{wg: sync.WaitGroup{}, cg: cg}
}

// GetForUsage позволяет получить consumer группу для использования
func (a *underCounterCG) GetForUsage() sarama.ConsumerGroup {
	a.wg.Add(1)
	return a.cg
}

// Done вызывается, когда необходимо просигнализировать о том, что взятая consumer группа не нужна
func (a *underCounterCG) Done() {
	a.wg.Done()
}

// CloseWhenIsNotUsed это неблокирующий запрос, который закроет consumer группу, когда она не нужна
func (a *underCounterCG) CloseWhenIsNotUsed() {
	go func() {
		a.wg.Wait()
		_ = a.cg.Close()
	}()
}
