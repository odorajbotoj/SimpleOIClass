# Simple OI Class

+ 一个极简的机房用文件下发/上传系统

---

## ver 4.0.0 更新日志

+ 现在可以通过REG模式让学生自行登记姓名了

+ 现已兼容 `UTF-8 BOM` ，Windows记事本也能用！

## 安装

+ 我预编译了在Windows Server2003/7/10上使用的32/64位可执行文件

+ 绿色软件，下载即可运行

+ 同级目录下需要存在 `send` `upld` `idmap`三个文件夹（或自定义。后文同）

## 使用

+ 教师机开启服务之后，同一网络下学生可以访问 `http://<教师机ip>:<设定的端口号>` 。

+ 当服务器处于REG模式时，学生可以提交自己的姓名完成注册。

+ 所有学生都可以下载您下发的文件（只需要放到 `send` 文件夹下）。

+ 当一台新的电脑连接的时候，会根据他的IP在 `idmap` 文件夹下创建一个文件。

+ 您需要在文件内部填写学生姓名（请使用 `UTF-8` 编码）。

+ 学生可以上传自己的文件，这些文件将保存至 `upld/<学生姓名>` 文件夹下。

+ 您可以在 `upld` 文件夹下选中学生，并打包测评。

## 配置

+ 您可在主程序同级目录下新建 `config.txt` 来进行配置

+ 配置文件为K=V格式 ( `Key=Value` )

+ 配置文件只需要按需填写，程序会自动识别。不填则使用默认配置。

+ `PORT` 为服务端口号，默认为 `:8080` （请勿漏写英文冒号！）

+ `IDMAP` 为学生-IP映射表存放目录，默认为 `idmap/`（请勿漏写英文正斜杠！） 

+ `SEND` 为下发文件存放目录，默认为 `send/`（请勿漏写英文正斜杠！） 

+ `UPLD` 为学生上传文件根目录，默认为 `upld/`（请勿漏写英文正斜杠！） 

+ `ACCEPT` 为允许上传的文件类型（后缀），默认为 `.cpp` ，书写遵循HTML input标签accept属性，这个字符串是一个以逗号为分隔的唯一文件类型说明符列表。[参考](https://developer.mozilla.org/zh-CN/docs/Web/HTML/Element/Input/file#accept)

+ `REG` 为REG模式的开关。默认为 `OFF` 。当值为 `ON` 时，服务器处于REG模式。此模式下学生仅能完成姓名注册，无法使用其他功能。

+ 所有内容仅在服务器启动时读取一次，若服务器运行中作了修改，请重启服务器。

## 本程序遵循GPL-v3协议

odorajbotoj @ ZJYZITLAB

2023.10.08
