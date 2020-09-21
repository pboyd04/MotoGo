package burst

import (
	"github.com/pboyd04/MotoGo/internal/util"
)

type ContentType byte

const (
	UnifiedData          ContentType = 0x00
	TCPHeaderCompression ContentType = 0x02
	UDPHeaderCompression ContentType = 0x03
	IPPacket             ContentType = 0x04
	ARP                  ContentType = 0x05
	ProprietaryData      ContentType = 0x09
	ShortData            ContentType = 0x0A
)

type DataHeaderBurst struct {
	DataType          DataType
	RSSIOk            bool
	RSParity          bool
	CRCParity         bool
	LCParity          bool
	Unknown           uint16
	HasRSSI           bool
	BurstSource       bool
	HardSync          bool
	HasSlotType       bool
	SyncType          byte
	ColorCode         byte
	SlotType          byte
	RSSI              float32
	IsGroup           bool
	ResponseRequested bool
	Compressed        bool
	HeaderDataType    byte
	PadOctectCount    byte
	ContentType       ContentType
	To                uint32
	From              uint32
	FullMessage       bool
	BlocksFollow      byte
	CRC               uint16
}

func NewDataHeaderBurstFromArray(data []byte) DataHeaderBurst {
	var p DataHeaderBurst
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
	if p.HasSlotType {
		p.ColorCode = data[offset+1] >> 4
		p.SlotType = data[offset+1] & 0x0F
		offset += 2
	}
	if p.HasRSSI {
		p.RSSI = util.CalcRSSI(data, offset)
	}
	p.IsGroup = (data[8] & 0x80) != 0
	p.ResponseRequested = (data[8] & 0x40) != 0
	p.Compressed = (data[8] & 0x20) != 0
	p.HeaderDataType = data[8] & 0x0F
	p.PadOctectCount = data[8] & 0x10
	p.PadOctectCount |= data[9] & 0x0F
	p.ContentType = ContentType(data[9] >> 4)
	p.To = util.ParseUint24(data, 10)
	p.From = util.ParseUint24(data, 13)
	p.FullMessage = (data[16] & 0x80) != 0
	p.BlocksFollow = data[16] & 0x7f
	p.CRC = util.ParseUint16(data, 17)
	return p
}

func (p DataHeaderBurst) GetBurstType() DataType {
	return p.DataType
}

func (p DataHeaderBurst) GetPayload() []byte {
	return make([]byte, 0)
}
