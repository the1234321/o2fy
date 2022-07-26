package icefile

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
	"time"

	"github.com/iafan/cwalk"
)

func writeToFiles(groupToWrite [][]byte, directory string) {

	for index := 0; index < len(groupToWrite); index++ {
		var int32_v int32 = BitConverter_ToInt32(groupToWrite[index], 16)
		var fileName string = string(groupToWrite[index][64 : 64+int32_v])
		fileName = strings.Replace(fileName, string([1]byte{}[0]), "", -1)

		if !strings.Contains(fileName, ".lua") {
			continue
		}

		if _, err := os.Stat(directory); os.IsNotExist(err) {
			os.MkdirAll(directory, 0777)
			// os.Chmod(directory, 0777)
		}

		iceHeaderSize := BitConverter_ToInt32(groupToWrite[index], 0xC)
		newLength := len(groupToWrite[index]) - int(iceHeaderSize)
		var file []byte = make([]byte, newLength)
		file = groupToWrite[index][iceHeaderSize : int(iceHeaderSize)+newLength]
		filePath := directory + fileName

		err := ioutil.WriteFile(filePath, file, 0644)
		if err != nil {
			fmt.Println(err)
		}

		fmt.Println("path: ", filePath, "file_name is: ", fileName, "file size:", newLength)
	}
}

func UnpackNGSFiles(filePath string, outputPath string) {
	file_to_open, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	defer file_to_open.Close()

	type_result, _ := LoadIceFile(file_to_open)
	// fmt.Println("result is:", type_result)

	switch type_result {
	case 3:
		// iceFile = new IceV3File(file_to_open);
		iceFile := new(IceV3File)
		iceFile.IceV3FileNew1(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	case 4:
		// iceFile = new IceV4File(file_to_open);
		iceFile := new(IceV4File)
		iceFile.IceV4FileNew1(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	case 5:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	case 6:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	case 7:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	case 8:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	case 9:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath)
		writeToFiles(iceFile.GroupTwoFiles, outputPath)
	default:
		fmt.Println("err type:", type_result)
	}

}

func DeIceFolder(folderPath string) {

	deIcefolder := "解压目录" + time.Now().Format("2006_01_0_15_04_05\\")
	if _, err := os.Stat(deIcefolder); os.IsNotExist(err) {
		os.MkdirAll(deIcefolder, 0777)
		// os.Chmod(deIcefolder, 0777)
	}

	cwalk.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {

			intputPath := folderPath + path
			outputPath := deIcefolder
			// fmt.Println(intputPath, " , ", outputPath)
			UnpackNGSFiles(intputPath, outputPath)
			return nil
		})
}
