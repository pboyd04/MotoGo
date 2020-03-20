package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// MasterKeepAliveRequestPacket a request to a Master to register
type MasterKeepAliveRequestPacket struct {
	Command  Command
	ID       RadioID
	Digital  bool
	CSBK     bool
	LinkType LinkType
}

//NewMasterKeepAliveRequestPacketByParam create registration packet from params
func NewMasterKeepAliveRequestPacketByParam(id RadioID, digital bool, csbk bool, linkType LinkType) MasterKeepAliveRequestPacket {
	var m MasterKeepAliveRequestPacket
	m.Command = MasterKeepAliveRequest
	m.ID = id
	m.Digital = digital
	m.CSBK = csbk
	m.LinkType = linkType
	return m
}

//NewMasterKeepAliveRequestPacketByArray create an registration packet from data array
func NewMasterKeepAliveRequestPacketByArray(data []byte) (MasterKeepAliveRequestPacket, error) {
	var p MasterKeepAliveRequestPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.Digital = ((data[5] & 0x20) != 0)
	p.CSBK = ((data[8] & 0x80) != 0)
	p.LinkType = LinkType(data[10])
	return p, nil
}

// GetCommand returns the command for the packet
func (p MasterKeepAliveRequestPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p MasterKeepAliveRequestPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p MasterKeepAliveRequestPacket) ToArray() []byte {
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
