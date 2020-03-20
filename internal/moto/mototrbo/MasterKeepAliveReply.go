package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// MasterKeepAliveReplyPacket a request to a Master to register
type MasterKeepAliveReplyPacket struct {
	Command  Command
	ID       RadioID
	Digital  bool
	CSBK     bool
	LinkType LinkType
}

//NewKeepAliveReplyPacketByParam create registration packet from params
func NewKeepAliveReplyPacketByParam(id RadioID, digital bool, csbk bool, linkType LinkType) MasterKeepAliveReplyPacket {
	var m MasterKeepAliveReplyPacket
	m.Command = RegistrationReply
	m.ID = id
	m.Digital = digital
	m.CSBK = csbk
	m.LinkType = linkType
	return m
}

//NewKeepAliveReplyPacketByArray create an registration packet from data array
func NewKeepAliveReplyPacketByArray(data []byte) (MasterKeepAliveReplyPacket, error) {
	var p MasterKeepAliveReplyPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.Digital = ((data[5] & 0x20) != 0)
	p.CSBK = ((data[8] & 0x80) != 0)
	p.LinkType = LinkType(data[10])
	return p, nil
}

// GetCommand returns the command for the packet
func (p MasterKeepAliveReplyPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p MasterKeepAliveReplyPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p MasterKeepAliveReplyPacket) ToArray() []byte {
	a := make([]byte, 14)
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
	a[10] = byte(p.LinkType)
	a[11] = 0x06
	a[12] = byte(p.LinkType)
	return a
}
