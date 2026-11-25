package route

import (
	"opennamu/route/tool"
)

func Api_calendar_event(config tool.Config, id string, code string) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !Calendar_view_allowed(db, config.IP) {
		return_data["response"] = "require auth"
		return return_data
	}

	id, id_ok := Calendar_id_normalize(id)
	code, code_ok := Calendar_code_normalize(code)
	if !id_ok || !code_ok {
		return_data["response"] = "not exist"
		return_data["data"] = "event"
		return return_data
	}

	event_data, event_exists := Calendar_event_get(db, id, code)
	if !event_exists {
		return_data["response"] = "not exist"
		return_data["data"] = "event"
		return return_data
	}

	calendar_data, calendar_exists := Calendar_meta_get(db, id)
	if !calendar_exists {
		return_data["response"] = "not exist"
		return_data["data"] = "calendar"
		return return_data
	}

	if event_data["user"] != "" {
		event_data["user_render"] = tool.Get_user_profile_image_ui(db, event_data["user"]) + tool.IP_parser(db, event_data["user"], config.IP)
	}

	return_data["response"] = "ok"
	return_data["data"] = event_data
	return_data["editable"] = Calendar_manage_allowed(db, config, calendar_data["owner"], event_data["user"])
	return return_data
}
