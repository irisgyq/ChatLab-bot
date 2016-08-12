package main

import (
    "log"
    "net/http"
    "os"

)

func main() {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }

    http.HandleFunc("/webhook", Handle)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}

func Handle(rw http.ResponseWriter, req *http.Request) {
    if req.Method == "GET" {
        query := req.URL.Query()
	if query.Get("hub.verify_token") != os.Getenv("VALIDATION_TOKEN") {
	    rw.WriteHeader(http.StatusUnauthorized)
	    return
        }
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte(query.Get("hub.challenge")))
    } else if req.Method == "POST" {
	HandlePOST(rw, req)
    } else {
	rw.WriteHeader(http.StatusMethodNotAllowed)
    }
}

func HandlePOST(rw http.ResponseWriter, req *http.Request) {

}
