package test

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
	"opennamu/route/tool"
)

const search_index_recovery_child_env = "OPENNAMU_SEARCH_INDEX_RECOVERY_CHILD"
const search_index_recovery_scenario_env = "OPENNAMU_SEARCH_INDEX_RECOVERY_SCENARIO"
const search_index_recovery_backend_env = "OPENNAMU_SEARCH_INDEX_RECOVERY_BACKEND"
const search_index_recovery_mysql_config_env = "OPENNAMU_SEARCH_INDEX_RECOVERY_MYSQL_CONFIG"

var search_index_recovery_scenarios = []string{
	"startup_missing",
	"version_markers",
	"metadata_corruption",
	"abandoned_staging",
	"staging_next_to_valid",
	"interrupted_install",
	"parent_retry",
	"create_retry",
	"staging_retry",
	"permission_retry",
	"storage_target_unavailable",
	"source_read_error",
	"install_failure",
	"runtime_directory_loss",
	"runtime_handle_error",
	"concurrent_changes",
}

func Test_Search_index_recovery_scenarios(t *testing.T) {
	for _, backend := range []string{"sqlite", "turso"} {
		for _, scenario := range search_index_recovery_scenarios {
			t.Run(backend+"/"+scenario, func(t *testing.T) {
				search_index_recovery_run_child(t, backend, scenario, "")
			})
		}
	}

	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		t.Run("mysql/startup_and_corrupt_metadata", func(t *testing.T) {
			mysql_config, _ := search_index_recovery_mysql_setup(t)
			search_index_recovery_run_child(t, "mysql", "metadata_corruption", mysql_config)
		})
	}
}

func search_index_recovery_run_child(t *testing.T, backend string, scenario string, mysql_config string) {
	t.Helper()
	root, err := os.MkdirTemp(filepath.Join("..", "tmp"), "search_index_recovery_")
	if err != nil {
		t.Fatalf("create isolated recovery root: %v", err)
	}
	t.Cleanup(func() {
		if err := os.RemoveAll(root); err != nil {
			t.Errorf("remove isolated recovery root: %v", err)
		}
	})
	root, err = filepath.Abs(root)
	if err != nil {
		t.Fatalf("resolve isolated recovery root: %v", err)
	}

	command_context, cancel := context.WithTimeout(context.Background(), 50*time.Second)
	defer cancel()
	command := exec.CommandContext(command_context, os.Args[0], "-test.run=^Test_Search_index_recovery_child$")
	command.Dir = root
	command.Env = search_index_recovery_environment(os.Environ(), search_index_recovery_child_env, "1")
	command.Env = search_index_recovery_environment(command.Env, search_index_recovery_scenario_env, scenario)
	command.Env = search_index_recovery_environment(command.Env, search_index_recovery_backend_env, backend)
	command.Env = search_index_recovery_environment(command.Env, "NAMU_DB", filepath.Join(root, "source_database"))
	command.Env = search_index_recovery_environment(command.Env, "NAMU_DB_TYPE", backend)
	if mysql_config != "" {
		command.Env = search_index_recovery_environment(command.Env, search_index_recovery_mysql_config_env, mysql_config)
	}
	output, err := command.CombinedOutput()
	if command_context.Err() != nil {
		t.Fatalf("search-index recovery child timed out (%s/%s): %s\n%s", backend, scenario, command_context.Err(), output)
	}
	if err != nil {
		t.Fatalf("search-index recovery child failed (%s/%s): %v\n%s", backend, scenario, err, output)
	}
}

func search_index_recovery_environment(environment []string, key string, value string) []string {
	prefix := key + "="
	filtered := make([]string, 0, len(environment)+1)
	for _, entry := range environment {
		if !strings.HasPrefix(entry, prefix) {
			filtered = append(filtered, entry)
		}
	}
	return append(filtered, prefix+value)
}

