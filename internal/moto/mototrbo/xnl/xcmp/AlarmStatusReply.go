package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// Alarm is the alarm number
type Alarm byte

const (
	AlarmTransmit     Alarm = 1
	AlarmReceive      Alarm = 2
	AlarmTemp         Alarm = 3
	AlarmAC           Alarm = 4
	AlarmFan          Alarm = 5
	AlarmVSWR         Alarm = 20
	AlarmTrasmitPower Alarm = 23
	AlarmUnknown1     Alarm = 81
)

// AlarmStatusReply a reply of the radio status
type AlarmStatusReply struct {
	OpCode   OpCode
	unknown  byte
	unknown1 byte
	Alarms   []AlarmStatus
}

// AlarmStatus an alarm status
type AlarmStatus struct {
	Severity byte
	State    byte
	Alarm    Alarm
	Unknown  []byte
}

// NewAlarmStatusReplyByParam create an data packet from data array
func NewAlarmStatusReplyByParam(EntityType EntityType, alarms []AlarmStatus) AlarmStatusReply {
	var p AlarmStatusReply
	p.OpCode = CombinedOpCode(Broadcast, OpCodeDeviceInitStatus)
	p.Alarms = alarms
	return p
}

// NewAlarmStatusReplyByArray create an data packet from data array
func NewAlarmStatusReplyByArray(data []byte) AlarmStatusReply {
	var p AlarmStatusReply
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	p.unknown = data[2]
	p.unknown1 = data[3]
	p.Alarms = make([]AlarmStatus, data[4])
	for i := 0; i < int(data[4]); i++ {
		p.Alarms[i].Severity = data[5+(i*7)]
		p.Alarms[i].State = data[6+(i*7)]
		p.Alarms[i].Alarm = Alarm(data[7+(i*7)])
		p.Alarms[i].Unknown = data[8+(i*7) : 12+(i*7)]
	}
	return p
}

// GetMessageType returns the message type for the packet
func (p AlarmStatusReply) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p AlarmStatusReply) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p AlarmStatusReply) ToArray() []byte {
	a := make([]byte, 4)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = p.unknown
	a[3] = p.unknown1
	for _, as := range p.Alarms {
		a = append(a, as.ToArray()...)
	}
	return a
}

// ToArray converts and AlarmStatus to an array
func (p AlarmStatus) ToArray() []byte {
	a := make([]byte, 3)
	a[0] = p.Severity
	a[1] = p.State
	a[2] = byte(p.Alarm)
	return append(a, p.Unknown...)
}
