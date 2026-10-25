package tool

import (
	"database/sql"
	"log"
	"sync"

	lang_data "opennamu/lang"
)

var global_lang_data = map[string]string{}
var global_lang_load = map[string]bool{}
var global_lang_lock sync.RWMutex

// lang_db_lock guards the filter-dedicated connection pool below. Readers
// (Get_language_shared) hold RLock for the duration of one translation, the
// swap path takes the write Lock, so a pool handed to a caller is never closed
// while that caller is still using it.
var lang_db_lock sync.RWMutex
var lang_db_pool *sql.DB
var lang_db_key string

// db_identity_key describes the DB the current connection parameters point at,
// so a cached pool is invalidated whenever the backend or DB file changes
// (e.g. NAMU_DB repointed in tests, or an sqlite/mysql switch).
func db_identity_key() string {
	key := Get_DB_type() + "|" + db_set["db_name"]
	if Get_DB_type() == "mysql" {
		key += "|" + db_set["db_mysql_user"] + "|" + db_set["db_mysql_pw"] +
			"|" + db_set["db_mysql_host"] + "|" + db_set["db_mysql_port"]
	}
	return key
}

// ensure_lang_db_locked returns the shared pool, rebuilding it when the DB
// identity changed. Callers must hold lang_db_lock for writing. The previous
// pool is closed only here, where no reader can still be using it.
func ensure_lang_db_locked(identity string) *sql.DB {
	if lang_db_pool != nil && lang_db_key == identity {
		return lang_db_pool
	}

	old := lang_db_pool
	lang_db_pool = DB_connect()
	lang_db_key = identity
	if old != nil {
		_ = old.Close()
	}

	return lang_db_pool
}

// Get_language_shared runs one translation against the filter-dedicated shared
// pool. The pool cannot be closed mid-call: readers hold lang_db_lock.RLock
// for the whole call and the swap path needs the write Lock. This is the
// entry point the load_lang pongo2 filter must use.
//
// The current language is still read from the DB on every call, so setting
// changes keep reflecting immediately. Failure semantics are unchanged:
// DB_connect/query errors still panic exactly as they did with per-call pools.
func Get_language_shared(key string, safe bool) string {
	identity := db_identity_key()

	lang_db_lock.RLock()
	if lang_db_pool != nil && lang_db_key == identity {
		defer lang_db_lock.RUnlock()
		return Get_language(lang_db_pool, key, safe)
	}
	lang_db_lock.RUnlock()

	// slow path: pool missing or DB identity changed
	lang_db_lock.Lock()
	defer lang_db_lock.Unlock()
	db := ensure_lang_db_locked(db_identity_key())
	return Get_language(db, key, safe)
}

// Get_lang_db returns the filter-dedicated shared *sql.DB pool. The returned
// pointer stays valid only until the next DB identity change, at which point
// the pool is closed; production code must use Get_language_shared, which
// holds a read lease for the whole call. This accessor exists for tests and
// diagnostics that exercise the swap behavior itself.
func Get_lang_db() *sql.DB {
	lang_db_lock.Lock()
	defer lang_db_lock.Unlock()
	return ensure_lang_db_locked(db_identity_key())
}

func Get_language(db *sql.DB, data string, safe bool) string {
	language := "ko-KR"
	QueryRow_DB(
		db,
		"select data from other where name = 'language'",
		[]any{&language},
	)

	lang_value, ok := Get_lang_cache(language, data)
	if !ok {
		Load_lang_data(language)
		lang_value, ok = Get_lang_cache(language, data)
	}

	if ok {
		if safe {
			return lang_value
		} else {
			return HTML_escape(lang_value)
		}
	}
	if language != "ko-KR" {
		fallback, fallback_ok := Get_lang_cache("ko-KR", data)
		if !fallback_ok {
			Load_lang_data("ko-KR")
			fallback, fallback_ok = Get_lang_cache("ko-KR", data)
		}
		if fallback_ok {
			if safe {
				return fallback
			}
			return HTML_escape(fallback)
		}
	}
	log.Default().Println(data + " (" + language + ")")
	return data + " (" + language + ")"
}

func Get_lang_cache(language string, data string) (string, bool) {
	global_lang_lock.RLock()
	defer global_lang_lock.RUnlock()

	lang_value, ok := global_lang_data[language+"_"+data]
	return lang_value, ok
}

func Load_lang_data(language string) {
	global_lang_lock.RLock()
	lang_load := global_lang_load[language]
	global_lang_lock.RUnlock()
	if lang_load {
		return
	}

	raw, err := lang_data.Read(language)
	if err != nil && language != "ko-KR" {
		raw, err = lang_data.Read("ko-KR")
	}
	if err != nil {
		panic(err)
	}

	lang_json := map[string]string{}
	if err := json.Unmarshal(raw, &lang_json); err != nil {
		panic(err)
	}

	global_lang_lock.Lock()
	defer global_lang_lock.Unlock()

	for k, v := range lang_json {
		global_lang_data[language+"_"+k] = v
	}
	global_lang_load[language] = true
}
