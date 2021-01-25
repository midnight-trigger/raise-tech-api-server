package http_client

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/imroc/req"
	"github.com/spf13/viper"
)

// HTTPClient struct : receiver
type HTTPClient struct {
	URI       string
	Params    *map[string]string
	Body      interface{}
	UseProxy  bool
	UserAgent string
}

// Get method : HTTP GET request, return response object
func (u *HTTPClient) Get() (resp *http.Response, err error) {
	r := req.New()

	// set request data
	params, err := u.setRequest(r)
	if err != nil {
		return
	}

	req, err := r.Get(u.URI, params)
	if err != nil {
		return
	}

	resp = req.Response()
	return
}

//Post method : HTTP POST request
func (u *HTTPClient) Post() (resp *http.Response, err error) {
	r := req.New()

	// set request data
	params, err := u.setRequest(r)
	if err != nil {
		return
	}

	req, err := r.Post(u.URI, params, req.BodyJSON(&u.Body))
	if err != nil {
		return
	}

	resp = req.Response()
	return
}

// setRequest private method : create HTTP request
func (u *HTTPClient) setRequest(r *req.Req) (params req.Param, err error) {
	// set params
	if u.Params != nil {
		for key, value := range *u.Params {
			params[key] = value
		}
	}

	// set proxy
	if u.UseProxy {
		proxyURL := viper.GetString("proxy.url")
		proxyPort := viper.GetString("proxy.port")
		r.SetProxy(func(r *http.Request) (*url.URL, error) {
			u, _ := url.ParseRequestURI(fmt.Sprintf("%s:%s", proxyURL, proxyPort))
			return u, nil
		})
	}
	return
}
