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

const server_update_remote = "https://github.com/opennamu/opennamu.git"

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

func find_server_repo_root() (string, error) {
	current, err := os.Getwd()
	if err != nil {
		return "", err
	}

	for {
		if _, err := os.Stat(filepath.Join(current, ".git")); err == nil {
			return current, nil
		}
		parent := filepath.Dir(current)
		if parent == current {
			break
		}
		current = parent
	}

	return "", fmt.Errorf("repository root not found")
}

func run_server_git(repo_root string, arguments ...string) ([]byte, error) {
	git_path, err := exec.LookPath("git")
	if err != nil {
		return nil, err
	}

	command := exec.Command(git_path, arguments...)
	command.Dir = repo_root
	output, err := command.CombinedOutput()
	if err != nil {
		return output, fmt.Errorf("git %s: %w: %s", strings.Join(arguments, " "), err, strings.TrimSpace(string(output)))
	}
	return output, nil
}

func clone_server_update(repo_root string, branch string) (string, error) {
	clone_root, err := os.MkdirTemp("", "opennamu-update-source-*")
	if err != nil {
		return "", err
	}

	_, err = run_server_git(
		repo_root,
		"clone",
		"--depth=1",
		"--no-tags",
		"--branch",
		branch,
		server_update_remote,
		clone_root,
	)
	if err != nil {
		os.RemoveAll(clone_root)
		return "", err
	}
	return clone_root, nil
}

func build_server_update(repo_root string) (string, error) {
	go_path, err := exec.LookPath("go")
	if err != nil {
		return "", err
	}

	temporary, err := os.CreateTemp("", "opennamu-update-*")
	if err != nil {
		return "", err
	}
	temporary_path := temporary.Name()
	if err := temporary.Close(); err != nil {
		os.Remove(temporary_path)
		return "", err
	}
	if err := os.Remove(temporary_path); err != nil {
		return "", err
	}

	command := exec.Command(go_path, "build", "-o", temporary_path, ".")
	command.Dir = repo_root
	output, err := command.CombinedOutput()
	if err != nil {
		os.Remove(temporary_path)
		return "", fmt.Errorf("go build: %w: %s", err, strings.TrimSpace(string(output)))
	}
	return temporary_path, nil
}

func update_server_source(db *sql.DB) (string, error) {
	repo_root, err := find_server_repo_root()
	if err != nil {
		return "", err
	}

	status, err := run_server_git(repo_root, "status", "--porcelain", "--untracked-files=all")
	if err != nil {
		return "", err
	}
	if strings.TrimSpace(string(status)) != "" {
		return "", fmt.Errorf("working tree is not clean")
	}

	branch := ""
	tool.QueryRow_DB(db, `select data from other where name = "update"`, []any{&branch})
	if branch != "stable" && branch != "beta" && branch != "dev" {
		branch = "stable"
	}

	clone_root, err := clone_server_update(repo_root, branch)
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(clone_root)

	executable, err := build_server_update(clone_root)
	if err != nil {
		return "", err
	}

	if _, err := run_server_git(repo_root, "fetch", "--depth=1", "--no-tags", server_update_remote, branch); err != nil {
		os.Remove(executable)
		return "", err
	}
	if _, err := run_server_git(repo_root, "merge", "--ff-only", "FETCH_HEAD"); err != nil {
		os.Remove(executable)
		return "", err
	}

	return executable, nil
}

func server_action_error(db *sql.DB, config tool.Config, action string, err error) string {
	log.Printf("server %s failed: %v", action, err)
	return tool.Get_error_page(db, config, "error")
}

func View_server_action(config tool.Config, action string, post bool) string {
	db := tool.DB_connect()
	defer tool.DB_close(db)

	if !tool.Check_acl(db, "", "", "owner_auth", config.IP) {
		return tool.Get_error_page(db, config, "auth")
	}
	if action != "restart" && action != "shutdown" && action != "update" {
		return tool.Get_error_page(db, config, "error")
	}

	if post {
		if !begin_server_action() {
			return tool.Get_error_page(db, config, "error")
		}
		tool.Do_insert_auth_history(db, config.IP, "server_"+action)

		if action == "shutdown" {
			schedule_server_exit()
			return tool.Get_language(db, "wiki_shutdown", true)
		}

		executable := ""
		var err error
		if action == "update" {
			executable, err = update_server_source(db)
		} else {
			executable, err = current_server_executable()
		}
		if err == nil {
			err = start_server_process(executable)
		}
		if err != nil {
			if action == "update" && executable != "" {
				os.Remove(executable)
			}
			cancel_server_action()
			return server_action_error(db, config, action, err)
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
