package mototrbo

import (
	"fmt"

	"github.com/pboyd04/MotoGo/internal/util"
)

// PeerListReplyPacket a request to a master for a list of peers
type PeerListReplyPacket struct {
	Command Command
	ID      RadioID
	Peers   []Peer
}

// Peer is a peer radio
type Peer struct {
	ID      RadioID
	Address string
	Port    uint16
	Mode    byte
}

// NewPeerListReplyPacketByParam create peer list request packet by params
func NewPeerListReplyPacketByParam(id RadioID, peers []Peer) PeerListReplyPacket {
	var m PeerListReplyPacket
	m.Command = PeerListReply
	m.ID = id
	m.Peers = peers
	return m
}

// NewPeerListReplyPacketByArray create an registration packet from data array
func NewPeerListReplyPacketByArray(data []byte) (PeerListReplyPacket, error) {
	var p PeerListReplyPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	length := util.ParseUint16(data, 5)
	count := length / 11
	p.Peers = make([]Peer, count)
	var i uint16
	for i = 0; i < count; i++ {
		array := data[7+(i*11):]
		p.Peers[i].ID = RadioID(util.ParseUint32(array, 0))
		p.Peers[i].Address = fmt.Sprintf("%d.%d.%d.%d", array[4], array[5], array[6], array[7])
		p.Peers[i].Port = util.ParseUint16(array, 8)
		p.Peers[i].Mode = array[10]
	}
	return p, nil
}

// GetCommand returns the command for the packet
func (p PeerListReplyPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p PeerListReplyPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p PeerListReplyPacket) ToArray() []byte {
	peerCount := len(p.Peers)
	a := make([]byte, 7+(peerCount*11))
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	a[5] = byte(peerCount >> 8)
	a[6] = byte(peerCount)
	index := 7
	for _, peer := range p.Peers {
		a[index] = byte(peer.ID >> 24)
		a[index+1] = byte(peer.ID >> 16)
		a[index+2] = byte(peer.ID >> 8)
		a[index+3] = byte(peer.ID)
		fmt.Scanf("%d.%d.%d.%d", &a[index+4], &a[index+5], &a[index+6], &a[index+7])
		a[index+8] = byte(peer.Port >> 8)
		a[index+9] = byte(peer.Port)
		a[index+10] = peer.Mode
		index += 11
	}
	return a
}
