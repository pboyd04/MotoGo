package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// RadioStatusReply a reply of the radio status
type RadioStatusReply struct {
	OpCode  OpCode
	unknown byte
	Status  StatusType
	Data    []byte
}

// NewRadioStatusReplyByParam create an data packet from data array
func NewRadioStatusReplyByParam(EntityType EntityType, status StatusType, data []byte) RadioStatusReply {
	var p RadioStatusReply
	p.OpCode = CombinedOpCode(Broadcast, OpCodeDeviceInitStatus)
	p.Status = status
	p.Data = data
	return p
}

// NewRadioStatusReplyByArray create an data packet from data array
func NewRadioStatusReplyByArray(data []byte) RadioStatusReply {
	var p RadioStatusReply
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	p.unknown = data[2]
	p.Status = StatusType(data[3])
	p.Data = data[4:]
	return p
}

// GetMessageType returns the message type for the packet
func (p RadioStatusReply) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p RadioStatusReply) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p RadioStatusReply) ToArray() []byte {
	a := make([]byte, 4)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = p.unknown
	a[3] = byte(p.Status)
	return append(a, p.Data...)
}
