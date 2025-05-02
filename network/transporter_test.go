package network

import (
	"strconv"
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

func TestSendMessageTo(t *testing.T) {
	ta := NewLocalTransport("A")
	tb := NewLocalTransport("B")

	assert.Nil(t, ta.Connect(tb))
	assert.Nil(t, tb.Connect(ta))

	assert.Equal(t, ta.peers["B"], tb)
	assert.Equal(t, tb.peers["A"], ta)

	payload := []byte("Hello from A")
	msg := Message{
		From:    "A",
		Payload: payload,
	}

	assert.Nil(t, ta.SendMessageTo("B", &msg))
	msgrv := <-tb.ReadChan()

	assert.Equal(t, msgrv, msg)
}

func TestBroadcastMessage(t *testing.T) {
	ts := []*LocalTransport{
		NewLocalTransport("A"),
		NewLocalTransport("B"),
		NewLocalTransport("C"),
		NewLocalTransport("D"),
	}

	for i := 0; i < 4; i++ {
		for j := i + 1; j < 4; j++ {
			assert.Nil(t, ts[i].Connect(ts[j]))
			assert.Nil(t, ts[j].Connect(ts[i]))
		}
	}

	numMsg := 100
	for i := 0; i < numMsg; i++ {
		payload := []byte("Hello No: " + strconv.Itoa(i))
		msg := NewMessage(ts[0].Address(), payload)

		assert.Nil(t, ts[0].BroadCastMessage(msg))

		for j := 1; j < 4; j++ {
			msgrv := <-ts[j].ReadChan()
			assert.Equal(t, &msgrv, msg)
		}
	}
}
