package route

import (
	"bufio"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"

	"opennamu/route/tool"
	"opennamu/route/tool/markup"
)

func emergency_print_menu() {
	fmt.Println("1. Backlink reset")
	fmt.Println("2. reCAPTCHA delete")
	fmt.Println("3. Ban delete")
	fmt.Println("4. Change host")
	fmt.Println("5. Change port")
	fmt.Println("6. Change skin")
	fmt.Println("7. Change password")
	fmt.Println("8. Set db version")
	fmt.Println("9. Delete set.json")
	fmt.Println("10. Change name")
	fmt.Println("11. Delete mysql.json")
	fmt.Println("14. Delete Main <HEAD>")
	fmt.Println("15. Give owner")
	fmt.Println("16. Delete 2FA password")
	fmt.Println("17. Change markup")
	fmt.Println("18. Change wiki access password")
	fmt.Println("19. Forced binary update")
	fmt.Println("20. Change domain")
	fmt.Println("21. Change TLS")
	fmt.Println("22. Delete body top")
	fmt.Println("23. Delete body bottom")
	fmt.Println("24. SQLite to MySQL (not supported)")
	fmt.Println("25. Recalc exist data_set")
	fmt.Println("26. Change update branch")
	fmt.Println("27. Change golang port")
}

func emergency_input(reader *bufio.Reader, message string) string {
	fmt.Print(message)
	data, _ := reader.ReadString('\n')
	return strings.TrimSpace(data)
}

func emergency_open_db() (db *sql.DB, err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("database open failed: %v", value)
		}
	}()

	tool.DB_boot()
	db = tool.DB_connect()
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, err
	}

	return db, nil
}

func emergency_exec(db *sql.DB, query string, values ...any) error {
	_, err := db.Exec(tool.DB_change(query), values...)
	return err
}

func emergency_set_other(db *sql.DB, name string, data string) error {
	if err := emergency_exec(db, "update other set data = ? where name = ?", data, name); err != nil {
		return err
	}

	var exists int
	err := db.QueryRow(
		tool.DB_change("select 1 from other where name = ? limit 1"),
		name,
	).Scan(&exists)
	if err == nil {
		return nil
	}
	if !errors.Is(err, sql.ErrNoRows) {
		return err
	}

	return emergency_exec(
		db,
		"insert into other (name, data, coverage) values (?, ?, '')",
		name,
		data,
	)
}

func emergency_delete_file(path string) error {
	err := os.Remove(path)
	if errors.Is(err, os.ErrNotExist) {
		return nil
	}
	return err
}

func emergency_reset_backlink(db *sql.DB) error {
	if err := emergency_exec(db, "delete from back"); err != nil {
		return err
	}
	if err := emergency_exec(db, "delete from data_set where set_name = 'link_count'"); err != nil {
		return err
	}

	rows, err := db.Query(
		tool.DB_change("select title, coalesce(data, '') from data order by title"),
	)
	if err != nil {
		return err
	}

	documents := make([][]string, 0)
	for rows.Next() {
		var title string
		var data string
		if err := rows.Scan(&title, &data); err != nil {
			rows.Close()
			return err
		}
		documents = append(documents, []string{title, data})
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	error_count := 0
	for index, document := range documents {
		if err := emergency_render_backlink(db, document[0], document[1]); err != nil {
			error_count++
			fmt.Println("backlink error:", document[0], err)
		}

		if (index+1)%100 == 0 {
			fmt.Printf("Backlink reset: %d/%d\n", index+1, len(documents))
		}
	}

	fmt.Printf("Backlink reset: %d documents", len(documents))
	if error_count > 0 {
		fmt.Printf(" (%d errors)", error_count)
	}
	fmt.Println()
	return nil
}

func emergency_render_backlink(db *sql.DB, title string, data string) (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = fmt.Errorf("%v", value)
		}
	}()

	markup.Get_render(db, title, data, "backlink")
	return nil
}

