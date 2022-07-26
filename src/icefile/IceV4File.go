package icefile

import (
	"os"
)

type IceV4File struct {
	IceFile

	keyconstant_1       uint
	keyconstant_2       uint
	keyconstant_3       uint
	groupOneCount       int
	groupTwoCount       int
	SecondPassThreshold int
}

func (iceV4 *IceV4File) InitIceV4File() {

	iceV4.InitIceFile(102400)
	iceV4.keyconstant_1 = 1129510338
	iceV4.keyconstant_2 = 3444586398
	iceV4.keyconstant_3 = 613566757

	iceV4.groupOneCount = 0
	iceV4.groupTwoCount = 0

}

func (iceV4 *IceV4File) IceV4FileNew1(inFile *os.File) {
	// iceV4.decryptShift = 16
	// iceV4.SecondPassThreshold = 102400
	iceV4.InitIceV4File()

	var numArray [][]byte = iceV4.splitGroups(inFile)

	iceV4.header = numArray[0]
	iceV4.GroupOneFiles = splitGroup(numArray[1], iceV4.groupOneCount)
	iceV4.GroupTwoFiles = splitGroup(numArray[2], iceV4.groupTwoCount)

}

func (iceV4 *IceV4File) IceV4FileNew2(headerData []byte, groupOneIn [][]byte, groupTwoIn [][]byte) {
	// iceV4.decryptShift = 16
	// iceV4.SecondPassThreshold = 102400
	iceV4.InitIceV4File()

	iceV4.header = headerData
	iceV4.GroupOneFiles = groupOneIn
	iceV4.GroupTwoFiles = groupTwoIn
	iceV4.groupOneCount = len(groupOneIn)
	iceV4.groupTwoCount = len(groupTwoIn)
}

func (iceV4 *IceV4File) splitGroups(inFile *os.File) [][]byte {
	// BinaryReader openReader = new BinaryReader(inFile);

	// inFile, _ = os.Open("9dc4e510eddebae273570fa5f5265eec")

	ReadBytes(inFile, 4)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	ReadInt32(inFile)
	var num1 int32 = ReadInt32(inFile)
	var compSize int32 = ReadInt32(inFile)

	// int num2 = num1 == 1 ? 288 : 272; //not used
	// var num2 int
	// if num1 == 1 {
	// 	num2 = 288
	// } else {
	// 	num2 = 272
	// }

	var blowfishKeys BlowfishKeys = iceV4.getBlowfishKeys(ReadBytes(inFile, 256), int(compSize))

	numArray1 := make([][]byte, 3)
	// var numArray2 [48]byte //not used
	var decryptedHeaderData []byte
	if num1 == 1 || num1 == 9 {
		inFile.Seek(0, SeekOriginBegin)
		var numArray3 []byte = ReadBytes(inFile, 288)
		var block []byte = ReadBytes(inFile, 48)

		tempBlewFish := new(BlewFish)
		tempBlewFish.BlewFishNew(blowfishKeys.groupHeadersKey)

		decryptedHeaderData = tempBlewFish.decryptBlock(block)

		numArray1[0] = make([]byte, 336)
		Array_Copy1(numArray3, &numArray1[0], 288)
		Array_Copy2(decryptedHeaderData, 0, &numArray1[0], 288, len(decryptedHeaderData))
	} else {
		switch num1 {
		case 8:
			inFile.Seek(288, SeekOriginBegin)
			decryptedHeaderData = ReadBytes(inFile, 48)
			inFile.Seek(0, SeekOriginBegin)
			numArray1[0] = ReadBytes(inFile, 336)

		case 327680:
			inFile.Seek(288, SeekOriginBegin)
			decryptedHeaderData = ReadBytes(inFile, 48)
			inFile.Seek(0, SeekOriginBegin)
			numArray1[0] = ReadBytes(inFile, 336)

		default:
			inFile.Seek(288, SeekOriginBegin)
			decryptedHeaderData = ReadBytes(inFile, 48)
			inFile.Seek(0, SeekOriginBegin)
			numArray1[0] = ReadBytes(inFile, 336)

		}
	}
	var groupHeaderArray []GroupHeader = readHeaders(decryptedHeaderData)
	iceV4.groupOneCount = int(groupHeaderArray[0].count)
	iceV4.groupTwoCount = int(groupHeaderArray[1].count)
	inFile.Seek(336, SeekOriginBegin)
	numArray1[1] = make([]byte, 0)
	numArray1[2] = make([]byte, 0)
	// #if DEBUG
	// if num1 == 8 || num1 == 9 {
	// 	fmt.Println("NGS Ice detected")
	// }
	// #endif
	if groupHeaderArray[0].decompSize > 0 {
		numArray1[1] = iceV4.extractGroup(groupHeaderArray[0], inFile, uint(num1&1) > 0, blowfishKeys.groupOneBlowfish[0], blowfishKeys.groupOneBlowfish[1], num1 == 8 || num1 == 9, false)
	}
	if groupHeaderArray[1].decompSize > 0 {
		numArray1[2] = iceV4.extractGroup(groupHeaderArray[1], inFile, uint(num1&1) > 0, blowfishKeys.groupTwoBlowfish[0], blowfishKeys.groupTwoBlowfish[1], num1 == 8 || num1 == 9, false)
	}
	return numArray1
}

