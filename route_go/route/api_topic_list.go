package route

import (
	"database/sql"
	"encoding/json"
	"log"
	"opennamu/route/tool"
	"strconv"
)

func Api_topic_list(call_arg []string) string {
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

	stmt, err := db.Prepare(tool.DB_change("select code, sub, stop, agree, date from rd where title = ? order by sub asc limit ?, 50"))
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows, err := stmt.Query(other_set["name"], page_int)
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close()

	var data_list [][]string
	ip_parser_temp := map[string][]string{}

	for rows.Next() {
		var code string
		var sub string
		var stop string
		var agree string
		var date string

		err := rows.Scan(&code, &sub, &stop, &agree, &date)
		if err != nil {
			log.Fatal(err)
		}

		stmt, err := db.Prepare(tool.DB_change("select ip, id from topic where code = ? order by id + 0 desc limit 1"))
		if err != nil {
			log.Fatal(err)
		}
		defer stmt.Close()

		var ip string
		var id string

		err = stmt.QueryRow(code).Scan(&ip, &id)
		if err != nil {
			if err == sql.ErrNoRows {
				ip = ""
			} else {
				log.Fatal(err)
			}
		}

		var ip_pre string
		var ip_render string

		if _, ok := ip_parser_temp[ip]; ok {
			ip_pre = ip_parser_temp[ip][0]
			ip_render = ip_parser_temp[ip][1]
		} else {
			ip_pre = tool.IP_preprocess(db, ip, other_set["ip"])[0]
			ip_render = tool.IP_parser(db, ip, other_set["ip"])

			ip_parser_temp[ip] = []string{ip_pre, ip_render}
		}

		data_list = append(data_list, []string{
			code,
			sub,
			stop,
			agree,
			ip_pre,
			ip_render,
			date,
			id,
		})
	}

	return_data := make(map[string]interface{})
	return_data["language"] = map[string]string{
		"closed":            tool.Get_language(db, "closed", false),
		"agreed_discussion": tool.Get_language(db, "agreed_discussion", false),
		"make_new_topic":    tool.Get_language(db, "make_new_topic", false),
		"stop":              tool.Get_language(db, "stop", false),
	}

	if len(data_list) == 0 {
		return_data["data"] = map[string]string{}
	} else {
		return_data["data"] = data_list
	}

	json_data, _ := json.Marshal(return_data)
	return string(json_data)
}
