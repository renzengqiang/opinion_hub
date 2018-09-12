package handler

import (
	"net/http"
	"runtime/debug"
	"github.com/go-martini/martini"
	"github.com/astaxie/beego/logs"
	"github.com/martini-contrib/render"
	"fmt"
)

func Recovery() martini.Handler {
	return func(res http.ResponseWriter, c martini.Context, logger *logs.BeeLogger, r render.Render) {
		defer func() {
			if err := recover(); err != nil {
				logger.Error("PANIC: %s\n%s", err, debug.Stack())
				type Resp struct {
					Env string
					Code int32
					Status string
					Message string
					Data []int64
				}
				resp := Resp{Env:martini.Env, Code:400, Status:"False", Message:fmt.Sprintf("Panic error. msg: %s", err)}
				r.JSON(200, resp)
			}
		}()
	}
}
