package tcp

import (
	"bytes"
	"github.com/tusharjoshi4531/block-chain.git/network"
	"io"
	"net"
)

type TcpTransportInterface struct {
	address string
}

func NewTcpTransportInterface(addr string) *TcpTransportInterface {
	return &TcpTransportInterface{
		address: addr,
	}
}

func (ti *TcpTransportInterface) Address() string {
	return ti.address
}

func (ti *TcpTransportInterface) SendMessage(msg *network.Message) error {
	conn, err := net.Dial("tcp", ti.address)
	if err != nil {
		return err
	}
	defer conn.Close()

	payload, err := msg.Bytes()
	if err != nil {
		return err
	}

	_, err = io.Copy(conn, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	return conn.(*net.TCPConn).CloseWrite()
}

type TcpTransportServer = network.DefaultTransport
