package main

import (
	"log"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"

	"github.com/Knetic/govaluate"
)

type calc struct {
	equation string //输出的字符串
	last     string //最后输出的字符串

	output  *widget.Label             // 输出框
	buttons map[string]*widget.Button // 按钮合集
	window  fyne.Window               // fyne对象
}

//设置输出值
func (c *calc) display(newtext string) {
	c.equation = newtext
	c.output.SetText(newtext)
}

//显示对应按钮值
func (c *calc) character(char rune, t int) {
	if t != 2 || string(char) != c.last {
		c.display(c.equation + string(char))
		c.last = string(char)
	}
}

//数字变为字符串显示对应按钮值
func (c *calc) digit(d int) {
	c.character(rune(d)+'0', 1)
}

//清除结果
func (c *calc) clear() {
	c.display("")
}

// 计算结果
func (c *calc) evaluate() {
	if c.output.Text == "error" {
		c.output.Text = ""
	}
	expression, err := govaluate.NewEvaluableExpression(c.output.Text)
	if err == nil {
		result, err := expression.Evaluate(nil)
		if err == nil {
			c.display(strconv.FormatFloat(result.(float64), 'f', -1, 64))
		}
	}

	if err != nil {
		log.Println("Error in calculation", err)
		c.display("error")
	}
}

//添加按钮和对应触发函数
func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button
	return button
}

//添加数字按钮和对应触发函数
func (c *calc) digitButton(number int) *widget.Button {
	str := strconv.Itoa(number)
	return c.addButton(str, func() {
		c.digit(number)
	})
}

//添加字符串按钮和对应触发函数
func (c *calc) charButton(char rune) *widget.Button {
	return c.addButton(string(char), func() {
		c.character(char, 2)
	})
}

//触发事件 根据键盘值触发
func (c *calc) onTypedRune(r rune) {
	if r == 'c' {
		r = 'C' // The button is using a capital C.
	}
	//在设置的按钮合集里触发按钮事件
	if button, ok := c.buttons[string(r)]; ok {
		button.OnTapped()
	}
}

//触发事件 根据键盘值触发
func (c *calc) onTypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter { //按键值为确认或回车则计算结果
		c.evaluate()
	} else if ev.Name == fyne.KeyBackspace && len(c.equation) > 0 { //按键值为删除且输入值大于0则删除一个
		c.display(c.equation[:len(c.equation)-1])
	}
}

//粘贴到输入框
func (c *calc) onPasteShortcut(shortcut fyne.Shortcut) {
	content := shortcut.(*fyne.ShortcutPaste).Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err != nil {
		return
	}

	c.display(c.equation + content)
}

//复制输入框内容
func (c *calc) onCopyShortcut(shortcut fyne.Shortcut) {
	shortcut.(*fyne.ShortcutCopy).Clipboard.SetContent(c.equation)
}

func (c *calc) loadUI(app fyne.App) {
	//构建输出框
	c.output = &widget.Label{Alignment: fyne.TextAlignTrailing}
	c.output.TextStyle.Monospace = true

	//构建等于按钮
	equals := c.addButton("=", c.evaluate)
	//突出显示
	equals.Importance = widget.HighImportance

	//新建窗口
	c.window = app.NewWindow("计算器")
	//设置内容 9 ** 2  n^{2}
	c.window.SetContent(container.NewGridWithColumns(1,
		c.output,
		container.NewGridWithColumns(4,
			c.addButton("C", c.clear),
			c.charButton('('),
			c.charButton(')'),
			c.charButton('/')),
		container.NewGridWithColumns(4,
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.charButton('*')),
		container.NewGridWithColumns(4,
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.charButton('-')),
		container.NewGridWithColumns(4,
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.charButton('+')),
		container.NewGridWithColumns(2,
			container.NewGridWithColumns(2,
				c.digitButton(0),
				c.charButton('.')),
			equals)),
	)

	//键盘输入事件
	c.window.Canvas().SetOnTypedRune(c.onTypedRune)
	//键盘确认删除事件
	c.window.Canvas().SetOnTypedKey(c.onTypedKey)
	//添加快捷操作复制 ctrl+c
	c.window.Canvas().AddShortcut(&fyne.ShortcutCopy{}, c.onCopyShortcut)
	//添加快捷操作粘贴 ctrl+v
	c.window.Canvas().AddShortcut(&fyne.ShortcutPaste{}, c.onPasteShortcut)
	//设置尺寸
	c.window.Resize(fyne.NewSize(200, 300))
	//运行
	c.window.Show()
}

func newCalculator() *calc {
	return &calc{
		buttons: make(map[string]*widget.Button, 19),
	}
}
