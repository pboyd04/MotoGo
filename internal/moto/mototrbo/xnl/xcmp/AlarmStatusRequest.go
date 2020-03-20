package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// AlarmStatusRequest a request for radio status
type AlarmStatusRequest struct {
	OpCode OpCode
}

// NewAlarmStatusRequestByParam create an data packet from data array
func NewAlarmStatusRequestByParam() AlarmStatusRequest {
	var p AlarmStatusRequest
	p.OpCode = CombinedOpCode(Request, OpCodeAlarm)
	return p
}

// NewAlarmStatusRequestByArray create an data packet from data array
func NewAlarmStatusRequestByArray(data []byte) AlarmStatusRequest {
	var p AlarmStatusRequest
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	return p
}

// GetMessageType returns the message type for the packet
func (p AlarmStatusRequest) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p AlarmStatusRequest) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p AlarmStatusRequest) ToArray() []byte {
	a := make([]byte, 3)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = 0x04
	return a
}
