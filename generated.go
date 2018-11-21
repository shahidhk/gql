package gql

// Author is an author type
type Author struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

// Author returns an array of authors
func (c *Client) Author() ([]Author, error) {

	var authors []Author
	var e interface{}

	err := c.Execute(Request{
		Query: `query {
			author {
				id
				name
			}
		}`,
	}, &authors, &e)

	if err != nil {
		return authors, err
	}

	return authors, nil
}