func Test_Search_index_recovery_child(t *testing.T) {
	if os.Getenv(search_index_recovery_child_env) != "1" {
		t.Skip("search-index recovery subprocess helper")
	}

	log_buffer := &search_index_recovery_log_buffer{}
	log.SetOutput(log_buffer)
	defer log.SetOutput(os.Stderr)

	backend := os.Getenv(search_index_recovery_backend_env)
	scenario := os.Getenv(search_index_recovery_scenario_env)
	if backend == "mysql" {
		if err := os.MkdirAll("data", 0o755); err != nil {
			t.Fatalf("create isolated MySQL config directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join("data", "mysql.json"), []byte(os.Getenv(search_index_recovery_mysql_config_env)), 0o600); err != nil {
			t.Fatalf("write isolated MySQL config: %v", err)
		}
	}
	tool.DB_boot()
	tool.DB_init()
	db := tool.DB_connect()
	defer tool.DB_close(db)
	search_index_recovery_seed_database(t, db)

	document_index_directory := filepath.Join("data", "bleve")
	document_version_file := filepath.Join("data", "bleve.version")
	bbs_index_directory := filepath.Join("data", "bleve_bbs")
	bbs_version_file := filepath.Join("data", "bleve_bbs.version")

	switch scenario {
	case "startup_missing":
	case "version_markers":
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		if err := os.WriteFile(bbs_version_file, []byte("obsolete"), 0o600); err != nil {
			t.Fatalf("write obsolete BBS marker: %v", err)
		}
	case "metadata_corruption":
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		if err := os.Remove(filepath.Join(document_index_directory, "index_meta.json")); err != nil {
			t.Fatalf("remove document metadata: %v", err)
		}
		if err := os.WriteFile(filepath.Join(bbs_index_directory, "index_meta.json"), []byte("corrupt metadata"), 0o600); err != nil {
			t.Fatalf("corrupt BBS metadata: %v", err)
		}
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("write document marker: %v", err)
		}
		if err := os.WriteFile(bbs_version_file, []byte("4"), 0o600); err != nil {
			t.Fatalf("write BBS marker: %v", err)
		}
	case "abandoned_staging":
		search_index_recovery_create_document_index(t, document_index_directory+".tmp")
		search_index_recovery_create_bbs_index(t, bbs_index_directory+".tmp")
	case "staging_next_to_valid":
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		search_index_recovery_create_document_index(t, document_index_directory+".tmp")
		search_index_recovery_create_bbs_index(t, bbs_index_directory+".tmp")
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("write document marker: %v", err)
		}
		if err := os.WriteFile(bbs_version_file, []byte("4"), 0o600); err != nil {
			t.Fatalf("write BBS marker: %v", err)
		}
	case "interrupted_install":
		search_index_recovery_create_document_index(t, document_index_directory)
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("write document marker: %v", err)
		}
		if err := os.Rename(document_index_directory, document_index_directory+".previous-123"); err != nil {
			t.Fatalf("simulate interrupted index replacement: %v", err)
		}
		if err := os.Rename(document_version_file, document_index_directory+".previous-123.version"); err != nil {
			t.Fatalf("simulate interrupted marker replacement: %v", err)
		}
	case "parent_retry":
		if err := os.WriteFile("data", []byte("temporary parent obstruction"), 0o600); err != nil {
			t.Fatalf("create temporary parent obstruction: %v", err)
		}
	case "create_retry":
		if err := os.Mkdir("data", 0o555); err != nil {
			t.Fatalf("create read-only index parent: %v", err)
		}
	case "staging_retry":
		if err := os.MkdirAll(document_index_directory+".tmp", 0o700); err != nil {
			t.Fatalf("create stale document staging directory: %v", err)
		}
		if err := os.MkdirAll(bbs_index_directory+".tmp", 0o700); err != nil {
			t.Fatalf("create stale BBS staging directory: %v", err)
		}
		t.Cleanup(func() {
			_ = os.Chmod(document_index_directory+".tmp", 0o700)
			_ = os.Chmod(bbs_index_directory+".tmp", 0o700)
		})
		if err := os.WriteFile(filepath.Join(document_index_directory+".tmp", "sentinel"), []byte("stale"), 0o600); err != nil {
			t.Fatalf("write document staging sentinel: %v", err)
		}
		if err := os.Chmod(document_index_directory+".tmp", 0o500); err != nil {
			t.Fatalf("make document staging directory unremovable: %v", err)
		}
		if err := os.WriteFile(filepath.Join(bbs_index_directory+".tmp", "sentinel"), []byte("stale"), 0o600); err != nil {
			t.Fatalf("write BBS staging sentinel: %v", err)
		}
		if err := os.Chmod(bbs_index_directory+".tmp", 0o500); err != nil {
			t.Fatalf("make BBS staging directory unremovable: %v", err)
		}
	case "permission_retry":
		if runtime.GOOS == "windows" {
			t.Skip("directory permission recovery is tested on Unix platforms")
		}
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("write document marker: %v", err)
		}
		if err := os.WriteFile(bbs_version_file, []byte("4"), 0o600); err != nil {
			t.Fatalf("write BBS marker: %v", err)
		}
		if err := os.Chmod(document_index_directory, 0); err != nil {
			t.Fatalf("deny document index access: %v", err)
		}
		if err := os.Chmod(bbs_index_directory, 0); err != nil {
			t.Fatalf("deny BBS index access: %v", err)
		}
	case "storage_target_unavailable":
		if runtime.GOOS == "windows" {
			t.Skip("index target symlink recovery is tested on Unix platforms")
		}
		for _, index := range []struct {
			active  string
			backing string
			version string
		}{
			{document_index_directory, filepath.Join("data", "document_backing"), "2"},
			{bbs_index_directory, filepath.Join("data", "bbs_backing"), "4"},
		} {
			if strings.Contains(index.active, "bleve_bbs") {
				search_index_recovery_create_bbs_index(t, index.backing)
			} else {
				search_index_recovery_create_document_index(t, index.backing)
			}
			if err := os.WriteFile(index.active+".version", []byte(index.version), 0o600); err != nil {
				t.Fatalf("write index marker: %v", err)
			}
			if err := os.Symlink(filepath.Base(index.backing), index.active); err != nil {
				t.Fatalf("link active index to backing directory: %v", err)
			}
			if err := os.Rename(index.backing, index.backing+".unavailable"); err != nil {
				t.Fatalf("simulate temporarily unavailable index target: %v", err)
			}
		}
	case "install_failure":
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		search_index_recovery_create_marker_directory(t, document_version_file)
		search_index_recovery_create_marker_directory(t, bbs_version_file)
	case "source_read_error":
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("write document marker: %v", err)
		}
		if err := os.WriteFile(bbs_version_file, []byte("4"), 0o600); err != nil {
			t.Fatalf("write BBS marker: %v", err)
		}
	case "runtime_directory_loss", "runtime_handle_error":
		search_index_recovery_create_document_index(t, document_index_directory)
		search_index_recovery_create_bbs_index(t, bbs_index_directory)
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("write document marker: %v", err)
		}
		if err := os.WriteFile(bbs_version_file, []byte("4"), 0o600); err != nil {
			t.Fatalf("write BBS marker: %v", err)
		}
	case "concurrent_changes":
	default:
		t.Fatalf("unknown recovery scenario %q", scenario)
	}

	tool.Search_index_start()
	tool.Search_bbs_index_start()

	switch scenario {
	case "parent_retry":
		search_index_recovery_wait_for_log(t, log_buffer, "document index recovery failed", "BBS index recovery failed")
		var wait_group sync.WaitGroup
		for worker := 0; worker < 12; worker++ {
			wait_group.Add(1)
			go func() {
				defer wait_group.Done()
				for request := 0; request < 20; request++ {
					_, _ = tool.Search_index_search("orchidtoken", "body", 0, 10)
					_, _ = tool.Search_bbs_index_search("recovery", "", 0, 10)
				}
			}()
		}
		wait_group.Wait()
		if err := os.Remove("data"); err != nil {
			t.Fatalf("remove temporary parent obstruction: %v", err)
		}
		if err := os.Mkdir("data", 0o755); err != nil {
			t.Fatalf("restore index parent directory: %v", err)
		}
	case "create_retry":
		search_index_recovery_wait_for_log(t, log_buffer, "create document index staging directory", "create BBS index staging directory")
		if err := os.Chmod("data", 0o755); err != nil {
			t.Fatalf("restore index parent write permission: %v", err)
		}
	case "staging_retry":
		search_index_recovery_wait_for_log(t, log_buffer, "data/bleve.tmp", "data/bleve_bbs.tmp")
		if err := os.Chmod(document_index_directory+".tmp", 0o700); err != nil {
			t.Fatalf("restore document staging directory permission: %v", err)
		}
		if err := os.Chmod(bbs_index_directory+".tmp", 0o700); err != nil {
			t.Fatalf("restore BBS staging directory permission: %v", err)
		}
	case "permission_retry":
		search_index_recovery_wait_for_log(t, log_buffer, "document index recovery failed", "BBS index recovery failed")
		if err := os.Chmod(document_index_directory, 0o700); err != nil {
			t.Fatalf("restore document index permission: %v", err)
		}
		if err := os.Chmod(bbs_index_directory, 0o700); err != nil {
			t.Fatalf("restore BBS index permission: %v", err)
		}
	case "storage_target_unavailable":
		search_index_recovery_wait_for_log(t, log_buffer, "document index recovery failed", "BBS index recovery failed")
		for _, index := range []struct {
			active  string
			backing string
		}{
			{document_index_directory, filepath.Join("data", "document_backing")},
			{bbs_index_directory, filepath.Join("data", "bbs_backing")},
		} {
			if info, err := os.Lstat(index.active); err != nil || info.Mode()&os.ModeSymlink == 0 {
				t.Fatalf("active index symlink was replaced during target outage: info=%v err=%v", info, err)
			}
			if err := os.Rename(index.backing+".unavailable", index.backing); err != nil {
				t.Fatalf("restore temporarily unavailable index target: %v", err)
			}
		}
	case "install_failure":
		search_index_recovery_wait_for_log(t, log_buffer, "install rebuilt document index", "install rebuilt BBS index")
		for _, index_directory := range []string{document_index_directory, bbs_index_directory} {
			if _, err := os.Stat(filepath.Join(index_directory, "index_meta.json")); err != nil {
				t.Fatalf("failed installation damaged original index %q: %v\n%s", index_directory, err, log_buffer.String())
			}
		}
		if !search_index_recovery_file_exists(filepath.Join(document_version_file, "sentinel")) || !search_index_recovery_file_exists(filepath.Join(bbs_version_file, "sentinel")) {
			t.Fatalf("failed installation changed original version marker directories\n%s", log_buffer.String())
		}
		for _, version_file := range []string{document_version_file, bbs_version_file} {
			if err := os.RemoveAll(version_file); err != nil {
				t.Fatalf("remove test marker obstruction: %v", err)
			}
		}
		if err := os.WriteFile(document_version_file, []byte("2"), 0o600); err != nil {
			t.Fatalf("repair document marker: %v", err)
		}
		if err := os.WriteFile(bbs_version_file, []byte("4"), 0o600); err != nil {
			t.Fatalf("repair BBS marker: %v", err)
		}
	case "runtime_directory_loss":
		search_index_recovery_wait_for_indexes(t, 10*time.Second)
		if err := os.RemoveAll(document_index_directory); err != nil {
			t.Fatalf("remove active document index: %v", err)
		}
		if err := os.RemoveAll(bbs_index_directory); err != nil {
			t.Fatalf("remove active BBS index: %v", err)
		}
		time.Sleep(search_index_runtime_wait)
		_ = tool.Search_index_ready()
		_ = tool.Search_bbs_index_ready()
	case "runtime_handle_error":
		search_index_recovery_wait_for_indexes(t, 10*time.Second)
		closed_document_directory := filepath.Join("data", "closed_document")
		closed_bbs_directory := filepath.Join("data", "closed_bbs")
		search_index_recovery_create_document_index(t, closed_document_directory)
		search_index_recovery_create_bbs_index(t, closed_bbs_directory)
		closed_document_index, err := bleve.Open(closed_document_directory)
		if err != nil {
			t.Fatalf("open closed document test handle: %v", err)
		}
		closed_bbs_index, err := bleve.Open(closed_bbs_directory)
		if err != nil {
			t.Fatalf("open closed BBS test handle: %v", err)
		}
		tool.Search_index_set_ready(closed_document_index)
		tool.Search_bbs_index_set_ready(closed_bbs_index)
		if err := closed_document_index.Close(); err != nil {
			t.Fatalf("close injected document index handle: %v", err)
		}
		if err := closed_bbs_index.Close(); err != nil {
			t.Fatalf("close injected BBS index handle: %v", err)
		}
		if _, ready := tool.Search_index_search("orchidtoken", "body", 0, 10); ready {
			t.Fatal("closed document index unexpectedly reported a successful search")
		}
		if _, ready := tool.Search_bbs_index_search("recovery", "", 0, 10); ready {
			t.Fatal("closed BBS index unexpectedly reported a successful search")
		}
	case "concurrent_changes":
		search_index_recovery_apply_concurrent_changes(t, db, log_buffer)
	case "source_read_error":
		search_index_recovery_apply_source_read_error(t, db)
	case "interrupted_install", "abandoned_staging", "staging_next_to_valid", "version_markers", "metadata_corruption", "startup_missing":
	default:
		t.Fatalf("scenario %q did not reach its recovery path", scenario)
	}

	if scenario != "concurrent_changes" && scenario != "source_read_error" {
		search_index_recovery_wait_for_indexes(t, 15*time.Second)
		search_index_recovery_assert_source_results(t)
	}
	if scenario == "abandoned_staging" || scenario == "staging_next_to_valid" {
		if search_index_recovery_file_exists(document_index_directory+".tmp") || search_index_recovery_file_exists(bbs_index_directory+".tmp") {
			t.Fatal("recovery left an abandoned staging directory behind")
		}
	}
	if scenario == "interrupted_install" && search_index_recovery_file_exists(document_index_directory+".previous-123") {
		t.Fatal("interrupted document index backup was not restored")
	}
}

