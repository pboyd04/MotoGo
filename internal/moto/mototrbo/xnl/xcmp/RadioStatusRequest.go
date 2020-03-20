package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// StatusType is the status type
type StatusType byte

const (
	RSSI         StatusType = 0x02
	ModelNumber  StatusType = 0x07
	SerialNumber StatusType = 0x0B
	RadioAlias   StatusType = 0x0F
)

// RadioStatusRequest a request for radio status
type RadioStatusRequest struct {
	OpCode OpCode
	Status StatusType
}

// NewRadioStatusRequestByParam create an data packet from data array
func NewRadioStatusRequestByParam(status StatusType) RadioStatusRequest {
	var p RadioStatusRequest
	p.OpCode = CombinedOpCode(Request, OpCodeRadioStatus)
	p.Status = status
	return p
}

// NewRadioStatusRequestByArray create an data packet from data array
func NewRadioStatusRequestByArray(data []byte) RadioStatusRequest {
	var p RadioStatusRequest
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	p.Status = StatusType(data[2])
	return p
}

// GetMessageType returns the message type for the packet
func (p RadioStatusRequest) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p RadioStatusRequest) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p RadioStatusRequest) ToArray() []byte {
	a := make([]byte, 3)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = byte(p.Status)
	return a
}
