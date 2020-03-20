package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// VersionInfoRequest a request for radio status
type VersionInfoRequest struct {
	OpCode OpCode
}

// NewVersionInfoRequestByParam create an data packet from data array
func NewVersionInfoRequestByParam() VersionInfoRequest {
	var p VersionInfoRequest
	p.OpCode = CombinedOpCode(Request, OpCodeVersionInfo)
	return p
}

// NewVersionInfoRequestByArray create an data packet from data array
func NewVersionInfoRequestByArray(data []byte) VersionInfoRequest {
	var p VersionInfoRequest
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	return p
}

// GetMessageType returns the message type for the packet
func (p VersionInfoRequest) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p VersionInfoRequest) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p VersionInfoRequest) ToArray() []byte {
	a := make([]byte, 3)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = 0x00
	return a
}