const search_index_runtime_wait = 5200 * time.Millisecond

type search_index_recovery_log_buffer struct {
	mu     sync.Mutex
	buffer bytes.Buffer
}

func (buffer *search_index_recovery_log_buffer) Write(data []byte) (int, error) {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.Write(data)
}

func (buffer *search_index_recovery_log_buffer) String() string {
	buffer.mu.Lock()
	defer buffer.mu.Unlock()
	return buffer.buffer.String()
}

func search_index_recovery_wait_for_log(t *testing.T, buffer *search_index_recovery_log_buffer, first string, second string) {
	t.Helper()
	deadline := time.Now().Add(10 * time.Second)
	for time.Now().Before(deadline) {
		logs := buffer.String()
		if strings.Contains(logs, first) && strings.Contains(logs, second) {
			return
		}
		time.Sleep(10 * time.Millisecond)
	}
	t.Fatalf("timed out waiting for recovery logs %q and %q:\n%s", first, second, buffer.String())
}

func search_index_recovery_seed_database(t *testing.T, db *sql.DB) {
	t.Helper()
	if err := tool.DB_transaction(db, func(tx *sql.Tx) error {
		queries := []struct {
			query string
			args  []any
		}{
			{"insert into data (title, data, type) values (?, ?, 'd')", []any{"recovery-doc", "orchidtoken original content"}},
			{"insert into data (title, data, type) values (?, ?, 'd')", []any{"deleted-doc", "deletetoken original content"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('title', ?, ?, ?)", []any{"91", "1", "recoverypost"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('data', ?, ?, ?)", []any{"91", "1", "posttoken original content"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('date', ?, ?, ?)", []any{"91", "1", "2026-01-01"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('title', ?, ?, ?)", []any{"92", "1", "deletedpost"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('data', ?, ?, ?)", []any{"92", "1", "deletetoken post"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment', ?, ?, ?)", []any{"91-1", "1", "nestedtoken original comment"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment_date', ?, ?, ?)", []any{"91-1", "1", "2026-01-01"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment', ?, ?, ?)", []any{"92-1", "1", "deletetoken deleted comment"}},
			{"insert into bbs_data (set_name, set_id, set_code, set_data) values ('comment_date', ?, ?, ?)", []any{"92-1", "1", "2026-01-01"}},
		}
		for _, query := range queries {
			if _, err := tx.Exec(tool.DB_change(query.query), query.args...); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		t.Fatalf("seed isolated recovery database: %v", err)
	}
}

func search_index_recovery_create_document_index(t *testing.T, directory string) {
	t.Helper()
	index, err := bleve.New(directory, tool.Search_index_mapping())
	if err != nil {
		t.Fatalf("create document Bleve index %q: %v", directory, err)
	}
	if err := index.Index("recovery-doc", tool.Search_document{Title_search: "recovery-doc", Data: "orchidtoken original content"}); err != nil {
		index.Close()
		t.Fatalf("seed document Bleve index %q: %v", directory, err)
	}
	if err := index.Close(); err != nil {
		t.Fatalf("close document Bleve index %q: %v", directory, err)
	}
}

func search_index_recovery_create_bbs_index(t *testing.T, directory string) {
	t.Helper()
	index, err := bleve.New(directory, tool.Bbs_search_index_mapping())
	if err != nil {
		t.Fatalf("create BBS Bleve index %q: %v", directory, err)
	}
	post := tool.Bbs_search_document{Set_id: "91", Set_code: "1", Title: "recoverypost", Data: "posttoken original content", Date: "2026-01-01"}
	comment := tool.Bbs_search_document{Set_id: "91", Set_code: "1", Comment_code: "1", Comment_data: "nestedtoken original comment", Date: "2026-01-01"}
	if err := index.Index(tool.Search_bbs_index_key("91", "1"), post); err != nil {
		index.Close()
		t.Fatalf("seed BBS post Bleve index %q: %v", directory, err)
	}
	if err := index.Index(tool.Search_bbs_index_comment_key("91", "1", "1"), comment); err != nil {
		index.Close()
		t.Fatalf("seed BBS comment Bleve index %q: %v", directory, err)
	}
	if err := index.Close(); err != nil {
		t.Fatalf("close BBS Bleve index %q: %v", directory, err)
	}
}

func search_index_recovery_create_marker_directory(t *testing.T, marker_path string) {
	t.Helper()
	if err := os.Mkdir(marker_path, 0o700); err != nil {
		t.Fatalf("create marker obstruction %q: %v", marker_path, err)
	}
	if err := os.WriteFile(filepath.Join(marker_path, "sentinel"), []byte("preserve"), 0o600); err != nil {
		t.Fatalf("write marker obstruction sentinel %q: %v", marker_path, err)
	}
}

func search_index_recovery_wait_for_indexes(t *testing.T, timeout time.Duration) {
	t.Helper()
	if !search_index_recovery_wait_for_indexes_result(timeout) {
		t.Fatal("document/BBS indexes did not recover")
	}
}

func search_index_recovery_wait_for_indexes_result(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		document_results, document_ready := tool.Search_index_search("orchidtoken", "body", 0, 20)
		bbs_results, bbs_ready := tool.Search_bbs_index_search_comment("nestedtoken", "91", 0, 20)
		if document_ready && bbs_ready && search_index_recovery_contains(document_results, "recovery-doc") && search_index_recovery_contains(bbs_results, tool.Search_bbs_index_comment_key("91", "1", "1")) {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

func search_index_recovery_wait_for_concurrent_results(timeout time.Duration) bool {
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) {
		document_results, document_ready := tool.Search_index_search("latesttoken", "body", 0, 20)
		comment_results, comments_ready := tool.Search_bbs_index_search_comment("latestcomment", "91", 0, 20)
		if document_ready && comments_ready && search_index_recovery_contains(document_results, "recovery-doc") && search_index_recovery_contains(comment_results, tool.Search_bbs_index_comment_key("91", "1", "1")) {
			return true
		}
		time.Sleep(20 * time.Millisecond)
	}
	return false
}

func search_index_recovery_assert_source_results(t *testing.T) {
	t.Helper()
	document_results, document_ready := tool.Search_index_search("orchidtoken", "body", 0, 20)
	if !document_ready || !search_index_recovery_contains(document_results, "recovery-doc") {
		t.Fatalf("document search after recovery = %v, ready=%v", document_results, document_ready)
	}
	bbs_results, bbs_ready := tool.Search_bbs_index_search("recovery", "91", 0, 20)
	if !bbs_ready || !search_index_recovery_contains(bbs_results, tool.Search_bbs_index_key("91", "1")) {
		t.Fatalf("BBS post search after recovery = %v, ready=%v", bbs_results, bbs_ready)
	}
	comment_results, comments_ready := tool.Search_bbs_index_search_comment("nestedtoken", "91", 0, 20)
	if !comments_ready || !search_index_recovery_contains(comment_results, tool.Search_bbs_index_comment_key("91", "1", "1")) {
		t.Fatalf("BBS comment search after recovery = %v, ready=%v", comment_results, comments_ready)
	}
}

func search_index_recovery_apply_concurrent_changes(t *testing.T, db *sql.DB, log_buffer *search_index_recovery_log_buffer) {
	t.Helper()
	var wait_group sync.WaitGroup
	for worker := 0; worker < 8; worker++ {
		wait_group.Add(1)
		go func() {
			defer wait_group.Done()
			for request := 0; request < 10; request++ {
				_, _ = tool.Search_index_search("orchidtoken", "body", 0, 20)
				_, _ = tool.Search_bbs_index_search("recovery", "91", 0, 20)
			}
		}()
	}

	if _, err := db.Exec(tool.DB_change("update data set data = ? where title = ?"), "latesttoken updated content", "recovery-doc"); err != nil {
		t.Fatalf("update document during rebuild: %v", err)
	}
	tool.Search_index_sync(db, "recovery-doc")
	if _, err := db.Exec(tool.DB_change("delete from data where title = ?"), "deleted-doc"); err != nil {
		t.Fatalf("delete document during rebuild: %v", err)
	}
	tool.Search_index_sync(db, "deleted-doc")

	if _, err := db.Exec(tool.DB_change("update bbs_data set set_data = ? where set_name = 'title' and set_id = ? and set_code = ?"), "updatedpost", "91", "1"); err != nil {
		t.Fatalf("update BBS title during rebuild: %v", err)
	}
	if _, err := db.Exec(tool.DB_change("update bbs_data set set_data = ? where set_name = 'data' and set_id = ? and set_code = ?"), "latesttoken post content", "91", "1"); err != nil {
		t.Fatalf("update BBS body during rebuild: %v", err)
	}
	tool.Search_bbs_index_update(db, "91", "1")
	if _, err := db.Exec(tool.DB_change("delete from bbs_data where set_id = ?"), "92"); err != nil {
		t.Fatalf("delete BBS post during rebuild: %v", err)
	}
	tool.Search_bbs_index_update(db, "92", "1")
	if _, err := db.Exec(tool.DB_change("update bbs_data set set_data = ? where set_name = 'comment' and set_id = ? and set_code = ?"), "latestcomment updated nested comment", "91-1", "1"); err != nil {
		t.Fatalf("update BBS comment during rebuild: %v", err)
	}
	tool.Search_bbs_index_update_comment(db, "91", "1", "1")
	if _, err := db.Exec(tool.DB_change("delete from bbs_data where set_id = ?"), "92-1"); err != nil {
		t.Fatalf("delete BBS comment during rebuild: %v", err)
	}
	tool.Search_bbs_index_delete_comment("92", "1", "1")
	wait_group.Wait()

	if !search_index_recovery_wait_for_concurrent_results(15 * time.Second) {
		t.Fatalf("latest concurrent document/BBS changes were not indexed:\n%s", log_buffer.String())
	}
	latest_document, document_ready := tool.Search_index_search("latesttoken", "body", 0, 20)
	if !document_ready || !search_index_recovery_contains(latest_document, "recovery-doc") {
		t.Fatalf("latest document edit was lost: %v, ready=%v", latest_document, document_ready)
	}
	deleted_document, document_ready := tool.Search_index_search("deletetoken", "body", 0, 20)
	if !document_ready || search_index_recovery_contains(deleted_document, "deleted-doc") {
		t.Fatalf("deleted document was resurrected: %v, ready=%v", deleted_document, document_ready)
	}
	latest_post, bbs_ready := tool.Search_bbs_index_search_data("latesttoken", "91", 0, 20)
	if !bbs_ready || !search_index_recovery_contains(latest_post, tool.Search_bbs_index_key("91", "1")) {
		t.Fatalf("latest BBS post edit was lost: %v, ready=%v", latest_post, bbs_ready)
	}
	deleted_post, bbs_ready := tool.Search_bbs_index_search("deletedpost", "92", 0, 20)
	if !bbs_ready || search_index_recovery_contains(deleted_post, tool.Search_bbs_index_key("92", "1")) {
		t.Fatalf("deleted BBS post was resurrected: %v, ready=%v", deleted_post, bbs_ready)
	}
	latest_comment, bbs_ready := tool.Search_bbs_index_search_comment("latestcomment", "91", 0, 20)
	if !bbs_ready || !search_index_recovery_contains(latest_comment, tool.Search_bbs_index_comment_key("91", "1", "1")) {
		t.Fatalf("latest nested comment edit was lost: %v, ready=%v", latest_comment, bbs_ready)
	}
	deleted_comment, bbs_ready := tool.Search_bbs_index_search_comment("deletetoken", "92", 0, 20)
	if !bbs_ready || search_index_recovery_contains(deleted_comment, tool.Search_bbs_index_comment_key("92", "1", "1")) {
		t.Fatalf("deleted nested comment was resurrected: %v, ready=%v", deleted_comment, bbs_ready)
	}
}

func search_index_recovery_apply_source_read_error(t *testing.T, db *sql.DB) {
	t.Helper()
	search_index_recovery_wait_for_indexes(t, 10*time.Second)
	if _, err := db.Exec(tool.DB_change("update data set data = ? where title = ?"), "freshsourcetoken current content", "recovery-doc"); err != nil {
		t.Fatalf("change document source before forcing index rebuild: %v", err)
	}
	if _, err := db.Exec(tool.DB_change("update bbs_data set set_data = ? where set_name = 'comment' and set_id = ? and set_code = ?"), "freshcommenttoken current comment", "91-1", "1"); err != nil {
		t.Fatalf("change comment source before forcing index rebuild: %v", err)
	}
	if err := db.Close(); err != nil {
		t.Fatalf("close source handle to inject read errors: %v", err)
	}
	tool.Search_index_sync(db, "recovery-doc")
	tool.Search_bbs_index_update_comment(db, "91", "1", "1")

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		document_results, document_ready := tool.Search_index_search("freshsourcetoken", "body", 0, 20)
		comment_results, comments_ready := tool.Search_bbs_index_search_comment("freshcommenttoken", "91", 0, 20)
		if document_ready && comments_ready && search_index_recovery_contains(document_results, "recovery-doc") && search_index_recovery_contains(comment_results, tool.Search_bbs_index_comment_key("91", "1", "1")) {
			old_document_results, _ := tool.Search_index_search("orchidtoken", "body", 0, 20)
			old_comment_results, _ := tool.Search_bbs_index_search_comment("nestedtoken", "91", 0, 20)
			if search_index_recovery_contains(old_document_results, "recovery-doc") || search_index_recovery_contains(old_comment_results, tool.Search_bbs_index_comment_key("91", "1", "1")) {
				t.Fatalf("source read failure reactivated stale indexes: document=%v comment=%v", old_document_results, old_comment_results)
			}
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("source read failures did not trigger full document and BBS index rebuilds")
}

func search_index_recovery_contains(values []string, expected string) bool {
	for _, value := range values {
		if value == expected {
			return true
		}
	}
	return false
}

func search_index_recovery_file_exists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

func search_index_recovery_mysql_setup(t *testing.T) (string, string) {
	t.Helper()
	dsn := os.Getenv("MYSQL_AUDIT_DSN")
	configuration, err := search_index_recovery_parse_mysql_dsn(dsn)
	if err != nil {
		t.Fatalf("parse MYSQL_AUDIT_DSN: %v", err)
	}
	database_name := fmt.Sprintf("opennamu_index_recovery_%d", time.Now().UnixNano()%1_000_000_000)
	server_db, err := sql.Open("mysql", dsn)
	if err != nil {
		t.Fatalf("open configured MySQL server: %v", err)
	}
	defer server_db.Close()
	if _, err := server_db.Exec("create database `" + database_name + "` character set utf8mb4"); err != nil {
		t.Fatalf("create isolated MySQL recovery schema: %v", err)
	}
	configuration["db_name"] = database_name
	configuration_json, err := json.Marshal(configuration)
	if err != nil {
		t.Fatalf("encode isolated MySQL connection settings: %v", err)
	}
	t.Cleanup(func() { search_index_recovery_drop_mysql_database(t, database_name) })
	return string(configuration_json), database_name
}

func search_index_recovery_drop_mysql_database(t *testing.T, database_name string) {
	t.Helper()
	server_db, err := sql.Open("mysql", os.Getenv("MYSQL_AUDIT_DSN"))
	if err != nil {
		t.Errorf("open MySQL for test-schema cleanup: %v", err)
		return
	}
	defer server_db.Close()
	if _, err := server_db.Exec("drop database if exists `" + database_name + "`"); err != nil {
		t.Errorf("drop isolated MySQL recovery schema: %v", err)
	}
}

func search_index_recovery_parse_mysql_dsn(dsn string) (map[string]string, error) {
	credential_end := strings.Index(dsn, "@")
	if credential_end <= 0 {
		return nil, fmt.Errorf("expected user:password@tcp(host:port)/ DSN")
	}
	credentials := dsn[:credential_end]
	credential_separator := strings.Index(credentials, ":")
	if credential_separator <= 0 {
		return nil, fmt.Errorf("expected user and password in DSN")
	}
	transport_start := strings.Index(dsn[credential_end+1:], "(")
	transport_end := strings.Index(dsn[credential_end+1:], ")")
	if transport_start < 0 || transport_end <= transport_start {
		return nil, fmt.Errorf("expected tcp(host:port) transport")
	}
	address := dsn[credential_end+1+transport_start+1 : credential_end+1+transport_end]
	host, port, err := net.SplitHostPort(address)
	if err != nil {
		return nil, err
	}
	return map[string]string{
		"user":     credentials[:credential_separator],
		"password": credentials[credential_separator+1:],
		"host":     host,
		"port":     port,
	}, nil
}
