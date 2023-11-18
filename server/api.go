package server

import (
	"net/http"
)

func Run() {
	r := newRouter()
	http.ListenAndServe(":3000", r)
}
