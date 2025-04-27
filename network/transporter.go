package network

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

type Transport interface {
	ReadChan() <-chan Message
	WriteChan() chan<- Message
	Connect(Transport) error
	SendMessage(string, *Message) error
	BroadCastMessage(*Message) error
	Address() string
}
