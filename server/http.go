package server

import (
	"fmt"
	"net/http"
)

// when user enteres long url save in db
func ShortURLrouting(w http.ResponseWriter, req *http.Request) {
	fmt.Fprintf(w, "This should render actual long url")
}

//TODO swaps shorturl link with acutal one
