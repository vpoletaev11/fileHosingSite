package errhand

import (
	"fmt"
	"log"
	"net/http"
)

// InternalError writes error in log and page.
func InternalError(pkg, fnc, username string, err error, w http.ResponseWriter) {
	if username == "" {
		log.Println(pkg+"."+fnc+"()", "INTERNAL ERROR:", err)
	} else {
		log.Println(pkg+"."+fnc+"()", "[user:", username+"]", "INTERNAL ERROR:", err)
	}

	w.WriteHeader(http.StatusInternalServerError)
	fmt.Fprintln(w, "INTERNAL ERROR. Please try later")
}
