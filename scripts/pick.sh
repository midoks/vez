#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin


VERSION=0.0.2
curPath=`pwd`
rootPath=$(dirname "$curPath")

PACK_NAME=vez

# go tool dist list
mkdir -p $rootPath/tmp/build
mkdir -p $rootPath/tmp/package

source ~/.bash_profile

cd $rootPath
LDFLAGS="-X \"github.com/midoks/vez/internal/conf.BuildTime=$(date -u '+%Y-%m-%d %I:%M:%S %Z')\""
LDFLAGS="${LDFLAGS} -X \"github.com/midoks/vez/internal/conf.BuildCommit=$(git rev-parse HEAD)\""


echo $LDFLAGS
build_app(){

	if [ -f $rootPath/tmp/build/vez ]; then
		rm -rf $rootPath/tmp/build/vez
		rm -rf $rootPath/vez
	fi

	if [ -f $rootPath/tmp/build/vez.exe ]; then
		rm -rf $rootPath/tmp/build/vez.exe
		rm -rf $rootPath/vez.exe
	fi

	echo "build_app" $1 $2

	echo "export CGO_ENABLED=0 GOOS=$1 GOARCH=$2"
	echo "cd $rootPath && go build vez.go"

	# export CGO_ENABLED=1 GOOS=linux GOARCH=amd64

	if [ $1 != "darwin" ]; then
		export CGO_ENABLED=0 GOOS=$1 GOARCH=$2
		export CGO_LDFLAGS="-static"
	fi

	cd $rootPath && go generate internal/assets/conf/conf.go
	cd $rootPath && go generate internal/assets/templates/templates.go
	cd $rootPath && go generate internal/assets/public/public.go



	if [ $1 == "windows" ]; then
		
		if [ $2 == "amd64" ]; then
			export CC=x86_64-w64-mingw32-gcc
			export CXX=x86_64-w64-mingw32-g++
		else
			export CC=i686-w64-mingw32-gcc
			export CXX=i686-w64-mingw32-g++
		fi

		cd $rootPath && go build -o vez.exe -ldflags "${LDFLAGS}" vez.go

		# -ldflags="-s -w"
		# cd $rootPath && go build vez.go && /usr/local/bin/strip vez
	fi

	if [ $1 == "linux" ]; then
		export CC=x86_64-linux-musl-gcc
		if [ $2 == "amd64" ]; then
			export CC=x86_64-linux-musl-gcc

		fi

		if [ $2 == "386" ]; then
			export CC=i486-linux-musl-gcc
		fi

		if [ $2 == "arm64" ]; then
			export CC=aarch64-linux-musl-gcc
		fi

		if [ $2 == "arm" ]; then
			export CC=arm-linux-musleabi-gcc
		fi

		cd $rootPath && go build -ldflags "${LDFLAGS}"  vez.go 
	fi

	if [ $1 == "darwin" ]; then
		echo "cd $rootPath && go build -v -ldflags '${LDFLAGS}'"
		cd $rootPath && go build -v -ldflags "${LDFLAGS}"
		
		cp $rootPath/vez $rootPath/tmp/build
	fi
	

	cp -r $rootPath/scripts $rootPath/tmp/build
	cp -r $rootPath/LICENSE $rootPath/tmp/build
	cp -r $rootPath/README.md $rootPath/tmp/build

	cd $rootPath/tmp/build && xattr -c * && rm -rf ./*/.DS_Store && rm -rf ./*/*/.DS_Store


	if [ $1 == "windows" ];then
		cp $rootPath/vez.exe $rootPath/tmp/build
	else
		cp $rootPath/vez $rootPath/tmp/build
	fi

	cd $rootPath/tmp/build && tar -zcvf ${PACK_NAME}_${VERSION}_$1_$2.tar.gz ./ && mv ${PACK_NAME}_${VERSION}_$1_$2.tar.gz $rootPath/tmp/package

}

golist=`go tool dist list`
echo $golist

build_app linux amd64
# build_app linux 386
# build_app linux arm64
# build_app linux arm
build_app darwin amd64
