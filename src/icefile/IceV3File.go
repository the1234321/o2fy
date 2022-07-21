package icefile

import "os"

type group struct {
	originalSize uint
	dataSize     uint
	fileCount    uint
	crc32        uint
}

type stGroup struct {
	group1     GroupHeader
	group2     GroupHeader
	group1Size uint
	group2Size uint
	key        uint
	reserve    uint
}

type stInfo struct {
	r1       uint
	crc32    uint
	r2       uint
	filesize uint
}

type IceV3File struct {
	IceFile

	groupOneCount int
	groupTwoCount int
	//102400
	SecondPassThreshold int

	//Structs based on ice.exe naming
	// group
	// stGroup
	// stInfo
}

func (iceV3 *IceV3File) InitIceV3File() {
	iceV3.InitIceFile(102400)
}

func (iceV3 *IceV3File) IceV3FileNew1(inFile *os.File) *IceV3File {
	var numArray [][]byte = iceV3.splitGroups(inFile)
	iceV3.header = numArray[0]
	iceV3.GroupOneFiles = splitGroup(numArray[1], iceV3.groupOneCount)
	iceV3.GroupTwoFiles = splitGroup(numArray[2], iceV3.groupTwoCount)

	// iceV3.decryptShift = 16
	iceV3.InitIceV3File()

	return iceV3
}

func (iceV3 *IceV3File) IceV3FileNew2(headerData []byte, groupOneIn [][]byte, groupTwoIn [][]byte) *IceV3File {
	iceV3.header = headerData
	iceV3.GroupOneFiles = groupOneIn
	iceV3.GroupTwoFiles = groupTwoIn

	// iceV3.decryptShift = 16
	iceV3.InitIceV3File()

	return iceV3
}

func (iceV3 *IceV3File) splitGroups(inFile *os.File) [][]byte {
	numArray1 := make([][]byte, 3)
	// BinaryReader openReader = new BinaryReader(inFile);
	numArray1[0] = ReadBytes(inFile, 128)

	inFile.Seek(0x10, SeekOriginBegin) //Skip the ICE header

	//Read group info
	var groupInfo stGroup
	// groupInfo.group1 = new GroupHeader();
	// groupInfo.group2 = new GroupHeader();
	iceV3.ReadGroupInfoGroup(inFile, groupInfo.group1)
	iceV3.ReadGroupInfoGroup(inFile, groupInfo.group2)
	groupInfo.group1Size = uint(ReadUInt32(inFile))
	groupInfo.group2Size = uint(ReadUInt32(inFile))
	groupInfo.key = uint(ReadUInt32(inFile))
	groupInfo.key = uint(ReadUInt32(inFile))

	//Read crypt info
	var info stInfo
	info.r1 = uint(ReadUInt32(inFile))
	info.crc32 = uint(ReadUInt32(inFile))
	info.r2 = uint(ReadUInt32(inFile))
	info.filesize = uint(ReadUInt32(inFile))

	//Seek past padding/unused data
	inFile.Seek(0x30, SeekOriginCurrent)

	//Generate key
	var key uint = groupInfo.group1Size
	if key > 0 {
		key = ReverseBytes(key)
	} else if info.r2 > 0 {
		key = iceV3.GetKey(groupInfo)
	}

	//Group 1
	if groupInfo.group1.decompSize > 0 {
		numArray1[1] = iceV3.extractGroup(groupInfo.group1, inFile, (info.r2&1) > 0, key, 0, info.r2 == 8 || info.r2 == 9, true)
	}

	//Group 2
	if groupInfo.group2.decompSize > 0 {
		numArray1[2] = iceV3.extractGroup(groupInfo.group2, inFile, (info.r2&1) > 0, key, 0, info.r2 == 8 || info.r2 == 9, true)
	}
	iceV3.groupOneCount = int(groupInfo.group1.count)
	iceV3.groupTwoCount = int(groupInfo.group2.count)

	return numArray1
}

func (iceV3 *IceV3File) ReadGroupInfoGroup(inFile *os.File, grp GroupHeader) {
	grp.decompSize = uint(ReadUInt32(inFile))
	grp.compSize = uint(ReadUInt32(inFile))
	grp.count = uint(ReadUInt32(inFile))
	grp.CRC = uint(ReadUInt32(inFile))
}

//uint reversal from ice.exe
func (iceV3 *IceV3File) bswap(v uint) uint {
	var r uint = v & 0xFF
	r <<= 8
	v >>= 8
	r |= v & 0xFF
	r <<= 8
	v >>= 8
	r |= v & 0xFF
	r <<= 8
	v >>= 8
	r |= v & 0xFF

	return r
}

func (iceV3 *IceV3File) GetKey(group stGroup) uint {
	return group.group1.decompSize ^ group.group2.decompSize ^ group.group2Size ^ group.key ^ 0xC8D7469A
}
