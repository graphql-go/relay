# graphql-relay-go

A Go/Golang library to help construct a [graphql-go](https://github.com/chris-ramon/graphql-go) server supporting react-relay.

See a live demo of here: http://bit.ly/try-graphql-go

Source code for demo can be found at https://github.com/sogko/golang-graphql-playground

### Tutorial
[Learn Golang + GraphQL + Relay Part 2: Your first Relay application]( https://wehavefaces.net/learn-golang-graphql-relay-2-a56cbcc3e341)

### Test
```bash
$ go get github.com/sogko/graphql-relay-go
$ go build && go test ./...
```

### TODO:
- [x] Starwars example
- [x] HTTP handler to easily create a Relay-compliant GraphQL server _(Moved to: [graphql-go-handler](https://github.com/sogko/graphql-go-handler))_
- [ ] In-code documentation (godocs)
- [ ] Usage guide / user documentation
- [x] Tutorial
- [ ] End-to-end example (graphql-relay-go + react-relay)
