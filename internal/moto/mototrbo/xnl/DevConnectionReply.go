package xnl

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// DevConnectionReplyPacket a request to initialize a XNL connection
type DevConnectionReplyPacket struct {
	Header     Header
	AssignedID Address
	AuthInfo   []byte
}

// NewDevConnectionReplyPacketByParam create deregistration packet from params
func NewDevConnectionReplyPacketByParam(dest Address, src Address, assignedID Address, authInfo []byte) DevConnectionReplyPacket {
	var pkt DevConnectionReplyPacket
	pkt.Header.OpCode = DeviceConnectionRequest
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = dest
	pkt.Header.Src = src
	pkt.Header.TransactionID = 0
	pkt.AssignedID = assignedID
	pkt.AuthInfo = authInfo
	return pkt
}

// NewDevConnectionReplyPacketByArray create an deregistration packet from data array
func NewDevConnectionReplyPacketByArray(data []byte) DevConnectionReplyPacket {
	var p DevConnectionReplyPacket
	p.Header.FromArray(data)
	p.AssignedID = Address(util.ParseUint16(data, 16))
	p.AuthInfo = data[16:]
	return p
}

// GetHeader returns the header for the packet
func (p DevConnectionReplyPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p DevConnectionReplyPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p DevConnectionReplyPacket) ToArray() []byte {
	payload := make([]byte, 4)
	payload[0] = 0x01
	payload[1] = 0x04
	payload[2] = byte(p.AssignedID >> 8)
	payload[3] = byte(p.AssignedID)
	return p.Header.ToArray(append(payload, p.AuthInfo...))
}
