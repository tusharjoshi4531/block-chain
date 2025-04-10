package network

type MessageHeader struct {
	From string
}

type Message struct {
	Header  MessageHeader
	Payload []byte
}

type Transporter interface {
	Receive() <-chan Message
	Connect(Transporter) error
	SendMessage(string, *Message) error
	Address() string
}
