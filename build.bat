go build -ldflags="-s -w"
upx -9 -vf --lzma --compress-icons=0 *.exe