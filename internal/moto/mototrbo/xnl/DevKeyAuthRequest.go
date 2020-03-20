package xnl

// DevAuthKeyRequestPacket a request to initialize a XNL connection
type DevAuthKeyRequestPacket struct {
	Header Header
}

// NewDevAuthKeyRequestPacket create deregistration packet from params
func NewDevAuthKeyRequestPacket(dest Address) DevAuthKeyRequestPacket {
	var pkt DevAuthKeyRequestPacket
	pkt.Header.OpCode = DeviceAuthKeyRequest
	pkt.Header.Protocol = ProtocolXNL
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = dest
	pkt.Header.Src = Address(0)
	pkt.Header.TransactionID = 0
	return pkt
}

// NewDevAuthKeyRequestPacketByArray create an deregistration packet from data array
func NewDevAuthKeyRequestPacketByArray(data []byte) DevAuthKeyRequestPacket {
	var p DevAuthKeyRequestPacket
	p.Header.FromArray(data)
	return p
}

// GetHeader returns the header for the packet
func (p DevAuthKeyRequestPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p DevAuthKeyRequestPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p DevAuthKeyRequestPacket) ToArray() []byte {
	return p.Header.ToArray([]byte{})
}
