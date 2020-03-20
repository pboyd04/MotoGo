package xcmp

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// DeviceInitStatusBroadcast a broadcast of the device init status
type DeviceInitStatusBroadcast struct {
	OpCode       OpCode
	MajorVersion byte
	MinorVersion byte
	RevVersion   byte
	EntityType   EntityType
	InitComplete bool
	Status       DeviceStatus
}

// DeviceStatus is the status of a child device
type DeviceStatus struct {
	DeviceType   byte
	DeviceStatus uint16
	Descriptor   map[byte]byte
}

// NewDeviceInitStatusBroadcastByParam create an data packet from data array
func NewDeviceInitStatusBroadcastByParam(EntityType EntityType, InitComplete bool, status DeviceStatus) DeviceInitStatusBroadcast {
	var p DeviceInitStatusBroadcast
	p.OpCode = CombinedOpCode(Broadcast, OpCodeDeviceInitStatus)
	p.MajorVersion = 1
	p.MinorVersion = 1
	p.RevVersion = 0
	p.EntityType = EntityType
	p.InitComplete = InitComplete
	p.Status = status
	return p
}

// NewDeviceInitStatusBroadcastByArray create an data packet from data array
func NewDeviceInitStatusBroadcastByArray(data []byte) DeviceInitStatusBroadcast {
	var p DeviceInitStatusBroadcast
	p.OpCode = OpCode(util.ParseUint16(data, 0))
	p.MajorVersion = data[2]
	p.MinorVersion = data[3]
	p.RevVersion = data[4]
	p.EntityType = EntityType(data[5])
	p.InitComplete = (data[6] == 0x01)
	if !p.InitComplete {
		p.Status.FromArray(data[7:])
	}
	return p
}

// GetMessageType returns the message type for the packet
func (p DeviceInitStatusBroadcast) GetMessageType() MessageType {
	return MessageType((p.OpCode >> 8) & 0xF0)
}

// GetOpCode returns the message opcode for the packet
func (p DeviceInitStatusBroadcast) GetOpCode() OpCode {
	return p.OpCode & 0x0FFF
}

// ToArray converts a packet to a byte array
func (p DeviceInitStatusBroadcast) ToArray() []byte {
	a := make([]byte, 8)
	a[0] = byte(p.OpCode >> 8)
	a[1] = byte(p.OpCode)
	a[2] = p.MajorVersion
	a[3] = p.MinorVersion
	a[4] = p.RevVersion
	a[5] = byte(p.EntityType)
	if p.InitComplete {
		a[6] = 0x01
	}
	return append(a, p.Status.ToArray()...)
}

// FromArray converts an array to a DeviceStatus
func (d *DeviceStatus) FromArray(data []byte) {
	d.DeviceType = data[0]
	d.DeviceStatus = util.ParseUint16(data, 1)
	length := data[3]
	if length > 0 {
		d.Descriptor = make(map[byte]byte)
		for index := 0; index < int(length); index += 2 {
			d.Descriptor[data[4+index]] = data[4+index+1]
		}
	}
}

// ToArray converts a DeviceStatus to an array
func (d *DeviceStatus) ToArray() []byte {
	a := make([]byte, 4)
	a[0] = d.DeviceType
	a[1] = byte(d.DeviceStatus >> 8)
	a[2] = byte(d.DeviceStatus)
	length := len(d.Descriptor) * 2
	a[3] = byte(length)
	if length > 0 {
		for key, value := range d.Descriptor {
			a = append(a, key, value)
		}
	}
	return a
}
