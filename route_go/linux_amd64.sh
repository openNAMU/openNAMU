echo "file_name : "
read file_name

echo "linux amd64"
export GOOS=linux
export GOARCH=amd64
CGO_ENABLED=0 go build $file_name.go
mv $file_name ./bin/$file_name.amd64.bin

echo "linux arm64"
export GOOS=linux
export GOARCH=arm64
CGO_ENABLED=0 go build $file_name.go
mv $file_name ./bin/$file_name.arm64.bin

echo "windows amd64"
export GOOS=windows
export GOARCH=amd64
CGO_ENABLED=0 go build $file_name.go
mv $file_name.exe ./bin/$file_name.amd64.exe

echo "windows arm64"
export GOOS=windows
export GOARCH=arm64
CGO_ENABLED=0 go build $file_name.go
mv $file_name.exe ./bin/$file_name.arm64.exe