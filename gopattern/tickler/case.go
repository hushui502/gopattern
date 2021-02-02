package main

import (
	"container/list"
	"context"
	"log"
	"math/rand"
	"sync"
	"time"
)

type Request int
var wg *sync.WaitGroup

type Service struct {
	mu sync.Mutex
	queue *list.List
	sema chan int
	loopSignal chan struct{}
}

func (s *Service) EnqueueRequest(request Request) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	// 当新来请求时，先入栈
	s.queue.PushBack(request)
	log.Printf("Added request to queue with length %d\n", s.queue.Len())
	// 触发循环事件
	s.tickleLoop()

	return nil
}

func (s *Service) loop(ctx context.Context) {
	log.Println("Starting service loop")
	for {
		select {
		case <-s.loopSignal:
			// 开始尝试处理请求
			s.tryDequeue()
		case <-ctx.Done():
			log.Println("Loop context cancelled")
			return
		}
	}
}

func (s *Service) tryDequeue() {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.queue.Len() == 0 {
		return
	}
	select {
	case s.sema <- 1:		// 这里是一个信号量，避免同时处理过多请求
		request := s.dequeue()
		log.Printf("Dequeued request %v\n", request)
		wg.Add(1)
		// 真正的请求处理逻辑
		go s.process(request)
	default:
		log.Printf("Received loop signal, but request limit is reached")
	}
}

func (s *Service) dequeue() Request {
	element := s.queue.Front()
	s.queue.Remove(element)
	return element.Value.(Request)
}

func (s *Service) process(request Request) {
	// 请求处理结束需要将信号量恢复
	defer s.replenish()
	log.Printf("Processing request %v\n", request)

	// simulate work
	<-time.After(time.Duration(rand.Intn(300)) * time.Microsecond)
}

func (s *Service) replenish() {
	wg.Done()
	<-s.sema
	log.Printf("Replenishing semaphore, now %d%d slots in use\n", len(s.sema), cap(s.sema))
	// 信号量恢复了再次尝试触发事件
	s.tickleLoop()
}

func (s *Service) tickleLoop() {
	select {
	// 发送循环事件信号
	case s.loopSignal <- struct{}{}:
	default:
	}
}

func NewService(ctx context.Context, requestLimit int) *Service {
	service := &Service{
		queue:      list.New(),
		sema:       make(chan int, requestLimit),
		loopSignal: make(chan struct{}, 1),
	}
	// 开启一个协程专门接受loop循环事件信号
	go service.loop(ctx)

	return service
}

// use
func main() {
	wg = &sync.WaitGroup{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	service := NewService(ctx, 3)
	for i := 0; i < 10; i++ {
		if err := service.EnqueueRequest(Request(i)); err != nil {
			log.Fatalf("error sending a request: %v\n", err)
			break
		}

		<-time.After(time.Duration(rand.Intn(100)) * time.Microsecond)
	}

	wg.Wait()
}
