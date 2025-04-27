package bcnetwork

import "github.com/google/uuid"

const (
	INIT int = iota
	OPEN
	CLOSED
)

type Request struct {
	Status   int
	ID       string
	Payload  *BCPayload
	OnData   func(*BCPayload)
	OnFinish func() 
	OnError  func(error)
}

func NewRequest(payload *BCPayload, onData func(*BCPayload), onFinish func(), onError func(error)) *Request {
	return &Request{
		Status:   INIT,
		ID:       uuid.New().String(),
		Payload:  payload,
		OnData:   onData,
		OnFinish: onFinish,
		OnError:  onError,
	}
}
