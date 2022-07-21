package windows

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/widget"
	"github.com/the1234321/o2fy/src/parse"
)

func get_pso2s_item() fyne.CanvasObject {

	// readme := widget.NewLabel("复制验证码图片, 会自动转换成文字")
	readme := canvas.NewText("该程序打开时，会自动将验证码图片转换成文字", color.Opaque)
	readme.TextSize = 20

	content := container.NewVBox(
		readme,
		container.NewHBox(
			widget.NewLabel("有问题请加群:21587709"),
		),
	)
	content.Add(container.NewHBox(
		widget.NewButton("打开充值网站1", func() {
			parse.Openbrowser("https://mcha.isao.net/profile_oem/OEMLogin.php?product_name=pso2&p_siteno=P00011")
		}),
		widget.NewButton("打开充值网站2", func() {
			parse.Openbrowser("https://gw.sega.jp/gw/login/")
		}),
	))

	result_str := binding.NewString()
	result_str.Set("")

	go func() { parse.StartParse(&result_str) }()

	return content
}
