package main

import (
	"github.com/matoous/visigo"
	"net/http"
	"fmt"
)

func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", visigo.Counter(finalHandler))
	http.ListenAndServe(":3000", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	count := visigo.Visits(r)
	response := fmt.Sprintf("This page was viewed %d times", count)
	w.Write([]byte(response))
}