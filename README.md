# graphql-relay-go

A Go/Golang library to help construct a graphql-go server supporting react-relay.

Note: Currently based on an experimental branch of [`graphql-go`](https://github.com/chris-ramon/graphql-go); pending merge for the following PRs:
- https://github.com/chris-ramon/graphql-go/pull/12
- ~~https://github.com/chris-ramon/graphql-go/pull/10~~
- ~~https://github.com/chris-ramon/graphql-go/pull/8~~

### HTTP Handler Usage

```go
package main

import (
	"net/http"
	"github.com/sogko/graphql-relay-go"

)

func main() {

  // define GraphQL schema using relay library helpers
  schema := types.NewGraphQLSchema(...)
  
	// simplest relay-compliant schema server
	h := graphql_relay.NewHandler(&graphql_relay.HandlerConfig{
  		Schema: &starwars.Schema,
  		Pretty: true,
  })
	
	// serve HTTP
	http.Handle("/graphql", h)
	http.ListenAndServe(":8080", nil)
}
```

`handler` will accept requests with
the parameters:

  * **`query`**: A string GraphQL document to be executed.

  * **`variables`**: The runtime values to use for any GraphQL query variables
    as a JSON object.

  * **`operationName`**: If the provided `query` contains multiple named
    operations, this specifies which operation should be executed. If not
    provided, an 400 error will be returned if the `query` contains multiple
    named operations.

GraphQL will first look for each parameter in the URL's query-string:

```
/graphql?query=query+getUser($id:ID){user(id:$id){name}}&variables={"id":"4"}
```

If not found in the query-string, it will look in the POST request body.
The `handler` will interpret it
depending on the provided `Content-Type` header.

  * **`application/json`**: the POST body will be parsed as a JSON
    object of parameters.

  * **`application/x-www-form-urlencoded`**: this POST body will be
    parsed as a url-encoded string of key-value pairs.

  * **`application/graphql`**: The POST body will be parsed as GraphQL
    query string, which provides the `query` parameter.


### Test
```bash
$ go get github.com/sogko/graphql-relay-go
$ go build && go test ./...
```

### TODO:
- [x] Starwars example
- [x] HTTP handler to easily create a Relay-compliant GraphQL server
- [ ] In-code documentation (godocs)
- [ ] Usage guide / user documentation
- [ ] End-to-end example (graphql-relay-go + react-relay)
