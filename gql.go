package gql

import (
	"bytes"
	"encoding/json"
	"net/http"
)

// GQLRequest is the GraphQL request containing Query and Variables
type GQLRequest struct {
	Query     string                 `json:"query"`
	Variables map[string]interface{} `json:"variables"`
}

// GQLResponse is the response from GraphQL server
type GQLResponse struct {
	Data   *json.RawMessage `json:"data"`
	Errors *json.RawMessage `json:"errors"`
}

// GQLError is a the GraphQL error from server
type GQLError struct {
	Message   string             `json:"message"`
	Locations []GQLErrorLocation `json:"locations"`
	Type      string             `json:"type"`
	Path      []interface{}      `json:"path"`
}

// Error returns the error message
func (e GQLError) Error() string {
	return e.Message
}

// GQLErrorLocation is the location of error in the query string
type GQLErrorLocation struct {
	Line   int `json:"line"`
	Column int `json:"column"`
}

// GQLClient can execute GraphQL queries against an endpoint
type GQLClient struct {
	Endpoint string
	Headers  map[string]string
	client   *http.Client
}

// NewGQLClient returns a GQLClient for given endpoint and headers
func NewGQLClient(endpoint string, headers map[string]string) *GQLClient {
	return &GQLClient{
		Endpoint: endpoint,
		Headers:  headers,
		client:   &http.Client{},
	}
}

// Execute executes the GQLRequest r using the GQLClient c and returns an error
// Response data and errors can be unmarshalled to the passed interfaces
func (c *GQLClient) Execute(r GQLRequest, data interface{}, errors interface{}) error {
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

	var response GQLResponse
	err = json.NewDecoder(res.Body).Decode(&response)
	if err != nil {
		return err
	}

	err = json.Unmarshal(*response.Data, data)
	if err != nil {
		return err
	}
	if response.Errors != nil {
		err = json.Unmarshal(*response.Errors, errors)
		if err != nil {
			return err
		}
	}

	return nil
}
