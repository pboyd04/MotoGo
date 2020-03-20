package xnl

// UnknownPacket an unknown packet
type UnknownPacket struct {
	Header  Header
	Payload []byte
}

// NewUnkownPacketByArray create an data packet from data array
func NewUnkownPacketByArray(data []byte) UnknownPacket {
	var p UnknownPacket
	p.Header.FromArray(data)
	p.Payload = data[14:]
	return p
}

// GetHeader returns the header for the packet
func (p UnknownPacket) GetHeader() Header {
	return p.Header
}

// SetHeader sets the header for the packet
func (p UnknownPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p UnknownPacket) ToArray() []byte {
	return p.Header.ToArray(p.Payload)
}
