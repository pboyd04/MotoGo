package xnl

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// DeviceSysmapBroadcastPacket a request to initialize a XNL connection
type DeviceSysmapBroadcastPacket struct {
	Header  Header
	Payload Payload
	Sysmaps []Sysmap
}

// Sysmap is an entry in the sysmap for a device
type Sysmap struct {
	DeviceType   byte
	DeviceNumber byte
	Address      Address
	AuthIndex    byte
}

// NewDeviceSysmapBroadcastPacket create device sysmap packet from params
func NewDeviceSysmapBroadcastPacket(src Address, maps []Sysmap) DeviceSysmapBroadcastPacket {
	var pkt DeviceSysmapBroadcastPacket
	pkt.Header.OpCode = DataMessage
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Src = src
	pkt.Header.TransactionID = 0
	pkt.Sysmaps = maps
	return pkt
}

// NewDeviceSysmapBroadcastPacketByArray create an data packet from data array
func NewDeviceSysmapBroadcastPacketByArray(data []byte) DeviceSysmapBroadcastPacket {
	var p DeviceSysmapBroadcastPacket
	p.Header.FromArray(data)
	length := util.ParseUint16(data, 14)
	p.Sysmaps = make([]Sysmap, length)
	for i := 0; i < int(length); i++ {
		p.Sysmaps[i] = SysmapFromArray(data[16+(i*5):])
	}
	return p
}

// GetHeader returns the header for the packet
func (p DeviceSysmapBroadcastPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p DeviceSysmapBroadcastPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p DeviceSysmapBroadcastPacket) ToArray() []byte {
	payload := make([]byte, 2)
	length := len(p.Sysmaps)
	payload[0] = byte(length >> 8)
	payload[1] = byte(length)
	for i := 0; i < int(length); i++ {
		payload = append(payload, p.Sysmaps[i].ToArray()...)
	}
	return p.Header.ToArray(payload)
}

// SysmapFromArray creates a Sysmap from an array
func SysmapFromArray(data []byte) Sysmap {
	var s Sysmap
	s.DeviceType = data[0]
	s.DeviceNumber = data[1]
	s.Address = Address(util.ParseUint16(data, 2))
	s.AuthIndex = data[4]
	return s
}

// ToArray converts a sysmap to an array
func (s Sysmap) ToArray() []byte {
	a := make([]byte, 5)
	a[0] = s.DeviceType
	a[1] = s.DeviceNumber
	a[2] = byte(s.Address >> 8)
	a[3] = byte(s.Address)
	a[4] = s.AuthIndex
	return a
}
