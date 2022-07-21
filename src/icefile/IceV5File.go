package icefile

import (
	"fmt"
	"os"
)

type DecryptBaseData struct {
	KeyStartPos   byte
	CrcStartPos   byte
	CrcEndPos     byte
	KeyConstTable []byte
	HeaderRol     byte
	Group2Rol     byte
}

type IceV5File struct {
	IceFile

	V5Decrypt         DecryptBaseData
	V6Decrypt         DecryptBaseData
	V7Decrypt         DecryptBaseData
	V8Decrypt         DecryptBaseData
	V9Decrypt         DecryptBaseData
	decryptionHeaders []DecryptBaseData

	groupHeaders  []GroupHeader
	allHeaderData []byte
	magicNumbers  []byte
	cryptHeaders  []byte
	iceType       int
	fileSize      int

	SecondPassThreshold int
}

func (this *IceV5File) InitIceV5File() {

	this.InitIceFile(153600)

	this.groupHeaders = make([]GroupHeader, 2)
	this.magicNumbers = make([]byte, 256)
	this.cryptHeaders = make([]byte, 48)

	// this.decryptShift = 16

	this.V5Decrypt = DecryptBaseData{
		KeyStartPos: 131,
		CrcStartPos: 10,
		CrcEndPos:   210,
		KeyConstTable: []byte{
			226,
			198,
			161,
			243,
		},
		HeaderRol: 25,
		Group2Rol: 17,
	}
	this.V6Decrypt = DecryptBaseData{
		KeyStartPos: 179,
		CrcStartPos: 80,
		CrcEndPos:   97,
		KeyConstTable: []byte{
			232,
			174,
			183,
			100,
		},
		HeaderRol: 15,
		Group2Rol: 4,
	}
	this.V7Decrypt = DecryptBaseData{
		KeyStartPos: 215,
		CrcStartPos: 23,
		CrcEndPos:   71,
		KeyConstTable: []byte{
			8,
			249,
			93,
			253,
		},
		HeaderRol: 10,
		Group2Rol: 7,
	}
	this.V8Decrypt = DecryptBaseData{
		KeyStartPos: 22,
		CrcStartPos: 84,
		CrcEndPos:   97,
		KeyConstTable: []byte{
			200,
			170,
			94,
			122,
		},
		HeaderRol: 28,
		Group2Rol: 5,
	}
	this.V9Decrypt = DecryptBaseData{
		KeyStartPos: 220,
		CrcStartPos: 189,
		CrcEndPos:   219,
		KeyConstTable: []byte{
			13,
			156,
			245,
			147,
		},
		HeaderRol: 8,
		Group2Rol: 14,
	}
	this.decryptionHeaders = []DecryptBaseData{
		this.V5Decrypt,
		this.V6Decrypt,
		this.V7Decrypt,
		this.V8Decrypt,
		this.V9Decrypt,
	}

	this.iceType = 5
	this.fileSize = 0

	// this.SecondPassThreshold = 153600
}

func (this *IceV5File) IceV5FileNew1(filename string) {
	this.InitIceV5File()
	// this.loadFile(File.OpenRead(filename))
	inFile, _ := os.Open(filename)
	this.loadFile(inFile)
}

func (this *IceV5File) IceV5FileNew2(inFile *os.File) {

	this.InitIceV5File()
	this.loadFile(inFile)
}

func (this *IceV5File) loadFile(inFile *os.File) {
	this.allHeaderData = ReadBytes(inFile, 352)
	this.iceType = int(BitConverter_ToInt32(this.allHeaderData, 8))
	this.decryptShift = this.iceType + 5
	this.fileSize = int(BitConverter_ToInt32(this.allHeaderData, 28))
	Array_Copy2(this.allHeaderData, 48, &(this.magicNumbers), 0, 256)
	Array_Copy2(this.allHeaderData, 304, &(this.cryptHeaders), 0, 48)
	inFile.Seek(0, SeekOriginBegin)
	var numArray [][]byte = this.splitGroups(inFile)
	this.header = numArray[0]
	var int32_1 int = int(BitConverter_ToInt32(this.header, 312))
	var int32_2 int = int(BitConverter_ToInt32(this.header, 328))
	this.GroupOneFiles = splitGroup(numArray[1], int32_1)
	this.GroupTwoFiles = splitGroup(numArray[2], int32_2)
}

