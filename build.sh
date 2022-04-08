#!/bin/bash

# var
WORKSPACE=$(cd `dirname $0` && pwd -P)
DATE=$(date "+%s")
GOPATH=/tmp/go-build${DATE}

# 项目名称，根据具体项目改动
PACKAGE_PATH=MyServer
MODULE_NAME=MyServer
MODULE_PATH=${GOPATH}/src/MyServer
VERSION=1.13.5

# env
export GOPATH
#export GOROOT=/usr/local/go 请在个人bash文件(~/.bash_profile) 中设置该变量
export PATH=${GOROOT}/bin:$GOPATH/bin:${PATH}:$GOBIN
#export GOPROXY=http://goproxy.intra.xiaojukeji.com,direct
export GOSUMDB=off
export GO111MODULE=on

if [ ! -d $GOROOT  ];then
  echo "EEROR !!! GO VERSION should more than 1.13.5, default is 1.13.5, please modify the VERSION in build.sh for a suitable go version"
  exit 1
fi

function build() {
    rm -rf $MODULE_PATH/$MODULE_NAME &> /dev/null
    #mkdir -p /tmp/xiaoju/ep/as/store/toggles
    mkdir -p $GOPATH/bin
    mkdir -p $MODULE_PATH
    ln -sf $WORKSPACE ${MODULE_PATH}/${MODULE_NAME}
    cd $MODULE_PATH/$MODULE_NAME

    #clean mod cache
    go clean -modcache

    # build
    echo "Building……" && make build

    if [[ $? != 0 ]];then
        echo -e "Build failed !"
        exit 1
    fi
    echo -e "Build success!"
}

build
