package redis

import (
	"fmt"
	"testing"
)

func TestCRC16(t *testing.T) {
	fmt.Println(crc16([]byte("123456789")))
	fmt.Println(crc16([]byte("123456789")))
	fmt.Println(crc16([]byte("123456789")))
	fmt.Println(crc16([]byte("123456789")))
	fmt.Println(crc16([]byte("123456789")))
	fmt.Println(crc16([]byte("123456789")))
	fmt.Println(crc16([]byte("123456789")))

	for _, v := range crcTable {
		fmt.Printf("0x%04x ", v)
	}
}
