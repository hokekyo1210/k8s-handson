package functions

import (
	"net/http"

	"git.dmm.com/tsuchida-yuki1/cloud-functions-go/function"
)

func TimeUTC(w http.ResponseWriter, r *http.Request) {
	function.TimeUTC(w, r)
}

func TimeJST(w http.ResponseWriter, r *http.Request) {
	function.TimeJST(w, r)
}
