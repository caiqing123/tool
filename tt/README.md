# tt
golang fyne开发得一个时间戳日期转换工具

打包
go build  -ldflags="-s -w -H windowsgui" main.go

对应桌面端的应用，在确定好目标平台后就可以直接打包：

fyne package -os darwin -icon myapp.png
fyne package -os linux -icon myapp.png
fyne package -os windows -icon myapp.png
上述命令分别对应 macOS、Linux 和 Windows平台的构建，而 myapp.png 是应用的图标文件。对于 macOS 平台，生成 app 应用文件；对于 Linux 平台，生成一个 tag.gz 文件，解包后可以放到 use/local/ 中使用；对于 Windows 平台，直接生成 exe 文件，可以直接运行。

而对于移动端平台，同样十分简单。在配置好了相应的环境后，运行

fyne package -os android -appID com.example.myapp -icon mobileIcon.png
fyne package -os ios - appID com.example.myapp -icon mobileIcon.png