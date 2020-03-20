package mototrbo

//Command represents a command byte for the Mototrbo radios
type Command byte

//RadioID represents a radio ID
type RadioID uint32

const (
	// XnlXcmpPacket represents an XNL data packet
	XnlXcmpPacket Command = 0x70
	// GroupVoiceCall represents a group voice call
	GroupVoiceCall Command = 0x80
	// PrivateVoiceCall represents a private voice call
	PrivateVoiceCall Command = 0x81
	// GroupDataCall represents a group data call
	GroupDataCall Command = 0x83
	// PrivateDataCall represents a private data call
	PrivateDataCall Command = 0x84
	// RegistrationRequest represents a request to register with the master
	RegistrationRequest Command = 0x90
	// RegistrationReply represents a reply to registeration request
	RegistrationReply Command = 0x91
	// PeerListRequest represnets a request for the peer list
	PeerListRequest Command = 0x92
	// PeerListReply represnets a reply of the peer list
	PeerListReply Command = 0x93
	// PeerRegisterRequest represents a request to register with a peer
	PeerRegisterRequest Command = 0x94
	// PeerRegisterReply represents a reply to registeration request
	PeerRegisterReply Command = 0x95
	// MasterKeepAliveRequest ask the master to send a keep alive reply
	MasterKeepAliveRequest Command = 0x96
	// MasterKeepAliveReply a reply from the master to a keep alive request
	MasterKeepAliveReply Command = 0x97
	// PeerKeepAliveRequest ask the peer to send a keep alive reply
	PeerKeepAliveRequest Command = 0x98
	// PeerKeepAliveReply a reply from a peer to a keep alive request
	PeerKeepAliveReply Command = 0x99
	// DeregisterRequest a request to deregister from the master
	DeregisterRequest Command = 0x9A
)

// Packet is a data structure sent over UDP to a radio system
type Packet interface {
	GetCommand() Command
	GetID() RadioID
	ToArray() []byte
}

// CreatePacketFromArray creates a Mototrbo packet from an array
func CreatePacketFromArray(array []byte) (Packet, error) {
	cmd := Command(array[0])
	switch cmd {
	case XnlXcmpPacket:
		return NewXNLXCMPPacketByArray(array)
	case GroupVoiceCall:
		return NewUserPacketByArray(array)
	case PrivateVoiceCall:
		return NewUserPacketByArray(array)
	case GroupDataCall:
		return NewUserPacketByArray(array)
	case PrivateDataCall:
		return NewUserPacketByArray(array)
	case RegistrationRequest:
		return NewRegistrationPacketByArray(array)
	case RegistrationReply:
		return NewRegistrationReplyPacketByArray(array)
	case PeerListRequest:
		return NewPeerListRequestPacketByArray(array)
	case PeerListReply:
		return NewPeerListReplyPacketByArray(array)
	case PeerRegisterRequest:
		return NewPeerRegistrationPacketByArray(array)
	case PeerRegisterReply:
		return NewPeerRegistrationReplyPacketByArray(array)
	case MasterKeepAliveRequest:
		return NewMasterKeepAliveRequestPacketByArray(array)
	case MasterKeepAliveReply:
		return NewKeepAliveReplyPacketByArray(array)
	case DeregisterRequest:
		return NewDeregistrationPacketByArray(array)
	case PeerKeepAliveRequest:
		return NewPeerKeepAliveRequestPacketByArray(array)
	case PeerKeepAliveReply:
		return NewPeerKeepAliveReplyPacketByArray(array)
	default:
		return NewUnknownPacket(array)
	}
}
