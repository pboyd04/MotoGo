package moto

import (
	"fmt"
	"net"
	"time"

	"github.com/pboyd04/MotoGo/internal/moto/mototrbo"
)

// RadioSystem represents a whole system of Mototrbo Radios
type RadioSystem struct {
	MyID       mototrbo.RadioID
	SystemType mototrbo.LinkType

	client       *mototrbo.Client
	master       *RemoteRadio
	peerRadios   map[mototrbo.RadioID]*RemoteRadio
	peerList     chan []mototrbo.Peer
	initializing *RemoteRadio
}

// NewRadioSystem create Radio System with the specified ID and the specified system type
func NewRadioSystem(id mototrbo.RadioID, systemType mototrbo.LinkType) (*RadioSystem, error) {
	sys := new(RadioSystem)
	sys.MyID = id
	sys.SystemType = systemType
	client, err := mototrbo.NewClient(":50001")
	if err != nil {
		return nil, err
	}
	sys.client = client
	sys.peerList = make(chan []mototrbo.Peer, 1)
	client.RegisterHandler(mototrbo.GroupVoiceCall, sys.gotUserPacket)
	client.RegisterHandler(mototrbo.PrivateVoiceCall, sys.gotUserPacket)
	client.RegisterHandler(mototrbo.GroupDataCall, sys.gotUserPacket)
	client.RegisterHandler(mototrbo.PrivateDataCall, sys.gotUserPacket)
	client.RegisterHandler(mototrbo.RegistrationReply, sys.gotRegisterReply)
	client.RegisterHandler(mototrbo.PeerListReply, sys.gotPeerListReply)
	client.RegisterHandler(mototrbo.PeerRegisterRequest, sys.gotRegisterRequest)
	client.RegisterHandler(mototrbo.PeerRegisterReply, sys.gotRegisterReply)
	client.RegisterHandler(mototrbo.MasterKeepAliveReply, sys.gotMasterKeepAliveReply)
	client.RegisterHandler(mototrbo.PeerKeepAliveRequest, sys.gotPeerKeepAliveRequest)
	client.RegisterHandler(mototrbo.PeerKeepAliveReply, sys.gotPeerKeepAliveReply)
	return sys, nil
}

// ConnectToMaster connects to a master repeater
func (sys *RadioSystem) ConnectToMaster(address string) error {
	if sys.master != nil {
		sys.Close()
	}
	master, err := NewRadio(address, true, sys)
	if err != nil {
		return err
	}
	sys.master = master
	sys.peerRadios = make(map[mototrbo.RadioID]*RemoteRadio)
	return nil
}

// Close closes the connection to all the radios in the RadioSystem and deregisters from the master
func (sys *RadioSystem) Close() {
	if sys.master != nil {
		sys.master.Deregister()
	}
	sys.client.Close()
}

// PeerList sends a peer list request to the master
func (sys *RadioSystem) PeerList() []*RemoteRadio {
	pkt := mototrbo.NewPeerListRequestPacketByParam(sys.MyID)
	sys.client.SendPacket(pkt, sys.master.Addr)
	peers := <-sys.peerList
	ret := make([]*RemoteRadio, 0)
	for _, peer := range peers {
		if peer.ID == sys.MyID || peer.ID == 0 {
			//This is either the master (already connected) or me
			continue
		}
		address := fmt.Sprintf("%s:%d", peer.Address, peer.Port)
		radio, err := NewRadio(address, false, sys)
		if err != nil {
			fmt.Printf("Error talking to peer %v\n", err)
		}
		sys.peerRadios[peer.ID] = radio
		ret = append(ret, radio)
	}
	return ret
}

// GetMasterID returns the master radio ID
func (sys *RadioSystem) GetMasterID() mototrbo.RadioID {
	if sys.master == nil {
		return sys.MyID
	}
	return sys.master.ID
}

// GetMasterXNLID returns the master radio XNL Protocol ID
func (sys *RadioSystem) GetMasterXNLID() uint16 {
	return sys.master.GetXNLID()
}

// GetMaster returns the remote radio representing the master radio
func (sys *RadioSystem) GetMaster() *RemoteRadio {
	return sys.master
}

// GetRadioByID returns the radio by its ID
func (sys *RadioSystem) GetRadioByID(id mototrbo.RadioID) *RemoteRadio {
	radio, ok := sys.peerRadios[id]
	if ok {
		return radio
	}
	if sys.master != nil && sys.master.ID == id {
		return sys.master
	}
	return nil
}

func (sys *RadioSystem) getRadioForPacket(pkt mototrbo.Packet) *RemoteRadio {
	return sys.GetRadioByID(pkt.GetID())
}

func (sys *RadioSystem) gotUserPacket(pkt mototrbo.Packet, _ net.Addr) bool {
	radio := sys.getRadioForPacket(pkt)
	if radio == nil {
		return false
	}
	return radio.gotUserPacket(pkt)
}

func (sys *RadioSystem) gotRegisterRequest(pkt mototrbo.Packet, addr net.Addr) bool {
	radio := sys.getRadioForPacket(pkt)
	pkt = mototrbo.NewPeerRegistrationReplyPacketByParam(sys.MyID, sys.SystemType)
	if radio == nil {
		sys.client.SendPacket(pkt, addr.(*net.UDPAddr))
	} else {
		sys.client.SendPacket(pkt, radio.Addr)
	}
	return true
}

func (sys *RadioSystem) gotRegisterReply(pkt mototrbo.Packet, _ net.Addr) bool {
	sys.initializing.ID = pkt.GetID()
	sys.initializing.ready <- true
	return true
}

func (sys *RadioSystem) gotPeerListReply(pkt mototrbo.Packet, _ net.Addr) bool {
	peerListPkt := pkt.(mototrbo.PeerListReplyPacket)
	sys.peerList <- peerListPkt.Peers
	return true
}

func (sys *RadioSystem) gotMasterKeepAliveReply(pkt mototrbo.Packet, _ net.Addr) bool {
	// Send a new request after 5 seconds...
	timer := time.NewTimer(5 * time.Second)
	radio := sys.getRadioForPacket(pkt)
	go func() {
		<-timer.C
		pkt := mototrbo.NewMasterKeepAliveRequestPacketByParam(sys.MyID, true, true, sys.SystemType)
		sys.client.SendPacket(pkt, radio.Addr)
	}()
	return true
}

func (sys *RadioSystem) gotPeerKeepAliveRequest(pkt mototrbo.Packet, addr net.Addr) bool {
	radio := sys.getRadioForPacket(pkt)
	pkt = mototrbo.NewPeerKeepAliveReplyPacketByParam(sys.MyID, true, true, sys.SystemType)
	if radio == nil {
		sys.client.SendPacket(pkt, addr.(*net.UDPAddr))
	} else {
		sys.client.SendPacket(pkt, radio.Addr)
	}
	return true
}

func (sys *RadioSystem) gotPeerKeepAliveReply(pkt mototrbo.Packet, _ net.Addr) bool {
	// Send a new request after 5 seconds...
	timer := time.NewTimer(5 * time.Second)
	radio := sys.getRadioForPacket(pkt)
	go func() {
		<-timer.C
		pkt := mototrbo.NewPeerKeepAliveRequestPacketByParam(sys.MyID, true, true, sys.SystemType)
		sys.client.SendPacket(pkt, radio.Addr)
	}()
	return true
}
