package windows

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
)

func GetMainWindow(o2fy_app fyne.App) fyne.Window {
	o2fy_main_window := o2fy_app.NewWindow("O2FY")

	tabs := container.NewAppTabs(
		container.NewTabItem("汉化启动器", get_launcher_item()),
		container.NewTabItem("我要充值！", get_pso2s_item()),
		container.NewTabItem("解包/打包", get_package_item(o2fy_main_window)),
		container.NewTabItem("追加模拟器", get_simulator_item()),
		container.NewTabItem("攻略网站", get_other_item()),
	)

	//tabs.Append(container.NewTabItemWithIcon("Home", theme.HomeIcon(), widget.NewLabel("Home tab")))

	tabs.SetTabLocation(container.TabLocationLeading)

	o2fy_main_window.SetContent(tabs)

	o2fy_main_window.Resize(fyne.NewSize(800, 600))

	return o2fy_main_window

}

func GetPSO2SWindow(o2fy_app fyne.App) fyne.Window {
	o2fy_main_window := o2fy_app.NewWindow("SEGA ID 验证码小程序")

	o2fy_main_window.SetContent(get_pso2s_item())

	// o2fy_main_window.Resize(fyne.NewSize(800, 600))

	return o2fy_main_window

}

func GetPackageWindow(o2fy_app fyne.App) fyne.Window {
	// o2fy_main_window := o2fy_app.NewWindow("SEGA ID 验证码小程序")
	o2fy_main_window := o2fy_app.NewWindow("PSO2 音乐提取器 (增加了cpk和acb格式提取)")

	o2fy_main_window.SetContent(get_package_item(o2fy_main_window))

	o2fy_main_window.Resize(fyne.NewSize(600, 80))

	return o2fy_main_window

}
