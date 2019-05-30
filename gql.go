package gql

import (
	"bytes"
	"encoding/json"
	"net/http"
	"strings"
)

// Request is the GraphQL request containing Query and Variables
type Request struct {
	Query         string                 `json:"query"`
	Variables     map[string]interface{} `json:"variables"`
	OperationName string                 `json:"operationName,omitempty"`
}

// Response is the response from GraphQL server
type Response struct {
	Data   *json.RawMessage `json:"data"`
	Errors *json.RawMessage `json:"errors"`
}

// Error is a the GraphQL error from server
type Error struct {
	Message    string           `json:"message"`
	Locations  []ErrorLocation  `json:"locations"`
	Type       string           `json:"type"`
	Path       []interface{}    `json:"path"`
	Extensions HasuraExtensions `json:"extensions"`
}

// Error returns the error message
func (e Error) Error() string {
	return e.Message
}

// HasuraExtensions is the error extension by Hasura
type HasuraExtensions struct {
	Path string `json:"path"`
	Code string `json:"code"`
}

// Errors are an array of GraphQL errors
type Errors []Error

// Error returns all error messages from Errors object
func (e Errors) Error() string {
	errors := []string{}
	for _, err := range e {
		errors = append(errors, err.Message)
	}
	return strings.Join(errors, ", ")
}

// ErrorLocation is the location of error in the query string
type ErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// Client can execute GraphQL queries against an endpoint
type Client struct {
	Endpoint string
	Headers  map[string]string
	client   *http.Client
}

// NewClient returns a Client for given endpoint and headers
func NewClient(endpoint string, headers map[string]string) *Client {
	return &Client{
		Endpoint: endpoint,
		Headers:  headers,
		client:   &http.Client{},
	}
}

// SetHeader sets a new header on the client
func (c *Client) SetHeader(key, value string) {
	c.Headers[key] = value
}

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
