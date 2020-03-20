package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// PeerListRequestPacket a request to a master for a list of peers
type PeerListRequestPacket struct {
	Command Command
	ID      RadioID
}

// NewPeerListRequestPacketByParam create peer list request packet by params
func NewPeerListRequestPacketByParam(id RadioID) PeerListRequestPacket {
	var m PeerListRequestPacket
	m.Command = PeerListRequest
	m.ID = id
	return m
}

// NewPeerListRequestPacketByArray create an registration packet from data array
func NewPeerListRequestPacketByArray(data []byte) (PeerListRequestPacket, error) {
	var p PeerListRequestPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	return p, nil
}

// GetCommand returns the command for the packet
func (p PeerListRequestPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p PeerListRequestPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p PeerListRequestPacket) ToArray() []byte {
	a := make([]byte, 5)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	return a
}
