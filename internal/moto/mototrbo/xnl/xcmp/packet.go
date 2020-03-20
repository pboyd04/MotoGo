package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// MessageType is the type of the XCMP message
type MessageType byte

// OpCode is the operation of the XCMP message
type OpCode uint16

// EntityType reprsents the kind of XCMP entity
type EntityType byte

const (
	// Request is an XCMP request for data
	Request MessageType = 0x00
	// Reply is an XCMP reply to a request
	Reply MessageType = 0x80
	// Broadcast is an XCMP message to everyone
	Broadcast MessageType = 0xB0

	OpCodeRadioStatus      OpCode = 0x00E
	OpCodeVersionInfo      OpCode = 0x00F
	OpCodeDeviceInitStatus OpCode = 0x400
	OpCodeDisplayText      OpCode = 0x401
	OpCodeAlarm            OpCode = 0x42E
)

// Packet is an XCMP packet
type Packet interface {
	GetMessageType() MessageType
	GetOpCode() OpCode
	ToArray() []byte
}

// CreatePacketFromArray creates a XCMP packet from an array
func CreatePacketFromArray(array []byte) Packet {
	opCode := OpCode(util.ParseUint16(array, 0))
	switch opCode {
	case OpCode(0xB400):
		return NewDeviceInitStatusBroadcastByArray(array)
	case OpCode(0x800E):
		return NewRadioStatusReplyByArray(array)
	case OpCode(0x800F):
		return NewVersionInfoReplyByArray(array)
	case OpCode(0x842e):
		return NewAlarmStatusReplyByArray(array)
	default:
		return NewUnkownPacketByArray(array)
	}
}

// CombinedOpCode combines the message type and opcode
func CombinedOpCode(mt MessageType, opCode OpCode) OpCode {
	m := OpCode(mt) << 8
	return m | opCode
}
