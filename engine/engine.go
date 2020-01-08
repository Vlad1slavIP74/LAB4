package engine

import (
	"sync"
)

type Command interface {
	Execute(handler Handler)
}

type Handler interface {
	Post(cmd Command)
}

type messageQueue struct {
	data struct {
		mx      sync.Mutex
		arr     []Command
		waiting bool
	}
	receivedSignal chan bool
}

func (mq *messageQueue) pull() Command {
	mq.data.mx.Lock()
	defer mq.data.mx.Unlock()
	if mq.size() == 0 {
		mq.data.waiting = true
		mq.data.mx.Unlock()
		<-mq.receivedSignal
		mq.data.mx.Lock()
	}
	res := mq.data.arr[0]
	mq.data.arr[0] = nil
	mq.data.arr = mq.data.arr[1:]
	return res
}

func (mq *messageQueue) push(cmd Command) {
	mq.data.mx.Lock()
	defer mq.data.mx.Unlock()
	mq.data.arr = append(mq.data.arr, cmd)
	if mq.data.waiting {
		mq.data.waiting = false
		mq.receivedSignal <- true
	}
}

func (mq *messageQueue) size() int {
	return len(mq.data.arr)
}

func (l *Loop) Post(cmd Command) {
	l.queue.push(cmd)
}

type commandFunc func(handler Handler)

func (c commandFunc) Execute(handler Handler) {
	c(handler)
}

type Loop struct {
	queue          *messageQueue
	receivedSignal bool
	stopSignal     chan bool
}

func (l *Loop) Start() {
	l.queue = new(messageQueue)
	l.stopSignal = make(chan bool)
	l.queue.receivedSignal = make(chan bool)
	go func() {
		for l.queue.size() != 0 || !l.receivedSignal {
			cmd := l.queue.pull()
			cmd.Execute(l)
		}
		l.stopSignal <- true
	}()
}

func (l *Loop) AwaitFinish() {
	l.Post(commandFunc(func(h Handler) { h.(*Loop).receivedSignal = true }))
	<-l.stopSignal
}
