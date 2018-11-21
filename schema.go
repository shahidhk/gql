package gql

// IntrospectionResponse is the response to an introspection query
type IntrospectionResponse struct {
	Schema `json:"__schema"`
}

// Schema is the GraphQL schema
type Schema struct {
	QueryType        RootType    `json:"queryType"`
	MutationType     RootType    `json:"mutationType"`
	SubscriptionType RootType    `json:"subscriptionType"`
	Types            []FullType  `json:"types"`
	Directives       []Directive `json:"directives"`
}

// RootType is the query/mutation root type
type RootType struct {
	Name string `json:"name"`
}

// FullType is a GraphQL type
type FullType struct {
	Kind          string       `json:"kind"`
	Name          string       `json:"name"`
	Description   string       `json:"description"`
	Fields        []Field      `json:"fields"`
	InputFields   []InputValue `json:"inputFields"`
	Interfaces    []TypeRef    `json:"interfaces"`
	EnumValues    []EnumValue  `json:"enumValues"`
	PossibleTypes []TypeRef    `json:"possibleTypes"`
}

// Field is a GraphQL field
type Field struct {
	Name              string       `json:"name"`
	Description       string       `json:"description"`
	Args              []InputValue `json:"args"`
	Type              TypeRef      `json:"type"`
	IsDeprecated      bool         `json:"isDeprecated"`
	DeprecationReason string       `json:"deprecationReason"`
}

// InputValue is a GraphQL input value
type InputValue struct {
	Name         string      `json:"name"`
	Description  string      `json:"description"`
	Type         TypeRef     `json:"type"`
	DefaultValue interface{} `json:"defaultValue"`
}

// TypeRef is a GraphQL type ref
type TypeRef struct {
	Kind   string   `json:"kind"`
	Name   string   `json:"name"`
	OfType *TypeRef `json:"ofType"`
}

// EnumValue is a GraphQL enum value
type EnumValue struct {
	Name              string `json:"name"`
	Description       string `json:"description"`
	Isdeprecated      bool   `json:"isdeprecated"`
	DeprecationReason string `json:"deprecationReason"`
}

// Directive is a GraphQL directive
type Directive struct {
	Name        string       `json:"name"`
	Description string       `json:"description"`
	Locations   []string     `json:"locations"`
	Args        []InputValue `json:"args"`
}
