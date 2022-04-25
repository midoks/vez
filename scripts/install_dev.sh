#!/bin/bash
PATH=/bin:/sbin:/usr/bin:/usr/sbin:/usr/local/bin:/usr/local/sbin:~/bin

APP_NAME=vez
TAGRT_DIR=/usr/local/${APP_NAME}_dev

mkdir -p $TAGRT_DIR
cd $TAGRT_DIR

export GIT_COMMIT=$(git rev-parse HEAD)
export BUILD_TIME=$(date -u '+%Y-%m-%d %I:%M:%S %Z')


if [ ! -d $TAGRT_DIR/${APP_NAME} ]; then
	git clone https://github.com/midoks/${APP_NAME}
	cd $TAGRT_DIR/${APP_NAME}
else
	cd $TAGRT_DIR/${APP_NAME}
	git pull https://github.com/midoks/${APP_NAME}
fi

go mod tidy
go mod vendor


rm -rf ${APP_NAME}
go build  -ldflags "-X \"github.com/midoks/${APP_NAME}/internal/conf.BuildTime=${BUILD_TIME}\" -X \"github.com/midoks/${APP_NAME}/internal/conf.BuildCommit=${GIT_COMMIT}\"" ./


cd $TAGRT_DIR/${APP_NAME}/scripts

sh make.sh

systemctl daemon-reload

service ${APP_NAME} restart

cd $TAGRT_DIR/${APP_NAME} && ./${APP_NAME} -v


