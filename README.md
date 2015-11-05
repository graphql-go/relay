# graphql-relay-go

A Go/Golang library to help construct a [graphql-go](https://github.com/graphql-go/graphql-go) server supporting react-relay.

See a live demo of here: http://bit.ly/try-graphql-go

Source code for demo can be found at https://github.com/graphql-go/golang-graphql-playground

### Notes:
This is based on alpha version of `graphql-go` and `graphql-relay-go`. 
Be sure to watch both repositories for latest changes.

### Tutorial
[Learn Golang + GraphQL + Relay Part 2: Your first Relay application]( https://wehavefaces.net/learn-golang-graphql-relay-2-a56cbcc3e341)

### Test
```bash
$ go get github.com/graphql-go/graphql-relay-go
$ go build && go test ./...
```

### TODO:
- [x] Starwars example
- [x] HTTP handler to easily create a Relay-compliant GraphQL server _(Moved to: [graphql-go-handler](https://github.com/graphql-go/graphql-go-handler))_
- [ ] In-code documentation (godocs)
- [ ] Usage guide / user documentation
- [x] Tutorial
- [ ] End-to-end example (graphql-relay-go + react-relay)
