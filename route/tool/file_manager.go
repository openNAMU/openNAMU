package tool

import (
	"crypto/sha256"
	"database/sql"
	"encoding/hex"
	"io"
	"path/filepath"
	"strings"

	"github.com/dlclark/regexp2"
)

func Get_file_max_size(db *sql.DB) int {
    data := "0"

    QueryRow_DB(
        db,
        "select data from other where name = 'upload'",
        []any{ &data },
    )

    file_max_size := Str_to_int(data)

    return file_max_size
}

func Get_file_main_dir(db *sql.DB) string {
    data := ""

    QueryRow_DB(
        db,
        "select data from other where name = 'image_where'",
        []any{ &data },
    )

    if data == "" {
        data = filepath.Join("..", "data", "images")
    } else {
        data = filepath.Clean(data)
    }

    return data
}

func Get_ext_allow_list(db *sql.DB) map[string]bool {
    rows := Query_DB(db, "select html from html_filter where kind = 'extension'")
    defer rows.Close()

    data_list := map[string]bool{}

    for rows.Next() {
        data := ""

        err := rows.Scan(&data)
        if err != nil {
            panic(err)
        }

        data = strings.ToLower(data)
        data = strings.TrimPrefix(data, ".")

        data_list[data] = true
    }

    return data_list
}

func Get_file_name_unallow_check(db *sql.DB, file_name string) bool {
    rows := Query_DB(db, "select html from html_filter where kind = 'file_name'")
    defer rows.Close()

    for rows.Next() {
        data := ""

        err := rows.Scan(&data)
        if err != nil {
            panic(err)
        }

        r, err := regexp2.Compile(data, 0)
        if err != nil {
            continue
        }

        m, err := r.MatchString(file_name)
        if err == nil && m {
            return true
        }
    }

    return false
}

func File_name_to_dir(file_name string, file_ext string) string {
    h := sha256.New224()
    io.WriteString(h, file_name)
    hash_hex := hex.EncodeToString(h.Sum(nil))

    return hash_hex + "." + file_ext
}