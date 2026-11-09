package test

import (
	"net/url"
	"testing"

	"opennamu/route"
	"opennamu/route/tool"
)

func Test_view_thread_post_flow(t *testing.T) {
	db := apiSlowDB(t)
	defer tool.DB_close(db)
	bbsSlowBoardSeed(db, t, "-1", 0)

	config := tool.Config{IP: "admin_ip", Session: &emptySession{}}
	new_topic_values := url.Values{
		"content": {"initial thread content"},
		"topic":   {"Thread view flow"},
		"title":   {"Thread view flow test"},
	}
	new_topic_result := route.View_thread(config, "0", "", "", new_topic_values)
	new_topic_target, ok := tool.Get_redirect_target(new_topic_result)
	if !ok || new_topic_target != "/bbs/w/-1/1#1" {
		t.Fatalf("new topic result target = %q (redirect=%v), want /bbs/w/-1/1#1", new_topic_target, ok)
	}

	comment_values := url.Values{"content": {"thread comment content"}}
	comment_result := route.View_thread(config, "1", "", "", comment_values)
	comment_target, ok := tool.Get_redirect_target(comment_result)
	if !ok || comment_target != "/bbs/w/-1/1#2" {
		t.Fatalf("thread comment result target = %q (redirect=%v), want /bbs/w/-1/1#2", comment_target, ok)
	}

	empty_topic_result := route.View_thread(config, "0", "", "", url.Values{})
	empty_topic_target, ok := tool.Get_redirect_target(empty_topic_result)
	if !ok || empty_topic_target != "/bbs/in/-1" {
		t.Fatalf("empty new-topic result target = %q (redirect=%v), want /bbs/in/-1", empty_topic_target, ok)
	}

	topic_page_result := route.View_thread(config, "1", "", "", nil)
	topic_page_target, ok := tool.Get_redirect_target(topic_page_result)
	if !ok || topic_page_target != "/bbs/w/-1/1" {
		t.Fatalf("topic page result target = %q (redirect=%v), want /bbs/w/-1/1", topic_page_target, ok)
	}
}
