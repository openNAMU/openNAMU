echo "file_name : "
read file_name

export GOOS=linux
export GOARCH=amd64
CGO_ENABLED=0 go build $file_name.go
mv $file_name ./bin/$file_name.amd64.bin

export GOOS=linux
export GOARCH=arm64
CGO_ENABLED=0 go build $file_name.go
mv $file_name ./bin/$file_name.arm64.bin

export GOOS=windows
export GOARCH=amd64
CGO_ENABLED=0 go build $file_name.go
mv $file_name.exe ./bin/$file_name.amd64.exe

export GOOS=windows
export GOARCH=arm64
CGO_ENABLED=0 go build $file_name.go
mv $file_name.exe ./bin/$file_name.arm64.exe