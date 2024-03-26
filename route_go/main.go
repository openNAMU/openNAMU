package main

import (
	"log"
	"opennamu/route"
	"os"
)

func main() {
	call_arg := os.Args[1:]

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	if call_arg[0] == "main_func_easter_egg" {
		route.Main_func_easter_egg()
	} else if call_arg[0] == "api_w_raw" {
		route.Api_w_raw(call_arg[1:])
	} else if call_arg[0] == "api_func_sha224" {
		route.Api_func_sha224(call_arg[1:])
	} else if call_arg[0] == "api_w_random" {
		route.Api_w_random(call_arg[1:])
	} else if call_arg[0] == "api_search" {
		route.Api_search(call_arg[1:])
	} else if call_arg[0] == "api_topic" {
		route.Api_thread(call_arg[1:])
	} else if call_arg[0] == "api_func_ip" {
		route.Api_func_ip(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_change" {
		route.Api_list_recent_change(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_edit_request" {
		route.Api_list_recent_edit_request(call_arg[1:])
	} else if call_arg[0] == "api_bbs" {
		route.Api_bbs(call_arg[1:])
	} else if call_arg[0] == "api_w_xref" {
		route.Api_w_xref(call_arg[1:])
	} else if call_arg[0] == "api_w_watch_list" {
		route.Api_w_watch_list(call_arg[1:])
	} else if call_arg[0] == "api_user_watch_list" {
		route.Api_user_watch_list(call_arg[1:])
	} else if call_arg[0] == "api_w_render" {
		route.Api_w_render(call_arg[1:])
	} else if call_arg[0] == "api_func_llm" {
		route.Api_func_llm(call_arg[1:])
	} else if call_arg[0] == "api_func_language" {
		route.Api_func_language(call_arg[1:])
	} else if call_arg[0] == "api_func_auth_list" {
		route.Api_func_auth_list(call_arg[1:])
	} else if call_arg[0] == "api_list_recent_discuss" {
		route.Api_list_recent_discuss(call_arg[1:])
	} else if call_arg[0] == "api_bbs_list" {
		route.Api_bbs_list(call_arg[1:])
	} else {
		log.Fatal("404")
	}
}
