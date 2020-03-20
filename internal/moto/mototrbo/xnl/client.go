package xnl

import (
	"fmt"
	"net"

	"github.com/pboyd04/MotoGo/internal/moto/mototrbo"
)

// PacketHandler is a function to handle a particular type of packet
type PacketHandler func(Packet) bool

// Client descripes a XNL Client
type Client struct {
	client  *mototrbo.Client
	id      mototrbo.RadioID
	xnlID   Address      //The remote radio's xnl ID
	myxnlID Address      //My xnl ID
	addr    *net.UDPAddr //The network address of the radio

	handlers      map[OpCode][]handler
	ready         chan bool
	transactionID uint16
	flag          byte

	PacketsIn chan Packet
}

type handler struct {
	handlerFunc PacketHandler
}

// NewClient creates a new Client instance
func NewClient(client *mototrbo.Client, id mototrbo.RadioID, addr *net.UDPAddr) *Client {
	c := new(Client)
	c.client = client
	c.id = id
	c.handlers = make(map[OpCode][]handler)
	c.PacketsIn = make(chan Packet, 5)
	c.ready = make(chan bool, 1)
	c.transactionID = 0x0100
	c.flag = 1
	c.addr = addr
	c.client.RegisterHandler(mototrbo.XnlXcmpPacket, c.gotXNLPacket)
	c.RegisterHandler(MasterStatusBroadcast, c.gotMasterStatusBroadcast)
	c.RegisterHandler(DeviceAuthKeyReply, c.gotDeviceAuthKeyReply)
	c.RegisterHandler(DeviceConnectionReply, c.gotDeviceConnectionReply)
	c.RegisterHandler(DeviceSysMapBroadcast, c.gotDeviceSystemBroadcast)
	c.RegisterHandler(DataMessage, c.ackDataPacket)
	c.RegisterHandler(DataMessageAck, c.gotAckPacket)
	c.init()
	return c
}

func (c *Client) gotXNLPacket(mp mototrbo.Packet, _ net.Addr) bool {
	motoPkt := mp.(mototrbo.XnlPacket)
	xnlPkt := CreatePacketFromArray(motoPkt.Payload)
	opCode := xnlPkt.GetHeader().OpCode
	generic := true
	handlers, ok := c.handlers[opCode]
	if ok {
		for _, handler := range handlers {
			ret := handler.handlerFunc(xnlPkt)
			if ret {
				generic = false
			}
		}
	}
	if generic {
		//Append to packet channel
		fmt.Printf("Got unhandled packet! %#v\n", xnlPkt)
		c.PacketsIn <- xnlPkt
	}
	return true
}

// SendPacket sends an XNL protocol packet to the radio
func (c *Client) SendPacket(pkt Packet) {
	header := pkt.GetHeader()
	if header.OpCode == DataMessage {
		dp := pkt.(DataPacket)
		header.Src = c.myxnlID
		header.Dest = c.xnlID
		header.TransactionID = c.transactionID
		header.Flags = c.flag
		c.flag++
		if c.flag > 0x07 {
			c.flag = 0
		}
		c.transactionID++
		dp.Header = header
		pkt = dp
	}
	mp := mototrbo.NewXNLXCMPPacketByParam(c.id, pkt)
	c.client.SendPacket(mp, c.addr)
}

// RegisterHandler registers a function to be called anytime a certain packet type is recieved
func (c *Client) RegisterHandler(opCode OpCode, handlerFunc PacketHandler) {
	var h handler
	h.handlerFunc = handlerFunc
	if cmdHandle, ok := c.handlers[opCode]; ok {
		cmdHandle = append(cmdHandle, h)
		c.handlers[opCode] = cmdHandle
	} else {
		c.handlers[opCode] = make([]handler, 1)
		c.handlers[opCode][0] = h
	}
}

// GetRadioXNLID gets the XNL ID of the remote radio
func (c *Client) GetRadioXNLID() Address {
	return c.xnlID
}

func (c *Client) init() {
	initPkt := NewInitPacket()
	c.SendPacket(initPkt)
	<-c.ready
}

// GetAuthKey authenticates with XNL on the radio and creates an XCMP connection
func (c *Client) GetAuthKey() {
	c.SendPacket(NewDevAuthKeyRequestPacket(c.xnlID))
}

func (c *Client) startConnection(tmpID uint16, authKey []byte) {
	pkt := NewDevConnectionRequestPacket(c.xnlID, Address(tmpID), Address(0), 0x0A, 0x01, authKey)
	c.SendPacket(pkt)
}

func (c *Client) gotMasterStatusBroadcast(pkt Packet) bool {
	c.xnlID = pkt.GetHeader().Src
	c.ready <- true
	return true
}

func (c *Client) gotDeviceAuthKeyReply(pkt Packet) bool {
	p := pkt.(DevKeyAuthReplyPacket)
	c.startConnection(p.TempID, p.AuthKey)
	return true
}

func (c *Client) gotDeviceConnectionReply(pkt Packet) bool {
	p := pkt.(DevConnectionReplyPacket)
	c.myxnlID = p.AssignedID
	return true
}

func (c *Client) gotDeviceSystemBroadcast(pkt Packet) bool {
	//NOOP for now
	return true
}

func (c *Client) ackDataPacket(pkt Packet) bool {
	header := pkt.GetHeader()
	if header.Dest == Address(0) || header.Dest == c.myxnlID {
		c.SendPacket(NewAckDataPacketByParam(header.Src, c.myxnlID, header.Protocol, header.Flags, header.TransactionID))
	}
	return false
}

func (c *Client) gotAckPacket(pkt Packet) bool {
	//NOOP for now
	return true
}
