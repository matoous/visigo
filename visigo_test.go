package visigo

import (
	"net/http"
	"testing"
	"time"
	"math/rand"
	"fmt"
	"net/http/httptest"
	"net/url"
	"log"
	"github.com/stretchr/testify/assert"
)

func TestPanics(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		}
	}()

	req, err := http.NewRequest("GET", "/", nil)
	if err != nil {
		t.Fatal(err)
	}
	Visits(req.URL)
}

func TestVisits(t *testing.T) {
	// seed random to use instead of slower crypto
	rand.Seed(time.Now().UnixNano())

	// random ip generating function
	randomIp := func() string {
		tokens := make([]byte, 4)
		rand.Read(tokens)
		return fmt.Sprintf("%v.%v.%v.%v", tokens[0], tokens[1], tokens[2], tokens[3])
	}

	// some url
	uri, err := url.Parse("/")
	if err != nil {
		log.Fatal(err)
	}

	closeTo := func(num uint64, to uint64) bool {
		return num > (to - ((to/100)*2)) && num < (to + ((to/100)*2))
	}

	// new request generating function
	newRequest := func() *http.Request {
		return &http.Request{
			RemoteAddr: randomIp(),
			URL: uri,
		}
	}

	handler := Counter(http.HandlerFunc(func (w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()

	// test
	for i := 0; i < 100000; i++ {
		handler.ServeHTTP(rec, newRequest())
	}
	cnt, err := Visits(uri)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, closeTo(cnt, 100000), true, fmt.Sprintf("Excpected: %v visits, got: %v", cnt, 100000))
}