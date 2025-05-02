package network

type LocalTransport struct {
	*DefaultTransport
}

func NewLocalTransport(address string) *LocalTransport {
	return &LocalTransport{
		DefaultTransport: NewDefaultTransport(address),
	}
}

func (t *LocalTransport) SendMessage(msg *Message) error {
	t.messageChanel <- *msg
	return nil
}
