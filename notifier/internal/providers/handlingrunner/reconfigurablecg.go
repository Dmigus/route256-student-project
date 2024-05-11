package handlingrunner

import (
	"context"
	"github.com/IBM/sarama"
	"sync"
)

const groupName = "notifier-group"

// reconfigurableConsumerGroup это ConsumerGroup, которую сожно переконфигурировать вызовом метода Init().
// Для получения актуальной ConsumerGroup необходимо использовать метод GetInitializedCG()
type reconfigurableConsumerGroup struct {
	mu            sync.Mutex
	isInitialized bool                 // доступен под мьютексом mu
	cgReady       chan struct{}        // канал, дающий право на использование cg
	cg            sarama.ConsumerGroup // можно читать, захватив cgReady, но менять только если захвачены оба mu и cgReady
}

func newReconfigurableConsumerGroup() *reconfigurableConsumerGroup {
	return &reconfigurableConsumerGroup{cgReady: make(chan struct{}, 1)}
}

// Init создаёт инициализирует новую consumer group. Если не удалось создать и подключиться к новой конфигурации, старая НЕ затирается
func (r *reconfigurableConsumerGroup) Init(addrs []string, config *sarama.Config) error {
	cg, err := sarama.NewConsumerGroup(addrs, groupName, config)
	if err != nil {
		return err
	}
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closeWhenReadyLocked()
	r.cg = cg
	// инвариант: r.cgReady может иметь ресурс, только если cg инициализирован
	r.cgReady <- struct{}{}
	r.isInitialized = true
	return nil
}

// GetInitializedCG возвращает канал для хранения одной ConsumerGroup, когда та будет инициализирована.
// Если контекст будет отменён, канал будет закрыт без отправки в него ConsumerGroup
func (r *reconfigurableConsumerGroup) GetInitializedCG(ctx context.Context) chan sarama.ConsumerGroup {
	res := make(chan sarama.ConsumerGroup, 1)
	select {
	case <-ctx.Done():
		close(res)
	default:
		go func() {
			defer close(res)
			select {
			case <-r.cgReady:
				res <- r.cg
				r.cgReady <- struct{}{}
			case <-ctx.Done():
			}
		}()
	}
	return res
}

func (r *reconfigurableConsumerGroup) closeWhenReadyLocked() {
	if !r.isInitialized {
		return
	}
	r.isInitialized = false
	// дожидаемся, что cg можно взять в эксклюзивное пользование
	<-r.cgReady
	_ = r.cg.Close()
	r.cg = nil
}

func (r *reconfigurableConsumerGroup) Shutdown() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.closeWhenReadyLocked()
}
