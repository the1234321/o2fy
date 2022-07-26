package main

import (
	"fyne.io/fyne/v2/app"
	"github.com/the1234321/o2fy/src/theme"
	"github.com/the1234321/o2fy/src/windows"
)

// func acbExample(path string) {
// 	a, err := acb.LoadCriAcbFile(path)
// 	if err != nil {
// 		panic(err)
// 	}
// 	for name, data := range a.Files() {
// 		fmt.Printf("Write: %s\n", name)

// 		f, err := os.Create(name)
// 		if err != nil {
// 			panic(err)
// 		}
// 		defer f.Close()
// 		f.Write(data)
// 	}
// }

func main() {
	o2fy_app := app.New()
	o2fy_app.Settings().SetTheme(&theme.O2fyTheme{})
	o2fy_main_window := windows.GetPackageWindow(o2fy_app)
	// o2fy_main_window := windows.GetPSO2SWindow(o2fy_app)
	// o2fy_main_window := windows.GetMainWindow(o2fy_app)
	o2fy_main_window.ShowAndRun()

	// acbExample("9d634e82ed7f43622e445a69fa9e3b.acb")

	// hca.DeIceMusic("E:\\PSO2_Fake\\PHANTASYSTARONLINE2_JP\\pso2_bin\\data\\win32reboot\\")

	// hca.DeIceMusic("E:\\PSO2_Fake\\PHANTASYSTARONLINE2_JP\\pso2_bin\\data\\win32reboot\\")
	// icefile.UnpackNGSFiles("3be9faae726c5ea5ec14c6c7ed9326fc", "")

	// icefile.DeIceFolder("E:\\PSO2_Fake\\PHANTASYSTARONLINE2_JP\\pso2_bin\\data\\")

}
