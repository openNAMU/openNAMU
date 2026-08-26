package main

import (
	"crypto/md5"
	_ "embed"
	"encoding/hex"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strconv"
	"time"

	"opennamu/route"
	"opennamu/route/tool"

	"github.com/flosch/pongo2/v6"
	"github.com/gin-gonic/gin"
)

var dev_mode = false

//go:embed version.json
var builtin_version_json []byte

func error_handler() gin.HandlerFunc {
	return func(c *gin.Context) {
		defer func() {
			if r := recover(); r != nil {
				err, ok := r.(error)
				if !ok {
					err = fmt.Errorf("%v", r)
				}

				stack_trace := debug.Stack()
				log.Printf("Panic recovered: %v\n%s", err, stack_trace)

				if dev_mode {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"response": "error",
						"error":    err.Error(),
						"stack":    string(stack_trace),
					})
				} else {
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"response": "error",
						"error":    "Internal Server Error",
					})
				}
			}
		}()

		c.Next()
	}
}

func pongo_init() {
	pongo2.RegisterFilter("md5_replace", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		h := md5.Sum([]byte(in.String()))

		return pongo2.AsValue(hex.EncodeToString(h[:])), nil
	})

	pongo2.RegisterFilter("load_lang", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		db := tool.DB_connect()
		defer tool.DB_close(db)

		return pongo2.AsValue(tool.Get_language(db, in.String(), false)), nil
	})

	pongo2.RegisterFilter("cut_100", func(in *pongo2.Value, param *pongo2.Value) (*pongo2.Value, *pongo2.Error) {
		data := in.String()
		data = tool.Get_slice(data, 0, 100)

		return pongo2.AsValue(data), nil
	})
}

func wait_startup_delay() {
	delay_text := os.Getenv("NAMU_START_DELAY_MS")
	os.Unsetenv("NAMU_START_DELAY_MS")

	delay, err := strconv.Atoi(delay_text)
	if err != nil || delay <= 0 {
		return
	}
	if delay > 10000 {
		delay = 10000
	}

	time.Sleep(time.Duration(delay) * time.Millisecond)
}

func main() {
	if len(os.Args) > 1 && os.Args[1] == "--opennamu-update" {
		os.Exit(route.Run_server_update(os.Args[2:]))
	}
	if len(os.Args) > 1 && os.Args[1] == "em" {
		os.Exit(route.Run_emergency_tool(os.Args[2:]))
	}

	wait_startup_delay()
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	port := "3000"
	if len(os.Args) > 1 {
		port = os.Args[1]
	}

	var r *gin.Engine
	if len(os.Args) > 2 && os.Args[2] == "dev" {
		dev_mode = true
		r = gin.Default()
	} else {
		gin.SetMode(gin.ReleaseMode)
		r = gin.New()
	}

	tool.Set_builtin_version_data(builtin_version_json)
	tool.Main_init()

	r.Use(error_handler())
	r.Use(tool.Session_middleware())
	r.Use(wiki_access_middleware())
	pongo_init()

	register_routes(r)

	log.Default().Println("Run in http://127.0.0.1:" + port)
	if err := r.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
