package route

import (
	"os"
	"path/filepath"
	"strings"

	"opennamu/route/tool"
)

func Api_image_exist(name string) map[string]any {
	if name == "" || filepath.Base(name) != name || strings.Contains(name, "..") {
		return map[string]any{}
	}

	db := tool.DB_connect()
	image_path := filepath.Join(tool.Get_image_url(db), name)
	tool.DB_close(db)
	if _, err := os.Stat(image_path); err != nil {
		return map[string]any{}
	}

	return map[string]any{"exist": "1"}
}
