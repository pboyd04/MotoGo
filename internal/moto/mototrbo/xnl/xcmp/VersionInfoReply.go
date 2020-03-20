package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// VersionInfoReply a reply of the radio version
type VersionInfoReply struct {
	OpCode  OpCode
	unknown byte
	Version string
}

// NewVersionInfoReplyReplyByParam create an data packet from data array
func NewVersionInfoReplyReplyByParam(version string) VersionInfoReply {
	var p VersionInfoReply
	p.OpCode = CombinedOpCode(Reply, OpCodeDeviceInitStatus)
	p.Version = version
	return p
}

// NewVersionInfoReplyByArray create an data packet from data array
func NewVersionInfoReplyByArray(data []byte) VersionInfoReply {
	var p VersionInfoReply
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	p.unknown = data[2]
	p.Version = string(data[3:])
	return p
}

// GetMessageType returns the message type for the packet
func (p VersionInfoReply) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p VersionInfoReply) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p VersionInfoReply) ToArray() []byte {
	a := make([]byte, 3)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = p.unknown
	return append(a, []byte(p.Version)...)
}
