# gql
A (WIP) client library for GraphQL written in Go with codegen for GraphQL types

```go
package main

import (
    "fmt"
    "log"

	"github.com/shahidhk/gql"
)

type Author struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}
type Authors struct {
	Authors []Authors `json:"authors"`
}

func main() {
	client := gql.NewClient("http://localhost:8080/v1/graphql", nil)
	var authors Authors
	err := client.Execute(gql.Request{Query: `query { authors {id name}}`}, &authors)
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println(authors)
}
```
