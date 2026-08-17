package route

import "opennamu/route/tool"

func Api_bbs_w_comment_all(config tool.Config, sub_code string, already_auth_check bool, do_type string) map[string]any {
	end_data := []map[string]string{}

	return_data := Api_bbs_w_comment_one(config, already_auth_check, do_type, sub_code)
	return_data_in := return_data["data"].([]map[string]string)

	for for_a := 0; for_a < len(return_data_in); for_a++ {
		end_data = append(end_data, return_data_in[for_a])

		temp_data := Api_bbs_w_comment_all(config, sub_code+"-"+return_data_in[for_a]["code"], already_auth_check, do_type)
		temp, _ := temp_data["data"].([]map[string]string)
		if len(temp) > 0 {
			for for_b := 0; for_b < len(temp); for_b++ {
				end_data = append(end_data, temp[for_b])
			}
		}
	}

	return map[string]any{
		"response": "ok",
		"data":     end_data,
	}
}
