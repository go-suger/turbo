package broadcast

import (
	"log"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

var brk *Broker[string, []byte]

func TestMain(t *testing.M) {
	brk = New[string, []byte]()

	t.Run()
}

func receiveMessages[K comparable, V any](br *Broker[K, V], key K) {
	ch, cancel, err := br.SubscribeWithCancel(key, 10)
	if err != nil {
		panic(err)
	}
	defer cancel()
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				log.Println("ch close")
				return
			}
			// log.Println(data)
		}
	}
}

func TestBroadcast(t *testing.T) {
	subNum := 1000
	var num atomic.Uint64

	wg := sync.WaitGroup{}

	for i := 0; i < subNum; i++ {
		go func() {
			wg.Add(1)
			ch, cancel, err := brk.SubscribeWithCancel(strconv.Itoa(i), 10)
			if err != nil {
				panic(err)
			}
			defer cancel()
			for {
				select {
				case _, ok := <-ch:
					if !ok {
						panic("ch close")
						return
					}
					num.Add(1)
				}
			}
		}()
	}

	go func() {
		var pustData = make([]byte, 20*1024*1024) // 20 KB = 20 * 1024 字节
		for i := range pustData {
			pustData[i] = 'a'
		}
		for {
			time.Sleep(10 * time.Microsecond)
			brk.Broadcast(pustData)
		}
	}()

	time.Sleep(time.Second)

	go func() {
		for i := 0; i < subNum; i++ {
			time.Sleep(50 * time.Microsecond)
			wg.Done()
		}
	}()

	wg.Wait()
}

type CustomKey struct {
	Id string
}

func TestCustomKey(t *testing.T) {
	_ = New[CustomKey, []byte]()
}

func BenchmarkBroadcast(b *testing.B) {
	var pustData = make([]byte, 20*1024*1024) // 20 KB = 20 * 1024 字节
	for i := range pustData {
		pustData[i] = 'a'
	}
	br := New[int, string]()
	for c := range 20 {
		go receiveMessages(br, c)
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		br.Broadcast(string(pustData))
	}
}
