package biblio

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/ugent-library/biblio-backend/internal/models"
)

type requestPayload struct {
	Data interface{} `json:"data"`
}

type responsePayload struct {
	Data json.RawMessage `json:"data"`
}

func (c *Client) get(path string, qp url.Values, responseData interface{}) (*http.Response, error) {
	req, err := c.newRequest("GET", path, qp, nil)
	if err != nil {
		return nil, err
	}
	if strings.Contains(path, "/restricted/") {
		req.SetBasicAuth(c.config.Username, c.config.Password)
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
	defer func() {
		io.Copy(ioutil.Discard, res.Body)
		res.Body.Close()
	}()

	if res.StatusCode == http.StatusNotFound {
		return res, models.HttpNotFound{
			Message: fmt.Sprintf("url %s not found", req.URL.String()),
		}
	}

	var p responsePayload
	if err = json.NewDecoder(res.Body).Decode(&p); err != nil {
		return res, err
	}

	if responseData != nil {
		return res, json.Unmarshal(p.Data, responseData)
	}
	return res, nil
}
