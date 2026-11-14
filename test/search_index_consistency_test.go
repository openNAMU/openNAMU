package test

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/blevesearch/bleve/v2"
	"opennamu/route/tool"
)

const search_index_consistency_restart_child_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_RESTART_CHILD"
const search_index_consistency_restart_phase_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_RESTART_PHASE"
const search_index_consistency_restart_backend_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_RESTART_BACKEND"
const search_index_consistency_pagination_child_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_PAGINATION_CHILD"
const search_index_consistency_pagination_backend_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_PAGINATION_BACKEND"
const search_index_consistency_escalation_child_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_ESCALATION_CHILD"
const search_index_consistency_escalation_backend_env = "OPENNAMU_SEARCH_INDEX_CONSISTENCY_ESCALATION_BACKEND"

func Test_Search_index_consistency_resumes_after_restart(t *testing.T) {
	backends := []struct {
		name         string
		mysql_config string
	}{
		{name: "sqlite"},
		{name: "turso"},
	}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		mysql_config, _ := search_index_recovery_mysql_setup(t)
		backends = append(backends, struct {
			name         string
			mysql_config string
		}{name: "mysql", mysql_config: mysql_config})
	}

	for _, backend := range backends {
		t.Run(backend.name, func(t *testing.T) {
			root, err := os.MkdirTemp(filepath.Join("..", "tmp"), "search_index_consistency_restart_")
			if err != nil {
				t.Fatalf("create isolated audit root: %v", err)
			}
			t.Cleanup(func() {
				if err := os.RemoveAll(root); err != nil {
					t.Errorf("remove isolated audit root: %v", err)
				}
			})
			root, err = filepath.Abs(root)
			if err != nil {
				t.Fatalf("resolve isolated audit root: %v", err)
			}

			for _, phase := range []string{"first", "second"} {
				command_context, cancel := context.WithTimeout(context.Background(), 45*time.Second)
				command := exec.CommandContext(command_context, os.Args[0], "-test.run=^Test_Search_index_consistency_restart_child$")
				command.Dir = root
				command.Env = search_index_recovery_environment(os.Environ(), search_index_consistency_restart_child_env, "1")
				command.Env = search_index_recovery_environment(command.Env, search_index_consistency_restart_phase_env, phase)
				command.Env = search_index_recovery_environment(command.Env, search_index_consistency_restart_backend_env, backend.name)
				command.Env = search_index_recovery_environment(command.Env, "NAMU_DB", filepath.Join(root, "source_database"))
				command.Env = search_index_recovery_environment(command.Env, "NAMU_DB_TYPE", backend.name)
				if backend.mysql_config != "" {
					command.Env = search_index_recovery_environment(command.Env, search_index_recovery_mysql_config_env, backend.mysql_config)
				}
				output, command_err := command.CombinedOutput()
				cancel()
				if command_err != nil {
					t.Fatalf("consistency audit child failed (%s/%s): %v\n%s", backend.name, phase, command_err, output)
				}
			}
		})
	}
}

