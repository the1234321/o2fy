package bgm

import (
	"encoding/binary"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/iafan/cwalk"
	"github.com/the1234321/o2fy/src/bgm/hca"
)

func checkIsCPK(intputPath string) bool {
	file, err := os.Open(intputPath)
	if err != nil {
		return false
	}
	defer file.Close()

	var type_num uint32
	binary.Read(file, binary.BigEndian, &type_num)

	return type_num == 1129335584
}

func checkIsAWB(intputPath string) bool {
	file, err := os.Open(intputPath)
	if err != nil {
		return false
	}
	defer file.Close()

	var type_num uint32
	binary.Read(file, binary.BigEndian, &type_num)

	return type_num == 1095127858

}

func checkIsNGS(path string) bool {
	return strings.Contains(path, "win32reboot")
}

func DeIceMusic(folderPath string) {

	// deIcefolder := "解压目录" + time.Now().Format("2006_01_0_15_04_05\\")
	basefolder := "解压目录" + time.Now().Format("2006_01_02_15_04_05\\")
	outputPathO2Normal := basefolder + "老O2音乐\\WAV格式音乐\\"
	outputPathO2Cpk := basefolder + "老O2音乐\\CPK格式音乐\\"
	outputPathO2Acb := basefolder + "老O2音乐\\ACB格式音乐\\"

	outputPathNGSNormal := basefolder + "NGS音乐\\WAV格式音乐\\"
	outputPathNGSCpk := basefolder + "NGS音乐\\CPK格式音乐\\"
	outputPathNGSAcb := basefolder + "NGS音乐\\ACB格式音乐\\"

	// os.MkdirAll(deIcefolder, 0777)
	os.MkdirAll(outputPathO2Normal, 0777)
	os.MkdirAll(outputPathO2Cpk, 0777)
	os.MkdirAll(outputPathO2Acb, 0777)

	os.MkdirAll(outputPathNGSNormal, 0777)
	os.MkdirAll(outputPathNGSCpk, 0777)
	os.MkdirAll(outputPathNGSAcb, 0777)

	cwalk.Walk(folderPath,
		func(path string, info os.FileInfo, err error) error {

			intputPath := folderPath + path

			if strings.Contains(intputPath, ".") {
				return nil
			}

			if info.IsDir() {
				return nil
			}

			var outputPath string

			/////////
			if checkIsCPK(intputPath) {
				// fmt.Println("found cpk")
				if checkIsNGS(intputPath) {
					outputPath = outputPathNGSCpk + filepath.Base(path) + ".cpk"
				} else {
					outputPath = outputPathO2Cpk + filepath.Base(path) + ".cpk"
				}
				bytesRead, _ := ioutil.ReadFile(intputPath)
				ioutil.WriteFile(outputPath, bytesRead, 0644)
				return nil
			}
			////////

			/////////
			if checkIsAWB(intputPath) {
				// fmt.Println("found acb")

				if checkIsNGS(intputPath) {
					outputPath = outputPathNGSAcb + filepath.Base(path) + ".acb"
				} else {
					outputPath = outputPathO2Acb + filepath.Base(path) + ".acb"
				}

				bytesRead, _ := ioutil.ReadFile(intputPath)
				ioutil.WriteFile(outputPath, bytesRead, 0644)
				return nil
			}
			////////

			// bytesRead, _ := ioutil.ReadFile(intputPath)
			// ioutil.WriteFile(outputPath, bytesRead, 0644)

			if checkIsNGS(intputPath) {
				outputPath = outputPathNGSNormal + filepath.Base(path) + ".wav"
			} else {
				outputPath = outputPathO2Normal + filepath.Base(path) + ".wav"
			}

			// fmt.Println(intputPath, " , ", outputPath)

			hca_decoder := hca.NewDecoder()
			ret := hca_decoder.DecodeFromFile(intputPath, outputPath)
			if ret == hca.RetSuccess {
				// fmt.Println(outputPath)
				// os.Remove(outputPath)
				// os.Remove(intputPath)
				// fmt.Println("successed!")
			} else if ret == hca.RetDecodeError {
				// bytesRead, _ := ioutil.ReadFile(intputPath)
				// ioutil.WriteFile(outputPath, bytesRead, 0644)
				// fmt.Println("key err, backing up!")
			}

			return nil
		})
}
