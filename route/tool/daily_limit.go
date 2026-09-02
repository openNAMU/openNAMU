package tool

import (
	"database/sql"
	"time"
)

func Check_edit_time(db *sql.DB, ip string) bool {
	auth_info := Get_auth_info(db, ip)
	day := auth_info["edit_day"]
	night := auth_info["edit_night"]
	if !day && !night {
		return true
	}

	hour := time.Now().Hour()
	if hour >= 8 && hour < 22 {
		return day
	}
	return night
}

func Get_daily_limit(auth_info map[string]bool, name string) int {
	for _, value := range []struct {
		name  string
		limit int
	}{
		{"unlimited", 0},
		{"100", 100},
		{"50", 50},
		{"10", 10},
	} {
		if auth_info[name+"_limit_"+value.name] {
			return value.limit
		}
	}
	return 50
}

func Check_daily_limit(db *sql.DB, ip string, name string) bool {
	limit := Get_daily_limit(Get_auth_info(db, ip), name)
	if limit == 0 {
		return true
	}

	query := ""
	switch name {
	case "edit":
		query = "select count(*) from history where ip = ? and date like ?"
	case "bbs_edit":
		query = `select count(*) from bbs_data user_data inner join bbs_data date_data on date_data.set_name = 'date' and date_data.set_id = user_data.set_id and date_data.set_code = user_data.set_code where user_data.set_name = 'user_id' and user_data.set_data = ? and date_data.set_data like ?`
	case "bbs_comment":
		query = `select count(*) from bbs_data user_data inner join bbs_data date_data on date_data.set_name = 'comment_date' and date_data.set_id = user_data.set_id and date_data.set_code = user_data.set_code where user_data.set_name = 'comment_user_id' and user_data.set_data = ? and date_data.set_data like ?`
	default:
		return true
	}

	count := 0
	QueryRow_DB(db, query, []any{&count}, ip, Get_date()+"%")
	return count < limit
}
