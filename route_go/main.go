package main

import (
	"fmt"
	"log"
	"opennamu/route"
	"os"
)

func main() {
	call_arg := os.Args[1:]

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	var route_data string
	if call_arg[0] == "main_func_easter_egg" {
		route_data = route.Main_func_easter_egg()
	} else if call_arg[0] == "api_w_raw" {
		route_data = route.Api_w_raw(call_arg[1:])
	} else if call_arg[0] == "api_func_sha224" {
		route_data = route.Api_func_sha224(call_arg[1:])
	} else if call_arg[0] == "api_w_random" {
		route_data = route.Api_w_random(call_arg[1:])
	} else if call_arg[0] == "api_search" {
		route_data = route.Api_search(call_arg[1:])
	} else if call_arg[0] == "api_topic" {
		route_data = route.Api_topic(call_arg[1:])
	} else if call_arg[0] == "api_func_ip" {
		route_data = route.Api_func_ip(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_change" {
		route_data = route.Api_list_recent_change(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_edit_request" {
		route_data = route.Api_list_recent_edit_request(call_arg[1:])
	} else if call_arg[0] == "api_bbs" {
		route_data = route.Api_bbs(call_arg[1:])
	} else if call_arg[0] == "api_w_xref" {
		route_data = route.Api_w_xref(call_arg[1:])
	} else if call_arg[0] == "api_w_watch_list" {
		route_data = route.Api_w_watch_list(call_arg[1:])
	} else if call_arg[0] == "api_user_watch_list" {
		route_data = route.Api_user_watch_list(call_arg[1:])
	} else if call_arg[0] == "api_w_render" {
		route_data = route.Api_w_render(call_arg[1:])
	} else if call_arg[0] == "api_func_llm" {
		route_data = route.Api_func_llm(call_arg[1:])
	} else if call_arg[0] == "api_func_language" {
		route_data = route.Api_func_language(call_arg[1:])
	} else if call_arg[0] == "api_func_auth_list" {
		route_data = route.Api_func_auth_list(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_discuss" {
		route_data = route.Api_list_recent_discuss(call_arg[1:])
	} else if call_arg[0] == "api_bbs_list" {
		route_data = route.Api_bbs_list(call_arg[1:])
	} else if call_arg[0] == "api_list_old_page" {
		route_data = route.Api_list_old_page(call_arg[1:])
	} else if call_arg[0] == "api_topic_list" {
		route_data = route.Api_topic_list(call_arg[1:])
	} else if call_arg[0] == "api_bbs_w_comment_n" {
		route_data = route.Api_bbs_w_comment(call_arg[1:])
	} else if call_arg[0] == "api_bbs_w_n" {
		route_data = route.Api_bbs_w(call_arg[1:])
	} else if call_arg[0] == "api_w_set_reset" {
		route_data = route.Api_w_set_reset(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_block" {
		route_data = route.Api_list_recent_block(call_arg[1:])
	} else if call_arg[0] == "api_list_title_index" {
		route_data = route.Api_list_title_index(call_arg[1:])
	} else if call_arg[0] == "api_user_setting_editor_post" {
		route_data = route.Api_user_setting_editor_post(call_arg[1:])
	} else if call_arg[0] == "api_user_setting_editor_delete" {
		route_data = route.Api_user_setting_editor_delete(call_arg[1:])
	} else if call_arg[0] == "api_user_setting_editor" {
		route_data = route.Api_user_setting_editor(call_arg[1:])
	} else if call_arg[0] == "api_setting" {
		route_data = route.Api_setting(call_arg[1:])
	} else if call_arg[0] == "api_setting_put" {
		route_data = route.Api_setting_put(call_arg[1:])
	} else if call_arg[0] == "api_func_ip_menu" {
		route_data = route.Api_func_ip_menu(call_arg[1:])
	} else if call_arg[0] == "api_func_ip_post" {
		route_data = route.Api_func_ip_post(call_arg[1:])
	} else if call_arg[0] == "api_func_acl_list" {
		route_data = route.Api_func_acl_list(call_arg[1:])
	} else {
		log.Fatal(call_arg[0] + " is 404")
	}

	fmt.Print(route_data)
}
