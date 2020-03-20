package xnl

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// DevConnectionRequestPacket a request to initialize a XNL connection
type DevConnectionRequestPacket struct {
	Header          Header
	connection      Address
	connectionType  byte
	connectionIndex byte
	key             []byte
}

// NewDevConnectionRequestPacket create deregistration packet from params
func NewDevConnectionRequestPacket(dest Address, src Address, connectionAddr Address, connectionType byte, connectionIndex byte, key []byte) DevConnectionRequestPacket {
	var pkt DevConnectionRequestPacket
	pkt.Header.OpCode = DeviceConnectionRequest
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = dest
	pkt.Header.Src = src
	pkt.Header.TransactionID = 0
	pkt.connection = connectionAddr
	pkt.connectionType = connectionType
	pkt.connectionIndex = connectionIndex
	pkt.key = Encrypt(key)
	return pkt
}

// NewDevConnectionRequestPacketByArray create an deregistration packet from data array
func NewDevConnectionRequestPacketByArray(data []byte) DevConnectionRequestPacket {
	var p DevConnectionRequestPacket
	p.Header.FromArray(data)
	p.connection = Address(util.ParseUint16(data, 14))
	p.connectionType = data[15]
	p.connectionIndex = data[16]
	p.key = data[17:]
	return p
}

// GetHeader returns the header for the packet
func (p DevConnectionRequestPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p DevConnectionRequestPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p DevConnectionRequestPacket) ToArray() []byte {
	payload := make([]byte, 4)
	payload[0] = byte(p.connection >> 8)
	payload[1] = byte(p.connection)
	payload[2] = p.connectionType
	payload[3] = p.connectionIndex
	return p.Header.ToArray(append(payload, p.key...))
}
