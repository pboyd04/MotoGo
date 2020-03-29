package burst

import "fmt"

//DataType describes the type of data in the packet
type DataType byte

const (
	//Taken from ETSI TS 102 361-1
	DataTypePIHeader         DataType = 0x00
	DataTypeVoiceLCHeader    DataType = 0x01
	DataTypeTerminatorWithLC DataType = 0x02
	DataTypeCSBK             DataType = 0x03
	DataTypeMBCHeader        DataType = 0x04
	DataTypeMBCContinuation  DataType = 0x05
	DataTypeDataHeader       DataType = 0x06
	DataTypeRateHalfData     DataType = 0x07
	DataTypeRateThreeQuarter DataType = 0x08
	DataTypeIdle             DataType = 0x09
	DataTypeRateFullData     DataType = 0x0a
	DataTypeUSBD             DataType = 0x0b
	//End values from ETSI TS 102 361-1

	//DataTypeUnknownSmall Don't know what this is... but it doesn't fit the regular format...
	DataTypeUnknownSmall DataType = 0x13
)

//Burst Represents the data/voice sent in this packet
type Burst interface {
	GetBurstType() DataType
	GetPayload() []byte
}

type UnknownBurst struct {
	DataType DataType
	Payload  []byte
}

func NewUnknownBurstFromArray(data []byte) UnknownBurst {
	var p UnknownBurst
	p.DataType = DataType(data[0])
	p.Payload = data[1:]
	fmt.Printf("Unknown packet created. Type = %02x\n", data[0])
	fmt.Printf(" Payload = %#v\n", p.Payload)
	return p
}

func (p UnknownBurst) GetBurstType() DataType {
	return p.DataType
}

func (p UnknownBurst) GetPayload() []byte {
	return p.Payload
}
