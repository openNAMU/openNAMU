package route

import (
	"net/url"
	"opennamu/route/tool"
	"strconv"
)

func View_challenge(config tool.Config, values url.Values) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)
	if !user_auth(db, config) {
		return tool.Get_redirect("/user")
	}

	if values != nil {
		edit_count := challenge_count(db, "select count(*) from history where ip = ?", config.IP)
		topic_count := challenge_count(db, "select count(*) from topic where ip = ?", config.IP)
		experience := 5 * edit_count

		for _, challenge := range []struct {
			name   string
			count  int
			reward int
		}{
			{"challenge_first_contribute", 1, 500},
			{"challenge_tenth_contribute", 10, 1000},
			{"challenge_hundredth_contribute", 100, 3000},
			{"challenge_thousandth_contribute", 1000, 10000},
		} {
			if edit_count >= challenge.count {
				user_save(db, config.IP, challenge.name, "1")
				experience += challenge.reward
			}
		}

		experience += 5 * topic_count
		for _, challenge := range []struct {
			name   string
			count  int
			reward int
		}{
			{"challenge_first_discussion", 1, 500},
			{"challenge_tenth_discussion", 10, 1000},
			{"challenge_hundredth_discussion", 100, 3000},
			{"challenge_thousandth_discussion", 1000, 10000},
		} {
			if topic_count >= challenge.count {
				user_save(db, config.IP, challenge.name, "1")
				experience += challenge.reward
			}
		}

		if tool.Check_acl(db, "", "", "all_admin_auth", config.IP) || challenge_is_complete(db, config.IP, "challenge_admin") {
			user_save(db, config.IP, "challenge_admin", "1")
			experience += 10000
		}

		level := 0
		for experience >= 500+level*50 {
			experience -= 500 + level*50
			level++
		}
		user_save(db, config.IP, "level", strconv.Itoa(level))
		user_save(db, config.IP, "experience", strconv.Itoa(experience))
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
