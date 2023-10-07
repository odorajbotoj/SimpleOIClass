# 编译环境 Linux x64
# 环境变量请自行替换
echo ==== win10 x86 ====
export GOROOT=$HOME/goroot/go1.21.1.linux-amd64/go && CGO_ENABLED=0 GOOS=windows GOARCH=386 GOROOT=$HOME/goroot/go1.21.1.linux-amd64/go $GOROOT/bin/go build -o SimpleOIClass_win10_x86.exe main.go

echo ==== win7 x86 ====
export GOROOT=$HOME/goroot/go1.20.8.linux-amd64/go && CGO_ENABLED=0 GOOS=windows GOARCH=386 GOROOT=$HOME/goroot/go1.20.8.linux-amd64/go $GOROOT/bin/go build -o SimpleOIClass_win7_x86.exe main.go

echo ==== oldwin x86 ====
export GOROOT=$HOME/goroot/go1.10.8.linux-amd64/go && CGO_ENABLED=0 GOOS=windows GOARCH=386 GOROOT=$HOME/goroot/go1.10.8.linux-amd64/go $GOROOT/bin/go build -o SimpleOIClass_oldwin_x86.exe main.go


echo ==== win10 x64 ====
export GOROOT=$HOME/goroot/go1.21.1.linux-amd64/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOROOT=$HOME/goroot/go1.21.1.linux-amd64/go $GOROOT/bin/go build -o SimpleOIClass_win10_x64.exe main.go

echo ==== win7 x64 ====
export GOROOT=$HOME/goroot/go1.20.8.linux-amd64/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOROOT=$HOME/goroot/go1.20.8.linux-amd64/go $GOROOT/bin/go build -o SimpleOIClass_win7_x64.exe main.go

echo ==== oldwin x64 ====
export GOROOT=$HOME/goroot/go1.10.8.linux-amd64/go && CGO_ENABLED=0 GOOS=windows GOARCH=amd64 GOROOT=$HOME/goroot/go1.10.8.linux-amd64/go $GOROOT/bin/go build -o SimpleOIClass_oldwin_x64.exe main.go
