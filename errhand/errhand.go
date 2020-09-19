package errhand

import (
	"fmt"
	"log"
	"net/http"
	"runtime"
	"strconv"
)

// InternalError writes error in log and page.
func InternalError(err error, w http.ResponseWriter) {
	pc := make([]uintptr, 15)
	n := runtime.Callers(2, pc)
	frames := runtime.CallersFrames(pc[:n])
	frame, _ := frames.Next()

	log.Println(frame.Function+"():"+strconv.Itoa(frame.Line), "INTERNAL ERROR:", err)

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "INTERNAL ERROR. Please try later")
}
