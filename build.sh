# 编译环境 Linux x64
# 环境变量请自行替换

mkdir SimpleOIClass_win10_x86/
mkdir SimpleOIClass_win10_x86/send/
mkdir SimpleOIClass_win10_x86/upld/
mkdir SimpleOIClass_win10_x86/idmap/
set GOROOT=$HOME/goroot/go1.21.1.linux-amd64/go
CGO_ENABLED=0 GOOS=windows GOARCH=386 $GOROOT/bin/go build -o SimpleOIClass_win10_x86/SimpleOIClass.exe main.go

mkdir SimpleOIClass_win7_x86/
mkdir SimpleOIClass_win7_x86/send/
mkdir SimpleOIClass_win7_x86/upld/
mkdir SimpleOIClass_win7_x86/idmap/
set GOROOT=$HOME/goroot/go1.20.8.linux-amd64/go
CGO_ENABLED=0 GOOS=windows GOARCH=386 $GOROOT/bin/go build -o SimpleOIClass_win7_x86/SimpleOIClass.exe main.go

mkdir SimpleOIClass_oldwin_x86/
mkdir SimpleOIClass_oldwin_x86/send/
mkdir SimpleOIClass_oldwin_x86/upld/
mkdir SimpleOIClass_oldwin_x86/idmap/
set GOROOT=$HOME/goroot/go1.10.8.linux-amd64/go
CGO_ENABLED=0 GOOS=windows GOARCH=386 $GOROOT/bin/go build -o SimpleOIClass_oldwin_x86/SimpleOIClass.exe main.go


mkdir SimpleOIClass_win10_x64/
mkdir SimpleOIClass_win10_x64/send/
mkdir SimpleOIClass_win10_x64/upld/
mkdir SimpleOIClass_win10_x64/idmap/
set GOROOT=$HOME/goroot/go1.21.1.linux-amd64/go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $GOROOT/bin/go build -o SimpleOIClass_win10_x64/SimpleOIClass.exe main.go

mkdir SimpleOIClass_win7_x64/
mkdir SimpleOIClass_win7_x64/send/
mkdir SimpleOIClass_win7_x64/upld/
mkdir SimpleOIClass_win7_x64/idmap/
set GOROOT=$HOME/goroot/go1.20.8.linux-amd64/go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $GOROOT/bin/go build -o SimpleOIClass_win7_x64/SimpleOIClass.exe main.go

mkdir SimpleOIClass_oldwin_x64/
mkdir SimpleOIClass_oldwin_x64/send/
mkdir SimpleOIClass_oldwin_x64/upld/
mkdir SimpleOIClass_oldwin_x64/idmap/
set GOROOT=$HOME/goroot/go1.10.8.linux-amd64/go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 $GOROOT/bin/go build -o SimpleOIClass_oldwin_x64/SimpleOIClass.exe main.go
