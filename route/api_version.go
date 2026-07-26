package route

import "opennamu/route/tool"

func Api_version(config tool.Config) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    version_list := tool.Get_last_version()

    up_data := ""
    tool.QueryRow_DB(
        db,
        `select data from other where name = "update"`,
        []any{ &up_data },
    )

    if up_data != "stable" && up_data != "beta" && up_data != "dev" {
        up_data = "stable"
    }

    return_data := make(map[string]any)
    return_data["version"] = version_list["r_ver"]
    return_data["db_version"] = version_list["c_ver"]
    return_data["skin_version"] = version_list["s_ver"]
    return_data["build"] = up_data

    return return_data
}
