package route

import (
	"opennamu/route/tool"
)

func View_edit_file_upload_post(config tool.Config, upload_files []map[string]string) string {
    resp := []any{}
    for _, v := range upload_files {
        data := Api_file_upload_post(
            config,
            v["file_name"],
            v["file_data"],
            v["file_ext"],
        )

        resp = append(resp, data["data"])
    }

    return ""
}
