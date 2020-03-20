package xnl

// AckDataPacket a request to acknowledge a prior data packet
type AckDataPacket struct {
	Header Header
}

// NewAckDataPacketByParam create ack data packet from params
func NewAckDataPacketByParam(dest Address, src Address, proto Protocol, flags byte, transactionID uint16) AckDataPacket {
	var pkt AckDataPacket
	pkt.Header.OpCode = DataMessageAck
	pkt.Header.Protocol = proto
	pkt.Header.Flags = flags
	pkt.Header.Dest = dest
	pkt.Header.Src = src
	pkt.Header.TransactionID = transactionID
	return pkt
}

// NewAckDataPacketByArray create an data packet from data array
func NewAckDataPacketByArray(data []byte) AckDataPacket {
	var p AckDataPacket
	p.Header.FromArray(data)
	return p
}

// GetHeader returns the header for the packet
func (p AckDataPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p AckDataPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p AckDataPacket) ToArray() []byte {
	return p.Header.ToArray([]byte{})
}
