package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// PeerKeepAliveReplyPacket a request to a Master to register
type PeerKeepAliveReplyPacket struct {
	Command Command
	ID      RadioID
	Digital bool
	CSBK    bool
}

//NewPeerKeepAliveReplyPacketByParam create registration packet from params
func NewPeerKeepAliveReplyPacketByParam(id RadioID, digital bool, csbk bool, linkType LinkType) PeerKeepAliveReplyPacket {
	var m PeerKeepAliveReplyPacket
	m.Command = PeerKeepAliveReply
	m.ID = id
	m.Digital = digital
	m.CSBK = csbk
	return m
}

//NewPeerKeepAliveReplyPacketByArray create an registration packet from data array
func NewPeerKeepAliveReplyPacketByArray(data []byte) (PeerKeepAliveReplyPacket, error) {
	var p PeerKeepAliveReplyPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.Digital = ((data[5] & 0x20) != 0)
	p.CSBK = ((data[8] & 0x80) != 0)
	return p, nil
}

// GetCommand returns the command for the packet
func (p PeerKeepAliveReplyPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p PeerKeepAliveReplyPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p PeerKeepAliveReplyPacket) ToArray() []byte {
	a := make([]byte, 10)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	a[5] = 0x45
	if p.Digital {
		a[5] |= 0x20
	} else {
		a[5] |= 0x10
	}
	if p.CSBK {
		a[8] |= 0x80
	}
	a[8] |= 0x20
	a[9] = 0x2C
	return a
}
