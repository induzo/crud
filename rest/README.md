# REST wrapper (JSON)

## Implementation

Consider the example in the mock folder as your entity manager.
In 5 lines, you can implement the REST handlers for your API:

```golang
package main

import (
    "github.com/go-chi/chi"
    "github.com/induzo/rest"
    "github.com/induzo/mock"
)

func main() {
    r := chi.NewRouter()

    // Subrouters:
    m := mock.NewMgr()
    r.Route("/e", func(r chi.Router) {
        r.Get("/", rest.GETListHandler(m))
        r.Post("/", rest.POSTHandler(m))
        r.Get("/{ID}", rest.GETHandler(m))
        r.Patch("/{ID}", rest.PATCHHandler(m))
    })


    srv := &http.Server{
       ReadTimeout:  10 * time.Second,
       WriteTimeout: 7200 * time.Second,
       IdleTimeout:  10 * time.Second,
       Addr:         8080,
       Handler:      r,
    }

    ListenAndServe(srv, conf.ForceInsecureTLS)
}
```

## Benchmarks (i7, 16GB)

```bash
    goos: linux
    goarch: amd64
    pkg: github.com/induzo/crud/rest
    BenchmarkPOSTHandler-8            300000              4190 ns/op            1912 B/op         20 allocs/op
    BenchmarkGETListHandler-8         300000              4694 ns/op            2131 B/op         31 allocs/op
    BenchmarkGETHandler-8             500000              2856 ns/op            1138 B/op         14 allocs/op
    BenchmarkDELETEHandler-8         3000000               596 ns/op              80 B/op          2 allocs/op
    BenchmarkPUTHandler-8             300000              4241 ns/op            1842 B/op         21 allocs/op
    BenchmarkPATCHHandler-8           500000              2772 ns/op            1544 B/op         17 allocs/op
```