func Test_Search_index_consistency_paginates_bounded_batches(t *testing.T) {
	backends := []struct {
		name         string
		mysql_config string
	}{
		{name: "sqlite"},
		{name: "turso"},
	}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		mysql_config, _ := search_index_recovery_mysql_setup(t)
		backends = append(backends, struct {
			name         string
			mysql_config string
		}{name: "mysql", mysql_config: mysql_config})
	}

	for _, backend := range backends {
		t.Run(backend.name, func(t *testing.T) {
			root, err := os.MkdirTemp(filepath.Join("..", "tmp"), "search_index_consistency_pagination_")
			if err != nil {
				t.Fatalf("create isolated pagination root: %v", err)
			}
			t.Cleanup(func() {
				if err := os.RemoveAll(root); err != nil {
					t.Errorf("remove isolated pagination root: %v", err)
				}
			})
			root, err = filepath.Abs(root)
			if err != nil {
				t.Fatalf("resolve isolated pagination root: %v", err)
			}

			command_context, cancel := context.WithTimeout(context.Background(), 45*time.Second)
			defer cancel()
			command := exec.CommandContext(command_context, os.Args[0], "-test.run=^Test_Search_index_consistency_pagination_child$")
			command.Dir = root
			command.Env = search_index_recovery_environment(os.Environ(), search_index_consistency_pagination_child_env, "1")
			command.Env = search_index_recovery_environment(command.Env, search_index_consistency_pagination_backend_env, backend.name)
			command.Env = search_index_recovery_environment(command.Env, "NAMU_DB", filepath.Join(root, "source_database"))
			command.Env = search_index_recovery_environment(command.Env, "NAMU_DB_TYPE", backend.name)
			if backend.mysql_config != "" {
				command.Env = search_index_recovery_environment(command.Env, search_index_recovery_mysql_config_env, backend.mysql_config)
			}
			output, err := command.CombinedOutput()
			if command_context.Err() != nil {
				t.Fatalf("pagination audit child timed out (%s): %v\n%s", backend.name, command_context.Err(), output)
			}
			if err != nil {
				t.Fatalf("pagination audit child failed (%s): %v\n%s", backend.name, err, output)
			}
		})
	}
}

