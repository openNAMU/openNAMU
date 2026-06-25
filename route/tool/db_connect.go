package tool

import (
	"database/sql"
	"os"
	"path/filepath"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"
	_ "modernc.org/sqlite"
)

var db_set = map[string]string{}

func Get_DB_set() map[string]string {
    new_db_set := map[string]string{}

    db_env, db_env_exist := os.LookupEnv("NAMU_DB")
    db_env_type, db_env_type_exist := os.LookupEnv("NAMU_DB_TYPE")
    if db_env_exist || db_env_type_exist {
        new_db_set["db_name"] = Choose(db_env, "data")
        new_db_set["db_type"] = Choose(db_env_type, "sqlite")

        return new_db_set
    }
    
    path_dir := filepath.Join("..", "data", "set.json")
    if File_exist_check(path_dir) {
        raw, err := os.ReadFile(path_dir)
        if err == nil {
            tmp := map[string]string{}
            if err := json.Unmarshal(raw, &tmp); err == nil {
                if v, ok := tmp["db"]; ok {
                    new_db_set["db_name"] = v
                } else {
                    new_db_set["db_name"] = "data"
                }
                
                if v, ok := tmp["db_type"]; ok {
                    new_db_set["db_type"] = v
                } else {
                    new_db_set["db_type"] = "sqlite"
                }

                return new_db_set
            }
        }
    }

    new_db_set["db_name"] = "data"
    new_db_set["db_type"] = "sqlite"
    
    return new_db_set
}

func Get_DB_set_MySQL(new_db_set map[string]string) map[string]string {
    path := filepath.Join("..", "data", "mysql.json")
    if !File_exist_check(path) {
        return new_db_set
    }

    raw, err := os.ReadFile(path)
    if err != nil {
        return new_db_set
    }

    tmp := new_db_set
    if err := json.Unmarshal(raw, &tmp); err != nil {
        return tmp
    }

    if host, ok := tmp["host"]; ok && host != "" {
        tmp["db_mysql_host"] = host
    } else {
        tmp["db_mysql_host"] = "127.0.0.1"
    }

    if port, ok := tmp["port"]; ok && port != "" {
        tmp["db_mysql_port"] = port
    } else {
        tmp["db_mysql_port"] = "3306"
    }

    if user, ok := tmp["user"]; ok {
        tmp["db_mysql_user"] = user
    }

    if pw, ok := tmp["password"]; ok {
        tmp["db_mysql_pw"] = pw
    }

    return tmp
}

func Exec_DB(db *sql.DB, query string, values ...any) {
    const retryDelay = 10 * time.Millisecond

    stmt, err := db.Prepare(DB_change(query))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    for {
        _, err = stmt.Exec(values...)
        if err == nil {
            return
        }

        if strings.Contains(err.Error(), "database is locked") {
            time.Sleep(retryDelay)
            continue
        }

        panic(err)
    }
}

func Query_DB(db *sql.DB, query string, values ...any) *sql.Rows {
    const retryDelay = 10 * time.Millisecond

    stmt, err := db.Prepare(DB_change(query))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    for {
        rows, err := stmt.Query(values...)
        if err == nil {
            return rows
        }

        if strings.Contains(err.Error(), "database is locked") {
            time.Sleep(retryDelay)
            continue
        }

        panic(err)
    }
}

// QueryRow_DB 이래서 포인터를 배우는구나...
func QueryRow_DB(db *sql.DB, query string, var_list []any, values ...any) bool {
    const retryDelay = 10 * time.Millisecond

    stmt, err := db.Prepare(DB_change(query))
    if err != nil {
        panic(err)
    }
    defer stmt.Close()

    for {
        row := stmt.QueryRow(values...)

        err := row.Scan(var_list...)
        switch err {
        case nil:
            return true
        case sql.ErrNoRows:
            return false
        }

        if strings.Contains(err.Error(), "database is locked") {
            time.Sleep(retryDelay)
            continue
        }

        panic(err)
    }
}

func DB_boot() map[string]string {
    new_db_set := Get_DB_set()
    if new_db_set["db_type"] == "mysql" {
        new_db_set = Get_DB_set_MySQL(new_db_set)
    }

    db_set = new_db_set

    return new_db_set
}

func DB_connect_init() (*sql.DB, error) {
    if db_set["db_type"] == "sqlite" {
        db, err := sql.Open("sqlite", filepath.Join("..", db_set["db_name"] + ".db") + "?_journal_mode=WAL&_busy_timeout=5000")
        if err != nil {
            return nil, err
        }

        if err := db.Ping(); err != nil {
            db.Close()
            return nil, err
        }

        return db, nil
    } else {
        db, err := sql.Open("mysql", db_set["db_mysql_user"] + ":" + db_set["db_mysql_pw"] + "@tcp(" + db_set["db_mysql_host"] + ":" + db_set["db_mysql_port"] + ")")
        if err != nil {
            return nil, err
        }

        if err := db.Ping(); err != nil {
            db.Close()
            return nil, err
        }

        return db, nil
    }
}

func DB_connect() *sql.DB {
    // log.Default().Println("DB open")

    if db_set["db_type"] == "sqlite" {
        db, err := sql.Open("sqlite", filepath.Join("..", db_set["db_name"] + ".db") + "?_journal_mode=WAL&_busy_timeout=5000")
        if err != nil {
            panic(err)
        }

        return db
    } else {
        db, err := sql.Open("mysql", db_set["db_mysql_user"] + ":" + db_set["db_mysql_pw"] + "@tcp(" + db_set["db_mysql_host"] + ":" + db_set["db_mysql_port"] + ")/" + db_set["db_name"])
        if err != nil {
            panic(err)
        }

        return db
    }
}

func DB_close(db *sql.DB) {
    db.Close()
    
    // log.Default().Println("DB close")
}

func Get_DB_type() string {
    return db_set["db_type"]
}

func DB_change(data string) string {
    if Get_DB_type() == "mysql" {
        data = strings.Replace(data, "random()", "rand()", -1)
        data = strings.Replace(data, "collate nocase", "collate utf8mb4_general_ci", -1)
    }

    return data
}