func emergency_recalc_data_set(db *sql.DB) error {
	rows, err := db.Query(
		tool.DB_change("select distinct doc_name from data_set where doc_rev = 'not_exist' or doc_rev = ''"),
	)
	if err != nil {
		return err
	}

	document_names := make([]string, 0)
	for rows.Next() {
		var document_name string
		if err := rows.Scan(&document_name); err != nil {
			rows.Close()
			return err
		}
		document_names = append(document_names, document_name)
	}
	if err := rows.Err(); err != nil {
		rows.Close()
		return err
	}
	rows.Close()

	for _, document_name := range document_names {
		var exists string
		err := db.QueryRow(
			tool.DB_change("select title from data where title = ? limit 1"),
			document_name,
		).Scan(&exists)
		if errors.Is(err, sql.ErrNoRows) {
			exists = "not_exist"
		} else if err != nil {
			return err
		} else {
			exists = ""
		}

		if err := emergency_exec(
			db,
			"update data_set set doc_rev = ? where doc_name = ? and (doc_rev = '' or doc_rev = 'not_exist')",
			exists,
			document_name,
		); err != nil {
			return err
		}
	}

	fmt.Printf("data_set recalculated: %d documents\n", len(document_names))
	return nil
}

func emergency_update_binary(branch string) error {
	if branch != "stable" && branch != "beta" {
		return fmt.Errorf("unsupported update branch: %s", branch)
	}

	executable, err := current_server_executable()
	if err != nil {
		return err
	}

	binary_url, err := get_server_update_url(branch)
	if err != nil {
		return err
	}

	temporary, err := os.CreateTemp(filepath.Dir(executable), ".opennamu-emergency-update-*")
	if err != nil {
		return err
	}
	temporary_path := temporary.Name()
	if err := temporary.Close(); err != nil {
		os.Remove(temporary_path)
		return err
	}

	if err := download_server_update(binary_url, temporary_path); err != nil {
		os.Remove(temporary_path)
		return err
	}

	if runtime.GOOS == "windows" {
		command := exec.Command(executable, server_update_mode, executable, temporary_path)
		detach_server_process(command)
		command.Stdout = os.Stdout
		command.Stderr = os.Stderr
		if err := command.Start(); err != nil {
			os.Remove(temporary_path)
			return err
		}
		fmt.Println("Update is scheduled. Close this process to replace the binary.")
		return nil
	}

	if err := os.Rename(temporary_path, executable); err != nil {
		os.Remove(temporary_path)
		return err
	}

	fmt.Println("Binary updated. Start openNAMU again.")
	return nil
}

