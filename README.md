# Visigo


[![Build Status](https://github.com/matoous/visigo/workflows/Tests/badge.svg)](https://github.com/matoous/visigo/actions) 
[![Build Status](https://github.com/matoous/visigo/workflows/Lint/badge.svg)](https://github.com/matoous/visigo/actions) 
[![GoDoc](https://godoc.org/github.com/matoous/visigo?status.svg)](https://godoc.org/github.com/matoous/visigo)
[![Go Report Card](https://goreportcard.com/badge/github.com/matoous/visigo)](https://goreportcard.com/report/github.com/matoous/visigo)
[![GitHub issues](https://img.shields.io/github/issues/matoous/visigo.svg)](https://github.com/matoous/visigo/issues)
[![License](https://img.shields.io/badge/license-MIT%20License-blue.svg)](https://github.com/matoous/visigo/LICENSE)


**Visigo** is http middleware for page unique visits counting. It uses HyperLogLog as 
a counter, so it's pretty fast.

**Warning:** Visigo stores HyperLogLog++ in *map*, so this implementation
should be used only on smaller sites.

## HyperLogLog++

[HyperLogLog++ paper](http://research.google.com/pubs/pub40671.html)  
[Google article about HyperLogLog++](https://research.neustar.biz/2013/01/24/hyperloglog-googles-take-on-engineering-hll/)

From [Wikipedia](https://en.wikipedia.org/wiki/HyperLogLog)  

> HyperLogLog is an algorithm for the count-distinct problem, approximating the number of distinct elements in a multiset.
Calculating the exact cardinality of a multiset requires an amount of memory proportional to the cardinality, which is impractical for very large data sets. Probabilistic cardinality estimators, such as the HyperLogLog algorithm, use significantly less memory than this, at the cost of obtaining only an approximation of the cardinality. The HyperLogLog algorithm is able to estimate cardinalities of > 109 with a typical accuracy of 2%, using 1.5 kB of memory.
 HyperLogLog is an extension of the earlier LogLog algorithm, itself deriving from the 1984 Flajolet–Martin algorithm.

## Install

Via go get tool

``` bash
$ go get github.com/matoous/visigo
```

## Usage


``` go
package main

import (
	"fmt"
	"net/http"

	"github.com/matoous/visigo"
)

func main() {
	http.Handle("/", visigo.Counter(http.HandlerFunc(final)))
	http.Handle("/total", visigo.Counter(http.HandlerFunc(total)))
	http.ListenAndServe(":3000", nil)
}

func final(w http.ResponseWriter, r *http.Request) {
	count, _ := visigo.Visits(r)
	response := fmt.Sprintf("This page was viewed by %d unique visitors", count)
	w.Write([]byte(response))
}

func total(w http.ResponseWriter, r *http.Request) {
	count, _ := visigo.TotalVisits()
	response := fmt.Sprintf("This website had %d unique visitors in total", count)
	w.Write([]byte(response))
}
```

## Testing

``` bash
$ go test ./...
```

## Notice

If you use **Visigo** on your site or in your project, please let me know!

If you have any issues, just feel free and open it in this repository, thanks!

## License

The MIT License (MIT). Please see [License File](LICENSE.md) for more information.
