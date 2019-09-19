package main

import (
	"io"
	"net/http"

	"github.com/pkg/errors"
)

type ReverseProxy struct {
	Director func(*http.Request) error
}

func (rp *ReverseProxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	err := rp.Director(r)
	if err != nil {
		err = errors.Wrap(err, "error changing request")
		w.Write([]byte(err.Error()))
		return
	}

	client := &http.Client{}
	if r.URL.Scheme == "" {
		r.URL.Scheme = "http://"
	}
	u := r.URL.Scheme + r.Host + r.URL.Path
	req, err := http.NewRequest(r.Method, u, r.Body)
	if err != nil {
		err = errors.Wrap(err, "error creating request")
		w.Write([]byte(err.Error()))
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		err = errors.Wrap(err, "error making request")
		w.Write([]byte(err.Error()))
		return
	}

	_, err = io.Copy(w, resp.Body)
	if err != nil {
		err = errors.Wrap(err, "error sending sending response to client")
		w.Write([]byte(err.Error()))
		return
	}

}
