package route

import (
	"encoding/json"
	"fmt"
	"opennamu/route/tool"
	"strconv"
)

func Api_list_history_tool(call_arg []string) {
	db_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[0]), &db_set)

	other_set := map[string]string{}
	json.Unmarshal([]byte(call_arg[1]), &other_set)

	db := tool.DB_connect(db_set)
	if db == nil {
		return
	}
	defer db.Close()

	main_dict := [][]string{}

	main_dict = append(
		main_dict,
		[]string{
			"/raw_rev/" + other_set["rev"] + "/" + tool.Url_parser(other_set["doc_name"]),
			tool.Get_language(db, db_set, "raw", false),
		},
	)
	main_dict = append(
		main_dict,
		[]string{
			"/revert/" + other_set["rev"] + "/" + tool.Url_parser(other_set["doc_name"]),
			tool.Get_language(db, db_set, "revert", false) + " (r" + other_set["rev"] + ")",
		},
	)

	rev_int, err := strconv.Atoi(other_set["rev"])
	if err != nil {
		rev_int = 0
	}

	if rev_int > 1 {
		rev_str := strconv.Itoa(rev_int - 1)

		main_dict = append(
			main_dict,
			[]string{
				"/revert/" + rev_str + "/" + tool.Url_parser(other_set["doc_name"]),
				tool.Get_language(db, db_set, "revert", false) + " (r" + rev_str + ")",
			},
		)
		main_dict = append(
			main_dict,
			[]string{
				"/diff/" + rev_str + "/" + other_set["rev"] + "/" + tool.Url_parser(other_set["doc_name"]),
				tool.Get_language(db, db_set, "compare", false),
			},
		)
	}

	main_dict = append(
		main_dict,
		[]string{
			"/history/" + tool.Url_parser(other_set["doc_name"]),
			tool.Get_language(db, db_set, "history", false),
		},
	)

	auth_name := tool.Get_user_auth(db, db_set, other_set["ip"])
	auth_info := tool.Get_auth_group_info(db, db_set, auth_name)

	_, ok := auth_info["hidel"]
	_, ok2 := auth_info["owner"]

	if ok || ok2 {
		main_dict = append(
			main_dict,
			[]string{
				"/history_hidden/" + other_set["rev"] + "/" + tool.Url_parser(other_set["doc_name"]),
				tool.Get_language(db, db_set, "hide", false),
			},
		)
	}

	if ok2 {
		main_dict = append(
			main_dict,
			[]string{
				"/history_delete/" + other_set["rev"] + "/" + tool.Url_parser(other_set["doc_name"]),
				tool.Get_language(db, db_set, "history_delete", false),
			},
		)
		main_dict = append(
			main_dict,
			[]string{
				"/history_send/" + other_set["rev"] + "/" + tool.Url_parser(other_set["doc_name"]),
				tool.Get_language(db, db_set, "send_edit", false),
			},
		)
	}

	json_data, _ := json.Marshal(main_dict)
	fmt.Print(string(json_data))
}
