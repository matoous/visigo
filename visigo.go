package visigo

import (
	"errors"
	"net/http"

	"github.com/axiomhq/hyperloglog"
	"github.com/tomasen/realip"
)

var counter map[string]*hyperloglog.Sketch

// ErrCount - error returned when you try to get count but didn't register middleware
var ErrCount = errors.New("count not found or error in HyperLogLog")

// Visits - get visits for given URL
func Visits(r *http.Request) (uint64, error) {
	if counter == nil {
		// no, you didn't ...
		panic("you need to register Visigo Counter first!")
	}
	if hll, found := counter[r.URL.String()]; found {
		return hll.Estimate(), nil
	}
	return 0, ErrCount
}

// TotalVisits gets total visits to all sites
func TotalVisits() (uint64, error) {
	hll := hyperloglog.New()
	for _, s := range counter {
		if err := hll.Merge(s); err != nil {
			return 0, err
		}
	}
	return hll.Estimate(), nil
}

// Counter - registers middleware for visits counting
func Counter(next http.Handler) http.Handler {
	counter = make(map[string]*hyperloglog.Sketch)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hll, found := counter[r.URL.String()]; found {
			hll.Insert([]byte(realip.RealIP(r)))
		} else {
			l := hyperloglog.New()
			l.Insert([]byte(realip.RealIP(r)))
			counter[r.URL.String()] = l
		}
		next.ServeHTTP(w, r)
	})
}
