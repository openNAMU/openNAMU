package route

import (
	"net/url"
	"opennamu/route/tool"
)

func View_challenge(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if values == nil && !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}

	if values != nil {
		result := Api_challenge_post(config)
		if result["response"] == "require auth" {
			return tool.Get_redirect("/user")
		}
		if result["response"] != "ok" {
			return tool.Get_error_page(db, config, "error")
		}
		return tool.Get_redirect("/challenge")
	}

	challenge_list := []struct {
		image    string
		name     string
		complete bool
	}{
		{"🌳", "register", true},
		{"🔰", "first_contribute", challenge_is_complete(db, config.IP, "challenge_first_contribute")},
		{"📝", "tenth_contribute", challenge_is_complete(db, config.IP, "challenge_tenth_contribute")},
		{"🖊️", "hundredth_contribute", challenge_is_complete(db, config.IP, "challenge_hundredth_contribute")},
		{"🏅", "thousandth_contribute", challenge_is_complete(db, config.IP, "challenge_thousandth_contribute")},
		{"💬", "first_discussion", challenge_is_complete(db, config.IP, "challenge_first_discussion")},
		{"💡", "tenth_discussion", challenge_is_complete(db, config.IP, "challenge_tenth_discussion")},
		{"📢", "hundredth_discussion", challenge_is_complete(db, config.IP, "challenge_hundredth_discussion")},
		{"📜", "thousandth_discussion", challenge_is_complete(db, config.IP, "challenge_thousandth_discussion")},
		{"☑️", "admin", challenge_is_complete(db, config.IP, "challenge_admin")},
	}

	green_html := ""
	red_html := ""
	for _, challenge := range challenge_list {
		design := challenge_design(
			challenge.image,
			tool.Get_language(db, "challenge_title_"+challenge.name, true),
			tool.Get_language(db, "challenge_info_"+challenge.name, true),
			challenge.complete,
		)
		if challenge.complete {
			green_html += design
		} else {
			red_html += design
		}
	}

	body := green_html + red_html + `<form method="post">
		<div id="opennamu_get_user_info">` + tool.HTML_escape(config.IP) + `</div>
		<hr class="main_hr">
		<button id="opennamu_save_button" type="submit">` + tool.Get_language(db, "reload", true) + `</button>
	</form>`
	return user_form_page(db, config, tool.Get_language(db, "challenge_and_level_manage", true), body)
}