func Run_emergency_tool(arguments []string) int {
	reader := bufio.NewReader(os.Stdin)
	choice := ""
	if len(arguments) > 0 {
		choice = arguments[0]
	} else {
		emergency_print_menu()
		choice = emergency_input(reader, "Insert selection number (EX : 9) : ")
	}

	if choice == "9" {
		if err := emergency_delete_file(filepath.Join("data", "set.json")); err != nil {
			fmt.Fprintln(os.Stderr, "Emergency tool failed:", err)
			return 1
		}
		fmt.Println("OK")
		return 0
	}
	if choice == "11" {
		if err := emergency_delete_file(filepath.Join("data", "mysql.json")); err != nil {
			fmt.Fprintln(os.Stderr, "Emergency tool failed:", err)
			return 1
		}
		fmt.Println("OK")
		return 0
	}
	if choice == "19" {
		branch := emergency_input(reader, "Branch (stable/beta) : ")
		if branch == "" {
			branch = "stable"
		}
		if err := emergency_update_binary(branch); err != nil {
			fmt.Fprintln(os.Stderr, "Emergency tool failed:", err)
			return 1
		}
		return 0
	}
	if choice == "24" {
		fmt.Fprintln(os.Stderr, "SQLite to MySQL conversion is not supported in the Go emergency tool.")
		return 1
	}

	db, err := emergency_open_db()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Emergency tool failed:", err)
		return 1
	}
	defer tool.DB_close(db)

	switch choice {
	case "1":
		err = emergency_reset_backlink(db)
	case "2":
		err = emergency_exec(db, "delete from other where name = 'recaptcha'")
		if err == nil {
			err = emergency_exec(db, "delete from other where name = 'sec_re'")
		}
	case "3":
		user_data := emergency_input(reader, "IP or Name : ")
		err = emergency_exec(
			db,
			"insert into rb (block, end, today, blocker, why, band, ongoing, login) values (?, 'release', ?, 'tool:emergency', '', '', '', '')",
			user_data,
			tool.Get_time(),
		)
		if err == nil {
			err = emergency_exec(db, "update rb set ongoing = '' where block = ?", user_data)
		}
	case "4":
		err = emergency_set_other(db, "host", emergency_input(reader, "Host : "))
	case "5":
		port := emergency_input(reader, "Port : ")
		port_num, port_err := strconv.Atoi(port)
		if port_err != nil || port_num < 1 || port_num > 65535 {
			err = fmt.Errorf("invalid port: %s", port)
		} else {
			err = emergency_set_other(db, "port", port)
		}
	case "6":
		err = emergency_set_other(db, "skin", emergency_input(reader, "Skin name : "))
	case "7":
		user_name := emergency_input(reader, "User name : ")
		user_password := emergency_input(reader, "User password : ")
		password_encode := tool.Get_user_encode(db, user_name)
		password := tool.Password_encode(db, user_password, password_encode)
		err = emergency_exec(
			db,
			"update user_set set data = ? where id = ? and name = 'pw'",
			password,
			user_name,
		)
	case "8":
		version := emergency_input(reader, "Insert version (0000000) : ")
		if version == "" {
			version = "0000000"
		}
		err = emergency_set_other(db, "ver", version)
	case "10":
		user_name := emergency_input(reader, "User name : ")
		new_name := emergency_input(reader, "New name : ")
		err = emergency_exec(db, "update user_set set id = ? where id = ?", new_name, user_name)
	case "14":
		err = emergency_exec(db, "delete from other where name = 'head'")
	case "15":
		user_name := emergency_input(reader, "User name : ")
		var user_exists string
		err = db.QueryRow(
			tool.DB_change("select id from user_set where id = ? limit 1"),
			user_name,
		).Scan(&user_exists)
		if errors.Is(err, sql.ErrNoRows) {
			err = fmt.Errorf("user not found: %s", user_name)
		}
		if err == nil {
			err = emergency_exec(db, "update user_set set data = 'owner' where id = ? and name = 'acl'", user_name)
		}
		if err == nil {
			var acl_exists int
			acl_err := db.QueryRow(
				tool.DB_change("select 1 from user_set where id = ? and name = 'acl' limit 1"),
				user_name,
			).Scan(&acl_exists)
			if errors.Is(acl_err, sql.ErrNoRows) {
				err = emergency_exec(db, "insert into user_set (name, id, data) values ('acl', ?, 'owner')", user_name)
			} else {
				err = acl_err
			}
		}
	case "16":
		user_name := emergency_input(reader, "User name : ")
		err = emergency_exec(db, "update user_set set data = '' where name = '2fa' and id = ?", user_name)
	case "17":
		err = emergency_set_other(db, "markup", emergency_input(reader, "Markup name : "))
	case "18":
		err = emergency_set_other(db, "wiki_access_password", emergency_input(reader, "Password : "))
	case "20":
		err = emergency_set_other(db, "domain", emergency_input(reader, "Domain (EX : 2du.pythonanywhere.com) : "))
	case "21":
		tls_value := emergency_input(reader, "TLS (http) [http, https] : ")
		if tls_value != "https" {
			tls_value = "http"
		}
		err = emergency_set_other(db, "http_select", tls_value)
	case "22":
		err = emergency_exec(db, "delete from other where name = 'body'")
	case "23":
		err = emergency_exec(db, "delete from other where name = 'bottom_body'")
	case "25":
		err = emergency_recalc_data_set(db)
	case "26":
		branch := emergency_input(reader, "Insert branch name (beta) [stable, beta] : ")
		if branch != "stable" && branch != "beta" {
			branch = "beta"
		}
		err = emergency_set_other(db, "update", branch)
	case "27":
		port := emergency_input(reader, "Port : ")
		port_num, port_err := strconv.Atoi(port)
		if port_err != nil || port_num < 1 || port_num > 65535 {
			err = fmt.Errorf("invalid port: %s", port)
		} else {
			err = emergency_set_other(db, "golang_port", port)
		}
	default:
		err = fmt.Errorf("unknown selection: %s", choice)
	}

	if err != nil {
		fmt.Fprintln(os.Stderr, "Emergency tool failed:", err)
		return 1
	}

	fmt.Println("OK")
	return 0
}
