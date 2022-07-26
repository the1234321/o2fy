package icefile

import (
	"fmt"
	"log"
	"os"

	oodle "github.com/new-world-tools/go-oodle"
)

type IceFile struct {
	// = 16;
	decryptShift int
	/// <summary>
	/// Group1 Files
	/// </summary>
	GroupOneFiles [][]byte
	/// <summary>
	/// Group2 Files
	/// </summary>
	GroupTwoFiles [][]byte
	/// <summary>
	/// Header of Ice file
	/// </summary>
	header []byte

	SecondPassThreshold int
}

type GroupHeader struct {
	decompSize uint
	compSize   uint
	count      uint
	CRC        uint
}

//public uint getStoredSize() => compSize > 0U ? compSize : decompSize;
func (this GroupHeader) getStoredSize() uint {
	if this.compSize > 0 {
		return this.compSize
	} else {
		return this.decompSize
	}
}

func ReverseBytes(x uint) uint {
	x = x>>16 | x<<16
	return (x&4278255360)>>8 | uint((int(x)&16711935)<<8)
}

func splitGroup(groupToSplit []byte, fileCount int) [][]byte {
	numArray := make([][]byte, fileCount)
	var sourceIndex int = 0
	for index := 0; index < fileCount && sourceIndex < len(groupToSplit); index++ {
		var int32_v int32 = BitConverter_ToInt32(groupToSplit, int(sourceIndex+4))
		numArray[index] = make([]byte, int32_v)
		Array_Copy2(groupToSplit, sourceIndex, &numArray[index], 0, int(int32_v))
		sourceIndex += int(int32_v)

	}
	return numArray
}

func (this *IceFile) InitIceFile(SecondPassThresholdValue int) {
	this.SecondPassThreshold = 102400
	this.decryptShift = 16
}

func (this *IceFile) decryptGroup(buffer []byte, key1 uint, key2 uint, v3Decrypt bool) []byte {
	var block1 []byte = make([]byte, len(buffer))
	if v3Decrypt == false {
		block1 = decrypt_block(buffer, uint(len(buffer)), key1, this.decryptShift)
	} else {
		Array_Copy2(buffer, 0, &block1, 0, len(buffer))
	}

	newBlewFish1 := new(BlewFish)
	newBlewFish1.BlewFishNew(ReverseBytes(key1))
	var block2 []byte = newBlewFish1.decryptBlock(block1)

	var numArray []byte = block2

	// fmt.Println("block2.Length: ", len(block2), ", SecondPassThreshold: ", this.SecondPassThreshold)

	if len(block2) <= this.SecondPassThreshold && !v3Decrypt {
		newBlewFish2 := new(BlewFish)
		newBlewFish2.BlewFishNew(ReverseBytes(key2))
		numArray = newBlewFish2.decryptBlock(block2)
	}
	return numArray
}

func combineGroup(filesToJoin [][]byte, headerLess bool) []byte {

	var outBytes []byte
	for i := 0; i < len(filesToJoin); i++ {
		for j := 0; j < len(filesToJoin[i]); j++ {
			outBytes = append(outBytes, filesToJoin[i][j])
		}
	}
	return outBytes
}

func (this *IceFile) extractGroup(
	header GroupHeader,
	inFile *os.File,
	encrypt bool,
	groupOneTempKey uint,
	groupTwoTempKey uint,
	ngsMode bool,
	v3Decrypt bool) []byte {

	var buffer []byte = ReadBytes(inFile, int(header.getStoredSize()))
	// var inData []byte= !encrypt ? buffer : decryptGroup(buffer, groupOneTempKey, groupTwoTempKey, v3Decrypt);
	var inData []byte

	// fmt.Println(header.compSize)
	if !encrypt {
		inData = buffer
	} else {
		inData = this.decryptGroup(buffer, groupOneTempKey, groupTwoTempKey, v3Decrypt)
	}
	// return header.compSize <= 0 ? inData : !ngsMode ? decompressGroup(inData, header.decompSize) : decompressGroupNgs(inData, header.decompSize);

	if header.compSize <= 0 {
		return inData
	} else {
		// !ngsMode ? decompressGroup(inData, header.decompSize) : decompressGroupNgs(inData, header.decompSize);
		if !ngsMode {
			return decompressGroup(inData, header.decompSize)
		} else {
			return decompressGroupNgs(inData, header.decompSize)
		}
	}
}

