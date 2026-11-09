package test

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"sync"
	"testing"

	"github.com/blevesearch/bleve/v2"
	"opennamu/route/tool"
)

const search_index_child_env = "OPENNAMU_PHASE6_SEARCH_INDEX_CHILD"

func Test_Search_index_pending_active_updates_and_query_contract(t *testing.T) {
	command := exec.Command(os.Args[0], "-test.run=^Test_Search_index_behavior_child$")
	command.Env = append(os.Environ(), search_index_child_env+"=1")
	if output, err := command.CombinedOutput(); err != nil {
		t.Fatalf("search-index behavior child failed: %v\n%s", err, output)
	}
}

func Test_Search_index_behavior_child(t *testing.T) {
	if os.Getenv(search_index_child_env) != "1" {
		t.Skip("search-index subprocess helper")
	}

	index_directory, err := os.MkdirTemp(filepath.Join("..", "tmp"), "phase6_search_index_")
	if err != nil {
		t.Fatalf("create temporary index directory: %v", err)
	}
	defer os.RemoveAll(index_directory)

	index, err := bleve.New(filepath.Join(index_directory, "index"), tool.Search_index_mapping())
	if err != nil {
		t.Fatalf("create search index: %v", err)
	}
	defer index.Close()

	if err := index.Index("queued-doc", tool.Search_document{Title_search: "queued-doc", Data: "obsolete orchid"}); err != nil {
		t.Fatalf("seed queued document: %v", err)
	}
	if err := index.Index("removed-doc", tool.Search_document{Title_search: "removed-doc", Data: "orchid"}); err != nil {
		t.Fatalf("seed pending-delete document: %v", err)
	}

	tool.Search_index_update("queued-doc", "violet unicorn")
	tool.Search_index_update("queued-doc", "violet unicorn final")
	tool.Search_index_delete("removed-doc")
	tool.Search_index_set_ready(index)
	if !tool.Search_index_ready() {
		t.Fatal("index is not ready after pending changes were applied")
	}

	assert_search := func(keyword string, search_type string, offset int, limit int, want []string) {
		t.Helper()
		got, ready := tool.Search_index_search(keyword, search_type, offset, limit)
		if !ready || !reflect.DeepEqual(got, want) {
			t.Fatalf("Search_index_search(%q, %q, %d, %d) = %v, ready %v; want %v, ready true", keyword, search_type, offset, limit, got, ready, want)
		}
	}
	assert_search("unicorn", "body", 0, 50, []string{"queued-doc"})
	assert_search("obsolete", "body", 0, 50, []string{})
	assert_search("orchid", "body", 0, 50, []string{})
	assert_search("queued", "title", 0, 50, []string{"queued-doc"})

	tool.Search_index_update("active-doc", "temporary cedar")
	tool.Search_index_delete("active-doc")
	tool.Search_index_update("active-doc", "fresh cedar")
	tool.Search_index_update("sort-b", "sharedtoken")
	tool.Search_index_update("sort-a", "sharedtoken")
	tool.Search_index_update("sort-c", "sharedtoken")
	assert_search("temporary", "body", 0, 50, []string{})
	assert_search("cedar", "body", 0, 50, []string{"active-doc"})
	assert_search("sharedtoken", "body", 0, 50, []string{"sort-a", "sort-b", "sort-c"})
	assert_search("sharedtoken", "body", 1, 1, []string{"sort-b"})

	var wait_group sync.WaitGroup
	for worker := 0; worker < 3; worker++ {
		worker := worker
		wait_group.Add(1)
		go func() {
			defer wait_group.Done()
			for update := 0; update < 8; update++ {
				tool.Search_index_update(fmt.Sprintf("concurrent-%d", worker), fmt.Sprintf("worker%d revision%d", worker, update))
				tool.Search_index_search("sharedtoken", "body", 0, 10)
			}
			tool.Search_index_update(fmt.Sprintf("concurrent-%d", worker), fmt.Sprintf("worker%d finalmarker", worker))
		}()
	}
	wait_group.Wait()
	for worker := 0; worker < 3; worker++ {
		assert_search(fmt.Sprintf("worker%d", worker), "body", 0, 50, []string{fmt.Sprintf("concurrent-%d", worker)})
	}
}
