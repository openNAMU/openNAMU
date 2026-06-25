package route

import (
	"opennamu/route/tool"
)

func View_main_other(config tool.Config) string {
    db := tool.DB_connect()
    defer tool.DB_close(db)
    
    out := tool.Get_template(
        db,
        config,
        tool.Get_language(db, "other_tool", true),
        `<h2>` + tool.Get_language(db, "user_tool", true) + `</h2>
        <ul>
            <li><a href="/manager/6">` + tool.Get_language(db, "user_tool", true) + `</a></li>
        </ul>
        <h2>` + tool.Get_language(db, "list", true) + `</h2>
        <h3>` + tool.Get_language(db, "admin", true) + `</h3>
        <ul>               
            <li><a href="/list/admin">` + tool.Get_language(db, "admin_list", true) + `</a></li>
            <li><a href="/list/admin/auth_use">` + tool.Get_language(db, "authority_use_list", true) + `</a></li>
        </ul>
        <h3>` + tool.Get_language(db, "discussion", true) + `</h3>
        <ul>
            <li><a href="/recent_discuss">` + tool.Get_language(db, "recent_discussion", true) + `</a></li>
        </ul>
        <h3>` + tool.Get_language(db, "document", true) + `</h3>
        <ul>
            <li><a href="/recent_change">` + tool.Get_language(db, "recent_change", true) + `</a></li>
            <li><a href="/list/random">` + tool.Get_language(db, "random_list", true) + `</a></li>
            <li><a href="/list/document/all">` + tool.Get_language(db, "all_document_list", true) + `</a></li>
            <li><a href="/list/document/acl">` + tool.Get_language(db, "acl_document_list", true) + `</a></li>
            <li><a href="/list/document/need">` + tool.Get_language(db, "need_document", true) + `</a></li>
            <li><a href="/list/document/long">` + tool.Get_language(db, "long_page", true) + `</a></li>
            <li><a href="/list/document/short">` + tool.Get_language(db, "short_page", true) + `</a></li>
            <li><a href="/list/document/old">` + tool.Get_language(db, "old_page", true) + `</a></li>
            <li><a href="/list/document/new">` + tool.Get_language(db, "new_page", true) + `</a></li>
            <li><a href="/list/document/no_link">` + tool.Get_language(db, "no_link_document_list", true) + `</a></li>
        </ul>
        <h3>` + tool.Get_language(db, "user", true) + `</h3>
        <ul>
            <li><a href="/recent_block">` + tool.Get_language(db, "recent_ban", true) + `</a></li>
            <li><a href="/list/user">` + tool.Get_language(db, "member_list", true) + `</a></li>
        </ul>
        <h3>` + tool.Get_language(db, "other", true) + `</h3>
        <ul>
            <li><a href="/list/file">` + tool.Get_language(db, "image_file_list", true) + `</a></li>
            <li><a href="/vote">` + tool.Get_language(db, "vote_list", true) + `</a></li>
            <li><a href="/bbs/main">` + tool.Get_language(db, "bbs_main", true) + `</a></li>
        </ul>
        <h2>` + tool.Get_language(db, "other", true) + `</h2>
        <ul>
            <li><a href="/upload">` + tool.Get_language(db, "upload", true) + `</a></li>
            <li><a href="/manager/10">` + tool.Get_language(db, "search", true) + `</a></li>
        </ul>
        <h2>` + tool.Get_language(db, "admin", true) + `</h2>
        <ul>
            <li><a href="/manager/1">` + tool.Get_language(db, "admin_tool", true) + `</a></li>
        </ul>`,
        []any{},
        [][]any{},
        map[string]string{},
    )

    return out
}
