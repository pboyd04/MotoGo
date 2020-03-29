package burst

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

type CSBKOpCode byte

const (
	//These CSBK Opcodes are from TS-102 361-2
	CSBKVoiceServiceRequest  CSBKOpCode = 0x04
	CSBKVoiceServiceAnswer   CSBKOpCode = 0x05
	CSBKNak                  CSBKOpCode = 0x26
	CSBKBSOutboundActivation CSBKOpCode = 0x38
	CSBKPreamble             CSBKOpCode = 0x3D

	//Other CSBK Opcodes (presumably Mototrbo specific)
	CSBKMototrboRadioCheck CSBKOpCode = 0x24
)

type CSBKBurst struct {
	DataType    DataType
	RSSIOk      bool
	RSParity    bool
	CRCParity   bool
	LCParity    bool
	Unknown     uint16
	HasRSSI     bool
	BurstSource bool
	HardSync    bool
	HasSlotType bool
	SyncType    byte
	ColorCode   byte
	SlotType    byte
	RSSI        float32
	LastBlock   bool
	ProtectFlag bool //This isn't used yet...
	CSBKOpCode  CSBKOpCode
	FeatureID   byte
	CRC         uint16
	Payload     []byte
}

func NewCSBKBurstFromArray(data []byte) CSBKBurst {
	var p CSBKBurst
	p.DataType = DataType(data[0])
	p.RSSIOk = (data[1] & 0x40) != 0
	p.RSParity = (data[1] & 0x04) != 0
	p.CRCParity = (data[1] & 0x02) != 0
	p.LCParity = (data[1] & 0x01) != 0
	p.Unknown = util.ParseUint16(data, 2)
	p.HasRSSI = (data[4] & 0x80) != 0
	p.BurstSource = (data[4] & 0x01) != 0
	p.HardSync = (data[5] & 0x40) != 0
	p.HasSlotType = (data[5] & 0x08) != 0
	p.SyncType = data[5] & 0x03
	offset := int(util.ParseUint16(data, 6)/8) + 8
	orig := offset
	if p.HasSlotType {
		p.ColorCode = data[offset+1] >> 4
		p.SlotType = data[offset+1] & 0x0F
		offset += 2
	}
	if p.HasRSSI {
		p.RSSI = util.CalcRSSI(data, offset)
	}
	//CSBK Data struct starts here...
	p.LastBlock = (data[8] & 0x80) != 0
	p.ProtectFlag = (data[8] & 0x40) != 0
	p.CSBKOpCode = CSBKOpCode(data[8] & 0x3F)
	p.FeatureID = data[9]
	orig -= 2
	p.CRC = util.ParseUint16(data, orig)
	p.Payload = data[10:orig]
	return p
}

func (p CSBKBurst) GetBurstType() DataType {
	return p.DataType
}

func (p CSBKBurst) GetPayload() []byte {
	return p.Payload
}
