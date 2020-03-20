package xnl

// Encrypt converts the auth data from the repeater to the format to send back in the handshake
func Encrypt(data []byte) []byte {
	dword1 := arrayToInt(data, 0)
	dword2 := arrayToInt(data, 4)
	var num1 uint32 = 0
	var num2 uint32 = EncrypterInput1
	var num3 uint32 = EncrypterInput2
	var num4 uint32 = EncrypterInput3
	var num5 uint32 = EncrypterInput4
	var num6 uint32 = EncrypterInput5
	for index := 0; index < 32; index++ {
		num1 += num2
		dword1 += (((dword2 << 4) + num3) ^ (dword2 + num1) ^ ((dword2 >> 5) + num4))
		dword2 += (((dword1 << 4) + num5) ^ (dword1 + num1) ^ ((dword1 >> 5) + num6))
	}
	a := make([]byte, 8)
	a[0] = byte(dword1 >> 24)
	a[1] = byte(dword1 >> 16)
	a[2] = byte(dword1 >> 8)
	a[3] = byte(dword1)
	a[4] = byte(dword2 >> 24)
	a[5] = byte(dword2 >> 16)
	a[6] = byte(dword2 >> 8)
	a[7] = byte(dword2)
	return a
}

func arrayToInt(data []byte, start int) uint32 {
	var ret uint32 = 0
	for index := 0; index < 4; index++ {
		ret = ret<<8 | uint32(data[index+start])
	}
	return ret
}
