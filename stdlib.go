package main

import (
	"net/http"
	"net/url"
)

func AppendSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || r.URL.Path[len(r.URL.Path)-1] != '/' {
			r2 := new(http.Request)
			*r2 = *r
			r2.URL = new(url.URL)
			*r2.URL = *r.URL
			r2.URL.Path += "/"
			next.ServeHTTP(w, r2)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}

func RedirectSlashMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "" || r.URL.Path[len(r.URL.Path)-1] != '/' {
			http.Redirect(w, r, r.URL.Path+"/", http.StatusMovedPermanently)
		} else {
			next.ServeHTTP(w, r)
		}
	})
}
