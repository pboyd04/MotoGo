package util

// ParseUint16 gets a uint16 from an array
func ParseUint16(data []byte, offset int) uint16 {
	return uint16(data[offset])<<8 | uint16(data[offset+1])
}

// ParseUint24 gets a uint32 consisting of 24-bits from an array
func ParseUint24(array []byte, offset int) uint32 {
	return uint32(array[offset])<<16 | uint32(array[offset+1])<<8 | uint32(array[offset+2])
}

// ParseUint32 gets a uint32 from an array
func ParseUint32(data []byte, offset int) uint32 {
	return uint32(data[offset])<<24 | uint32(data[offset+1])<<16 | uint32(data[offset+2])<<8 | uint32(data[offset+3])
}

// CalcRSSI converts the data sent by the repeater to floating point at the correct offset
func CalcRSSI(data []byte, offset int) float32 {
	return -1.0*float32(data[offset]) + float32((float64(data[offset+1])*1000.0+128.0)/256000.0)
}
