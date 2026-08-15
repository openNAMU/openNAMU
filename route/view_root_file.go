package route

import (
	"os"
	"path/filepath"
	"strings"
)

func Read_root_file(name string) ([]byte, string, bool) {
	name = strings.TrimPrefix(name, "/")
	if name == "" || filepath.Base(name) != name {
		return nil, "", false
	}

	dot := strings.LastIndex(name, ".")
	if dot <= 0 || strings.Contains(name[:dot], ".") {
		return nil, "", false
	}
	extension := strings.ToLower(name[dot+1:])
	mime_type := ""
	switch extension {
	case "txt":
		mime_type = "text/plain; charset=utf-8"
	case "xml":
		mime_type = "application/xml; charset=utf-8"
	case "ico":
		mime_type = "image/x-icon"
	default:
		return nil, "", false
	}

	data, err := os.ReadFile(name)
	if err != nil {
		if os.IsNotExist(err) {
			return []byte{}, mime_type, true
		}
		return nil, "", false
	}
	return data, mime_type, true
}
