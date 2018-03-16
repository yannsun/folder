#!/bin/bash
cd `dirname $0`/.. || exit 1;

. build/env.sh
build_date=$(date +"%Y-%m-%d %H:%M:%S")
echo "start to install project..." $build_date
echo $GOPATH
${gobin} version
pwd

# 单独编译安装,若为调试开发模式,建议开启 并发检测 -race 函数

ldflags=()
#value=`svnversion`
value=$1
ldflags+=("-X" "\"main.SvnVersion=r${value}\"")
ldflags+=("-X" "\"main.BuildDate=${build_date}\"")
#去掉二进制段信息中的gdb调试信息,也就是strip的作用
ldflags+=("-s")
ldflags+=("-w")

cd ${projectpath}/src/gourd/launcher/swallow/

#GOOS=darwin  GOARCH=amd64 ${gobin} build -v -ldflags "${ldflags[*]}" -o ${projectpath}/bin/swallow || exit 1
#echo "build darwin amd64 done."

GOBIN=${GOBIN} GOOS=linux  GOARCH=amd64 ${gobin} install || exit 1
echo "install [swallow] [linux amd64] done."

cd ${projectpath}/src/gourd/launcher/gourd/

#GOOS=darwin  GOARCH=amd64 ${gobin} build -v -ldflags "${ldflags[*]}" -o ${projectpath}/bin/gourd_darwin_amd64 || exit 1
#echo "build darwin amd64 done."

#GOOS=darwin  GOARCH=386 ${gobin} build -v -ldflags "${ldflags[*]}" -o ${projectpath}/bin/w_gourd_darwin_386_r${gitver} || exit 1
#echo "build darwin 386 done."
#
GOBIN=${GOBIN} GOOS=linux  GOARCH=amd64 ${gobin} install || exit 1
echo "install [gourd] [linux amd64] done."

#GOBIN=${GOBIN} GOOS=linux  GOARCH=386 ${gobin} build -v -ldflags "${ldflags[*]}" -o ${projectpath}/bin/w_gourd_linux_386_r${gitver} || exit 1
#echo "build linux 386 done."
#cd ${projectpath}

echo "intsall project done..." $(date +"%Y-%m-%d %H:%M:%S")
exit 0
