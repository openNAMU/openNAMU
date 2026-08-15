package route

import (
	stdjson "encoding/json"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"opennamu/route/tool"
)

var skin_info_link = map[string]string{
	"ACME":        "https://raw.githubusercontent.com/openNAMU/openNAMU-Skin-ACME/master/info.json",
	"Liberty":     "https://raw.githubusercontent.com/openNAMU/openNAMU-Skin-Liberty/master/info.json",
	"Before Namu": "https://raw.githubusercontent.com/openNAMU/openNAMU-Skin-Before_Namu/master/info.json",
}

func get_skin_latest_version(info map[string]any) (string, bool) {
	info_link, _ := info["info_link"].(string)
	if info_link == "" {
		name, _ := info["name"].(string)
		info_link = skin_info_link[name]
	}
	if info_link == "" {
		return "", false
	}

	parsed, err := url.Parse(info_link)
	if err != nil || parsed.Host == "" || (parsed.Scheme != "http" && parsed.Scheme != "https") {
		return "", false
	}

	client := http.Client{Timeout: 3 * time.Second}
	response, err := client.Get(info_link)
	if err != nil {
		return "", false
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return "", false
	}

	remote_data := map[string]any{}
	decoder := stdjson.NewDecoder(io.LimitReader(response.Body, 1<<20))
	if err := decoder.Decode(&remote_data); err != nil {
		return "", false
	}
	version, ok := remote_data["skin_ver"].(string)
	return version, ok && version != ""
}

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

func Api_skin_info_all(config tool.Config) ([]byte, bool) {
	db := tool.DB_connect()
	current := tool.Get_use_skin_name_session(db, config.IP, config.Session)
	skin_list := tool.Get_skin_list(current, true)
	tool.DB_close(db)

	data := map[string]map[string]any{}
	for _, name := range skin_list {
		if name == "" || filepath.Base(name) != name || strings.Contains(name, "..") {
			continue
		}

		raw, err := tool.Read_view_file(name + "/info.json")
		if err != nil {
			continue
		}
		info := map[string]any{}
		if err := stdjson.Unmarshal(raw, &info); err != nil {
			continue
		}
		if name == current {
			info["main"] = "true"
		}
		if version, ok := get_skin_latest_version(info); ok {
			info["lastest_version"] = map[string]string{"skin_ver": version}
		}
		data[name] = info
	}

	result, err := stdjson.Marshal(data)
	if err != nil {
		return nil, false
	}
	return result, true
}
