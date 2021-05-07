package main

import (
	"net/http"
	"regexp"
)

type Router struct {
	routes map[*regexp.Regexp]func(http.ResponseWriter, *http.Request)
}

func (self *Router) registerRoute(rxp *regexp.Regexp, handler func(http.ResponseWriter, *http.Request)) {
	self.routes[rxp] = handler
}

func (self *Router) ServeHTTP(response http.ResponseWriter, request *http.Request) {
	request_path := request.URL.Path
	for route, handler := range self.routes {
		if route.MatchString(request_path) {
			handler(response, request)
			return
		}
	}

	response.WriteHeader(404)
	response.Write([]byte("not found"))
}

func createRouter() *Router {
	var new_router *Router = new(Router)
	new_router.routes = make(map[*regexp.Regexp]func(http.ResponseWriter, *http.Request))

	return new_router
}
