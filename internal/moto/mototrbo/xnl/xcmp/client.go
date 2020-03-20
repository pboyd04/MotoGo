package xcmp

import (
	"fmt"

	"github.com/pboyd04/MotoGo/internal/moto/mototrbo/xnl"
)

// PacketHandler is a function to handle a particular type of packet
type PacketHandler func(Packet) bool

// Client descripes a XCMP Client
type Client struct {
	client *xnl.Client

	handlers map[MessageType]map[OpCode][]handler
	ready    chan bool

	PacketsIn  chan Packet
	Version    string
	EntityType EntityType
}

type handler struct {
	handlerFunc PacketHandler
}

// NewClient creates a new Client instance
func NewClient(client *xnl.Client) *Client {
	c := new(Client)
	c.client = client
	c.handlers = make(map[MessageType]map[OpCode][]handler)
	c.PacketsIn = make(chan Packet, 5)
	c.ready = make(chan bool, 1)
	c.handlers[Request] = make(map[OpCode][]handler)
	c.handlers[Reply] = make(map[OpCode][]handler)
	c.handlers[Broadcast] = make(map[OpCode][]handler)
	c.client.RegisterHandler(xnl.DataMessage, c.gotXNLPacket)
	c.RegisterHandler(Broadcast, OpCodeDeviceInitStatus, c.gotDeviceInitBroadcast)
	c.client.GetAuthKey()
	<-c.ready
	return c
}

func (c *Client) gotXNLPacket(pkt xnl.Packet) bool {
	dp := pkt.(xnl.DataPacket)
	if dp.Header.Protocol == xnl.ProtocolXCMP {
		xcmpPkt := CreatePacketFromArray(dp.Payload.ToArray())
		mt := xcmpPkt.GetMessageType()
		opcode := xcmpPkt.GetOpCode()
		generic := true
		handlers, ok := c.handlers[mt][opcode]
		if ok {
			for _, handler := range handlers {
				ret := handler.handlerFunc(xcmpPkt)
				if ret {
					generic = false
				}
			}
		}
		if generic {
			//Append to packet channel
			//fmt.Printf("Got unhandled XCMP message! %#v\n", xcmpPkt)
			c.PacketsIn <- xcmpPkt
		}
		return true
	}
	return false
}

// RegisterHandler registers a function to be called anytime a certain packet type is recieved
func (c *Client) RegisterHandler(messagetype MessageType, opcode OpCode, handlerFunc PacketHandler) {
	var h handler
	h.handlerFunc = handlerFunc
	if cmdHandle, ok := c.handlers[messagetype][opcode]; ok {
		cmdHandle = append(cmdHandle, h)
		c.handlers[messagetype][opcode] = cmdHandle
	} else {
		c.handlers[messagetype][opcode] = make([]handler, 1)
		c.handlers[messagetype][opcode][0] = h
	}
}

// SendPacket sends an XCMP protocol packet to the radio
func (c *Client) SendPacket(pkt Packet) {
	mp := xnl.NewDataPacket(xnl.Address(0), xnl.Address(0), xnl.ProtocolXCMP, pkt)
	c.client.SendPacket(mp)
}

func (c *Client) gotDeviceInitBroadcast(pkt Packet) bool {
	initBroadcast := pkt.(DeviceInitStatusBroadcast)
	c.Version = fmt.Sprintf("%d.%d.%d", initBroadcast.MajorVersion, initBroadcast.MinorVersion, initBroadcast.RevVersion)
	c.EntityType = initBroadcast.EntityType
	if initBroadcast.InitComplete {
		c.ready <- true
	}
	return true
}
