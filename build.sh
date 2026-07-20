#!/bin/bash
set -e

cd "$(dirname "$0")"

echo "=== 编译前端 ==="
cd frontend && npm run build && cd ..

echo "=== 编译 Windows exe ==="
export PATH="$HOME/go/bin:/usr/lib/go-1.26/bin:$PATH"
export GOOS=windows GOARCH=amd64
export CGO_ENABLED=1
export CC=x86_64-w64-mingw32-gcc
export CXX=x86_64-w64-mingw32-g++

wails build -tags webkit2_41 -o OptiWin.exe -ldflags="-s -w" -trimpath

echo "=== 复制到桌面 ==="
cp -f build/bin/OptiWin.exe /mnt/c/Users/Administrator/Desktop/OptiWin.exe

ls -lh build/bin/OptiWin.exe
echo "OK"
