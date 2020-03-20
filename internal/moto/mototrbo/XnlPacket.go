package mototrbo

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// XnlPacket an XNL data packet
type XnlPacket struct {
	Command Command
	ID      RadioID
	Payload []byte
}

// XnlPayload an Xnl Packet Payload
type XnlPayload interface {
	ToArray() []byte
}

//NewXNLXCMPPacketByParam create an XNL/XCMP from params and payload
func NewXNLXCMPPacketByParam(id RadioID, payload XnlPayload) XnlPacket {
	var p XnlPacket
	p.Command = XnlXcmpPacket
	p.ID = id
	p.Payload = payload.ToArray()
	return p
}

//NewXNLXCMPPacketByArray create an XNL/XCMP from data array
func NewXNLXCMPPacketByArray(data []byte) (XnlPacket, error) {
	var p XnlPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.Payload = data[5:]
	return p, nil
}

// GetCommand returns the command for the packet
func (p XnlPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p XnlPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p XnlPacket) ToArray() []byte {
	a := make([]byte, 5)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	return append(a, p.Payload...)
}
