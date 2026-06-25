package route

import (
	"opennamu/route/tool"
)

func Api_bbs_w_comment_all(config tool.Config, sub_code string, already_auth_check bool, do_type string) []map[string]string {
    end_data := []map[string]string{}
    
    return_data := Api_bbs_w_comment_one(config, already_auth_check, do_type, sub_code)
    return_data_in := return_data["data"].([]map[string]string)

    for for_a := 0; for_a < len(return_data_in); for_a++ {
        end_data = append(end_data, return_data_in[for_a])

        temp := Api_bbs_w_comment_all(config, sub_code + "-" + return_data_in[for_a]["code"], already_auth_check, do_type)
        if len(temp) > 0 {
            for for_b := 0; for_b < len(temp); for_b++ {
                end_data = append(end_data, temp[for_b])
            }
        }
    }

    return end_data
}

func Api_bbs_w_comment(config tool.Config, do_type string, sub_code string) map[string]any {
    db := tool.DB_connect()
    defer tool.DB_close(db)

    if do_type == "length" {
        bbs_and_post_num := sub_code

        comment_length := "0"
        tool.QueryRow_DB(
            db,
            "select count(*) from bbs_data where set_name = 'comment_date' and set_id = ? order by set_code + 0 desc",
            []any{ &comment_length },
            bbs_and_post_num,
        )

        reply_length := "0"
        tool.QueryRow_DB(
            db,
            "select count(*) from bbs_data where set_name = 'comment_date' and set_id = ? order by set_code + 0 desc",
            []any{ &reply_length },
            bbs_and_post_num + "-%",
        )

        comment_length_int := tool.Str_to_int(comment_length)
        reply_length_int := tool.Str_to_int(reply_length)
        
        length_int := comment_length_int + reply_length_int

        data_list := map[string]any{
            "response" : "ok",
            "comment" : comment_length,
            "reply" : reply_length,
            "data" : length_int,
        }

        return data_list
    } else {
        return_data := make(map[string]any)
        
        temp := []map[string]string{}
        if !tool.Check_acl(db, "", "", "bbs_comment", config.IP) {
            return_data["response"] = "require auth"
        } else {
            temp = Api_bbs_w_comment_all(config, sub_code, true, do_type)
        }

        return_data["response"] = "ok"
        return_data["data"] = temp

        return return_data
    }
}
