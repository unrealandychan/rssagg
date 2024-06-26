package main

import "net/http"

func handlerError(w http.ResponseWriter, r *http.Request) {
	responseWithError(w, http.StatusInternalServerError, "An error occurred")
}
