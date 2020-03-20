package xnl

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

// OpCode An XNL OpCode
type OpCode uint16

// Address An XNL Address
type Address uint16

// Protocol XCMP or XNL
type Protocol byte

const (
	// MasterStatusBroadcast returned for the init of the connection
	MasterStatusBroadcast OpCode = 0x02
	// DeviceMasterQuery starts the init of the connection
	DeviceMasterQuery OpCode = 0x03
	// DeviceAuthKeyRequest asks the master for an authkey
	DeviceAuthKeyRequest OpCode = 0x04
	// DeviceAuthKeyReply returns from the master with the authkey
	DeviceAuthKeyReply OpCode = 0x05
	// DeviceConnectionRequest asks for a connection with the master's encrypted authkey
	DeviceConnectionRequest OpCode = 0x06
	// DeviceConnectionReply response to a connection request
	DeviceConnectionReply OpCode = 0x07
	// DeviceSysMapBroadcast gets a list of devices seen in XNL
	DeviceSysMapBroadcast OpCode = 0x09
	// DataMessage sends XNL or XCMP data
	DataMessage OpCode = 0x0b
	// DataMessageAck acknowledges a DataMessage packet
	DataMessageAck OpCode = 0x0c

	// ProtocolXNL the packet is just XNL
	ProtocolXNL Protocol = 0
	// ProtocolXCMP the payload is in XCMP format
	ProtocolXCMP Protocol = 1
)

// Header for the Packet
type Header struct {
	OpCode        OpCode
	Protocol      Protocol
	Flags         byte
	Dest          Address
	Src           Address
	TransactionID uint16
}

// Packet represents an XNL packet
type Packet interface {
	ToArray() []byte
	GetHeader() Header
}

// CreatePacketFromArray creates a XNL packet from an array
func CreatePacketFromArray(array []byte) Packet {
	opCode := OpCode(util.ParseUint16(array, 2))
	switch opCode {
	case MasterStatusBroadcast:
		return NewMasterStatusBroadcastPacketByArray(array)
	case DeviceAuthKeyReply:
		return NewDevKeyAuthReplyPacketByArray(array)
	case DeviceConnectionReply:
		return NewDevConnectionReplyPacketByArray(array)
	case DeviceSysMapBroadcast:
		return NewDeviceSysmapBroadcastPacketByArray(array)
	case DataMessage:
		return NewDataPacketByArray(array)
	case DataMessageAck:
		return NewAckDataPacketByArray(array)
	default:
		return NewUnkownPacketByArray(array)
	}
}

// ToArray converts an XNL header to an array
func (h *Header) ToArray(payload []byte) []byte {
	pktlength := 12 //Minimum pkt length
	a := make([]byte, pktlength+2)
	pktlength += len(payload)
	a[0] = byte(pktlength >> 8)
	a[1] = byte(pktlength)
	a[2] = byte(h.OpCode >> 8)
	a[3] = byte(h.OpCode)
	a[4] = byte(h.Protocol)
	a[5] = h.Flags
	a[6] = byte(h.Dest >> 8)
	a[7] = byte(h.Dest)
	a[8] = byte(h.Src >> 8)
	a[9] = byte(h.Src)
	a[10] = byte(h.TransactionID >> 8)
	a[11] = byte(h.TransactionID)
	a[12] = byte(len(payload) >> 8)
	a[13] = byte(len(payload))
	return append(a, payload...)
}

// FromArray converts an array to an XNL header
func (h *Header) FromArray(data []byte) {
	h.OpCode = OpCode(util.ParseUint16(data, 2))
	h.Protocol = Protocol(data[4])
	h.Flags = data[5]
	h.Dest = Address(util.ParseUint16(data, 6))
	h.Src = Address(util.ParseUint16(data, 8))
	h.TransactionID = util.ParseUint16(data, 10)
}
