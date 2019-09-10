package main

import (
	"fmt"
	"net/http"

	"github.com/matoous/visigo"
)

func main() {
	http.Handle("/", visigo.Counter(http.HandlerFunc(final)))
	http.Handle("/total", visigo.Counter(http.HandlerFunc(total)))
	_ = http.ListenAndServe(":3000", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	count, _ := visigo.Visits(r)
	response := fmt.Sprintf("This page was viewed by %d unique visitors", count)
	_, _ = w.Write([]byte(response))
}

func total(w http.ResponseWriter, _ *http.Request) {
	count, _ := visigo.TotalVisits()
	response := fmt.Sprintf("This website had %d unique visitors in total", count)
	_, _ = w.Write([]byte(response))
}
