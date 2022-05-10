package librecat

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/ugent-library/biblio-backend/internal/jsonapi"
)

type requestPayload struct {
	Data interface{} `json:"data"`
}

type responsePayload struct {
	Data   json.RawMessage `json:"data"`
	Errors jsonapi.Errors  `json:"errors"`
}

func (c *Client) get(path string, qp url.Values, responseData interface{}) (*http.Response, error) {
	req, err := c.newRequest("GET", path, qp, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req, responseData)
}

func (c *Client) post(path string, requestData, responseData interface{}) (*http.Response, error) {
	req, err := c.newRequest("POST", path, nil, requestData)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req, responseData)
}

func (c *Client) put(path string, requestData, responseData interface{}) (*http.Response, error) {
	req, err := c.newRequest("PUT", path, nil, requestData)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req, responseData)
}

func (c *Client) delete(path string, qp url.Values, responseData interface{}) (*http.Response, error) {
	req, err := c.newRequest("DELETE", path, qp, nil)
	if err != nil {
		return nil, err
	}
	return c.doRequest(req, responseData)
}

func (c *Client) newRequest(method, path string, vals url.Values, requestData interface{}) (*http.Request, error) {
	var buf io.ReadWriter
	if requestData != nil {
		buf = new(bytes.Buffer)
		err := json.NewEncoder(buf).Encode(&requestPayload{Data: requestData})
		if err != nil {
			return nil, err
		}
	}

	u := c.config.URL + path
	if vals != nil {
		u = u + "?" + vals.Encode()
	}

	req, err := http.NewRequest(method, u, buf)
	if err != nil {
		return nil, err
	}
	req.SetBasicAuth(c.config.Username, c.config.Password)
	if requestData != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "application/json")

	return req, nil
}

func (c *Client) doRequest(req *http.Request, responseData interface{}) (*http.Response, error) {
	res, err := c.http.Do(req)
	if err != nil {
		return res, err
	}
	defer res.Body.Close()

	var p responsePayload
	if err = json.NewDecoder(res.Body).Decode(&p); err != nil {
		return res, err
	}
	if len(p.Errors) > 0 {
		return res, p.Errors
	}
	if responseData != nil {
		return res, json.Unmarshal(p.Data, responseData)
	}
	return res, nil
}