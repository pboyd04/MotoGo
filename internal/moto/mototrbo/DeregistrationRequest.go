package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// DeregistrationRequestPacket a request to a master to deregister
type DeregistrationRequestPacket struct {
	Command Command
	ID      RadioID
}

//NewDeregistrationPacketByParam create deregistration packet from params
func NewDeregistrationPacketByParam(id RadioID) DeregistrationRequestPacket {
	var m DeregistrationRequestPacket
	m.Command = DeregisterRequest
	m.ID = id
	return m
}

//NewDeregistrationPacketByArray create an deregistration packet from data array
func NewDeregistrationPacketByArray(data []byte) (DeregistrationRequestPacket, error) {
	var p DeregistrationRequestPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	return p, nil
}

// GetCommand returns the command for the packet
func (p DeregistrationRequestPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p DeregistrationRequestPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p DeregistrationRequestPacket) ToArray() []byte {
	a := make([]byte, 5)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	return a
}
