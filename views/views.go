package views

import (
	"embed"
	"io/fs"
)

//go:embed main_css ringo
var files embed.FS

var BuiltinSkinList = []string{
	"main_css",
	"ringo",
}

func Read(name string) ([]byte, error) {
	return files.ReadFile(name)
}

func ReadDir(name string) ([]fs.DirEntry, error) {
	return files.ReadDir(name)
}
