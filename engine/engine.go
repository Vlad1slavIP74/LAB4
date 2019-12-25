package engine

import (
	"sync"
)

type Command interface {
	Execute(handler Handler)
}

type Handler interface {
	Post(cmd Command, wg *sync.WaitGroup)
}

type messageQueue struct {
	data []Command
}

type MessageQueue struct {
	data []Command
}

type commandFunc func(handler Handler)

func (c commandFunc) Execute(handler Handler) {
	c(handler)
}

func (mq *MessageQueue) pull() Command {
	res := mq.data[0]
	mq.data[0] = nil
	mq.data = mq.data[1:]
	return res
}

func (mq *MessageQueue) push(cmd Command, wg *sync.WaitGroup) {
	mq.data = append(mq.data, cmd)
	wg.Done()
}

func (mq *MessageQueue) size() int {
	return len(mq.data)
}

func (l *Loop) Post(cmd Command, wg *sync.WaitGroup) {
	l.queue.push(cmd, wg)
}

type Loop struct {
	queue *MessageQueue
}

func (l *Loop) Start(wg *sync.WaitGroup) {
	l.queue = new(MessageQueue)

	go func() {
		for l.queue.size() != 0 {
			cmd := l.queue.pull()
			cmd.Execute(l)
			wg.Done()
		}
	}()
}
