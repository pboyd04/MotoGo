package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// MasterRegistrationReplyPacket a request to a Master to register
type MasterRegistrationReplyPacket struct {
	Command  Command
	ID       RadioID
	Digital  bool
	CSBK     bool
	LinkType LinkType
}

//NewRegistrationReplyPacketByParam create registration packet from params
func NewRegistrationReplyPacketByParam(id RadioID, digital bool, csbk bool, linkType LinkType) MasterRegistrationReplyPacket {
	var m MasterRegistrationReplyPacket
	m.Command = MasterKeepAliveReply
	m.ID = id
	m.Digital = digital
	m.CSBK = csbk
	m.LinkType = linkType
	return m
}

//NewRegistrationReplyPacketByArray create an registration packet from data array
func NewRegistrationReplyPacketByArray(data []byte) (MasterRegistrationReplyPacket, error) {
	var p MasterRegistrationReplyPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.Digital = ((data[5] & 0x20) != 0)
	p.CSBK = ((data[8] & 0x80) != 0)
	p.LinkType = LinkType(data[10])
	return p, nil
}

// GetCommand returns the command for the packet
func (p MasterRegistrationReplyPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p MasterRegistrationReplyPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p MasterRegistrationReplyPacket) ToArray() []byte {
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