func (this *IceV5File) IceV5FileNew3(headerFilename string, group1 string, group2 string) {
	// throw new NotSupportedException();
	fmt.Println("NotSupportedException")
}

func (this *IceV5File) calculateKeyStep1() int {
	var keyStartPos int = int(this.decryptionHeaders[this.iceType-5].KeyStartPos)
	var crcStartPos int = int(this.decryptionHeaders[this.iceType-5].CrcStartPos)
	var count int = int(this.decryptionHeaders[this.iceType-5].CrcEndPos) - crcStartPos

	reverse_result := reverse_byte(ComputeHash(this.magicNumbers, crcStartPos, count))
	var temp_key uint = this.calcBlowfishKeys(this.magicNumbers, this.getKey(this.magicNumbers, (uint)((int)(BitConverter_ToUInt32(reverse_result, 0))^(int)(BitConverter_ToUInt32(this.magicNumbers, keyStartPos))^this.fileSize^1129510338)))

	//not used
	// var key uint = this.getKey(this.magicNumbers, temp_key)

	BitConverter_ToUInt32(reverse_byte(BitConverter_GetBytes(temp_key<<this.decryptionHeaders[this.iceType-5].HeaderRol|temp_key>>(32-uint(this.decryptionHeaders[this.iceType-5].HeaderRol)))), 0)

	//not used
	// var num1 uint = temp_key>>this.decryptionHeaders[this.iceType-5].HeaderRol | temp_key<<32 - uint(this.decryptionHeaders[this.iceType-5].HeaderRol)
	// var num2 uint = key>>this.decryptionHeaders[this.iceType-5].HeaderRol | key<<32 - uint(this.decryptionHeaders[this.iceType-5].HeaderRol)
	return 0
}

func (this *IceV5File) getKey(keys []byte, temp_key uint) uint {
	// var num1 uint = (uint)(((int)(temp_key) & byteMaxValue) + int(this.decryptionHeaders[this.iceType-5].KeyConstTable[0])&byteMaxValue)
	// var num2 uint = (uint)((int)(temp_key>>8) + int(this.decryptionHeaders[this.iceType-5].KeyConstTable[1])&byteMaxValue)
	// var num3 uint = (uint)((int)(temp_key>>16) + int(this.decryptionHeaders[this.iceType-5].KeyConstTable[2])&byteMaxValue)
	// var num4 uint = (uint)((int)(temp_key>>24) + int(this.decryptionHeaders[this.iceType-5].KeyConstTable[3])&byteMaxValue)
	// var num5 byte = (byte)(this.decryptionHeaders[this.iceType-5].KeyConstTable[1] & 7)
	// var num6 byte = (byte)(this.decryptionHeaders[this.iceType-5].KeyConstTable[3] & 7)
	// var num7 byte = (byte)(this.decryptionHeaders[this.iceType-5].KeyConstTable[0] & 7)
	// var num8 byte = (byte)(this.decryptionHeaders[this.iceType-5].KeyConstTable[2] & 7)
	// return (uint)((uint)((keys[(int)(num3)]<<num8|keys[(int)(num3)]>>8-num8)&byteMaxValue)<<24|(uint)((keys[(int)(num1)]<<num7|keys[(int)(num1)]>>8-num7)&byteMaxValue)<<16|(uint)((keys[(int)(num2)]<<num5|keys[(int)(num2)]>>8-num5)&byteMaxValue)<<8) | (uint)((keys[(int)(num4)]<<num6|keys[(int)(num4)]>>8-num6)&byteMaxValue)
	return 0
}

