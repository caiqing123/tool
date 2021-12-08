package main

import (
	_ "embed"
	"fmt"
	"os"
	"strconv"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

const timeTemplate1 = "2006-01-02 15:04:05"

//go:embed img/time_icon.jpeg
var icon []byte

func main() {
	os.Setenv("FYNE_FONT", "./msyh.ttc")
	myApp := app.New()
	setting := myApp.Settings()
	setting.SetTheme(theme.LightTheme())
	resource := fyne.NewStaticResource("time_icon", icon)
	myApp.SetIcon(resource)
	//设置名称
	myWindow := myApp.NewWindow("时间戳转换工具")

	//开始画布
	timerGridElems3, ticker := TimeNow()
	defer ticker.Stop()
	gridElems := TimestampToDate()
	gridElems2 := DateToTimestamp()
	timerGridElems3 = append(timerGridElems3, gridElems...)
	timerGridElems3 = append(timerGridElems3, gridElems2...)

	//New 返回一个新的 Container 实例，其中包含指定的 CanvasObjects，它将根据指定的 Layout 进行布局。 layout.NewGridLayout(4) 传入每行的列数,这里为4列
	grid := container.New(layout.NewGridLayout(4), timerGridElems3...)
	//设置内容
	myWindow.SetContent(grid)
	//设置窗口尺寸大小
	myWindow.Resize(fyne.NewSize(500, 50))
	//启动
	myWindow.ShowAndRun()
	os.Unsetenv("FYNE_FONT")
}

func TimeNow() ([]fyne.CanvasObject, *time.Ticker) {
	timeNow := 1
	now := time.Now().Unix()
	// 转换为字符串
	millionSec := strconv.Itoa(int(now))
	//设置时间戳标签默认值
	timeStampInp := widget.NewLabel(millionSec)
	//赋值时间戳标签的文本
	timeStampInp.SetText(strconv.Itoa(int(now)))
	//设置日期标签默认值
	text3 := widget.NewLabel("DATE")
	nums, err := strconv.Atoi(timeStampInp.Text)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	//转换时间戳变为正常格式
	date, err := timeStampToDate(nums, timeNow)
	if err != nil {
		fmt.Println(err)
		return nil, nil
	}
	//赋值日期文本
	text3.SetText(date)
	ticker := time.NewTicker(time.Second)
	//开启协程更新时间
	go func(label, text3 *widget.Label) {
		for t := range ticker.C {
			if timeNow == 1 {
				millionSec = strconv.Itoa(int(t.Unix()))
			} else {
				millionSec = strconv.Itoa(int(t.UnixNano() / 1e6))
			}
			label.SetText(millionSec)
			nums, err := strconv.Atoi(timeStampInp.Text)
			if err != nil {
				fmt.Println(err)
				return
			}
			date, err := timeStampToDate(nums, timeNow)
			if err != nil {
				return
			}
			text3.SetText(date)
		}
	}(timeStampInp, text3)
	//设置按钮 添加点击触发事件
	click1 := widget.NewButton("复制", func() {
		//复制到剪切版
		copyClipBoard(text3.Text)
	})
	//添加按钮图标
	click1.SetIcon(theme.ContentCopyIcon())

	provinceSelect := widget.NewSelect([]string{"秒", "毫秒"}, func(value string) {
		if value == "毫秒" {
			timeNow = 2
		}
		if value == "秒" {
			timeNow = 1
		}
	})
	provinceSelect.Selected = "秒"

	//返回画布 timeStampInp 时间戳 复制按钮 正常时间
	return []fyne.CanvasObject{
		timeStampInp, provinceSelect, click1, text3,
	}, ticker
}

func TimestampToDate() []fyne.CanvasObject {
	timeNow := 1
	//创建输入框
	timeStampInp := widget.NewEntry()
	now := time.Now().Unix()
	//设置值
	timeStampInp.SetText(strconv.Itoa(int(now)))
	//默认值
	timeStampInp.SetPlaceHolder("TIMESTAMP")
	text3 := widget.NewLabel("DATE")
	click1 := widget.NewButton("转换", func() {
		nums, err := strconv.Atoi(timeStampInp.Text)
		if err != nil {
			fmt.Println(err)
			return
		}
		date, err := timeStampToDate(nums, timeNow)
		if err != nil {
			return
		}
		text3.SetText(date)
		//复制到剪切版
		copyClipBoard(date)
		fmt.Printf("%d--->%s\n", nums, date)
	})
	click1.SetIcon(theme.VisibilityIcon())

	provinceSelect := widget.NewSelect([]string{"秒", "毫秒"}, func(value string) {
		if value == "毫秒" {
			timeStampInp.SetText(strconv.Itoa(int(time.Now().UnixNano() / 1e6)))
			timeNow = 2
		}
		if value == "秒" {
			timeStampInp.SetText(strconv.Itoa(int(time.Now().Unix())))
			timeNow = 1
		}
	})
	provinceSelect.Selected = "秒"

	return []fyne.CanvasObject{
		timeStampInp, provinceSelect, click1, text3,
	}
}

func DateToTimestamp() []fyne.CanvasObject {
	timeNow := 1
	timeStampInp := widget.NewEntry()
	now := time.Now().Unix()
	date, _ := timeStampToDate(int(now), timeNow)
	timeStampInp.SetText(date)
	timeStampInp.SetPlaceHolder("DATE")
	text3 := widget.NewLabel("TIMESTAMP")

	click1 := widget.NewButton("转换", func() {
		stamp, err := time.ParseInLocation(timeTemplate1, timeStampInp.Text, time.Local)
		if err != nil {
			fmt.Println(err)
			return
		}
		millionSec := ""
		if timeNow == 1 {
			millionSec = strconv.Itoa(int(stamp.Unix()))
		} else {
			millionSec = strconv.Itoa(int(stamp.UnixNano() / 1e6))
		}
		text3.SetText(millionSec)
		copyClipBoard(millionSec)
		fmt.Printf("%s--->%s\n", stamp, millionSec)
	})
	click1.SetIcon(theme.VisibilityIcon())

	provinceSelect := widget.NewSelect([]string{"秒", "毫秒"}, func(value string) {
		if value == "毫秒" {
			timeNow = 2
		}
		if value == "秒" {
			timeNow = 1
		}
	})
	provinceSelect.Selected = "秒"
	return []fyne.CanvasObject{
		timeStampInp, provinceSelect, click1, text3,
	}
}

func timeStampToDate(t int, types int) (date string, err error) {
	nums, err := strconv.Atoi(strconv.Itoa(t))
	if err != nil {
		return
	}
	if types == 1 {
		date = time.Unix(int64(nums), 0).Format(timeTemplate1)
	} else {
		date = time.Unix(int64(nums/1000), 0).Format(timeTemplate1)
	}
	return
}

func copyClipBoard(context string) {
	clipboard := fyne.CurrentApp().Driver().AllWindows()[0].Clipboard()
	clipboard.SetContent(context)
	fmt.Println("success clipboard", context)
}
