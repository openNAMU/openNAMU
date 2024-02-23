$file_name = Read-Host "file_name : "

$env:GOOS = "linux"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = 0
go build $file_name.go
Remove-Item ".\bin\$file_name.amd64.bin"
Move-Item "opennamu" ".\bin\$file_name.amd64.bin"

$env:GOOS = "linux"
$env:GOARCH = "arm64"
$env:CGO_ENABLED = 0
go build $file_name.go
Remove-Item ".\bin\$file_name.arm64.bin"
Move-Item "opennamu" ".\bin\$file_name.arm64.bin"

$env:GOOS = "windows"
$env:GOARCH = "amd64"
$env:CGO_ENABLED = 0
go build $file_name.go
Remove-Item ".\bin\$file_name.amd64.exe"
Move-Item "opennamu.exe" ".\bin\$file_name.amd64.exe"

$env:GOOS = "windows"
$env:GOARCH = "arm64"
$env:CGO_ENABLED = 0
go build $file_name.go
Remove-Item ".\bin\$file_name.arm64.exe"
Move-Item "opennamu.exe" ".\bin\$file_name.arm64.exe"