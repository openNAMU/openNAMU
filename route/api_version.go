package route

import (
	"database/sql"
	stdjson "encoding/json"
	"io"
	"net/http"
	"time"

	"opennamu/route/tool"
)

func get_version_branch(db *sql.DB) string {
	branch := ""
	tool.QueryRow_DB(db, `select data from other where name = "update"`, []any{&branch})
	if branch != "stable" && branch != "beta" && branch != "dev" {
		branch = "stable"
	}
	return branch
}

func get_remote_version(branch string) string {
	if branch != "stable" && branch != "beta" && branch != "dev" {
		return ""
	}

	client := http.Client{Timeout: 3 * time.Second}
	response, err := client.Get("https://raw.githubusercontent.com/openNAMU/openNAMU/" + branch + "/version.json")
	if err != nil {
		return ""
	}
	defer response.Body.Close()
	if response.StatusCode != http.StatusOK {
		return ""
	}

	version_data := map[string]any{}
	decoder := stdjson.NewDecoder(io.LimitReader(response.Body, 1<<20))
	if err := decoder.Decode(&version_data); err != nil {
		return ""
	}
	if beta_data, ok := version_data["beta"].(map[string]any); ok {
		if version, ok := beta_data["r_ver"].(string); ok {
			return version
		}
	}
	version, _ := version_data["r_ver"].(string)
	return version
}

func Api_version(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	version_list := tool.Get_last_version()
	up_data := get_version_branch(db)

	return_data := make(map[string]any)
	return_data["version"] = version_list["r_ver"]
	return_data["db_version"] = version_list["c_ver"]
	return_data["skin_version"] = version_list["s_ver"]
	return_data["build"] = up_data

	return return_data
}
