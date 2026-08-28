package route

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"opennamu/route/tool"
)

var server_action_mutex sync.Mutex
var server_action_started bool

func begin_server_action() bool {
	server_action_mutex.Lock()
	defer server_action_mutex.Unlock()

	if server_action_started {
		return false
	}
	server_action_started = true
	return true
}

func cancel_server_action() {
	server_action_mutex.Lock()
	server_action_started = false
	server_action_mutex.Unlock()
}

func schedule_server_exit() {
	time.AfterFunc(500*time.Millisecond, func() {
		os.Exit(0)
	})
}

func server_process_environment() []string {
	environment := []string{}
	for _, value := range os.Environ() {
		if strings.HasPrefix(value, "NAMU_START_DELAY_MS=") {
			continue
		}
		environment = append(environment, value)
	}
	return append(environment, "NAMU_START_DELAY_MS=1000")
}

func start_server_process(executable string) error {
	working_dir, err := os.Getwd()
	if err != nil {
		return err
	}

	command := exec.Command(executable, os.Args[1:]...)
	command.Dir = working_dir
	command.Env = server_process_environment()
	detach_server_process(command)
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	if err := command.Start(); err != nil {
		return fmt.Errorf("start server process: %w", err)
	}
	return nil
}

func current_server_executable() (string, error) {
	executable, err := os.Executable()
	if err != nil {
		return "", err
	}
	return filepath.Abs(executable)
}

func server_action_error(db *sql.DB, config tool.Config, action string, err error) string {
	log.Printf("server %s failed: %v", action, err)
	return tool.Get_error_page(db, config, "error")
}

func View_server_action(config tool.Config, action string, post bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !post && !tool.Check_permission(db, "server_action", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if action != "restart" && action != "shutdown" && action != "update" {
		return tool.Get_error_page(db, config, "error")
	}

	if post {
		result := Api_server_action_post(config, action)
		if result["response"] == "require auth" {
			return tool.Get_error_page(db, config, "auth")
		}
		if result["response"] != "ok" {
			return server_action_error(db, config, action, fmt.Errorf("server action rejected"))
		}
		if !begin_server_action() {
			return tool.Get_error_page(db, config, "error")
		}
		if action == "shutdown" {
			schedule_server_exit()
			return tool.Get_language(db, "wiki_shutdown", true)
		}

		if action == "update" {
			err := start_server_update(get_version_branch(db))
			if err != nil {
				cancel_server_action()
				return server_action_error(db, config, action, err)
			}
		} else {
			executable, err := current_server_executable()
			if err != nil {
				cancel_server_action()
				return server_action_error(db, config, action, err)
			}
			err = start_server_process(executable)
			if err != nil {
				cancel_server_action()
				return server_action_error(db, config, action, err)
			}
		}

		schedule_server_exit()
		return tool.Get_language(db, "warning_restart", true)
	}

	title := action
	button := action
	warning := ""
	switch action {
	case "restart":
		title = tool.Get_language(db, "wiki_restart", true)
		button = tool.Get_language(db, "restart", true)
	case "shutdown":
		title = tool.Get_language(db, "wiki_shutdown", true)
		button = tool.Get_language(db, "shutdown", true)
	case "update":
		title = tool.Get_language(db, "update", true)
		button = tool.Get_language(db, "update", true)
		warning = tool.Get_language(db, "update_warning", true) + `<hr class="main_hr">`
	}
	version_data := ""
	if action == "update" {
		version_list := tool.Get_last_version()
		latest_version := get_remote_version(get_version_branch(db))
		version_data = `<ul><li>` + tool.Get_language(db, "version", true) + ` : ` + tool.HTML_escape(version_list["r_ver"]) + `</li><li>` + tool.Get_language(db, "lastest", true) + ` : ` + tool.HTML_escape(latest_version) + `</li></ul><hr class="main_hr">`
	}
	data := warning + version_data + `<form method="post"><button type="submit">` + button + `</button></form>`
	return tool.Get_template(db, config, title, data, []any{}, [][]any{{"manager", tool.Get_language(db, "return", true)}}, map[string]string{})
}
