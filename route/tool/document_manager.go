package tool

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	"github.com/dlclark/regexp2"
)

func Do_edit_filter(db *sql.DB, config Config, doc_name string, data string) bool {
    if !Check_acl(db, "", "", "edit_filter_pass", config.IP) {
        rows := Query_DB(
            db, 
            `select plus, plus_t from html_filter where kind = 'regex_filter' and plus != ''`,
        )
        defer rows.Close()

        for rows.Next() {
            var plus string
            var plus_t string

            err := rows.Scan(&plus, &plus_t)
            if err != nil {
                panic(err)
            }

            r, err := regexp2.Compile(plus, 0)
            if err != nil {
                continue
            }

            m, err := r.MatchString(data)
            if err == nil && m {
                return false
            }
        }
    }

    return true
}

func Do_edit_send_require_check(db *sql.DB, config Config, data string) bool {
    if !Check_acl(db, "", "", "edit_bottom_compulsion", config.IP) {
        var check string

        QueryRow_DB(
            db,
            `select data from other where name = "edit_bottom_compulsion"`,
            []any{ &check },
        )
        if check != "" && data == "" {
            return false
        }
    }

    return true
}

func Do_edit_text_checkbox_check(db *sql.DB, config Config, data string) bool {
    var check string

    QueryRow_DB(
        db,
        `select data from other where name = "copyright_checkbox_text"`,
        []any{ &check },
    )

    if check != "" && data != "yes" {
        return false
    }

    return true
}

func Do_edit_slow_check(db *sql.DB, config Config, do_type string) bool {
    if !Check_acl(db, "", "", "slow_edit", config.IP) {
        var check string

        if do_type == "edit" {
            QueryRow_DB(
                db,
                `select data from other where name = 'slow_edit'`,
                []any{ &check },
            )
        } else {
            // do_type == "thread"
            QueryRow_DB(
                db,
                `select data from other where name = 'slow_thread'`,
                []any{ &check },
            )
        }

        if check != "" {
            slow_edit := Str_to_int(check)

            var last_edit string

            if do_type == "edit" {
                QueryRow_DB(
                    db,
                    `select date from history where ip = ? order by date desc limit 1`,
                    []any{ &last_edit },
                    config.IP,
                )
            } else {
                QueryRow_DB(
                    db,
                    `select date from topic where ip = ? order by date desc limit 1`,
                    []any{ &last_edit },
                    config.IP,
                )
            }

            if last_edit != "" {
                last_edit_compact := strings.NewReplacer(
                    " ", "",
                    ":", "",
                    "-", "",
                ).Replace(last_edit)

                last_edit_num := Str_to_int(last_edit_compact)
                now_edit_num := Str_to_int(
                    time.Now().Add(-time.Duration(slow_edit) * time.Second).Format("20060102150405"),
                )

                if last_edit_num > now_edit_num {
                    return false
                }
            }
        }
    }

    return true
}

// Do_edit_max_length_check : over is false, under is true
func Do_edit_max_length_check(db *sql.DB, config Config, data string) bool {
    var check string

    exist := QueryRow_DB(
        db,
        "select data from other where name = 'document_content_max_length'",
        []any{ &check },
    )

    if !exist || check == "" {
        return true
    }

    return len(data) <= Str_to_int(check)
}

func Get_edit_length_diff(A string, B string) string {
    A_len := len(A)
    B_len := len(B)

    if A_len > B_len {
        diff_len := A_len - B_len
        diff_len_str := strconv.Itoa(diff_len)

        return "-" + diff_len_str
    } else if B_len > A_len {
        diff_len := B_len - A_len
        diff_len_str := strconv.Itoa(diff_len)

        return "+" + diff_len_str
    } else {
        return "0"
    }
}

func Do_watchlist_alarm_send(db *sql.DB, config Config, doc_name string) {
    rows := Query_DB(
        db,
        `select id from user_set where name = 'watchlist' and data = ?`,
        doc_name,
    )
    defer rows.Close()

    for rows.Next() {
        var id string

        err := rows.Scan(&id)
        if err != nil {
            panic(err)
        }
    
        Send_alarm(
            db,
            id,
            config.IP,
            `<a href="/w/` + Url_parser(doc_name) + `">` + HTML_escape(doc_name) + `</a>`,
        )
    }
}

