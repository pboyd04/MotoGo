package mototrbo

import "github.com/pboyd04/MotoGo/internal/util"

// UserPacket a request to a Master to register
type UserPacket struct {
	Command     Command
	ID          RadioID
	Source      RadioID
	Destination RadioID
	Encrypted   bool //0x80
	End         bool //0x40
	TimeSlot    byte //0x20
	PhoneCall   bool //0x10
	Payload     []byte
}

//NewUserPacketByParam create user packet from params
func NewUserPacketByParam(cmd Command, id RadioID, src RadioID, dest RadioID, encrypted bool, end bool, timeslot byte, phoneCall bool, payload []byte) UserPacket {
	var m UserPacket
	m.Command = cmd
	m.ID = id
	m.Source = src
	m.Destination = dest
	m.Encrypted = encrypted
	m.End = end
	m.TimeSlot = timeslot
	m.PhoneCall = phoneCall
	m.Payload = payload
	return m
}

//NewUserPacketByArray create an user packet from data array
func NewUserPacketByArray(data []byte) (UserPacket, error) {
	var p UserPacket
	p.Command = Command(data[0])
	p.ID = RadioID(util.ParseUint32(data, 1))
	p.Source = RadioID(util.ParseUint24(data, 6))
	p.Destination = RadioID(util.ParseUint24(data, 9))
	//Unknown data[12:16]
	p.Encrypted = ((data[17] & 0x80) != 0)
	p.End = ((data[17] & 0x40) != 0)
	if (data[17] & 0x20) != 0 {
		p.TimeSlot = 2
	} else {
		p.TimeSlot = 1
	}
	p.PhoneCall = ((data[17] & 0x10) != 0)
	p.Payload = data[18:]
	return p, nil
}

// GetCommand returns the command for the packet
func (p UserPacket) GetCommand() Command {
	return p.Command
}

// GetID returns the RadioID for the packet
func (p UserPacket) GetID() RadioID {
	return p.ID
}

// ToArray converts a packet to a byte array
func (p UserPacket) ToArray() []byte {
	a := make([]byte, 18)
	a[0] = byte(p.Command)
	a[1] = byte(p.ID >> 24)
	a[2] = byte(p.ID >> 16)
	a[3] = byte(p.ID >> 8)
	a[4] = byte(p.ID)
	a[5] = 0x00
	a[6] = byte(p.Source >> 16)
	a[7] = byte(p.Source >> 8)
	a[8] = byte(p.Source)
	a[9] = byte(p.Destination >> 16)
	a[10] = byte(p.Destination >> 8)
	a[11] = byte(p.Destination)
	if p.Encrypted {
		a[17] |= 0x80
	}
	if p.End {
		a[17] |= 0x40
	}
	if p.TimeSlot == 2 {
		a[17] |= 0x20
	}
	if p.PhoneCall {
		a[17] |= 0x10
	}
	return append(a, p.Payload...)
}
