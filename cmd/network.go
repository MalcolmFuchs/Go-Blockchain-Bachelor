package cmd

import (
	"fmt"
	"net"
	"sync"
)

type Network struct {
	Nodes map[string]net.Conn
	mutex sync.Mutex
}

func NewNetwork() *Network {
	return &Network{
		Nodes: make(map[string]net.Conn),
	}
}

func (n *Network) AddNode(address string, conn net.Conn) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	n.Nodes[address] = conn
	fmt.Printf("Node %s added to the network\n", address)
}

func (n *Network) RemoveNode(address string) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	if _, exists := n.Nodes[address]; exists {
		delete(n.Nodes, address)
		fmt.Printf("Node %s removed from the network\n", address)
	}
}

func (n *Network) BroadcastMessage(message []byte) {
	n.mutex.Lock()
	defer n.mutex.Unlock()
	for address, conn := range n.Nodes {
		_, err := conn.Write(message)
		if err != nil {
			fmt.Printf("Failed to send message to node %s: %v\n", address, err)
		}
	}
}
