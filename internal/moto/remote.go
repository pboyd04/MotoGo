package moto

import (
	"fmt"
	"net"
	"time"
	"unicode/utf16"

	"github.com/pboyd04/MotoGo/internal/moto/mototrbo"
	"github.com/pboyd04/MotoGo/internal/moto/mototrbo/xnl"
	"github.com/pboyd04/MotoGo/internal/moto/mototrbo/xnl/xcmp"
	"github.com/pboyd04/MotoGo/internal/util"
)

// RemoteRadio represents a radio on the RadioSystem
type RemoteRadio struct {
	ID         mototrbo.RadioID
	IsMaster   bool
	EntityType xcmp.EntityType
	Addr       *net.UDPAddr

	system        *RadioSystem
	xnlClient     *xnl.Client
	xcmpClient    *xcmp.Client
	activeCalls   map[mototrbo.RadioID]*RadioCall
	ready         chan bool
	callChannel   chan *RadioCall
	callCountChan chan int
	alarms        map[string]bool
}

// NewRadio creates a new RemoteRadio instance
func NewRadio(address string, isMaster bool, sys *RadioSystem) (*RemoteRadio, error) {
	radio := new(RemoteRadio)
	addr, err := net.ResolveUDPAddr("udp", address)
	if err != nil {
		return nil, err
	}
	radio.Addr = addr
	radio.activeCalls = make(map[mototrbo.RadioID]*RadioCall)
	radio.IsMaster = isMaster
	radio.system = sys
	radio.ready = make(chan bool, 1)
	radio.callChannel = make(chan *RadioCall, 10)
	radio.alarms = make(map[string]bool)
	sys.initializing = radio
	radio.register()
	if radio.IsMaster {
		// Start the Keep Alive...
		pkt := mototrbo.NewMasterKeepAliveRequestPacketByParam(sys.MyID, true, true, sys.SystemType)
		radio.SendPacket(pkt)
	}
	sys.initializing = nil
	return radio, nil
}

// SendPacket sends a packet to the radio
func (r *RemoteRadio) SendPacket(pkt mototrbo.Packet) {
	r.system.client.SendPacket(pkt, r.Addr)
}

// GetXNLID gets the radios XNL ID
func (r *RemoteRadio) GetXNLID() uint16 {
	if r.xnlClient == nil {
		return 0
	}
	return uint16(r.xnlClient.GetRadioXNLID())
}

// GetXCMPVersion returns the XCMP version of the remote radio
func (r *RemoteRadio) GetXCMPVersion() string {
	return r.xcmpClient.Version
}

// GetActiveCallCount returns the number of active calls on the remote radio
func (r *RemoteRadio) GetActiveCallCount() int {
	return len(r.activeCalls)
}

// register registers with the master or peer as appropriate
func (r *RemoteRadio) register() {
	if r.IsMaster == false {
		pkt := mototrbo.NewPeerRegistrationPacketByParam(r.system.MyID, r.system.SystemType)
		r.SendPacket(pkt)
	} else {
		pkt := mototrbo.NewRegistrationPacketByParam(r.system.MyID, true, true, r.system.SystemType)
		r.SendPacket(pkt)
	}
	select {
	case <-r.ready:
		return
	case <-time.After(5 * time.Second):
		r.register()
	}
}

// Deregister tells the master radio we are no longer part of the system
func (r *RemoteRadio) Deregister() {
	if r.IsMaster {
		pkt := mototrbo.NewDeregistrationPacketByParam(r.ID)
		r.SendPacket(pkt)
	}
}

// InitXNL initializes XNL communication with the radio
func (r *RemoteRadio) InitXNL() {
	r.xnlClient = xnl.NewClient(r.system.client, r.system.MyID, r.Addr)
	r.xcmpClient = xcmp.NewClient(r.xnlClient)
}

// ListenForCalls starts a loop waiting for call packets and keeping the connection alive
func (r *RemoteRadio) ListenForCalls(calls chan *RadioCall, callCount chan int) {
	for {
		r.callCountChan = callCount
		call := <-r.callChannel
		calls <- call
	}
}

func (r *RemoteRadio) processXnlAsyncPacket(pkt mototrbo.XnlPacket) {
	/*
		fmt.Printf("Got an unrequested XNL packet...\n")
		xnlPkt := pkt.XNL
		if xnlPkt.Protocol == xnl.ProtocolXCMP {
			fmt.Printf("Is an XCMP Packet\n")
			xcmpPkt := xcmp.CreatePacketFromArray(xnlPkt.Payload)
			fmt.Printf("XCMP Packet = %#v\n", xcmpPkt)
			switch xcmpPkt.OpCode {
			case xcmp.OpCode_Alarm:
				r.processAlarm(xcmpPkt.Payload)
			default:
				fmt.Printf("Unknown op code in packet %#v\n", xcmpPkt)
			}
		}*/
}

func (r *RemoteRadio) processAlarm(payload []byte) {
	fmt.Printf("Alarm Severity = %x\n", payload[0])
	if payload[1] == 0 {
		fmt.Printf("Alarm Deactivated\n")
	} else {
		fmt.Printf("Alarm Activated\n")
	}
	switch payload[2] {
	case 5:
		fmt.Printf("Fan Alarm\n")
	default:
		fmt.Printf("Unknown Alarm %x\n", payload[2])
	}
	fmt.Printf("Alarm Length = %x\n", uint32(payload[3])<<24|uint32(payload[4])<<16|uint32(payload[5])<<8|uint32(payload[6]))
}

