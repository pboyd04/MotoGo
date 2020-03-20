package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// PeerRegistrationReplyPacket a reply to a peer registeration
type PeerRegistrationReplyPacket struct {
	Command  Command
	ID       RadioID
	LinkType LinkType
}

//NewPeerRegistrationReplyPacketByParam create registration packet from params
func NewPeerRegistrationReplyPacketByParam(id RadioID, linkType LinkType) PeerRegistrationReplyPacket {
	var m PeerRegistrationReplyPacket
	m.Command = PeerRegisterReply
	m.ID = id
	m.LinkType = linkType
	return m
}

//NewPeerRegistrationReplyPacketByArray create an registration packet from data array
func NewPeerRegistrationReplyPacketByArray(data []byte) (PeerRegistrationReplyPacket, error) {
	var p PeerRegistrationReplyPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.LinkType = LinkType(data[5])
	return p, nil
}

// GetCommand returns the command for the packet
func (p PeerRegistrationReplyPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p PeerRegistrationReplyPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p PeerRegistrationReplyPacket) ToArray() []byte {
	a := make([]byte, 9)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	a[5] = byte(p.LinkType)
	a[6] = 0x06
	a[7] = byte(p.LinkType)
	return a
}
