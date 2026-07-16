package main

import "github.com/gin-gonic/gin"

func register_routes(r *gin.Engine) {
	register_api_routes(r)
	register_watch_routes(r)
	register_list_routes(r)
	register_auth_routes(r)
	register_wiki_routes(r)
	register_recent_routes(r)
	register_bbs_routes(r)
	register_history_edit_routes(r)
	register_search_routes(r)
}
