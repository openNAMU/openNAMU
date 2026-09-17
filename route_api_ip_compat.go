package main

import (
	"strconv"

	"opennamu/route"

	"github.com/gin-gonic/gin"
)

func Compat_ip_post(c *gin.Context) {
	Compat_api_data(c, route.Api_func_ip_post(Make_route_config(c), Compat_ip_post_data(c)))
}

func Compat_ip_post_data(c *gin.Context) []string {
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
