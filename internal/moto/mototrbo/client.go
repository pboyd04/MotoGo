package mototrbo

import (
	"fmt"
	"net"
)

// LinkType the type of link to the system
type LinkType byte

// PacketHandler is a function to handle a particular type of packet
type PacketHandler func(Packet, net.Addr) bool

const (
	// IPSiteConnect the system is an IP Site Connect System
	IPSiteConnect LinkType = 0x04
	// CapacityPlus the system is a Capacity Plus System
	CapacityPlus LinkType = 0x08
)

// Client descripes a Mototrbo Network Client
type Client struct {
	conn *net.UDPConn

	handlers map[Command][]handler

	PacketsIn chan Packet
}

type handler struct {
	handlerFunc PacketHandler
}

// NewClient creates a new Client instance
func NewClient(address string) (*Client, error) {
	c := new(Client)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}
	c.conn = conn
	c.handlers = make(map[Command][]handler)
	c.PacketsIn = make(chan Packet, 5)
	go c.startListener()
	return c, nil
}

// Close ends the UDP connecton
func (c *Client) Close() {
	if c.conn != nil {
		c.conn.Close()
	}
}

// SendPacket sends a MotoPacket to the radio
func (c *Client) SendPacket(pkt Packet, addr *net.UDPAddr) {
	c.conn.WriteTo(pkt.ToArray(), addr)
}

// RegisterHandler registers a function to be called anytime a certain packet type is recieved
func (c *Client) RegisterHandler(cmd Command, handlerFunc PacketHandler) {
	var h handler
	h.handlerFunc = handlerFunc
	if cmdHandle, ok := c.handlers[cmd]; ok {
		cmdHandle = append(cmdHandle, h)
		c.handlers[cmd] = cmdHandle
	} else {
		c.handlers[cmd] = make([]handler, 1)
		c.handlers[cmd][0] = h
	}
}

func (c *Client) startListener() {
	for {
		// Get packet
		raw := make([]byte, 1024)
		size, addr, _ := c.conn.ReadFrom(raw)
		raw = raw[0 : size+1]
		pkt, err := CreatePacketFromArray(raw)
		if err != nil {
			fmt.Printf("Error: %#v from %v\n", err, addr)
			continue
		}
		// Check if we have a handler
		cmd := pkt.GetCommand()
		generic := true
		handlers, ok := c.handlers[cmd]
		if ok {
			for _, handler := range handlers {
				ret := handler.handlerFunc(pkt, addr)
				if ret {
					generic = false
				}
			}
		}
		if generic {
			//Append to packet channel
			fmt.Printf("Got unhandled packet! %#v\n", pkt)
			c.PacketsIn <- pkt
		}
	}
}
