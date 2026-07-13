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
	} else {
		log.Default().Println(data + " (" + language + ")")
		return data + " (" + language + ")"
	}
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
		global_lang_data[language + "_" + k] = v
	}
	global_lang_load[language] = true
}
