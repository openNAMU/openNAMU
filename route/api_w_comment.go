package route

import (
	"opennamu/route/tool"
)

func Api_bbs_w_comment_make(config tool.Config, doc_name string) string {
    config_copy := config
    config_copy.IP = "Tool:System"

    data_api := Api_bbs_w_post(config_copy, "0", doc_name, "")
    data_api_in := data_api["data"]

    return data_api_in
}

func Api_w_comment_ui(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    other_set := map[string]string{}
    json.Unmarshal([]byte(config.Other_set), &other_set)

    return_data := Api_w_comment(config, other_set["doc_name"])
    db_code_str := return_data["data"]

    return_data = make(map[string]string)
    return_data["response"] = "ok"
    return_data["data"] = View_bbs_in_w_comment(db, config, "", "0", db_code_str)

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}

func Api_w_comment(config tool.Config, doc_name string) map[string]string {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    db_code := tool.Get_document_setting(db, doc_name, "document_comment_code", "")
    
    db_code_str := ""
    if len(db_code) >= 1 {
        db_code_str = db_code[0][0]
    }

    if db_code_str == "" {
        db_code_str = Api_bbs_w_comment_make(config, doc_name)

        tool.Exec_DB(
            db,
            "insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'document_comment_code', ?)",
            doc_name,
            db_code_str,
        )
    }

    return_data := make(map[string]string)
    return_data["response"] = "ok"
    return_data["data"] = db_code_str

    return return_data
}