// GetSerialNumber retrieves the remote radio serial number
func (r *RemoteRadio) GetSerialNumber() string {
	statusPkt := r.xcmpClient.SendAndWaitForRadioStatusReply(xcmp.SerialNumber)
	return string(statusPkt.Data)
}

// GetModelNumber retrieves the remote radio model number
func (r *RemoteRadio) GetModelNumber() string {
	statusPkt := r.xcmpClient.SendAndWaitForRadioStatusReply(xcmp.ModelNumber)
	x := len(statusPkt.Data)
	return string(statusPkt.Data[:x-2])
}

// GetModelName retrieves the remote radio model name
func (r *RemoteRadio) GetModelName() string {
	modelNumber := r.GetModelNumber()
	switch modelNumber {
	case "M27TRR9JA7BN":
		return "XPR 8400 (UHF 450-512MHz)"
	case "M27QPR9JA7BN":
		return "XPR 8400 (UHF 403-470MHz)"
	case "M27JQR9JA7BN":
		return "XPR 8400 (VHF 136-174MHz)"
	default:
		return fmt.Sprintf("Unknown Model Number (%s)", modelNumber)
	}
}

// GetRSSI retrieves the remote radio Received Signal Strength Indicator for both time slots
func (r *RemoteRadio) GetRSSI() (float32, float32) {
	statusPkt := r.xcmpClient.SendAndWaitForRadioStatusReply(xcmp.RSSI)
	return util.CalcRSSI(statusPkt.Data, 0), util.CalcRSSI(statusPkt.Data, 2)
}

// GetFirmwareVersion retrieves the remote radio firmware version
func (r *RemoteRadio) GetFirmwareVersion() string {
	r.xcmpClient.SendPacket(xcmp.NewVersionInfoRequestByParam())
	pkt := <-r.xcmpClient.PacketsIn
	verInfo := pkt.(xcmp.VersionInfoReply)
	return verInfo.Version
}

// GetRadioAlias retrieves the remote radio alias
func (r *RemoteRadio) GetRadioAlias() string {
	statusPkt := r.xcmpClient.SendAndWaitForRadioStatusReply(xcmp.RadioAlias)
	length := len(statusPkt.Data) / 2
	newU := make([]uint16, length)
	for i := 0; i < length; i++ {
		newU[i] = uint16(statusPkt.Data[i*2])<<8 | uint16(statusPkt.Data[(i*2)+1])
	}
	s := utf16.Decode(newU)
	return string(s)
}

// GetAlarmStatus retrieves the remote radio alarm status info
func (r *RemoteRadio) GetAlarmStatus() map[string]bool {
	r.xcmpClient.SendPacket(xcmp.NewAlarmStatusRequestByParam())
	pkt := <-r.xcmpClient.PacketsIn
	almStatus := pkt.(xcmp.AlarmStatusReply)
	for _, as := range almStatus.Alarms {
		switch as.Alarm {
		case xcmp.AlarmTransmit:
			r.alarms["Transmit"] = (as.State == 0x01)
		case xcmp.AlarmReceive:
			r.alarms["Receive"] = (as.State == 0x01)
		case xcmp.AlarmTemp:
			r.alarms["Temperature"] = (as.State == 0x01)
		case xcmp.AlarmAC:
			r.alarms["AC"] = (as.State == 0x01)
		case xcmp.AlarmFan:
			r.alarms["Fan"] = (as.State == 0x01)
		case xcmp.AlarmVSWR:
			r.alarms["VSWR"] = (as.State == 0x01)
		case xcmp.AlarmTrasmitPower:
			r.alarms["Transmit Power"] = (as.State == 0x01)
		}
	}
	return r.alarms
}

func (r *RemoteRadio) gotUserPacket(pkt mototrbo.Packet) bool {
	upkt := pkt.(mototrbo.UserPacket)
	to := upkt.Destination
	call := new(RadioCall)
	if _, ok := r.activeCalls[to]; ok {
		call = r.activeCalls[to]
	} else {
		call.StartTime = time.Now()
	}
	if pkt.GetCommand() == mototrbo.GroupVoiceCall || pkt.GetCommand() == mototrbo.GroupDataCall {
		call.Group = true
	} else {
		call.Group = false
	}
	if pkt.GetCommand() == mototrbo.GroupDataCall || pkt.GetCommand() == mototrbo.PrivateDataCall {
		call.Audio = false
	} else {
		call.Audio = true
	}
	call.From = upkt.Source
	call.To = to
	call.Encrypted = upkt.Encrypted
	call.End = upkt.End
	call.Timeslot = upkt.TimeSlot
	call.IsPhone = upkt.PhoneCall
	//if !call.Audio {
	//	fmt.Printf("%#v\n", upkt.Payload)
	//}
	if call.End {
		call.EndTime = time.Now()
		delete(r.activeCalls, to)
		r.callChannel <- call
	} else {
		r.activeCalls[to] = call
	}
	if r.callCountChan != nil {
		r.callCountChan <- len(r.activeCalls)
	}
	return true
}

func (r *RemoteRadio) gotRegisterReply(pkt mototrbo.Packet) bool {
	r.ID = pkt.GetID()
	r.ready <- true
	return true
}
