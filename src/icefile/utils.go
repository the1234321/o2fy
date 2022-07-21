package icefile

import (
	"bytes"
	"encoding/binary"
	"os"
	"reflect"
)

const (
	SeekOriginBegin   = 0
	SeekOriginCurrent = 1
	SeekOriginEnd     = 2

	byteMaxValue byte = 255

	keyconstant_1 uint = 1129510338
	keyconstant_2 uint = 3444586398
	keyconstant_3 uint = 613566757
	keyconstant_4 uint = 1321528399
)

func ReadBytes(file *os.File, num int) []byte {
	result := make([]byte, num)
	file.Read(result)
	return result
}

func ReadByte(file *os.File) (result byte) {
	binary.Read(file, binary.LittleEndian, &result)
	return
}

func ReadUInt32(file *os.File) (result uint32) {
	binary.Read(file, binary.LittleEndian, &result)
	return
}

func ReadInt32(file *os.File) (result int32) {
	binary.Read(file, binary.LittleEndian, &result)
	return
}

func BitConverter_ToInt32(buf []byte, index int) (result int32) {
	buffer := bytes.NewBuffer(buf[index : index+4])
	binary.Read(buffer, binary.LittleEndian, &result)
	return

}

func BitConverter_ToUInt32(buf []byte, index int) (result uint32) {
	buffer := bytes.NewBuffer(buf[index : index+4])
	binary.Read(buffer, binary.LittleEndian, &result)
	return result

}

func BitConverter_ToBoolean(buf []byte, index int) (result bool) {
	buffer := bytes.NewBuffer(buf[index : index+4])
	binary.Read(buffer, binary.LittleEndian, &result)
	return result

}

func BitConverter_GetBytes(n uint) []byte {
	data := int32(n)
	bytebuf := bytes.NewBuffer([]byte{})
	binary.Write(bytebuf, binary.LittleEndian, data)
	return bytebuf.Bytes()
}

func reverse_byte(s []byte) []byte {
	n := reflect.ValueOf(s).Len()
	swap := reflect.Swapper(s)
	for i, j := 0, n-1; i < j; i, j = i+1, j-1 {
		swap(i, j)
	}
	return s
}

func Array_Copy1(src []byte, dest *[]byte, length int) {
	if len(*dest) < length {
		*dest = make([]byte, length)
	}
	dest_tmp := src[0:length]
	copy(*dest, dest_tmp)
}

func Array_Copy1_uint32(src *[]uint, dest *[]uint32, length int) {
	// dest_tmp := (*src)[0:length]
	// copy(*dest, dest_tmp)
	destIndex := 0
	srcIndex := 0
	for length > 0 {
		(*dest)[destIndex] = uint32((*src)[srcIndex])
		destIndex++
		srcIndex++
		length--
	}
}

func Array_Copy2(src []byte, srcIndex int, dest *[]byte, destIndex int, length int) {
	for length > 0 {
		(*dest)[destIndex] = src[srcIndex]
		destIndex++
		srcIndex++
		length--
	}
}
