package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// UnknownPacket an unknown packet
type UnknownPacket struct {
	OpCode  OpCode
	Payload []byte
}

// NewUnkownPacketByArray create an data packet from data array
func NewUnkownPacketByArray(data []byte) UnknownPacket {
	var p UnknownPacket
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	p.Payload = data[2:]
	return p
}

// GetMessageType returns the message type for the packet
func (p UnknownPacket) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p UnknownPacket) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p UnknownPacket) ToArray() []byte {
	a := make([]byte, 2)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	return append(a, p.Payload...)
}