func (iceV4 *IceV4File) getBlowfishKeys(magicNumbers []byte, compSize int) BlowfishKeys {
	blowfishKeys := new(BlowfishKeys)
	blowfishKeys.BlowfisKeysNew()

	reverse_result := reverse_byte(ComputeHash(magicNumbers, 124, 96))
	var temp_key uint = uint((int(BitConverter_ToUInt32(reverse_result, 0)) ^ int(BitConverter_ToUInt32(magicNumbers, 108)) ^ compSize ^ 1129510338))
	var key uint = iceV4.getKey(magicNumbers, temp_key)
	blowfishKeys.groupOneBlowfish[0] = iceV4.calcBlowfishKeys(magicNumbers, key)
	blowfishKeys.groupOneBlowfish[1] = iceV4.getKey(magicNumbers, blowfishKeys.groupOneBlowfish[0])
	blowfishKeys.groupTwoBlowfish[0] = blowfishKeys.groupOneBlowfish[0]>>15 | blowfishKeys.groupOneBlowfish[0]<<17
	blowfishKeys.groupTwoBlowfish[1] = blowfishKeys.groupOneBlowfish[1]>>15 | blowfishKeys.groupOneBlowfish[1]<<17
	var x uint = blowfishKeys.groupOneBlowfish[0]<<13 | blowfishKeys.groupOneBlowfish[0]>>19
	blowfishKeys.groupHeadersKey = ReverseBytes(x)

	return *blowfishKeys
}

func (iceV4 *IceV4File) getKey(keys []byte, temp_key uint) uint {
	var num1 uint = uint((byte(int(temp_key)) & byteMaxValue) + 93&byteMaxValue)
	var num2 uint = uint(byte(int(temp_key>>8)) + 63&byteMaxValue)
	var num3 uint = uint(byte(int(temp_key>>16)) + 69&byteMaxValue)
	var num4 uint = uint(byte(int(temp_key>>24)) - 58&byteMaxValue)

	// fmt.Println(num1, " ", num2, " ", num3, " ", num4, "")

	ret := (uint)((int((keys[int(num2)]<<7|keys[int(num2)]>>1)&byteMaxValue)<<24 | int((keys[int(num4)]<<6|keys[int(num4)]>>2)&byteMaxValue)<<16 | int((keys[int(num1)]<<5|keys[int(num1)]>>3)&byteMaxValue)<<8) | int((keys[int(num3)]<<5|keys[int(num3)]>>3)&byteMaxValue))
	// fmt.Println("keys:", keys, ", temp_key", temp_key)
	// fmt.Println("getKey ", ret)
	return ret
}

