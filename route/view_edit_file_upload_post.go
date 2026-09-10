package route

import "opennamu/route/tool"

type Upload_file_data struct {
	File_name    string
	File_data    []byte
	File_ext     string
	License      string
	License_text string
	Replace      string
	Captcha      string
}

func View_edit_file_upload_post(config tool.Config, upload_files []Upload_file_data) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	captcha := ""
	if len(upload_files) > 0 {
		captcha = upload_files[0].Captcha
	}
	if !tool.Captcha_check(db, config.Session, config.IP, captcha) {
		return tool.Get_error_page(db, config, "recaptcha")
	}

	last_doc := ""
	for _, v := range upload_files {
		data := api_file_upload_post(
			config,
			v.File_name,
			v.File_data,
			v.File_ext,
			v.License,
			v.License_text,
			"",
			false,
			len(upload_files) > 1,
			v.Replace == "1",
		)
		if data["response"] != "ok" {
			error_name, _ := data["data"].(string)
			if error_name == "" {
				error_name = "upload error"
			}
			return tool.Get_error_page(db, config, error_name)
		}
		last_doc, _ = data["data"].(string)
	}

	if last_doc == "" {
		return tool.Get_error_page(db, config, "invalid data")
	}
	return tool.Get_redirect("/w/" + tool.Url_parser(last_doc))
}