func Test_Search_index_consistency_pagination_child(t *testing.T) {
	if os.Getenv(search_index_consistency_pagination_child_env) != "1" {
		t.Skip("search-index consistency pagination subprocess helper")
	}
	if os.Getenv(search_index_consistency_pagination_backend_env) == "mysql" {
		if err := os.MkdirAll("data", 0o755); err != nil {
			t.Fatalf("create isolated MySQL configuration directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join("data", "mysql.json"), []byte(os.Getenv(search_index_recovery_mysql_config_env)), 0o600); err != nil {
			t.Fatalf("write isolated MySQL connection settings: %v", err)
		}
	}
	tool.DB_boot()
	tool.DB_init()
	db := tool.DB_connect()
	defer tool.DB_close(db)
	search_index_recovery_seed_database(t, db)
	for document_number := 0; document_number < 105; document_number++ {
		title := fmt.Sprintf("bounded-%03d", document_number)
		if _, err := db.Exec(tool.DB_change("insert into data (title, data, type) values (?, ?, 'd')"), title, "bounded pagination content"); err != nil {
			t.Fatalf("seed document %q: %v", title, err)
		}
	}
	tool.Search_index_start()
	tool.Search_bbs_index_start()
	search_index_recovery_wait_for_indexes(t, 15*time.Second)

	state_dir := filepath.Join("data", "pagination_audit")
	phases := []string{"sql", "index", "index", "complete"}
	processed_counts := []int{100, 7, 100, 7}
	for batch, expected_phase := range phases {
		result, err := tool.Search_index_consistency_audit_batch(db, state_dir, time.Now())
		if err != nil {
			t.Fatalf("run bounded document audit batch %d: %v", batch+1, err)
		}
		if result.Phase != expected_phase || result.Processed != processed_counts[batch] || result.Rebuild_requested {
			t.Fatalf("document audit batch %d violated bounded paging: got %+v, phase=%q processed=%d", batch+1, result, expected_phase, processed_counts[batch])
		}
		if (batch == len(phases)-1) != result.Cycle_completed {
			t.Fatalf("document audit completion state is wrong after batch %d: %+v", batch+1, result)
		}
	}
	result, err := tool.Search_index_consistency_audit_batch(db, state_dir, time.Now())
	if err != nil {
		t.Fatalf("check completed audit interval: %v", err)
	}
	if !result.Skipped || result.Processed != 0 {
		t.Fatalf("completed audit did not respect its schedule interval: %+v", result)
	}
}

func Test_Search_index_consistency_escalates_widespread_divergence(t *testing.T) {
	backends := []struct {
		name         string
		mysql_config string
	}{
		{name: "sqlite"},
		{name: "turso"},
	}
	if os.Getenv("MYSQL_AUDIT_DSN") != "" {
		mysql_config, _ := search_index_recovery_mysql_setup(t)
		backends = append(backends, struct {
			name         string
			mysql_config string
		}{name: "mysql", mysql_config: mysql_config})
	}

	for _, backend := range backends {
		t.Run(backend.name, func(t *testing.T) {
			root, err := os.MkdirTemp(filepath.Join("..", "tmp"), "search_index_consistency_escalation_")
			if err != nil {
				t.Fatalf("create isolated escalation root: %v", err)
			}
			t.Cleanup(func() {
				if err := os.RemoveAll(root); err != nil {
					t.Errorf("remove isolated escalation root: %v", err)
				}
			})
			root, err = filepath.Abs(root)
			if err != nil {
				t.Fatalf("resolve isolated escalation root: %v", err)
			}

			command_context, cancel := context.WithTimeout(context.Background(), 45*time.Second)
			defer cancel()
			command := exec.CommandContext(command_context, os.Args[0], "-test.run=^Test_Search_index_consistency_escalation_child$")
			command.Dir = root
			command.Env = search_index_recovery_environment(os.Environ(), search_index_consistency_escalation_child_env, "1")
			command.Env = search_index_recovery_environment(command.Env, search_index_consistency_escalation_backend_env, backend.name)
			command.Env = search_index_recovery_environment(command.Env, "NAMU_DB", filepath.Join(root, "source_database"))
			command.Env = search_index_recovery_environment(command.Env, "NAMU_DB_TYPE", backend.name)
			if backend.mysql_config != "" {
				command.Env = search_index_recovery_environment(command.Env, search_index_recovery_mysql_config_env, backend.mysql_config)
			}
			output, err := command.CombinedOutput()
			if command_context.Err() != nil {
				t.Fatalf("escalation audit child timed out (%s): %v\n%s", backend.name, command_context.Err(), output)
			}
			if err != nil {
				t.Fatalf("escalation audit child failed (%s): %v\n%s", backend.name, err, output)
			}
		})
	}
}

func Test_Search_index_consistency_escalation_child(t *testing.T) {
	if os.Getenv(search_index_consistency_escalation_child_env) != "1" {
		t.Skip("search-index consistency escalation subprocess helper")
	}
	if os.Getenv(search_index_consistency_escalation_backend_env) == "mysql" {
		if err := os.MkdirAll("data", 0o755); err != nil {
			t.Fatalf("create isolated MySQL configuration directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join("data", "mysql.json"), []byte(os.Getenv(search_index_recovery_mysql_config_env)), 0o600); err != nil {
			t.Fatalf("write isolated MySQL connection settings: %v", err)
		}
	}
	tool.DB_boot()
	tool.DB_init()
	db := tool.DB_connect()
	defer tool.DB_close(db)
	search_index_recovery_seed_database(t, db)
	for document_number := 0; document_number < 8; document_number++ {
		title := fmt.Sprintf("escalation-%03d", document_number)
		if _, err := db.Exec(tool.DB_change("insert into data (title, data, type) values (?, ?, 'd')"), title, "authoritativetoken document source"); err != nil {
			t.Fatalf("seed escalation document %q: %v", title, err)
		}
	}
	for post_number := 0; post_number < 3; post_number++ {
		set_id := fmt.Sprintf("%d", 94+post_number)
		queries := []struct {
			name string
			data string
		}{
			{"title", "escalationpost" + set_id},
			{"data", "authoritativeposttoken source"},
			{"date", "2026-01-02"},
			{"comment", "authoritativecommenttoken source"},
			{"comment_date", "2026-01-02"},
		}
		for _, query := range queries {
			storage_id := set_id
			if query.name == "comment" || query.name == "comment_date" {
				storage_id += "-1"
			}
			if _, err := db.Exec(tool.DB_change("insert into bbs_data (set_name, set_id, set_code, set_data) values (?, ?, ?, ?)"), query.name, storage_id, "1", query.data); err != nil {
				t.Fatalf("seed escalation BBS %s record for %s: %v", query.name, set_id, err)
			}
		}
	}

	document_index, err := bleve.New(filepath.Join("data", "bleve"), tool.Search_index_mapping())
	if err != nil {
		t.Fatalf("create divergent document index: %v", err)
	}
	for document_number := 0; document_number < 10; document_number++ {
		title := "recovery-doc"
		if document_number == 1 {
			title = "deleted-doc"
		} else if document_number >= 2 {
			title = fmt.Sprintf("escalation-%03d", document_number-2)
		}
		if err := document_index.Index(title, tool.Search_index_document(title, "wrong indexed content")); err != nil {
			document_index.Close()
			t.Fatalf("create stale document %q: %v", title, err)
		}
	}
	if err := document_index.Close(); err != nil {
		t.Fatalf("close divergent document index: %v", err)
	}

	bbs_index, err := bleve.New(filepath.Join("data", "bleve_bbs"), tool.Bbs_search_index_mapping())
	if err != nil {
		t.Fatalf("create divergent BBS index: %v", err)
	}
	post_ids := []string{"91", "92", "94", "95", "96"}
	for _, set_id := range post_ids {
		post := tool.Bbs_search_document{Set_id: set_id, Set_code: "1", Title: "stale-post", Data: "wrong indexed post", Date: "2026-01-01"}
		comment := tool.Bbs_search_document{Set_id: set_id, Set_code: "1", Comment_code: "1", Comment_data: "wrong indexed comment", Date: "2026-01-01"}
		if err := bbs_index.Index(tool.Search_bbs_index_key(set_id, "1"), post); err != nil {
			bbs_index.Close()
			t.Fatalf("create stale BBS post %s: %v", set_id, err)
		}
		if err := bbs_index.Index(tool.Search_bbs_index_comment_key(set_id, "1", "1"), comment); err != nil {
			bbs_index.Close()
			t.Fatalf("create stale BBS comment %s: %v", set_id, err)
		}
	}
	if err := bbs_index.Close(); err != nil {
		t.Fatalf("close divergent BBS index: %v", err)
	}
	if err := os.WriteFile(filepath.Join("data", "bleve.version"), []byte("3"), 0o600); err != nil {
		t.Fatalf("write current document marker: %v", err)
	}
	if err := os.WriteFile(filepath.Join("data", "bleve_bbs.version"), []byte("5"), 0o600); err != nil {
		t.Fatalf("write current BBS marker: %v", err)
	}

	tool.Search_index_start()
	tool.Search_bbs_index_start()
	if !Search_index_recovery_wait_for_indexes_ready(10 * time.Second) {
		t.Fatal("divergent indexes were not activated")
	}
	document_result, err := tool.Search_index_consistency_audit_batch(db, filepath.Join("data", "escalation_state"), time.Now())
	if err != nil {
		t.Fatalf("audit widespread document divergence: %v", err)
	}
	if !document_result.Rebuild_requested || document_result.Total_repaired != 10 {
		t.Fatalf("widespread document divergence did not escalate after targeted repairs: %+v", document_result)
	}
	bbs_state_dir := filepath.Join("data", "escalation_state")
	bbs_result, err := tool.Search_bbs_index_consistency_audit_batch(db, bbs_state_dir, time.Now())
	if err != nil {
		t.Fatalf("audit widespread BBS post divergence: %v", err)
	}
	if bbs_result.Rebuild_requested || bbs_result.Total_repaired != 5 {
		t.Fatalf("BBS posts should receive targeted repairs before threshold: %+v", bbs_result)
	}
	bbs_result, err = tool.Search_bbs_index_consistency_audit_batch(db, bbs_state_dir, time.Now())
	if err != nil {
		t.Fatalf("audit widespread BBS comment divergence: %v", err)
	}
	if !bbs_result.Rebuild_requested || bbs_result.Total_repaired != 10 {
		t.Fatalf("widespread BBS divergence did not escalate after targeted repairs: %+v", bbs_result)
	}

	deadline := time.Now().Add(15 * time.Second)
	for time.Now().Before(deadline) {
		document_results, document_ready := tool.Search_index_search("authoritativetoken", "body", 0, 30)
		post_results, post_ready := tool.Search_bbs_index_search_data("authoritativeposttoken", "94", 0, 20)
		comment_results, comment_ready := tool.Search_bbs_index_search_comment("authoritativecommenttoken", "94", 0, 20)
		if document_ready && post_ready && comment_ready && search_index_recovery_contains(document_results, "escalation-000") && search_index_recovery_contains(post_results, tool.Search_bbs_index_key("94", "1")) && search_index_recovery_contains(comment_results, tool.Search_bbs_index_comment_key("94", "1", "1")) {
			return
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatal("existing recovery mechanisms did not rebuild widespread document, BBS post, and nested-comment divergence")
}

func Test_Search_index_consistency_restart_child(t *testing.T) {
	if os.Getenv(search_index_consistency_restart_child_env) != "1" {
		t.Skip("search-index consistency restart subprocess helper")
	}

	if os.Getenv(search_index_consistency_restart_backend_env) == "mysql" {
		if err := os.MkdirAll("data", 0o755); err != nil {
			t.Fatalf("create isolated MySQL configuration directory: %v", err)
		}
		if err := os.WriteFile(filepath.Join("data", "mysql.json"), []byte(os.Getenv(search_index_recovery_mysql_config_env)), 0o600); err != nil {
			t.Fatalf("write isolated MySQL connection settings: %v", err)
		}
	}
	tool.DB_boot()
	tool.DB_init()
	db := tool.DB_connect()
	defer tool.DB_close(db)
	phase := os.Getenv(search_index_consistency_restart_phase_env)
	if phase == "first" {
		search_index_recovery_seed_database(t, db)
		Search_index_recovery_create_divergent_indexes(t, filepath.Join("data", "bleve"), filepath.Join("data", "bleve_bbs"))
		if err := os.WriteFile(filepath.Join("data", "bleve.version"), []byte("3"), 0o600); err != nil {
			t.Fatalf("write current document index marker: %v", err)
		}
		if err := os.WriteFile(filepath.Join("data", "bleve_bbs.version"), []byte("5"), 0o600); err != nil {
			t.Fatalf("write current BBS index marker: %v", err)
		}
	} else if phase != "second" {
		t.Fatalf("unknown restart test phase %q", phase)
	}

	tool.Search_index_start()
	tool.Search_bbs_index_start()
	if !Search_index_recovery_wait_for_indexes_ready(10 * time.Second) {
		t.Fatal("indexes were not ready for the restart test")
	}
	state_dir := filepath.Join("data", "audit_progress")

	if phase == "first" {
		if err := os.MkdirAll(state_dir, 0o755); err != nil {
			t.Fatalf("create audit checkpoint directory: %v", err)
		}
		if err := os.WriteFile(tool.Search_index_consistency_state_path(state_dir, "document"), []byte("{"), 0o600); err != nil {
			t.Fatalf("seed a corrupt document checkpoint: %v", err)
		}
		tool.Search_index_consistency_audit_scheduled(db, state_dir, time.Now())
		document_state, err := os.ReadFile(tool.Search_index_consistency_state_path(state_dir, "document"))
		if err != nil {
			t.Fatalf("read document audit checkpoint before restart: %v", err)
		}
		bbs_state, err := os.ReadFile(tool.Search_index_consistency_state_path(state_dir, "bbs"))
		if err != nil {
			t.Fatalf("read BBS audit checkpoint before restart: %v", err)
		}
		if !strings.Contains(string(document_state), `"phase":"index"`) || !strings.Contains(string(bbs_state), `"phase":"comments"`) {
			t.Fatalf("audit checkpoints did not persist independent phases: document=%s BBS=%s", document_state, bbs_state)
		}
		for _, index_name := range []string{"document", "bbs"} {
			last_run_path := filepath.Join(state_dir, ".opennamu_"+index_name+"_search_index_audit_last_run")
			if _, err := os.Stat(last_run_path); err != nil {
				t.Fatalf("scheduled %s audit did not persist its independent due state: %v", index_name, err)
			}
		}
		return
	}

	document_state_before, err := os.ReadFile(tool.Search_index_consistency_state_path(state_dir, "document"))
	if err != nil {
		t.Fatalf("read document checkpoint after process restart: %v", err)
	}
	tool.Search_index_consistency_audit_scheduled(db, state_dir, time.Now().Add(-time.Hour))
	document_state_after, err := os.ReadFile(tool.Search_index_consistency_state_path(state_dir, "document"))
	if err != nil {
		t.Fatalf("read document checkpoint after scheduled not-due check: %v", err)
	}
	if string(document_state_before) != string(document_state_after) {
		t.Fatal("scheduled audit ignored its independent one-minute due state")
	}

	document_result, err := tool.Search_index_consistency_audit_batch(db, state_dir, time.Now())
	if err != nil {
		t.Fatalf("resume document index audit: %v", err)
	}
	if document_result.Phase != "complete" || !document_result.Cycle_completed || document_result.Rebuild_requested {
		t.Fatalf("document audit did not resume from its saved index cursor: %+v", document_result)
	}
	bbs_result, err := tool.Search_bbs_index_consistency_audit_batch(db, state_dir, time.Now())
	if err != nil {
		t.Fatalf("resume BBS comment audit: %v", err)
	}
	if bbs_result.Phase != "index" || bbs_result.Cycle_completed || bbs_result.Rebuild_requested {
		t.Fatalf("BBS audit did not resume from its saved comment phase: %+v", bbs_result)
	}
	bbs_result, err = tool.Search_bbs_index_consistency_audit_batch(db, state_dir, time.Now())
	if err != nil {
		t.Fatalf("finish resumed BBS index audit: %v", err)
	}
	if bbs_result.Phase != "complete" || !bbs_result.Cycle_completed || bbs_result.Rebuild_requested {
		t.Fatalf("resumed BBS audit did not complete: %+v", bbs_result)
	}

	document_results, document_ready := tool.Search_index_search("orchidtoken", "body", 0, 20)
	if !document_ready || !search_index_recovery_contains(document_results, "recovery-doc") {
		t.Fatalf("resumed audit did not repair the document index: results=%v ready=%v", document_results, document_ready)
	}
	comment_results, comment_ready := tool.Search_bbs_index_search_comment("nestedtoken", "91", 0, 20)
	if !comment_ready || !search_index_recovery_contains(comment_results, tool.Search_bbs_index_comment_key("91", "1", "1")) {
		t.Fatalf("resumed audit did not repair the nested-comment index: results=%v ready=%v", comment_results, comment_ready)
	}
	orphan_results, orphan_ready := tool.Search_bbs_index_search_comment("orphan", "93", 0, 20)
	if !orphan_ready || search_index_recovery_contains(orphan_results, tool.Search_bbs_index_comment_key("93", "1", "1")) {
		t.Fatalf("resumed audit left an orphan nested comment: results=%v ready=%v", orphan_results, orphan_ready)
	}

	var source_document string
	if err := db.QueryRow("select data from data where title = ?", "recovery-doc").Scan(&source_document); err != nil || source_document != "orchidtoken original content" {
		t.Fatalf("resumed audit changed authoritative SQL content: data=%q err=%v", source_document, err)
	}
}
