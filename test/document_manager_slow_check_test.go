package test

import (
	"testing"
	"time"

	"opennamu/route/tool"
)

func Test_Do_edit_slow_check_modes(t *testing.T) {
	for _, do_type := range []string{"edit", "thread"} {
		t.Run(do_type, func(t *testing.T) {
			db := apiSlowDB(t)
			defer tool.DB_close(db)

			const user_id = "slow_check_user"
			seedBbsNonExempt(db, t, user_id)
			setting_name := "slow_thread"
			if do_type == "edit" {
				setting_name = "slow_edit"
			}
			if _, err := db.Exec(
				`insert into other (name, data, coverage) values (?, '60', '')`,
				setting_name,
			); err != nil {
				t.Fatalf("seed %s setting: %v", setting_name, err)
			}

			recent_date := tool.Get_time()
			if do_type == "edit" {
				if _, err := db.Exec(
					`insert into history (id, title, data, date, ip, send, leng, hide, type) values ('1', 'slow_check', '', ?, ?, '', '0', '', '')`,
					recent_date,
					user_id,
				); err != nil {
					t.Fatalf("seed edit history: %v", err)
				}
			} else {
				for _, value := range []struct {
					name string
					data string
				}{
					{name: "comment_date", data: recent_date},
					{name: "comment_user_id", data: user_id},
				} {
					if _, err := db.Exec(
						`insert into bbs_data (set_name, set_code, set_id, set_data) values (?, '1', '-1-1', ?)`,
						value.name,
						value.data,
					); err != nil {
						t.Fatalf("seed thread %s: %v", value.name, err)
					}
				}
			}

			config := tool.Config{IP: user_id}
			if tool.Do_edit_slow_check(db, config, do_type) {
				t.Fatal("recent write should be blocked by the 60-second limit")
			}

			old_date := time.Now().Add(-2 * time.Minute).Format("2006-01-02 15:04:05")
			if do_type == "edit" {
				if _, err := db.Exec(`update history set date = ? where ip = ?`, old_date, user_id); err != nil {
					t.Fatalf("age edit history: %v", err)
				}
			} else {
				if _, err := db.Exec(
					`update bbs_data set set_data = ? where set_name = 'comment_date' and set_id = '-1-1'`,
					old_date,
				); err != nil {
					t.Fatalf("age thread history: %v", err)
				}
			}
			if !tool.Do_edit_slow_check(db, config, do_type) {
				t.Fatal("write older than the limit should be allowed")
			}
		})
	}
}