func Do_add_history(db *sql.DB, doc_name string, data string, date string, ip string, send string, length string, mode string, type_check string) {
    var history_recording_off_check string
    
    QueryRow_DB(
        db,
        `select data from other where name = "history_recording_off"`,
        []any{ &history_recording_off_check },
    )
    if history_recording_off_check != "" {
        return
    }

    var id_data string

    if mode == "add" || mode == "setting" {
        old_id_data := "1"
        
        QueryRow_DB(
            db,
            `select id from history where title = ? order by id + 0 asc limit 1`,
            []any{ &old_id_data },
            doc_name,
        )

        old_id_data_int := Str_to_int(old_id_data)
        old_id_data_int -= 1

        id_data = strconv.Itoa(old_id_data_int)
    } else {
        old_id_data := "0"
        
        QueryRow_DB(
            db,
            `select id from history where title = ? order by id + 0 desc limit 1`,
            []any{ &old_id_data },
            doc_name,
        )

        old_id_data_int := Str_to_int(old_id_data)
        old_id_data_int += 1

        id_data = strconv.Itoa(old_id_data_int)

        if id_data == "1" {
            mode = "r1"
        }

        if strings.HasPrefix(doc_name, "user:") {
            mode = "user"
        } else if strings.HasPrefix(doc_name, "category:") {
            mode = "category"
        } else if strings.HasPrefix(doc_name, "file:") {
            mode = "file"
        }
    }

    send = strings.ReplaceAll(send, "<", "")
    send = strings.ReplaceAll(send, ">", "")

    if len(send) > 512 {
        send = send[:512]
    }

    if type_check != "" {
        send = send + " (" + type_check + ")"
    }

    if mode != "add" && mode != "setting" && mode != "user" {
        Do_add_recent_history(db, "normal", id_data, doc_name, date)
    }

    if mode != "add" && mode != "setting" {
        Do_add_recent_history(db, mode, id_data, doc_name, date)

        var document_count int

        QueryRow_DB(
            db, 
            `select count(*) from data`,
            []any{ &document_count },
        )

        document_count_str := strconv.Itoa(document_count)

        Exec_DB(
            db,
            `delete from other where name = "count_all_title"`,
        )
        Exec_DB(
            db,
            `insert into other (name, data, coverage) values ("count_all_title", ?, "")`,
            document_count_str,
        )

        data_set_exist := ""
        if mode == "delete" {
            data_set_exist = "not_exist"
        }

        Exec_DB(
            db,
            `delete from data_set where doc_name = ? and set_name = "edit_request_doing"`,
            doc_name,
        )

        Exec_DB(
            db,
            `delete from data_set where doc_name = ? and set_name = "last_edit"`,
            doc_name,
        )
        Exec_DB(
            db,
            `insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'last_edit', ?)`,
            doc_name,
            date,
        )

        Exec_DB(
            db,
            `delete from data_set where doc_name = ? and set_name = "length"`,
            doc_name,
        )
        Exec_DB(
            db,
            `insert into data_set (doc_name, doc_rev, set_name, set_data) values (?, '', 'length', ?)`,
            doc_name,
            len(data),
        )

        Exec_DB(
            db,
            `update data_set set doc_rev = ? where doc_name = ? and (doc_rev = '' or doc_rev = 'not_exist')`,
            data_set_exist,
            doc_name,
        )
    }

    Exec_DB(
        db,
        `insert into history (id, title, data, date, ip, send, leng, hide, type) values (?, ?, ?, ?, ?, ?, ?, '', ?)`,
        id_data,
        doc_name,
        data,
        date,
        ip,
        send,
        length,
        mode,
    )
}

func Do_add_recent_history(db *sql.DB, mode string, id string, title string, date string) {
    var length int

    QueryRow_DB(
        db,
        `select count(*) from rc where type = ?`,
        []any{ &length },
        mode,
    )
    if length >= 200 {
        var id string
        var title string

        exist := QueryRow_DB(
            db,
            `select id, title from rc where type = ? order by date asc limit 1`,
            []any{ &id, &title },
            mode,
        )

        if exist {
            Exec_DB(
                db,
                `delete from rc where id = ? and title = ? and type = ?`,
                id,
                title,
                mode,
            )
        }
    }

    Exec_DB(
        db,
        `insert into rc (id, title, date, type) values (?, ?, ?, ?)`,
        id,
        title,
        date,
        mode,
    )
}

func Get_document_markup(db *sql.DB, doc_name string, do_type string) string {
    markup := ""

    if do_type == "document" {
        QueryRow_DB(
            db,
            "select set_data from data_set where doc_name = ? and set_name = 'document_markup'",
            []any{ &markup },
            doc_name,
        )

        if markup == "" {
            QueryRow_DB(
                db,
                "select data from other where name = 'markup'",
                []any{ &markup },
            )
        }
    }

    if markup == "" || markup == "namumark_beta" {
        markup = "namumark"
    }

    return markup
}

func Get_document_is_redirect(db *sql.DB, doc_name string) bool {
    var is_redirect string

    exist := QueryRow_DB(
        db,
        "select distinct type from back where link = ? and type = 'redirect'",
        []any{ &is_redirect },
        doc_name,
    )

    return exist
}