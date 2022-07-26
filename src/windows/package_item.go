package windows

import (
	"fmt"
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/the1234321/o2fy/src/bgm"
)

func PathExists(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	return false
}

func get_package_item(o2fy_main_window fyne.Window) fyne.CanvasObject {
	input := widget.NewEntry()
	input.SetPlaceHolder("..输入路径 (例子: X:\\pso2_bin\\)")

	var ButtonStart *widget.Button

	var ended bool = false

	ButtonStart = widget.NewButton("开始解包音乐", func() {

		if ended {
			o2fy_main_window.Close()
		} else {
			var path string = input.Text + "\\pso2.exe"

			if !PathExists(path) {
				input.SetPlaceHolder("\"" + input.Text + "\"" + " 路径不对，请重新输入")
				input.SetText("")
			} else {
				path := input.Text + "\\data\\"
				fmt.Println("输入正确: ", path)

				ButtonStart.SetText("解压中...")

				input.Disable()
				ButtonStart.Disable()

				bgm.DeIceMusic(path)

				ended = true
				ButtonStart.Enable()
				ButtonStart.SetText("全部音乐解压完成，点击退出")
			}
		}

	})

	content := container.NewVBox(input, ButtonStart)

	return content
}
