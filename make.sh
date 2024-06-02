#!/bin/bash

# 确认参数是否存在
if [ ! -n "$1" ] || [ ! -n "$2" ]; then
    echo "参数错误: "
    echo "示例: ./make.sh -u/-w/-a/-o [读卡器选项]"
    echo "示例: ./make.sh -u 0:none/1:YMC60/2:CX522"
    exit 1
fi

export GOPATH=/Users/gao/Develop/golang/

# 设置默认读卡器选项
READER_TAG=""

# 根据第二个参数设置读卡器选项
case $2 in
    "0")
        READER_TAG="-tags 老崔定制读卡器"
        ;;
    "1")
        READER_TAG="-tags CX522读卡器"
        ;;
    "2")
        READER_TAG="-tags YMC60读卡器"
        ;;
    *)
        echo "不支持的读卡器选项: $2"
        echo "支持的选项: 0:none, 1:YMC60，2:CX522"
        exit 1
        ;;
esac

# 根据第一个参数编译对应系统的目标
case $1 in
    "-u")
        echo "正在编译 Ubuntu 版本..." $READER_TAG
        CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build $READER_TAG -o ecds main.go
        ;;
    "-w")
        echo "正在编译 Windows 版本..." $READER_TAG
        CGO_ENABLED=0 GOOS=windows GOARCH=386 go build $READER_TAG -o ecds.exe main.go
        ;;
    "-a")
        echo "正在编译 ARM 版本..." $READER_TAG
        CGO_ENABLED=0 GOOS=linux GOARCH=arm go build $READER_TAG -o ecds main.go
        ;;
    "-r")
        echo "正在编译 树莓派 版本..." $READER_TAG
        CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build $READER_TAG -o ecds main.go
        ;;
    "-o")
        echo "正在编译 OSX 版本..." $READER_TAG
        CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build $READER_TAG -o ecds main.go
        ;;
    *)
        echo "不支持的系统选项: $1"
        echo "支持的选项: -u (Ubuntu), -w (Windows), -a (ARM), -o (OSX)"
        exit 1
        ;;
esac

# 编译成功提示
echo "编译完成"
