package visigo

import (
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

// newRequest is helper function for testing, it creates new request from given IP to given URL.
func newRequest(ip string, url *url.URL) *http.Request {
	return &http.Request{
		RemoteAddr: ip,
		URL:        url,
	}
}

func TestPanics(t *testing.T) {
	uri, err := url.Parse("/1")
	if err != nil {
		assert.NoError(t, err, "must parse url")
	}

	assert.Panics(t, func() {
		_, _ = Visits(uri)
	}, "should panic if the middleware is not registered and count function is called")
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

	// accuracy better than 2%
	closeTo := func(num uint64, to uint64) bool {
		a := float32(num)
		b := float32(to)
		return a > (b-(b*0.02)) && a < (b+(b*0.02))
	}

	handler := Counter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
	rec := httptest.NewRecorder()

	// test up to 1 000 000
	var limit uint64 = 90
	for i := uint64(0); limit < 1000000; {
		for ; i < limit; i++ {
			handler.ServeHTTP(rec, newRequest(randomIp(), uri))
		}
		cnt, err := Visits(uri)
		if err != nil {
			log.Fatal(err)
		}
		assert.Equal(t, closeTo(cnt, limit), true, fmt.Sprintf("Excpected: %v visits, got: %v", cnt, limit))
		limit *= 9
	}
}

func TestTotalVisits(t *testing.T) {
	uri1, err := url.Parse("/1")
	if err != nil {
		assert.NoError(t, err, "must parse url")
	}
	uri2, err := url.Parse("/2")
	if err != nil {
		assert.NoError(t, err, "must parse url")
	}

	t.Run("counts multiple paths", func(t *testing.T) {
		handler := Counter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, newRequest("1.1.1.1", uri1))
		handler.ServeHTTP(rec, newRequest("2.2.2.2", uri2))
		total, err := TotalVisits()
		assert.NoError(t, err, "should not return error")
		assert.Equal(t, uint64(2), total, "should count distinct IPs over all paths")
	})
	t.Run("merges hyperloglogs", func(t *testing.T) {
		handler := Counter(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {}))
		rec := httptest.NewRecorder()
		handler.ServeHTTP(rec, newRequest("1.1.1.1", uri1))
		handler.ServeHTTP(rec, newRequest("1.1.1.1", uri2))
		total, err := TotalVisits()
		assert.NoError(t, err, "should not return error")
		assert.Equal(t, uint64(1), total, "should merge the logs and count only distinct IPs")
	})
}