func (iceV4 *IceV4File) calcBlowfishKeys(keys []byte, temp_key uint) uint {

	var temp_key1 uint = 2382545500 ^ temp_key
	var num1 uint = (uint)(613566757 * temp_key1 >> 32)
	var num2 uint = ((((temp_key1 - num1) >> 1) + num1) >> 2) * 7

	for index := int(temp_key1) - int(num2) + 2; index > 0; index-- {
		temp_key1 = iceV4.getKey(keys, temp_key1)
	}

	ret := uint32(int(temp_key1) ^ 1129510338 ^ -850380898)
	return uint(ret)
}

func (iceV4 *IceV4File) getRawData(compress bool, forceUnencrypted bool) []byte {
	return iceV4.packFile(iceV4.header, combineGroup(iceV4.GroupOneFiles, true), combineGroup(iceV4.GroupTwoFiles, true), iceV4.groupOneCount, iceV4.groupTwoCount, compress, forceUnencrypted)
}

func (iceV4 *IceV4File) packFile(
	headerData []byte,
	groupOneIn []byte,
	groupTwoIn []byte,
	groupOneCount int,
	groupTwoCount int,
	compress bool,
	forceUnencrypted bool) []byte {
	//Setup ICE header
	Array_Copy2(BitConverter_GetBytes(1), 0, &headerData, 0x18, 0x4)

	//Set group data in ICE header
	Array_Copy2(BitConverter_GetBytes(uint(len(groupOneIn))), 0, &headerData, 0x120, 0x4)
	Array_Copy2(BitConverter_GetBytes(uint(len(groupTwoIn))), 0, &headerData, 0x130, 0x4)
	Array_Copy2(BitConverter_GetBytes(uint(groupOneCount)), 0, &headerData, 0x128, 0x4)
	Array_Copy2(BitConverter_GetBytes(uint(groupTwoCount)), 0, &headerData, 0x138, 0x4)
	Array_Copy2(BitConverter_GetBytes(uint(len(groupOneIn))), 0, &headerData, 0x140, 0x4)
	Array_Copy2(BitConverter_GetBytes(uint(len(groupTwoIn))), 0, &headerData, 0x144, 0x4)

	var compressedContents1 []byte = getCompressedContents(groupOneIn, compress)
	var compressedContents2 []byte = getCompressedContents(groupTwoIn, compress)
	var compSize int = len(headerData) + len(compressedContents1) + len(compressedContents2)

	//Set main CRC (Should be done after potential compression, but before encryption)

	// var mainCrc = new Crc32Alt().GetCrc32(compressedContents2, new Crc32Alt().GetCrc32(compressedContents1));
	tempCrc32Alt1 := new(Crc32Alt)
	tempCrc32Alt1.Crc32Alt()
	tempCrc32Alt2 := new(Crc32Alt)
	tempCrc32Alt2.Crc32Alt()
	mainCrc := tempCrc32Alt1.GetCrc32(compressedContents2, tempCrc32Alt2.GetCrc32(compressedContents1, 0))

	Array_Copy2(BitConverter_GetBytes(mainCrc), 0, &headerData, 0x14, 0x4)

	if forceUnencrypted {
		//Set encryption flag to 0
		headerData[24] = 0
		headerData[25] = 0
		headerData[26] = 0
		headerData[27] = 0

		//Set array to 0
		for index := 0; index < 256; index++ {
			headerData[32+index] = 0
		}

		//Set encrypted group 1 size to 0
		headerData[320] = 0
		headerData[321] = 0
		headerData[322] = 0
		headerData[323] = 0

		//Set encrypted group 2 size to 0
		headerData[324] = 0
		headerData[325] = 0
		headerData[326] = 0
		headerData[327] = 0
		compress = false
	}
	var boolean bool = BitConverter_ToBoolean(headerData, 24)
	magicNumbers := make([]byte, 256)
	Array_Copy2(headerData, 32, &magicNumbers, 0, 256)
	var blowfishKeys BlowfishKeys = iceV4.getBlowfishKeys(magicNumbers, compSize)
	// numArray1 := make([]byte, 0) //not used
	// numArray2 := make([]byte, 0)
	outBytes := make([]byte, compSize)
	var destinationIndex1 int = 336
	var destinationIndex2 int = 288
	Array_Copy2(BitConverter_GetBytes(uint(len(groupOneIn))), 0, &headerData, destinationIndex2, 4)
	Array_Copy2(BitConverter_GetBytes(uint(len(groupTwoIn))), 0, &headerData, destinationIndex2+16, 4)
	if compress {
		if uint(len(groupOneIn)) > 0 {
			Array_Copy2(BitConverter_GetBytes(uint(len(compressedContents1))), 0, &headerData, destinationIndex2+4, 4)
			var num int = 4
			if (uint)(len(groupTwoIn)) > 0 {
				num = 2
			}
			Array_Copy2(BitConverter_GetBytes(uint(len(groupOneIn)-num)), 0, &headerData, 320, 4)
		}
		if uint(len(groupTwoIn)) > 0 {
			Array_Copy2(BitConverter_GetBytes(uint(len(compressedContents2))), 0, &headerData, destinationIndex2+20, 4)
			var num int = 3
			if uint(len(groupOneIn)) > 0 {
				num = 5
			}
			Array_Copy2(BitConverter_GetBytes(uint(len(groupTwoIn)-num)), 0, &headerData, 324, 4)
		}
	} else {
		headerData[destinationIndex2+4] = 0
		headerData[destinationIndex2+5] = 0
		headerData[destinationIndex2+6] = 0
		headerData[destinationIndex2+7] = 0
		headerData[destinationIndex2+20] = 0
		headerData[destinationIndex2+21] = 0
		headerData[destinationIndex2+22] = 0
		headerData[destinationIndex2+23] = 0
	}
	if uint(len(compressedContents1)) > 0 {
		var numArray4 []byte = iceV4.packGroup(compressedContents1, blowfishKeys.groupOneBlowfish[0], blowfishKeys.groupOneBlowfish[1], boolean)
		Array_Copy2(numArray4, 0, &outBytes, destinationIndex1, len(numArray4))
		destinationIndex1 += len(numArray4)
	}
	if uint(len(compressedContents2)) > 0 {
		var numArray4 []byte = iceV4.packGroup(compressedContents2, blowfishKeys.groupTwoBlowfish[0], blowfishKeys.groupTwoBlowfish[1], boolean)
		Array_Copy2(numArray4, 0, &outBytes, destinationIndex1, len(numArray4))
		// var num int = destinationIndex1 + len(numArray4) //not used
	}

	//CRC32 for groups
	// Array.Copy(BitConverter.GetBytes(new Crc32Alt().GetCrc32(compressedContents1)), 0, headerData, 0x12C, 0x4);
	tempCrc32Alt1 = new(Crc32Alt)
	tempCrc32Alt1.Crc32Alt()
	Array_Copy2(BitConverter_GetBytes(tempCrc32Alt1.GetCrc32(compressedContents1, 0)), 0, &headerData, 0x12C, 0x4)

	// Array.Copy(BitConverter.GetBytes(new Crc32Alt().GetCrc32(compressedContents2)), 0, headerData, 0x13C, 0x4);
	tempCrc32Alt2 = new(Crc32Alt)
	tempCrc32Alt2.Crc32Alt()
	Array_Copy2(BitConverter_GetBytes(tempCrc32Alt2.GetCrc32(compressedContents2, 0)), 0, &headerData, 0x13C, 0x4)

	Array_Copy1(headerData, &outBytes, 336)
	if boolean {

		// var blewFish BlewFish= new BlewFish(blowfishKeys.groupHeadersKey);
		blewFish := new(BlewFish)
		blewFish.BlewFishNew(blowfishKeys.groupHeadersKey)

		block := make([]byte, 48)
		Array_Copy2(headerData, 288, &block, 0, 48)
		Array_Copy2(blewFish.encryptBlock(block), 0, &outBytes, 288, 48)
	}
	Array_Copy2(BitConverter_GetBytes(uint(compSize)), 0, &outBytes, 28, 4)
	return outBytes
}

type BlowfishKeys struct {
	groupHeadersKey  uint
	groupOneBlowfish []uint
	groupTwoBlowfish []uint
}

func (this *BlowfishKeys) BlowfisKeysNew() {
	this.groupOneBlowfish = make([]uint, 2)
	this.groupTwoBlowfish = make([]uint, 2)
}
