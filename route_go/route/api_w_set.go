package route

import (
    "opennamu/route/tool"

    jsoniter "github.com/json-iterator/go"
)

func Document_set_list() map[string]string {
    set_list := map[string]string{}

    set_list["document_markup"] = ""
    set_list["document_top"] = ""
    set_list["document_editor_top"] = ""
    set_list["document_comment_code"] = ""

    return set_list
}

func Api_w_set(call_arg []string) string {
    var json = jsoniter.ConfigCompatibleWithStandardLibrary

    other_set := map[string]string{}
    json.Unmarshal([]byte(call_arg[0]), &other_set)

    db := tool.DB_connect()
    defer db.Close()

    set_list := Document_set_list()
    return_data := make(map[string]interface{})

    if _, ok := set_list[other_set["set_name"]]; ok {
        doc_rev := ""
        if val, ok := other_set["doc_rev"]; ok {
            doc_rev = val
        }

        return_data["data"] = tool.Get_document_setting(db, other_set["doc_name"], other_set["set_name"], doc_rev)
        return_data["response"] = "ok"
    } else {
        return_data["response"] = "not exist"
    }

    json_data, _ := json.Marshal(return_data)
    return string(json_data)
}
