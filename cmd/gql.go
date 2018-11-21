package main

import (
	"fmt"
	"log"

	"github.com/shahidhk/gql"
)

type hasuraError struct {
	Path  string `json:"path"`
	Error string `json:"error"`
	Code  string `json:"code"`
}

var c = gql.NewClient("http://localhost:8080/v1alpha1/graphql", nil)
var introspectionResult gql.IntrospectionResponse
var e []hasuraError

func main() {

	err := c.Execute(gql.Request{
		Query: gql.IntrospectionQuery,
	}, &introspectionResult, &e)

	if err != nil {
		log.Fatal(err)
	}

	if len(e) > 0 {
		log.Fatalf("error: %v", e)
	}

	fmt.Printf("%v", introspectionResult)
}
