package network

import (
	"fmt"
	"sync"
)

type LocalTransport struct {
	address       string
	peers         map[string]Transport
	lock          sync.RWMutex
	messageChanel chan Message
}

func NewLocalTransport(address string) *LocalTransport {
	return &LocalTransport{
		address:       address,
		peers:         make(map[string]Transport),
		messageChanel: make(chan Message, 1024),
	}
}

func (t *LocalTransport) ReadChan() <-chan Message {
	return t.messageChanel
}

func (t *LocalTransport) WriteChan() chan<- Message {
	return t.messageChanel
}

func (t *LocalTransport) Connect(otherTransport Transport) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[otherTransport.Address()] = otherTransport

	return nil
}

func (t *LocalTransport) SendMessage(to string, msg *Message) error {
	t.lock.RLock()
	peer, ok := t.peers[to]
	t.lock.RUnlock()

	if !ok {
		return fmt.Errorf("sender (%s) is not connected to receiver (%s)", t.Address(), to)
	}

	peer.WriteChan() <- *msg
	return nil
}

func (t *LocalTransport) BroadCastMessage(msg *Message) error {
	for k := range t.peers {
		t.SendMessage(k, msg)
	}
	return nil
}

func (t *LocalTransport) Address() string {
	return t.address
}
