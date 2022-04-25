#!/bin/bash

_os=`uname`
_path=`pwd`
_dir=`dirname $_path`

APP_NAME=vez

sed "s:{APP_PATH}:${_dir}:g" $_dir/scripts/init.d/${APP_NAME}.tpl > $_dir/scripts/init.d/${APP_NAME}
chmod +x $_dir/scripts/init.d/${APP_NAME}


if [ -d /etc/init.d ];then
	cp $_dir/scripts/init.d/${APP_NAME} /etc/init.d/${APP_NAME}
	chmod +x /etc/init.d/${APP_NAME}
fi

echo `dirname $_path`