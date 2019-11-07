// +build !js !wasm

package gql

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// ExecuteWithHeaders makes a request with extra headers
func (c *Client) ExecuteWithHeaders(r Request, headers map[string]string, data interface{}) error {
	for k, v := range headers {
		c.SetHeader(k, v)
	}
	err := c.Execute(r, data)
	for k := range headers {
		delete(c.Headers, k)
	}
	return err
}

// Execute executes the Request r using the Client c and returns an error
// Response data and errors can be unmarshalled to the passed interfaces
func (c *Client) Execute(r Request, data interface{}) error {
	payload, err := json.Marshal(r)
	if err != nil {
		return err
	}
	req, err := http.NewRequest("POST", c.Endpoint, bytes.NewBuffer(payload))
	if err != nil {
		return err
	}
	for k, v := range c.Headers {
		req.Header.Set(k, v)
	}
	res, err := c.client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	var response Response
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	if response.Data != nil {
		err = json.Unmarshal(*response.Data, data)
		if err != nil {
			return err
		}
	}

	if response.Errors != nil {
		var errors Errors
		err = json.Unmarshal(*response.Errors, &errors)
		if err != nil {
			return err
		}
		return errors
	}
	return nil
}
