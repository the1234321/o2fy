package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"strings"

	"github.com/the1234321/o2fy/src/icefile"
)

func writeToFiles(groupToWrite [][]byte, directory string) {

	if _, err := os.Stat(directory); os.IsNotExist(err) {
		os.Mkdir(directory, 0777)
		os.Chmod(directory, 0777)
	}

	for index := 0; index < len(groupToWrite); index++ {
		var int32_v int32 = icefile.BitConverter_ToInt32(groupToWrite[index], 16)
		var fileName string = string(groupToWrite[index][64 : 64+int32_v])
		fileName = strings.Replace(fileName, string([1]byte{}[0]), "", -1)

		iceHeaderSize := icefile.BitConverter_ToInt32(groupToWrite[index], 0xC)
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

func unpackNGSFiles(filePath string, outputPath string) {
	file_to_open, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
	}
	type_result, _ := icefile.LoadIceFile(file_to_open)
	fmt.Println("result is:", type_result)

	postFix := "_解压后/"

	switch type_result {
	case 3:
		// iceFile = new IceV3File(file_to_open);
		iceFile := new(icefile.IceV3File)
		iceFile.IceV3FileNew1(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	case 4:
		// iceFile = new IceV4File(file_to_open);
		iceFile := new(icefile.IceV4File)
		iceFile.IceV4FileNew1(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	case 5:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(icefile.IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	case 6:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(icefile.IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	case 7:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(icefile.IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	case 8:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(icefile.IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	case 9:
		// iceFile = new IceV5File(file_to_open);
		iceFile := new(icefile.IceV5File)
		iceFile.IceV5FileNew2(file_to_open)
		writeToFiles(iceFile.GroupOneFiles, outputPath+path.Base(filePath)+postFix)
		writeToFiles(iceFile.GroupTwoFiles, outputPath+path.Base(filePath)+postFix)
	default:
		fmt.Println("err type:", type_result)
	}

}

func main() {
	// o2fy_app := app.New()
	// o2fy_app.Settings().SetTheme(&theme.O2fyTheme{})
	// o2fy_main_window := windows.GetPSO2SWindow(o2fy_app)
	// // o2fy_main_window := windows.GetMainWindow(o2fy_app)
	// o2fy_main_window.ShowAndRun()

	unpackNGSFiles("4f02512ea12c3db46207c63ac5f550", "")

}
