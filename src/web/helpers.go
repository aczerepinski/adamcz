package web

import (
	"net/http"
)

func getParam(param string, r *http.Request) string {
	r.ParseForm()
	return r.Form.Get(param)
}
