package mod_http

import (
	"net/http"
	"time"
	"fmt"
	"os"
	"reflect"
	"strings"
)

type http_access_log struct {
	startTime time.Time
}

func NewHttpServing() *http_access_log {
	return &http_access_log{time.Now()}
}

func (log *http_access_log) End(w http.ResponseWriter, r *http.Request) {
	endTime := time.Now()
	duration := endTime.Sub(log.startTime)

	var addr string
	pos := strings.LastIndex(r.RemoteAddr, ":")
	if pos >= 0 {
		addr = r.RemoteAddr[:pos]
	} else {
		addr = r.RemoteAddr
	}
	user := "-"
	if r.URL.User != nil {
		u := r.URL.User.Username()
		if u != "" {
			user = fmt.Sprintf("\"%s\"", u)
		}
	}
	referer := r.Header.Get("Referer")
	if referer == "" {
		referer = "-"
	}
	userAgent := r.Header.Get("User-Agent")
	if userAgent == "" {
		userAgent = "-"
	}
	xForwardFor := r.Header.Get("X-Forwarded-For")
	if xForwardFor == "" {
		xForwardFor = "-"
	}

	respP := reflect.ValueOf(w)
	resp := respP.Elem()
	status := resp.FieldByName("status")
	written := resp.FieldByName("written")
	fmt.Fprintf(os.Stderr, "%s - %s [%s] \"%s %s %s\" %v %v \"%s\" \"%s\" \"%s\" %v\n",
        addr,
		user,
        log.startTime.Format("2/Jan/2006:15:04:05 -0700"),
        r.Method, r.RequestURI, r.Proto,
        status,
        written,
		referer,
		userAgent,
		xForwardFor,
        duration,
    )
}

func (log *http_access_log) formatTime(format string) string {
	if format == "" {
		format = "2006-01-02 15:04:05.000"
	}
	return log.startTime.Format(format)
}
