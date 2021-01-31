package function

import (
	"fmt"
	"net/http"
	"time"
)

func TimeJST(w http.ResponseWriter, r *http.Request) {
	jst, _ := time.LoadLocation("Asia/Tokyo")
	response := time.Now().In(jst)
	fmt.Fprint(w, response)
}
