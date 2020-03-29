package burst

import "fmt"

type VoiceBurst struct {
	Slot        byte
	DataType    DataType
	Unknown     byte
	Flags       byte
	Frames      [][]byte
	FrameErrors []bool
	LCHardBits  []byte
}

func NewVoiceBurstFromArray(data []byte) VoiceBurst {
	var p VoiceBurst
	p.Slot = 1
	if (data[0] & 0x80) != 0 {
		p.Slot = 2
	}
	p.DataType = DataType(data[0] & 0x3F)
	p.Unknown = data[1]
	p.Flags = data[2]
	p.Frames = make([][]byte, 3)
	p.FrameErrors = make([]bool, 3)
	x := 0
	y := 3
	for i := 0; i < 3; i++ {
		frame := make([]byte, 7)
		for j := 0; j < 7; j++ {
			tmp := data[y]
			if x > 0 {
				if j > 0 {
					frame[j-1] |= byte(uint16(tmp) >> (8 - x))
				}
				frame[j] = byte(uint16(tmp) << x)
			} else {
				frame[j] = tmp
			}
			y++
		}
		frame[6] &= 0x80
		x += 2
		y--
		p.Frames[i] = frame
	}
	p.FrameErrors[0] = (data[2] & 0x01) != 0
	p.FrameErrors[1] = (data[9] & 0x40) != 0
	p.FrameErrors[2] = (data[15] & 0x10) != 0
	if (p.Flags & 0x02) != 0 {
		p.LCHardBits = make([]byte, 4)
		p.LCHardBits[0] = data[22]
		p.LCHardBits[1] = data[23]
		p.LCHardBits[2] = data[24]
		p.LCHardBits[3] = data[25]
	} else if (p.Flags & 0x10) != 0 {
		fmt.Printf("Embed LC bits...\n")
	} else if (p.Flags & 0x04) != 0 {
		fmt.Printf("EMB...\n")
	}
	fmt.Printf("Voice packet created. Packet = %#v\n", p)
	return p
}

func (p VoiceBurst) GetBurstType() DataType {
	return p.DataType
}

func (p VoiceBurst) GetPayload() []byte {
	return make([]byte, 0)
}
