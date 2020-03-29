package moto

import (
	"fmt"
	"strings"
	"time"

	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/pboyd04/MotoGo/internal/moto/mototrbo"
	"github.com/pboyd04/MotoGo/internal/moto/mototrbo/burst"
)

// RadioCall describes a call on the Radio
type RadioCall struct {
	Group     bool
	Audio     bool
	From      mototrbo.RadioID
	To        mototrbo.RadioID
	Encrypted bool
	End       bool
	Timeslot  byte
	IsPhone   bool
	StartTime time.Time
	EndTime   time.Time
	Payload   []burst.Burst
}

type PacketData interface {
}

func (r RadioCall) ConsolidateData() PacketData {
	b := r.Payload[0]
	switch b.GetBurstType() {
	case burst.DataTypeDataHeader:
		dh := b.(burst.DataHeaderBurst)
		a := make([]byte, 0)
		for _, burst := range r.Payload[1:] {
			a = append(a, burst.GetPayload()...)
		}
		if dh.ContentType == burst.IPPacket {
			packet := gopacket.NewPacket(a, layers.LayerTypeIPv4, gopacket.Default)
			return packet
		}
		fmt.Printf("Unknown data header type %x\n", dh.ContentType)
		return make([]byte, 0)
	default:
		fmt.Printf("Unknown packet type %x\n", b.GetBurstType())
		return make([]byte, 0)
	}
}

func dumpHex(data []byte) string {
	builder := new(strings.Builder)
	length := len(data)
	for row := 0; row < length/8; row++ {
		builder.WriteString(fmt.Sprintf("%06x", row*8))
		rowData := data[row*8:]
		rowLength := len(rowData)
		if rowLength > 8 {
			rowLength = 8
		}
		for i := 0; i < rowLength; i++ {
			builder.WriteString(fmt.Sprintf(" %02x", rowData[i]))
		}
		builder.WriteString("\n")
	}
	return builder.String()
}
