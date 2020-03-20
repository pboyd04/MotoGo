package xnl

// DataPacket a request to initialize a XNL connection
type DataPacket struct {
	Header  Header
	Payload Payload
}

// Payload an Xnl Packet Payload
type Payload interface {
	ToArray() []byte
}

// NewDataPacket create data packet from params
func NewDataPacket(dest Address, src Address, proto Protocol, payload Payload) DataPacket {
	var pkt DataPacket
	pkt.Header.OpCode = DataMessage
	pkt.Header.Protocol = proto
	pkt.Header.Flags = 0x00
	pkt.Header.Dest = dest
	pkt.Header.Src = src
	pkt.Payload = payload
	return pkt
}

// NewDataPacketByArray create an data packet from data array
func NewDataPacketByArray(data []byte) DataPacket {
	var p DataPacket
	p.Header.FromArray(data)
	var g GenericPayload
	g.Payload = data[14:]
	p.Payload = g
	return p
}

// GetHeader returns the header for the packet
func (p DataPacket) GetHeader() Header {
	return p.Header
}

// SetHeader returns the header for the packet
func (p DataPacket) SetHeader(h Header) {
	p.Header = h
}

// ToArray converts a packet to a byte array
func (p DataPacket) ToArray() []byte {
	return p.Header.ToArray(p.Payload.ToArray())
}

// GenericPayload a generic payload struct
type GenericPayload struct {
	Payload []byte
}

// ToArray converts a generic payload to an array
func (g GenericPayload) ToArray() []byte {
	return g.Payload
}
