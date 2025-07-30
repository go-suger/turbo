package broadcast

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Subscriber[V any] struct {
	Ch       chan V
	msgCh    chan V
	stopCh   chan struct{}
	isClosed atomic.Bool
	wg       sync.WaitGroup
}

func (s *Subscriber[V]) run() {
	s.wg.Add(1)
	go s.runner()
}

// 订阅器处理协程
func (s *Subscriber[V]) runner() {
	defer s.wg.Done()
	defer close(s.Ch)
	defer close(s.msgCh)
	defer close(s.stopCh)
	for {
		select {
		case msg, ok := <-s.msgCh:
			if !ok {
				return
			}
			select {
			case s.Ch <- msg:
			case <-s.stopCh:
				return
			}
		case <-s.stopCh:
			return
		}
	}
}

func (s *Subscriber[V]) unsubscribe() {
	s.isClosed.Store(true)
	select {
	case s.stopCh <- struct{}{}:
	default:
	}
	s.wg.Wait()
}

// Broker Thread-safe subscription Broker
// no message backpressure, messages will be discarded when the buffer is full, and only one broadcast is guaranteed.
type Broker[K comparable, V any] struct {
	subscribers sync.Map
}

// New K is comparable, Equal(T) bool
func New[K comparable, V any]() *Broker[K, V] {
	return &Broker[K, V]{
		subscribers: sync.Map{},
	}
}

// SubList Subscriber List
func (b *Broker[K, V]) SubList() []K {
	var list []K
	b.subscribers.Range(func(key, value interface{}) bool {
		sub := value.(*Subscriber[V])
		if !sub.isClosed.Load() {
			list = append(list, key.(K))
		}
		return true
	})
	return list
}

// Broadcast
// Returns errors that were not sent successfully, optionally ignored
func (b *Broker[K, V]) Broadcast(message V) error {
	var failedKeys []K
	b.subscribers.Range(func(key, value interface{}) bool {
		sub := value.(*Subscriber[V])
		if sub.isClosed.Load() {
			failedKeys = append(failedKeys, key.(K))
			return true
		}
		select {
		case sub.msgCh <- message:
		default:
			failedKeys = append(failedKeys, key.(K))
		}
		return true
	})
	if len(failedKeys) > 0 {
		return fmt.Errorf("broadcast failed for keys: %v", failedKeys)
	}
	return nil
}

// Subscribe broadcast
func (b *Broker[K, V]) Subscribe(key K, bufSize int) (chan V, error) {
	sub := &Subscriber[V]{
		Ch:       make(chan V, bufSize),
		msgCh:    make(chan V, bufSize+10),
		stopCh:   make(chan struct{}),
		isClosed: atomic.Bool{},
	}
	actualValue, loaded := b.subscribers.LoadOrStore(key, sub)
	if loaded {
		return nil, fmt.Errorf("subscriber id %v already exists", key)
	}
	actualSub := actualValue.(*Subscriber[V])
	actualSub.run()
	return actualSub.Ch, nil
}

// SubscribeWithCancel subscribe
func (b *Broker[K, V]) SubscribeWithCancel(key K, bufSize int) (ch chan V, cancel func(), err error) {
	ch, err = b.Subscribe(key, bufSize)
	if err != nil {
		return nil, nil, err
	}
	return ch, func() {
		b.Unsubscribe(key)
	}, nil
}

// Unsubscribe unsubscribe broadcast
func (b *Broker[K, V]) Unsubscribe(key K) {
	if value, loaded := b.subscribers.LoadAndDelete(key); loaded {
		sub := value.(*Subscriber[V])
		sub.unsubscribe()
	}
}
