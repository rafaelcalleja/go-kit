package transaction

import "sync"

type Broker struct {
	stopCh    chan struct{}
	publishCh chan interface{}
	subCh     chan chan interface{}
	unsubCh   chan chan interface{}
	wg        *sync.WaitGroup
	subs      map[chan interface{}]struct{}
}

func NewBroker() *Broker {
	return &Broker{
		stopCh:    make(chan struct{}),
		publishCh: make(chan interface{}, 1),
		subCh:     make(chan chan interface{}, 1),
		unsubCh:   make(chan chan interface{}, 1),
		wg:        &sync.WaitGroup{},
	}
}

func (b *Broker) Start() {
	b.subs = map[chan interface{}]struct{}{}
	for {
		select {
		case <-b.stopCh:
			for msgCh := range b.subs {
				close(msgCh)
			}
			return
		case msgCh := <-b.subCh:
			b.subs[msgCh] = struct{}{}
		case msgCh := <-b.unsubCh:
			delete(b.subs, msgCh)
		case msg := <-b.publishCh:
			for msgCh := range b.subs {
				// msgCh is buffered, use non-blocking send to protect the broker:
				select {
				case msgCh <- msg:
				default:
				}
			}
		}
	}
}

func (b *Broker) Stop() {
	close(b.stopCh)
}

func (b *Broker) Subscribe() chan interface{} {
	msgCh := make(chan interface{}, 1)
	b.subCh <- msgCh
	return msgCh
}

func (b *Broker) Unsubscribe(msgCh chan interface{}) {
	b.unsubCh <- msgCh
	close(msgCh)
}

func (b *Broker) MsgReceived() {
	b.wg.Done()
}

func (b *Broker) Publish(msg interface{}) {
	b.wg.Add(len(b.subs))

	if nil != msg {
		msgId := msg.(*BrokerConnClient)
		msgId.Id = len(b.subs)
		b.publishCh <- msgId
	} else {
		b.publishCh <- msg
	}

	b.wg.Wait()
}
