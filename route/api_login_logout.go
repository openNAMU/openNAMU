package route

import "opennamu/route/tool"

func Api_login_logout(config tool.Config) map[string]any {
	config.Session.Delete("id")

	return_data := make(map[string]any)

	if err := config.Session.Save(); err != nil {
		return_data["response"] = "error"
		return_data["data"] = "session delete error"

		return return_data
	}

	return_data["response"] = "ok"

	return return_data
}
