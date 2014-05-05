package router

import (
	"fmt"
	"io"
	"net/http"
	"strings"
)

type Resquestor struct {
	Header http.Header
	host   string
	routes Routes
}

func NewRequestor(host string, routes Routes) *Resquestor {
	return &Resquestor{
		host:   host,
		routes: routes,
	}
}

func (r *Resquestor) RequestForHandler(
	handler string,
	params Params,
	body io.Reader,
) (*http.Request, error) {
	route, ok := r.routes.RouteForHandler(handler)
	if !ok {
		return &http.Request{}, fmt.Errorf("No route exists for handler %", handler)
	}
	path, err := route.PathWithParams(params)
	if err != nil {
		return &http.Request{}, err
	}

	url := r.host + "/" + strings.TrimLeft(path, "/")

	req, err := http.NewRequest(route.Method, url, body)
	if err != nil {
		return &http.Request{}, err
	}

	for key, values := range r.Header {
		req.Header[key] = []string{}
		copy(req.Header[key], values)
	}
	return req, nil
}