func (this *IceV5File) calcBlowfishKeys(keys []byte, temp_key uint) uint {
	var temp_key1 uint = 2382545500 ^ temp_key
	var num1 uint = (uint)((1321528399 * temp_key1) >> 32)
	// var num2 uint = temp_key1 - num1>>1 //not used
	var num3 uint = (num1 >> 2) * 13
	for index := (int)(temp_key1) - (int)(num3) + 3; index > 0; index-- {
		temp_key1 = this.getKey(keys, temp_key1)
	}
	return (uint)((int)(temp_key1) ^ 1129510338 ^ -850380898)
}

func (this *IceV5File) splitGroups(inFile *os.File) [][]byte {
	// BinaryReader openReader = new BinaryReader(inFile);
	ReadBytes(inFile, 4)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	inFile.Seek(48, SeekOriginBegin)
	var numArray1 []byte = ReadBytes(inFile, 256)
	var block []byte = ReadBytes(inFile, 48)
	inFile.Seek(0, SeekOriginBegin)
	var numArray2 []byte = ReadBytes(inFile, 304)
	var keyStartPos int = int(this.decryptionHeaders[this.iceType-5].KeyStartPos)
	var crcStartPos int = int(this.decryptionHeaders[this.iceType-5].CrcStartPos)
	var count int = int(this.decryptionHeaders[this.iceType-5].CrcEndPos) - crcStartPos

	reverse_result := reverse_byte(ComputeHash(numArray1, crcStartPos, count))
	var temp_key uint = (uint)((int)(BitConverter_ToUInt32(reverse_result, 0)) ^ (int)(BitConverter_ToUInt32(numArray1, keyStartPos)) ^ this.fileSize ^ 1129510338)

	var key1 uint = this.getKey(numArray1, temp_key)
	var num uint = this.calcBlowfishKeys(numArray1, key1)
	var key2 uint = this.getKey(numArray1, num)
	var key3 uint = ReverseBytes(uint(num<<this.decryptionHeaders[this.iceType-5].Group2Rol) | num>>(32-uint(this.decryptionHeaders[this.iceType-5].Group2Rol)))
	var groupOneTempKey uint = num<<this.decryptionHeaders[this.iceType-5].HeaderRol | num>>(32-uint(this.decryptionHeaders[this.iceType-5].HeaderRol))
	var groupTwoTempKey uint = key2<<this.decryptionHeaders[this.iceType-5].HeaderRol | key2>>(32-uint(this.decryptionHeaders[this.iceType-5].HeaderRol))

	tmpBlewFish := new(BlewFish)
	tmpBlewFish.BlewFishNew(key3)
	var decryptedHeaderData []byte = tmpBlewFish.decryptBlock(block)

	this.groupHeaders = readHeaders(decryptedHeaderData)
	inFile.Seek(352, SeekOriginBegin)
	// numArray3 := [3][]byte{
	// 	make([]byte, 352), //byte[352],
	// 	make([]byte, 0),
	// 	make([]byte, 0),
	// }
	numArray3 := make([][]byte, 3)
	numArray3[0] = make([]byte, 352)
	numArray3[1] = make([]byte, 0)
	numArray3[2] = make([]byte, 0)

	Array_Copy1(numArray2, &numArray3[0], 304)
	Array_Copy2(decryptedHeaderData, 0, &numArray3[0], 304, 48)
	if this.groupHeaders[0].decompSize > 0 {
		numArray3[1] = this.extractGroup(this.groupHeaders[0], inFile, true, num, key2, false, false)
	}
	if this.groupHeaders[1].decompSize > 0 {
		numArray3[2] = this.extractGroup(this.groupHeaders[1], inFile, true, groupOneTempKey, groupTwoTempKey, false, false)
	}
	return numArray3
}
