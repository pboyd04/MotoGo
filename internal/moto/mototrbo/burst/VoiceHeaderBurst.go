package burst

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

type FullLinkControlOpcode byte

const (
	FLCOGroupVoiceChannelUser      FullLinkControlOpcode = 0x00
	FLCOUnitToUnitVoiceChannelUser FullLinkControlOpcode = 0x03
)

type VoiceHeaderBurst struct {
	DataType              DataType
	RSSIOk                bool
	RSParity              bool
	CRCParity             bool
	LCParity              bool
	Unknown               uint16
	HasRSSI               bool
	BurstSource           bool
	HardSync              bool
	HasSlotType           bool
	SyncType              byte
	ColorCode             byte
	SlotType              byte
	RSSI                  float32
	Protected             bool
	FullLinkControlOpcode FullLinkControlOpcode
	FeatureID             byte
	SerivceOptions        byte
	DestAddress           uint32
	SourceAddress         uint32
	CRC                   uint16
}

func NewVoiceHeaderBurstFromArray(data []byte) VoiceHeaderBurst {
	var p VoiceHeaderBurst
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
	orig -= 2
	p.CRC = util.ParseUint16(data, orig)
	p.Protected = (data[8] & 0x80) != 0
	p.FullLinkControlOpcode = FullLinkControlOpcode(data[8] & 0x3F)
	p.FeatureID = data[9]
	p.SerivceOptions = data[10]
	p.DestAddress = util.ParseUint24(data, 11)
	p.SourceAddress = util.ParseUint24(data, 14)
	return p
}

func (p VoiceHeaderBurst) GetBurstType() DataType {
	return p.DataType
}

func (p VoiceHeaderBurst) GetPayload() []byte {
	return make([]byte, 0)
}
