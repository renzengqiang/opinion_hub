package handler

import (
	"github.com/go-martini/martini"
	"github.com/astaxie/beego/logs"
	"time"
	"net/http"
)

func Logger() martini.Handler {
	return func(res http.ResponseWriter, req *http.Request, c martini.Context, logger *logs.BeeLogger) {
		start := time.Now()

		addr := req.Header.Get("X-Real-IP")
		if addr == "" {
			addr = req.Header.Get("X-Forwarded-For")
			if addr == "" {
				addr = req.RemoteAddr
			}
		}

		logger.Info("Started %s %s for %s", req.Method, req.URL.Path, addr)
		rw := res.(martini.ResponseWriter)
		c.Next()
		logger.Info("Completed %v %s in %v\n", rw.Status(), http.StatusText(rw.Status()), time.Since(start))
	}
}
