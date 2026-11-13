package route

import (
	"database/sql"

	"opennamu/route/tool"
)

func Bbs_vote_controls_html(db *sql.DB, action string, my_vote string, up_count string, down_count string) string {
	if up_count == "" {
		up_count = "0"
	}
	if down_count == "" {
		down_count = "0"
	}

	return `<form method="post" action="` + tool.HTML_escape(action) + `" style="display: inline;">` +
		Bbs_vote_button_html(db, "up", my_vote, up_count) +
		Bbs_vote_button_html(db, "down", my_vote, down_count) +
		`</form>`
}

func Bbs_vote_button_html(db *sql.DB, vote_type string, my_vote string, count string) string {
	language_key := "upvote"
	cancel_language_key := "cancel_upvote"
	if vote_type == "down" {
		language_key = "downvote"
		cancel_language_key = "cancel_downvote"
	}

	vote_label := tool.Get_language(db, language_key, true)
	accessible_label := vote_label
	selected := my_vote == vote_type || my_vote == "both"
	if selected {
		accessible_label = tool.Get_language(db, cancel_language_key, true)
	} else if my_vote != "" && my_vote != "none" {
		accessible_label = tool.Get_language(db, "change_vote", true) + ": " + vote_label
	}

	selected_attributes := ` aria-pressed="false"`
	if selected {
		selected_attributes = ` aria-pressed="true" style="font-weight: bold; outline: 2px solid currentColor;"`
	}

	return `<button class="__ON_BUTTON__" type="submit" name="vote_type" value="` + vote_type + `" aria-label="` + tool.HTML_escape(accessible_label) + `" title="` + tool.HTML_escape(accessible_label) + `"` + selected_attributes + `>` + vote_label + ` ` + tool.HTML_escape(count) + `</button>`
}
