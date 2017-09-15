package visigo

import (
	"github.com/clarkduvall/hyperloglog"
	"github.com/go-errors/errors"
	"github.com/tomasen/realip"
	"hash/fnv"
	"net/http"
	"net/url"
)

const (
	defaultPrecision = 18
)

type hashableIp struct {
	realIp []byte
}

func (hip *hashableIp) Sum64() uint64 {
	h := fnv.New64a()
	h.Write(hip.realIp)
	return h.Sum64()
}

var counter map[*url.URL]*hyperloglog.HyperLogLogPlus

// VisigoError - error returned when you try to get count but didn't register middleware
var ErrCount = errors.New("Count not found or error in HyperLogLog")

// Visits - get visits for given URL
func Visits(u *url.URL) (uint64, error) {
	if counter == nil {
		// no, you didn't ...
		panic("You need to register Visigo Counter first!")
	}
	if hll, found := counter[u]; found {
		return hll.Count(), nil
	}
	return 0, ErrCount
}

// Counter - registers middleware for visits counting
func Counter(next http.Handler) http.Handler {
	counter = make(map[*url.URL]*hyperloglog.HyperLogLogPlus)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if hll, found := counter[r.URL]; !found {
			// get hyperloglog or fail silently
			if l, err := hyperloglog.NewPlus(defaultPrecision); err == nil {
				ip := &hashableIp{
					realIp: []byte(realip.RealIP(r)),
				}
				l.Add(ip)
				counter[r.URL] = l
			}
		} else {
			// it's perfectly fine to omit map assignment since it is a pointer
			ip := &hashableIp{
				realIp: []byte(realip.RealIP(r)),
			}
			hll.Add(ip)
		}
		// serve
		next.ServeHTTP(w, r)
	})
}
