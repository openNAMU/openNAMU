package route

import (
	"database/sql"
	"log"
	"opennamu/route/tool"
	"strconv"

	jsoniter "github.com/json-iterator/go"
)

func Api_list_recent_block(call_arg []string) string {
	var json = jsoniter.ConfigCompatibleWithStandardLibrary

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &other_set)

	db := tool.DB_connect()
	defer db.Close()

	page_int, err := strconv.Atoi(other_set["num"])
	if err != nil {
		log.Fatal(err)
	}

	if page_int > 0 {
		page_int = (page_int * 50) - 50
	} else {
		page_int = 0
	}

	// private 공개 안되도록 조심할 것
	var stmt *sql.Stmt
	var rows *sql.Rows
	if other_set["set_type"] == "all" {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band != 'private' order by today desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(page_int)
		if err != nil {
			log.Fatal(err)
		}
	} else if other_set["set_type"] == "ongoing" {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where ongoing = '1' and band != 'private' order by end desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(page_int)
		if err != nil {
			log.Fatal(err)
		}
	} else if other_set["set_type"] == "regex" {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'regex' order by today desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(page_int)
		if err != nil {
			log.Fatal(err)
		}
	} else if other_set["set_type"] == "private" {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'private' order by today desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(page_int)
		if err != nil {
			log.Fatal(err)
		}
	} else if other_set["set_type"] == "user" {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where block = ? and band != 'private' order by today desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(other_set["user_name"], page_int)
		if err != nil {
			log.Fatal(err)
		}
	} else if other_set["set_type"] == "cidr" {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where band = 'cidr' order by today desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(page_int)
		if err != nil {
			log.Fatal(err)
		}
	} else {
		stmt, err = db.Prepare(tool.DB_change("select why, block, blocker, end, today, band, ongoing from rb where blocker = ? and band != 'private' order by today desc limit ?, 50"))
		if err != nil {
			log.Fatal(err)
		}

		rows, err = stmt.Query(other_set["user_name"], page_int)
		if err != nil {
			log.Fatal(err)
		}
	}

	defer stmt.Close()
	defer rows.Close()

	data_list := [][]string{}
	ip_parser_temp := map[string][]string{}

	for rows.Next() {
		var why string
		var block string
		var blocker string
		var end string
		var today string
		var band string
		var ongoing string

		err := rows.Scan(&why, &block, &blocker, &end, &today, &band, &ongoing)
		if err != nil {
			log.Fatal(err)
		}

		var ip_pre_blocker string
		var ip_render_blocker string

		if _, ok := ip_parser_temp[blocker]; ok {
			ip_pre_blocker = ip_parser_temp[blocker][0]
			ip_render_blocker = ip_parser_temp[blocker][1]
		} else {
			ip_pre_blocker = tool.IP_preprocess(db, blocker, other_set["ip"])[0]
			ip_render_blocker = tool.IP_parser(db, blocker, other_set["ip"])

			ip_parser_temp[blocker] = []string{ip_pre_blocker, ip_render_blocker}
		}

		var ip_pre_block string
		var ip_render_block string

		if band == "" {
			if _, ok := ip_parser_temp[block]; ok {
				ip_pre_block = ip_parser_temp[block][0]
				ip_render_block = ip_parser_temp[block][1]
			} else {
				ip_pre_block = tool.IP_preprocess(db, block, other_set["ip"])[0]
				ip_render_block = tool.IP_parser(db, block, other_set["ip"])

				ip_parser_temp[block] = []string{ip_pre_block, ip_render_block}
			}
		} else {
			ip_pre_block = block
			ip_render_block = block
		}

		data_list = append(data_list, []string{
			why,
			ip_pre_block,
			ip_render_block,
			ip_pre_blocker,
			ip_render_blocker,
			end,
			today,
			band,
			ongoing,
		})
	}

	if other_set["set_type"] == "private" {
		if !tool.Check_acl(db, "", "", "owner_auth", other_set["ip"]) {
			data_list = [][]string{}
		}
	}

	return_data := make(map[string]interface{})
	return_data["language"] = map[string]string{
		"all":         tool.Get_language(db, "all", false),
		"regex":       tool.Get_language(db, "regex", false),
		"cidr":        tool.Get_language(db, "cidr", false),
		"private":     tool.Get_language(db, "private", false),
		"in_progress": tool.Get_language(db, "in_progress", false),
		"admin":       tool.Get_language(db, "admin", false),
		"blocked":     tool.Get_language(db, "blocked", false),
		"limitless":   tool.Get_language(db, "limitless", false),
		"release":     tool.Get_language(db, "release", false),
		"start":       tool.Get_language(db, "start", false),
		"end":         tool.Get_language(db, "end", false),
		"ban":         tool.Get_language(db, "ban", false),
	}
	return_data["data"] = data_list

	auth_name := tool.Get_user_auth(db, other_set["ip"])
	auth_info := tool.Get_auth_group_info(db, auth_name)

	return_data["auth"] = auth_info

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
