package xnl

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// DevKeyAuthReplyPacket a request to initialize a XNL connection
type DevKeyAuthReplyPacket struct {
	Header  Header
	TempID  uint16
	AuthKey []byte
}

// NewDevKeyAuthReplyPacketByParam create data packet from params
func NewDevKeyAuthReplyPacketByParam(dest Address, src Address, tempID uint16, authKey []byte) DevKeyAuthReplyPacket {
	var pkt DevKeyAuthReplyPacket
	pkt.Header.OpCode = DataMessage
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = dest
	pkt.Header.Src = src
	pkt.Header.TransactionID = 0
	pkt.TempID = tempID
	pkt.AuthKey = authKey
	return pkt
}

// NewDevKeyAuthReplyPacketByArray create an data packet from data array
func NewDevKeyAuthReplyPacketByArray(data []byte) DevKeyAuthReplyPacket {
	var p DevKeyAuthReplyPacket
	p.Header.FromArray(data)
	p.TempID = util.ParseUint16(data, 14)
	p.AuthKey = data[16:]
	return p
}

// GetHeader returns the header for the packet
func (p DevKeyAuthReplyPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p DevKeyAuthReplyPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p DevKeyAuthReplyPacket) ToArray() []byte {
	payload := make([]byte, 2)
	payload[0] = byte(p.TempID >> 8)
	payload[1] = byte(p.TempID)
	return p.Header.ToArray(append(payload, p.AuthKey...))
}
