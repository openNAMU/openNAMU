package lang

import "embed"

//go:embed *.json
var fs embed.FS

func Read(name string) ([]byte, error) {
	return fs.ReadFile(name + ".json")
}
