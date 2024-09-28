$file_name = Read-Host "file_name "

$to = $args[0]
if($to -eq $null) {
    $to = "all"
}

if($to -eq "linux_amd64" -or $to -eq "all") {
    Write-Host "linux amd64"
    $env:GOOS = "linux"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = 0
    go build $file_name.go
    Remove-Item ".\bin\$file_name.amd64.bin"
    Move-Item "opennamu" ".\bin\$file_name.amd64.bin"
}

if($to -eq "linux_arm64" -or $to -eq "all") {
    Write-Host "linux arm64"
    $env:GOOS = "linux"
    $env:GOARCH = "arm64"
    $env:CGO_ENABLED = 0
    go build $file_name.go
    Remove-Item ".\bin\$file_name.arm64.bin"
    Move-Item "opennamu" ".\bin\$file_name.arm64.bin"
}

if($to -eq "windows_amd64" -or $to -eq "all") {
    Write-Host "windows amd64"
    $env:GOOS = "windows"
    $env:GOARCH = "amd64"
    $env:CGO_ENABLED = 0
    go build $file_name.go
    Remove-Item ".\bin\$file_name.amd64.exe"
    Move-Item "opennamu.exe" ".\bin\$file_name.amd64.exe"
}

if($to -eq "windows_arm64" -or $to -eq "all") {
    Write-Host "windows arm64"
    $env:GOOS = "windows"
    $env:GOARCH = "arm64"
    $env:CGO_ENABLED = 0
    go build $file_name.go
    Remove-Item ".\bin\$file_name.arm64.exe"
    Move-Item "opennamu.exe" ".\bin\$file_name.arm64.exe"
}

if($to -eq "mac_arm64" -or $to -eq "all") {
    Write-Host "mac arm64"
    $env:GOOS = "darwin"
    $env:GOARCH = "arm64"
    $env:CGO_ENABLED = 0
    go build $file_name.go
    Remove-Item ".\bin\$file_name.mac.arm64.bin"
    Move-Item "opennamu" ".\bin\$file_name.mac.arm64.bin"
}