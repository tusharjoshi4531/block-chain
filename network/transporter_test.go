package network

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConnect(t *testing.T) {
	ta := NewLocalTransport("A")
	tb := NewLocalTransport("B")
	
	assert.Nil(t, ta.Connect(tb))
	assert.Nil(t, tb.Connect(ta))

	assert.Equal(t, ta.peers["B"], tb)
	assert.Equal(t, tb.peers["A"], ta)
}

func TestSendMessage(t *testing.T) {
	ta := NewLocalTransport("A")
	tb := NewLocalTransport("B")
	
	assert.Nil(t, ta.Connect(tb))
	assert.Nil(t, tb.Connect(ta))

	assert.Equal(t, ta.peers["B"], tb)
	assert.Equal(t, tb.peers["A"], ta)

	payload := []byte("Hello from A")
	msg := Message {
		Header: MessageHeader{
			From: "A",
		},
		Payload: payload,
	}

	assert.Nil(t, ta.SendMessage("B", &msg))
	msgrv := <-tb.Receive()

	assert.Equal(t, msgrv, msg)
}

