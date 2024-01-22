package sol

import (
	"encoding/binary"
	"math/big"
)

func convertUint64ToBuffer(num uint64) []byte {
	buffer := make([]byte, 16)
	binary.LittleEndian.PutUint64(buffer, num)
	return buffer
}

func bigIntToLittleEndianBytes(num *big.Int, arraySize int) []byte {
	numBytes := num.Bytes()

	byteLength := (num.BitLen() + 7) / 8

	littleEndianBytes := make([]byte, arraySize)

	for i := 0; i < byteLength; i++ {
		littleEndianBytes[i] = numBytes[byteLength-1-i]
	}

	return littleEndianBytes
}
