package main

import (
	"encoding/json"
	"fmt"
	"os"

	"opennamu/tool"
)

func main() {
	call_arg := os.Args[1:]
	db_set := map[string]string{}

	json.Unmarshal([]byte(call_arg[0]), &db_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	var title string
	err := db.QueryRow(tool.DB_change(db_set, "select title from data where title not like 'user:%' and title not like 'category:%' and title not like 'file:%' order by random() limit 1")).Scan(&title)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Print(title)
}
