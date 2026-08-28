package route

import (
	"net"
	"strings"
	"time"

	"opennamu/route/tool"

	"github.com/dlclark/regexp2"
)

func Api_give_auth_patch(config tool.Config, auth string, change_auth string, user_name string, end_date string, target_type string, reason string, release bool) map[string]any {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	new_data := make(map[string]any)

	end_date = strings.TrimSpace(end_date)
	if end_date != "" {
		end_time, err := time.Parse("2006-01-02", end_date)
		if err != nil {
			end_time, err = time.Parse("2006-01-02 15:04:05", end_date)
		}
		if err != nil {
			new_data["response"] = "error"
			new_data["data"] = "invalid end date"
			return new_data
		}
		end_date = end_time.Format("2006-01-02 15:04:05")
	}

	ip := config.IP
	before_auth := auth
	band := ""
	switch target_type {
	case "", "normal":
		target_type = "normal"
	case "regex":
		band = "regex"
		if user_name != "" {
			if _, err := regexp2.Compile(user_name, 0); err != nil {
				new_data["response"] = "error"
				new_data["data"] = "invalid regex"
				return new_data
			}
		}
	case "cidr":
		band = "cidr"
		if user_name != "" {
			if _, _, err := net.ParseCIDR(user_name); err != nil {
				new_data["response"] = "error"
				new_data["data"] = "invalid cidr"
				return new_data
			}
		}
	case "private":
		band = "private"
		if !tool.Check_acl(db, "", "", "owner_auth", ip) {
			new_data["response"] = "require auth"
			return new_data
		}
	default:
		new_data["response"] = "error"
		new_data["data"] = "invalid target type"
		return new_data
	}

	required_auth := "give_auth"
	if band == "regex" || band == "cidr" {
		required_auth = "give_range_auth"
	}
	if band == "private" {
		required_auth = "owner_auth"
	}
	if !tool.Check_acl(db, "", "", required_auth, ip) {
		new_data["response"] = "require auth"
		return new_data
	}

	if user_name != "" {
		before_auth = tool.Get_auth_target_group(db, user_name, target_type)
	}

	can_change := false
	if release {
		can_change = user_name != "" && tool.Auth_can_change_auth(db, ip, before_auth, "ip")
	} else {
		can_change = tool.Auth_can_change_auth(db, ip, before_auth, change_auth)
	}
	if !can_change {
		new_data["response"] = "require auth"
		return new_data
	}

	if user_name != "" {
		tool.Do_auth_insert(db, user_name, end_date, reason, change_auth, ip, band, release)
		action := change_auth
		if release {
			action = "release"
		}
		tool.Do_insert_auth_history(db, ip, "give_auth ("+user_name+") -> "+action)
	} else {
		if auth != change_auth {
			type auth_target_data struct {
				user_name string
				end_date  string
			}

			target_list := []auth_target_data{}
			rows := tool.Query_DB(
				db,
				"select a.id, coalesce((select e.data from user_set as e where e.id = a.id and e.name = 'acl_end' limit 1), '') from user_set as a where a.name = 'acl' and a.data = ?",
				auth,
			)
			for rows.Next() {
				data := auth_target_data{}
				if err := rows.Scan(&data.user_name, &data.end_date); err == nil {
					target_list = append(target_list, data)
				}
			}
			rows.Close()

			for _, data := range target_list {
				tool.Do_auth_insert(db, data.user_name, data.end_date, "", change_auth, ip, "", false)
			}
		}
		tool.Do_insert_auth_history(db, ip, "give_auth ("+auth+") -> "+change_auth)
	}

	new_data["response"] = "ok"
	return new_data
}
