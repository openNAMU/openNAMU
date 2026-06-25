package route

import (
	"opennamu/route/tool"
)

func Api_w_raw(config tool.Config, doc_name string, exist_check string, rev string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    new_data := make(map[string]any)

    if !tool.Check_acl(db, doc_name, "", "render", config.IP) {
        new_data["response"] = "require auth"
    } else if exist_check != "" {
        title := ""
        exist := tool.QueryRow_DB(
            db,
            "select title from data where title = ? collate nocase",
            []any{ &title },
            doc_name,
        )

        if !exist {
            new_data["exist"] = false
        } else {
            new_data["exist"] = true
        }

        new_data["data"] = title
        new_data["response"] = "ok"
    } else {        
        exist := false

        data := ""
        hide := ""
        if rev != "" {
            exist = tool.QueryRow_DB(
                db,
                "select data, hide from history where title = ? and id = ?",
                []any{ &data, &hide },
                doc_name, rev,
            )
        } else {
            exist = tool.QueryRow_DB(
                db,
                "select data from data where title = ?",
                []any{ &data },
                doc_name,
            )
        }

        if !exist {
            new_data["response"] = "not exist"
        } else {
            check_pass := false
            if hide != "" {
                if tool.Check_acl(db, "", "", "hidel_auth", config.IP) {
                    check_pass = true
                } else {
                    new_data["response"] = "require auth"
                }
            } else {
                check_pass = true
            }

            if check_pass {
                new_data["title"] = doc_name
                new_data["data"] = data

                new_data["response"] = "ok"
            }
        }
    }

    return new_data
}
