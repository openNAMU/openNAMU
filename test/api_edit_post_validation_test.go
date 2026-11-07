package test

import (
	"strings"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_Api_edit_post_validation_responses(t *testing.T) {
	db := Setup_rebuild_db(t)
	defer tool.DB_close(db)

	if _, err := db.Exec("insert into other (name, data, coverage) values ('title_max_length', '10', '')"); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name     string
		config   tool.Config
		doc_name string
		want     string
	}{
		{name: "empty title", config: tool.Config{IP: "admin_ip"}, doc_name: "", want: "empty title"},
		{name: "title length", config: tool.Config{IP: "admin_ip"}, doc_name: strings.Repeat("x", 11), want: "title length"},
		{name: "permission", config: tool.Config{IP: "192.0.2.44"}, doc_name: "Example", want: "permission denied"},
	}

	for _, test_case := range cases {
		t.Run(test_case.name, func(t *testing.T) {
			response := route.Api_edit_post(test_case.config, test_case.doc_name, "content", "", "", "")
			if response["response"] != "error" || response["data"] != test_case.want {
				t.Fatalf("response=%v, want response=error data=%q", response, test_case.want)
			}
		})
	}
}
