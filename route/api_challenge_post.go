package route

import (
	"strconv"

	"opennamu/route/tool"
)

func Api_challenge_post(config tool.Config) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	return_data := make(map[string]any)
	if !user_auth(db, config) {
		return_data["response"] = "require auth"
		return return_data
	}

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

	if tool.Check_permission(db, "treat_as_admin", config.IP) || challenge_is_complete(db, config.IP, "challenge_admin") {
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

	return_data["response"] = "ok"
	return return_data
}
