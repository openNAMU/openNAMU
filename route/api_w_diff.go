package route

import (
	"database/sql"
	"opennamu/route/tool"
	"strings"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func Api_w_diff(config tool.Config, doc_name string, before_rev string, after_rev string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    return_data := make(map[string]any)

    if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
        return_data["response"] = "require auth"

        return return_data
    }

    hidden_auth := tool.Check_acl(db, "", "", "hidel_auth", config.IP)

    before_data, before_response := Get_diff_revision(db, doc_name, before_rev, hidden_auth)
    if before_response != "ok" {
        return_data["response"] = before_response
        return_data["data"] = before_rev

        return return_data
    }

    after_data, after_response := Get_diff_revision(db, doc_name, after_rev, hidden_auth)
    if after_response != "ok" {
        return_data["response"] = after_response
        return_data["data"] = after_rev

        return return_data
    }

    diff_tool := diffmatchpatch.New()
    
    before_data = strings.ReplaceAll(before_data, "\r", "")
    after_data = strings.ReplaceAll(after_data, "\r", "")

    diff_list := diff_tool.DiffMain(before_data, after_data, true)
    diff_list = diff_tool.DiffCleanupSemantic(diff_list)

    data_list := make([]map[string]string, 0, len(diff_list))
    for _, diff_data := range diff_list {
        if diff_data.Text == "" {
            continue
        }

        diff_type := "equal"
        switch diff_data.Type {
        case diffmatchpatch.DiffDelete:
            diff_type = "delete"
        case diffmatchpatch.DiffInsert:
            diff_type = "insert"
        }

        data_list = append(data_list, map[string]string{
            "type": diff_type,
            "data": diff_data.Text,
        })
    }

    return_data["response"] = "ok"
    return_data["before_rev"] = before_rev
    return_data["after_rev"] = after_rev
    return_data["data"] = data_list

    return return_data
}

func Get_diff_revision(db *sql.DB, doc_name string, rev string, hidden_auth bool) (string, string) {
    if rev == "0" {
        return "", "ok"
    }

    data := ""
    hide := ""
    exist := tool.QueryRow_DB(
        db,
        "select data, hide from history where title = ? and id = ?",
        []any{ &data, &hide },
        doc_name,
        rev,
    )

    if !exist {
        return "", "not exist"
    }
    if hide != "" && !hidden_auth {
        return "", "require auth"
    }

    return data, "ok"
}
