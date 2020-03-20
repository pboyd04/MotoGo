package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// PeerRegistrationRequestPacket a request to a peer to register
type PeerRegistrationRequestPacket struct {
	Command  Command
	ID       RadioID
	LinkType LinkType
}

//NewPeerRegistrationPacketByParam create registration packet from params
func NewPeerRegistrationPacketByParam(id RadioID, linkType LinkType) PeerRegistrationRequestPacket {
	var m PeerRegistrationRequestPacket
	m.Command = PeerRegisterRequest
	m.ID = id
	m.LinkType = linkType
	return m
}

//NewPeerRegistrationPacketByArray create an registration packet from data array
func NewPeerRegistrationPacketByArray(data []byte) (PeerRegistrationRequestPacket, error) {
	var p PeerRegistrationRequestPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.LinkType = LinkType(data[5])
	return p, nil
}

// GetCommand returns the command for the packet
func (p PeerRegistrationRequestPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p PeerRegistrationRequestPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p PeerRegistrationRequestPacket) ToArray() []byte {
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
