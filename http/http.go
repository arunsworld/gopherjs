package http

import (
	"github.com/gopherjs/gopherjs/js"
	"honnef.co/go/js/xhr"
)

// Options passed to HTTP calls
type Options struct {
	ResponseType string
	Headers      map[string]string
}

// Response is the response from a call
type Response struct {
	Status     int
	XhrRequest *xhr.Request
	Error      error
}

// GET HTTP call
func GET(url string, params map[string]string, options Options, response chan Response) {
	if params != nil {
		u := js.Global.Get("URLSearchParams").New()
		for k, v := range params {
			u.Call("set", k, v)
		}
		url += "?" + u.String()
	}
	call("GET", url, "", options, response)
}

// POST HTTP call
func POST(url string, data string, options Options, response chan Response) {
	call("POST", url, data, options, response)
}

func call(method string, url string, data string, options Options, response chan Response) {
	req := xhr.NewRequest(method, url)
	responseType := options.ResponseType
	if responseType == "" {
		responseType = xhr.Text
	}
	req.ResponseType = responseType
	if options.Headers != nil {
		for k, v := range options.Headers {
			req.SetRequestHeader(k, v)
		}
	}
	err := req.Send(data)
	if err != nil {
		response <- Response{Status: req.Status, Error: err}
		return
	}
	response <- Response{Status: req.Status, XhrRequest: req}
}