func decompressGroup(inData []byte, bufferLength uint) []byte {
	input := make([]byte, len(inData))
	Array_Copy1(inData, &input, len(input))
	for index := 0; index < len(input); index++ {
		input[index] ^= 149
	}
	return Decompress(input, bufferLength)
}

func decompressGroupNgs(inData []byte, bufferLength uint) []byte {

	if !oodle.IsDllExist() {
		err := oodle.Download()
		if err != nil {
			log.Fatalf("no oo2core_9_win64.dll")
		}
	}

	decompressedData, err := oodle.Decompress(inData, int64(bufferLength))
	if err != nil {
		fmt.Println("oodle err: ", err)
	}
	return decompressedData

	// return Oodle.Decompress(inData, bufferLength)
	return []byte{}
}

func readHeaders(decryptedHeaderData []byte) []GroupHeader {

	// fmt.Println(decryptedHeaderData)

	groupHeaderArray := make([]GroupHeader, 2)

	groupHeaderArray[0].decompSize = uint(BitConverter_ToUInt32(decryptedHeaderData, 0))
	groupHeaderArray[0].compSize = uint(BitConverter_ToUInt32(decryptedHeaderData, 4))
	groupHeaderArray[0].count = uint(BitConverter_ToUInt32(decryptedHeaderData, 8))
	groupHeaderArray[0].CRC = uint(BitConverter_ToUInt32(decryptedHeaderData, 12))
	// groupHeaderArray[1] = new GroupHeader();
	groupHeaderArray[1].decompSize = uint(BitConverter_ToUInt32(decryptedHeaderData, 16))
	groupHeaderArray[1].compSize = uint(BitConverter_ToUInt32(decryptedHeaderData, 20))
	groupHeaderArray[1].count = uint(BitConverter_ToUInt32(decryptedHeaderData, 24))
	groupHeaderArray[1].CRC = uint(BitConverter_ToUInt32(decryptedHeaderData, 28))
	return groupHeaderArray
}

func (this *IceFile) packGroup(buffer []byte, key1 uint, key2 uint, encrypt bool) []byte {
	if !encrypt {
		return buffer
	}
	var block []byte = buffer
	if len(buffer) <= this.SecondPassThreshold {
		tempBlewFish := new(BlewFish)
		tempBlewFish.BlewFishNew(ReverseBytes(key2))
		block = tempBlewFish.encryptBlock(buffer)
	}
	tempBlewFish2 := new(BlewFish)
	tempBlewFish2.BlewFishNew(ReverseBytes(key1))
	var data_block []byte = tempBlewFish2.encryptBlock(block)
	return decrypt_block2(data_block, (uint)(len(data_block)), key1)
}

func getCompressedContents(buffer []byte, compress bool) []byte {
	if !compress || (uint)(len(buffer)) <= 0 {
		return buffer
	}

	//TODO
	// var numArray []byte = PrsCompDecomp.compress(buffer)
	var numArray []byte
	for index := 0; index < len(numArray); index++ {
		numArray[index] ^= 149
	}
	return numArray
}

func LoadIceFile(inFile *os.File) (int /* *IceFile*/, error) {
	inFile.Seek(8, SeekOriginBegin)
	var num int = int(ReadByte(inFile))
	inFile.Seek(0, SeekOriginBegin)
	// 	var iceFile *IceFile = nil
	// 	switch num {
	// 	case 3:
	// 		iceFile = new(IceV3File)
	// 		iceFile = IceV3FileNew1(inFile)
	// 		break
	// 	case 4:
	// 		iceFile = IceV4FileNew1(inFile)
	// 		break
	// 	case 5:
	// 		iceFile = IceV5FileNew2(inFile)
	// 		break
	// 	case 6:
	// 		iceFile = IceV5FileNew2(inFile)
	// 		break
	// 	case 7:
	// 		iceFile = IceV5FileNew2(inFile)
	// 		break
	// 	case 8:
	// 		iceFile = IceV5FileNew2(inFile)
	// 		break
	// 	case 9:
	// 		iceFile = IceV5FileNew2(inFile)
	// 		break
	// 	default:
	// 		return nil, fmt.Errorf("Invalid version: %i", num)
	// 	}
	// 	inStream.Dispose()
	return num, nil
}
