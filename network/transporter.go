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
	Receive() <-chan Message
	Chan() chan<- Message
	Connect(Transport) error
	SendMessage(string, *Message) error
	BroadCastMessage(*Message) error
	Address() string
}
