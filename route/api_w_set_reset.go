package route

import (
	"opennamu/route/tool"
)

func Api_w_set_reset(config tool.Config, doc_name string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    ip := config.IP

    if tool.Check_acl(db, "", "", "owner_auth", ip) {
        tool.Exec_DB(
            db,
            "delete from acl where title = ?",
            doc_name,
        )

        tool.Exec_DB(
            db,
            "delete from data_set where doc_name = ? and set_name = 'acl_date'",
            doc_name,
        )

        set_list := []string{
            "document_markup",
            "document_top",
            "document_editor_top",
        }

        for for_a := 0; for_a < len(set_list); for_a++ {
            tool.Exec_DB(
                db,
                "delete from data_set where doc_name = ? and set_name = ?",
                doc_name, set_list[for_a],
            )
        }

        return_data := make(map[string]any)
        return_data["response"] = "ok"
        return_data["language"] = map[string]string{
            "reset": tool.Get_language(db, "reset", false),
        }

        return return_data
    } else {
        return_data := make(map[string]any)
        return_data["response"] = "require auth"
        return_data["language"] = map[string]string{
            "authority_error": tool.Get_language(db, "authority_error", false),
        }

        return return_data
    }
}
