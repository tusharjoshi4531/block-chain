package network

import (
	"encoding/gob"
	"fmt"
	"io"
	"sync"

	"github.com/tusharjoshi4531/block-chain.git/util"
)

type Message struct {
	// Header
	From    string
	Payload []byte
}

func NewMessage(from string, payload []byte) *Message {
	return &Message{
		From:    from,
		Payload: payload,
	}
}

func (msg *Message) Encode(w io.Writer) error {
	return gob.NewEncoder(w).Encode(msg)
}

func (msg *Message) Decode(r io.Reader) error {
	return gob.NewDecoder(r).Decode(msg)
}

func (msg *Message) Bytes() ([]byte, error) {
	return util.EncodeToBytes(msg)
}

type TransportInterface interface {
	Address() string
	SendMessage(*Message) error
}

type Transport interface {
	BroadCastMessage(*Message) error
	Address() string
	SendMessageTo(string, *Message) error
	ReadChan() <-chan Message
	WriteChan() chan<- Message
	Connect(TransportInterface) error
}

type DefaultTransport struct {
	address       string
	peers         map[string]TransportInterface
	lock          sync.RWMutex
	messageChanel chan Message
}

func NewDefaultTransport(address string) *DefaultTransport {
	return &DefaultTransport{
		address:       address,
		peers:         make(map[string]TransportInterface),
		messageChanel: make(chan Message, 1024),
	}
}

func (t *DefaultTransport) ReadChan() <-chan Message {
	return t.messageChanel
}

func (t *DefaultTransport) WriteChan() chan<- Message {
	return t.messageChanel
}

func (t *DefaultTransport) Connect(otherTransport TransportInterface) error {
	t.lock.Lock()
	defer t.lock.Unlock()

	t.peers[otherTransport.Address()] = otherTransport

	return nil
}

func (t *DefaultTransport) SendMessageTo(to string, msg *Message) error {
	t.lock.RLock()
	peer, ok := t.peers[to]
	t.lock.RUnlock()

	if !ok {
		return fmt.Errorf("sender (%s) is not connected to receiver (%s)", t.Address(), to)
	}

	peer.SendMessage(msg)
	return nil
}

func (t *DefaultTransport) BroadCastMessage(msg *Message) error {
	for k := range t.peers {
		t.SendMessageTo(k, msg)
	}
	return nil
}

func (t *DefaultTransport) Address() string {
	return t.address
}
