package mototrbo

import (
	"fmt"

	"github.com/pboyd04/MotoGo/internal/moto/mototrbo/burst"
	"github.com/pboyd04/MotoGo/internal/util"
)

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
	RTP         RTPData
	Payload     burst.Burst
}

// RTPData represents the Real Time Protocol data sent by the radio
type RTPData struct {
	Version        byte
	Padding        bool
	Extension      bool
	CSRCCount      byte
	Marker         bool
	PayloadType    byte
	SequenceNumber uint16
	Timestamp      uint32
	SSRCID         uint32
}

//NewUserPacketByParam create user packet from params
func NewUserPacketByParam(cmd Command, id RadioID, src RadioID, dest RadioID, encrypted bool, end bool, timeslot byte, phoneCall bool, rtp RTPData, payload burst.Burst) UserPacket {
	var m UserPacket
	m.Command = cmd
	m.ID = id
	m.Source = src
	m.Destination = dest
	m.Encrypted = encrypted
	m.End = end
	m.TimeSlot = timeslot
	m.PhoneCall = phoneCall
	m.RTP = rtp
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
	p.RTP = RTPFromArray(data[18:])
	if !p.RTP.Extension {
		p.Payload = BurstFromArray(data[30:])
	} else {
		fmt.Printf("Have a header extenstion! Don't know how to process packet! %#v\n", data[18:])
	}
	return p, nil
}

//RTPFromArray converts RTP Data from an array
func RTPFromArray(data []byte) RTPData {
	var p RTPData
	p.Version = data[0] >> 6
	p.Padding = (data[0] & 0x20) != 0
	p.Extension = (data[0] & 0x10) != 0
	p.CSRCCount = data[0] & 0x0F
	p.Marker = (data[1] & 0x80) != 0
	p.PayloadType = data[1] & 0x7F
	p.SequenceNumber = util.ParseUint16(data, 2)
	p.Timestamp = util.ParseUint32(data, 4)
	p.SSRCID = util.ParseUint32(data, 8)
	return p
}

//BurstFromArray convers Burst Data from an array
func BurstFromArray(data []byte) burst.Burst {
	dt := burst.DataType(data[0] & 0x3F)
	switch dt {
	case burst.DataTypeVoiceLCHeader:
		return burst.NewVoiceHeaderBurstFromArray(data)
	case burst.DataTypeTerminatorWithLC:
		return burst.NewVoiceTerminatorBurstFromArray(data)
	case burst.DataTypeCSBK:
		return burst.NewCSBKBurstFromArray(data)
	case burst.DataTypeDataHeader:
		return burst.NewDataHeaderBurstFromArray(data)
	case burst.DataTypeRateThreeQuarter:
		return burst.NewDataBurstFromArray(data)
	case burst.DataTypeRateFullData:
		//This only seems to happen for voice transmitions...
		return burst.NewVoiceBurstFromArray(data)
	default:
		return burst.NewUnknownBurstFromArray(data)
	}
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
	//TODO convert RTP and Burst to array...
	return a
	//return append(a, p.Payload...)
}
