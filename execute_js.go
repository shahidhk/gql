// +build js,wasm

package gql

import (
	"bytes"
	"encoding/json"
	"net/http"
	"syscall/js"
)

// ExecuteWithHeaders makes a request with extra headers
func (c *Client) ExecuteWithHeaders(r Request, headers map[string]string, callBack js.Value) error {
	for k, v := range headers {
		c.SetHeader(k, v)
	}
	c.Execute(r, callBack)
	for k := range headers {
		delete(c.Headers, k)
	}
	return nil
}

// Execute executes the Request r using the Client c and returns an error
// Response data and errors can be unmarshalled to the passed interfaces
func (c *Client) Execute(r Request, callBack js.Value) {
	payload, err := json.Marshal(r)
	if err != nil {
		callBack.Invoke(js.Null(), err.Error())
		return
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		callBack.Invoke(js.Null(), err.Error())
		return
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	go func() {
		res, err := c.client.Do(req)
		if err != nil {
			callBack.Invoke(js.Null(), err.Error())
			return
		}
		defer res.Body.Close()
		var response Response
		err = json.NewDecoder(res.Body).Decode(&response)
		if err != nil {
			callBack.Invoke(js.Null(), err.Error())
			return
		}

		if response.Data != nil {
			var data map[string]interface{}
			err = json.Unmarshal(*response.Data, &data)
			if err != nil {
				callBack.Invoke(js.Null(), err.Error())
				return
			}
			callBack.Invoke(data, js.Null())
		}

		if response.Errors != nil {
			var errors Errors
			err = json.Unmarshal(*response.Errors, &errors)
			if err != nil {
				callBack.Invoke(js.Null(), err.Error())
				return
			}
			callBack.Invoke(js.Null(), errors.Error())
		}
	}()
}
