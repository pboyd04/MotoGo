package xnl

// InitPacket a request to initialize a XNL connection
type InitPacket struct {
	Header Header
}

// NewInitPacket create deregistration packet from params
func NewInitPacket() InitPacket {
	var pkt InitPacket
	pkt.Header.OpCode = DeviceMasterQuery
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = Address(0)
	pkt.Header.Src = Address(0)
	pkt.Header.TransactionID = 0
	return pkt
}

// NewInitPacketByArray create an deregistration packet from data array
func NewInitPacketByArray(data []byte) InitPacket {
	var p InitPacket
	p.Header.FromArray(data)
	return p
}

// GetHeader returns the header for the packet
func (p InitPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p InitPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p InitPacket) ToArray() []byte {
	return p.Header.ToArray([]byte{})
}
