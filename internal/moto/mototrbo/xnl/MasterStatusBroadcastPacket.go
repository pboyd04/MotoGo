package xnl

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// MasterStatusBroadcastPacket a response from the XNL master
type MasterStatusBroadcastPacket struct {
	Header        Header
	ProtocolRev   uint16
	ProtocolMinor uint16
	DeviceType    byte
	DeviceNumber  byte
	Unknown       byte
}

// NewMasterStatusBroadcastPacket create master status broadcast packet from params
func NewMasterStatusBroadcastPacket(src Address) MasterStatusBroadcastPacket {
	var pkt MasterStatusBroadcastPacket
	pkt.Header.OpCode = MasterStatusBroadcast
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = Address(0)
	pkt.Header.Src = src
	pkt.Header.TransactionID = 0
	pkt.ProtocolRev = 0
	pkt.ProtocolMinor = 1
	pkt.DeviceType = 1
	pkt.DeviceNumber = 1
	pkt.Unknown = 1
	return pkt
}

// NewMasterStatusBroadcastPacketByArray create an master status broadcast packet from data array
func NewMasterStatusBroadcastPacketByArray(data []byte) MasterStatusBroadcastPacket {
	var p MasterStatusBroadcastPacket
	p.Header.FromArray(data)
	p.ProtocolRev = util.ParseUint16(data, 14)
	p.ProtocolMinor = util.ParseUint16(data, 16)
	p.DeviceType = data[18]
	p.DeviceNumber = data[19]
	p.Unknown = data[20]
	return p
}

// GetHeader returns the header for the packet
func (p MasterStatusBroadcastPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p MasterStatusBroadcastPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p MasterStatusBroadcastPacket) ToArray() []byte {
	return p.Header.ToArray([]byte{})
}
