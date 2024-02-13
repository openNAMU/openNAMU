package main

import (
	"opennamu/route"
	"os"
)

func main() {
	call_arg := os.Args[1:]

	if call_arg[0] == "main_func_easter_egg" {
		route.Main_func_easter_egg()
	} else if call_arg[0] == "api_w_raw" {
		route.Api_w_raw(call_arg[1:])
	} else if call_arg[0] == "api_func_sha224" {
		route.Api_func_sha224(call_arg[1:])
	} else if call_arg[0] == "view_random" {
		route.View_random(call_arg[1:])
	} else if call_arg[0] == "api_search" {
		route.Api_search(call_arg[1:])
	}
}
