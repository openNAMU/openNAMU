package tool

func Namumark() map[string]interface{} {
	end_data := make(map[string]interface{})
	end_data["data"] = ""
	end_data["js_data"] = ""
	end_data["backlink"] = [][]string{}
	end_data["link_count"] = 0

	return end_data
}
