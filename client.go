package todoist

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"gopkg.in/resty.v1"
)

// An endpoint represents a single call of the Todoist API. The PathParams will be a name -> description map of the expected path params, passed in
// at call time.
type endpoint struct {
	Path       string
	PathParams map[string]string
	Method     string
}

type todoistResponse struct {
	StatusCode int
	Body       []byte
}

func makeCall(token string, endpointName string, pathParams map[string]string, data interface{}) (todoistResponse, error) {
	result := todoistResponse{}

	// first, find the endpoint
	ep, epOK := endpoints[endpointName]
	if !epOK {
		return result, errors.New("endpoint not found")
	}

	if token == "" && config.AuthToken != "" {
		token = config.AuthToken
	}

	client := resty.New()
	r := client.R().SetAuthToken(token)

	// build the URL
	url := "https://api.todoist.com/rest/v1" + ep.Path
	for k, v := range pathParams {
		url = strings.Replace(url, ":"+k, v, -1)
	}

	// data depends entirely on the call
	if data != nil {
		if ep.Method == http.MethodGet || ep.Method == http.MethodDelete {
			// try to convert
			p, pOK := data.(map[string]string)
			if !pOK {
				return result, errors.New("for GET and DELETE, the body must be a map[string]string{}")
			}
			r.SetQueryParams(p)
		} else if ep.Method == http.MethodPost || ep.Method == http.MethodPut || ep.Method == http.MethodPatch {
			r.SetBody(data)
		}
	}
	var err error
	var resp *resty.Response
	switch ep.Method {
	case http.MethodGet:
		resp, err = r.Get(url)
	case http.MethodPost:
		resp, err = r.Post(url)
	case http.MethodPatch:
		resp, err = r.Patch(url)
	case http.MethodDelete:
		resp, err = r.Delete(url)
	case http.MethodPut:
		resp, err = r.Put(url)
	}

	if err != nil {
		return result, err
	}
	result.StatusCode = resp.StatusCode()

	if resp.StatusCode() == http.StatusBadRequest || resp.StatusCode() == http.StatusUnauthorized || resp.StatusCode() == http.StatusForbidden || resp.StatusCode() == http.StatusNotFound {
		// body is a string, but we need to flag an error
		err = fmt.Errorf(strings.TrimRight(string(resp.Body()), "\n"))
	}

	result.Body = resp.Body()
	return result, err
}

const (
	EndpointNameGetProjects = "GetProjects"
)

var endpoints = map[string]endpoint{
	EndpointNameGetProjects: {
		Path:       "/projects",
		PathParams: map[string]string{},
		Method:     http.MethodGet,
	},
}