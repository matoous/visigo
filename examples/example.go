package main

import (
	"fmt"
	"github.com/matoous/visigo"
	"net/http"
)

func main() {
	finalHandler := http.HandlerFunc(final)

	http.Handle("/", visigo.Counter(finalHandler))
	http.ListenAndServe(":3000", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	count, _ := visigo.Visits(r.URL)
	response := fmt.Sprintf("This page was viewed %d times", count)
	w.Write([]byte(response))
}
