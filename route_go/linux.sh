#!/bin/bash

read -p "file_name: " file_name
to="${1:-all}"

build() {
    local os=$1
    local arch=$2
    local ext=$3

    echo "$os $arch"
    
    GOOS=$os
    GOARCH=$arch
    
    if [[ "$arch" == "darwin" ]]; then
        local arch="mac.$arch"
    fi
    
    CGO_ENABLED=0
    go build -o "bin/$file_name.$arch$ext" "$file_name.go"
}

if [[ "$to" == "linux_amd64" || "$to" == "all" ]]; then
    build "linux" "amd64" ".bin"
fi

if [[ "$to" == "linux_arm64" || "$to" == "all" ]]; then
    build "linux" "arm64" ".bin"
fi

if [[ "$to" == "windows_amd64" || "$to" == "all" ]]; then
    build "windows" "amd64" ".exe"
fi

if [[ "$to" == "windows_arm64" || "$to" == "all" ]]; then
    build "windows" "arm64" ".exe"
fi

if [[ "$to" == "mac_arm64" || "$to" == "all" ]]; then
    build "darwin" "arm64" ".bin"
fi