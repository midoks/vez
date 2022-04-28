#!/bin/bash

check_go_environment() {
	if test ! -x "$(command -v go)"; then
		printf "\e[1;31mmissing go running environment\e[0m\n"
		exit 1
	fi
}

load_vars() {
	OS=$(uname | tr '[:upper:]' '[:lower:]')

	VERSION=$(get_latest_release "midoks/webp_server_go")

	TARGET_DIR="/opt/webps"
}

get_latest_release() {
    curl -sL "https://api.github.com/repos/$1/releases/latest" | grep '"tag_name":' | cut -d'"' -f4
}

get_arch() {
	echo "package main
import (
	\"fmt\"
	\"runtime\"
)
func main() { fmt.Println(runtime.GOARCH) }" > /tmp/go_arch.go

	ARCH=$(go run /tmp/go_arch.go)
}

get_download_url() {
	DOWNLOAD_URL="https://github.com/midoks/webp_server_go/releases/download/$VERSION/webp-server-${OS}-${ARCH}"
}

# download file
download_file() {
    url="${1}"
    destination="${2}"

    printf "Fetching ${url} \n\n"

    if test -x "$(command -v curl)"; then
        code=$(curl --connect-timeout 15 -w '%{http_code}' -L "${url}" -o "${destination}")
    elif test -x "$(command -v wget)"; then
        code=$(wget -t2 -T15 -O "${destination}" --server-response "${url}" 2>&1 | awk '/^  HTTP/{print $2}' | tail -1)
    else
        printf "\e[1;31mNeither curl nor wget was available to perform http requests.\e[0m\n"
        exit 1
    fi

    if [ "${code}" != 200 ]; then
        printf "\e[1;31mRequest failed with code %s\e[0m\n" $code
        exit 1
    else 
	    printf "\n\e[1;33mDownload succeeded\e[0m\n"
    fi
}


main() {
	check_go_environment

	load_vars

	get_arch

	get_download_url

	DOWNLOAD_FILE="$(mktemp)"
	download_file $DOWNLOAD_URL $DOWNLOAD_FILE

	if [ ! -d "$TARGET_DIR" ]; then
		mkdir -p "$TARGET_DIR"
	fi

	cp -rf $DOWNLOAD_FILE $TARGET_DIR/webp-server
	chmod 755 $TARGET_DIR/webp-server

	if [ ! -f "/usr/lib/systemd/system/webps.service" ];then
		wget  -t2 -T15 -O "/usr/lib/systemd/system/webps.service" https://raw.githubusercontent.com/midoks/webp_server_go/master/scripts/webps.service
	fi


	if [ ! -f "$TARGET_DIR/config.json" ];then
		wget  -t2 -T15 -O "$TARGET_DIR/config.json" https://raw.githubusercontent.com/midoks/webp_server_go/master/config.json
	fi

	

	systemctl daemon-reload
	service webps restart

	$TARGET_DIR/webp-server -V
}

main "$@" || exit 1
