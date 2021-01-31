package function

import (
	"fmt"
	"net/http"
	"time"
)

func TimeUTC(w http.ResponseWriter, r *http.Request) {
	response := time.Now()
	fmt.Fprint(w, response)
}
