package mototrbo

import (
	"fmt"

	"github.com/pboyd04/MotoGo/internal/util"
)

// UnknownPacket a Mototrbo packet of unknown type
type UnknownPacket struct {
	Command Command
	ID      RadioID
	Payload []byte
}

// NewUnknownPacket create an unknown packet from a byte array
func NewUnknownPacket(array []byte) (UnknownPacket, error) {
	var p UnknownPacket
	if len(array) < 5 {
		return p, fmt.Errorf("error packet isn't big enough %#v", array)
	}
	p.Command = Command(array[0])
	p.ID = RadioID(util.ParseUint32(array, 1))
	p.Payload = array[5:]
	return p, nil
}

// GetCommand returns the command for the packet
func (p UnknownPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p UnknownPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p UnknownPacket) ToArray() []byte {
	a := make([]byte, 5)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	return append(a, p.Payload...)
}
