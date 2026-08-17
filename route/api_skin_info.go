package route

import (
	"opennamu/route/tool"
	"path/filepath"
	"strings"
)

func Api_skin_info(config tool.Config, name string) ([]byte, bool) {
	db := tool.DB_connect()
	if name == "" {
		name = tool.Get_use_skin_name_session(db, config.IP, config.Session)
	}
	tool.DB_close(db)

	if name == "" || filepath.Base(name) != name || strings.Contains(name, "..") {
		return nil, false
	}

	data, err := tool.Read_view_file(name + "/info.json")
	if err != nil {
		return nil, false
	}

	return data, true
}
