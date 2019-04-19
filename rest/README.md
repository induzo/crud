# REST wrapper HTTP/JSON

Consider the example in mock_test.go as your entity manager.
In 5 lines, you can implement the REST handlers for your API:

```golang
package main

import "github.com/go-chi/chi"

func main() {
    r := chi.NewRouter()

    // Subrouters:
    m := newMgrMock()
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
