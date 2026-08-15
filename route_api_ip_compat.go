package main

import (
	"strconv"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func compat_ip_post(c *gin.Context) {
	compat_api_data(c, route.Api_func_ip_post(make_route_config(c), compat_ip_post_data(c)))
}

func compat_ip_post_data(c *gin.Context) []string {
	data := c.PostFormArray("data")
	if len(data) == 0 {
		data = c.PostFormArray("data[]")
	}
	if len(data) > 0 {
		return data
	}

	for index := 1; ; index++ {
		value := c.PostForm("data_" + strconv.Itoa(index))
		if value == "" {
			break
		}
		data = append(data, value)
	}

	return data
}
