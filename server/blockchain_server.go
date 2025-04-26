package server

import "github.com/tusharjoshi4531/block-chain.git/network"

type BlockChainServer interface {
	Listen()
	ConnectPeer(network.Transport) error
	Kill()
}
