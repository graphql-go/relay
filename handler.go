package graphql_relay

import (
	"encoding/json"
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/gorilla/schema"
	"github.com/unrolled/render"
	"io/ioutil"
	"net/http"
)

const (
	ContentTypeJSON           = "application/json"
	ContentTypeGraphQL        = "application/graphql"
	ContentTypeFormUrlEncoded = "application/x-www-form-urlencoded"
)

var decoder = schema.NewDecoder()

type Handler struct {
	Schema *types.GraphQLSchema
	render *render.Render
}
type requestOptions struct {
	Query         string `json:"query" url:"query" schema:"query"`
	Variables     string `json:"variables" url:"variables" schema:"variables"`
	OperationName string `json:"operationName" url:"operationName" schema:"operationName"`
}

func getRequestOptions(r *http.Request) *requestOptions {

	query := r.URL.Query().Get("query")
	if query != "" {
		return &requestOptions{
			Query:         query,
			Variables:     r.URL.Query().Get("variables"),
			OperationName: r.URL.Query().Get("operationName"),
		}
	}
	if r.Method != "POST" {
		return &requestOptions{}
	}
	if r.Body == nil {
		return &requestOptions{}
	}

	switch r.Header.Get("Content-Type") {
	case ContentTypeGraphQL:
		body, err := ioutil.ReadAll(r.Body)
		r.Body.Close()
		if err != nil {
			return &requestOptions{}
		}
		return &requestOptions{
			Query: string(body),
		}
	case ContentTypeFormUrlEncoded:
		var opts requestOptions
		err := r.ParseForm()
		if err != nil {
			return &requestOptions{}
		}
		err = decoder.Decode(&opts, r.PostForm)
		if err != nil {
			return &requestOptions{}
		}
		return &opts
	case ContentTypeJSON:
		fallthrough
	default:
		jsonDecoder := json.NewDecoder(r.Body)
		var opts requestOptions
		err := jsonDecoder.Decode(&opts)
		if err != nil {
			return &requestOptions{}
		}
		return &opts
	}
}

// ServeHTTP provides an entry point into executing graphQL queries
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	// get query
	opts := getRequestOptions(r)

	// execute graphql query
	resultChannel := make(chan *types.GraphQLResult)
	params := gql.GraphqlParams{
		Schema:        *h.Schema,
		RequestString: opts.Query,
	}
	go gql.Graphql(params, resultChannel)
	result := <-resultChannel

	// render result
	h.render.JSON(w, http.StatusOK, result)
}

type HandlerConfig struct {
	Schema *types.GraphQLSchema
	Pretty bool
}

func NewHandlerConfig() *HandlerConfig {
	return &HandlerConfig{
		Schema: nil,
		Pretty: true,
	}
}

func NewHandler(p *HandlerConfig) *Handler {
	if p == nil {
		p = NewHandlerConfig()
	}
	if p.Schema == nil {
		panic("undefined graphQL schema")
	}
	r := render.New(render.Options{
		IndentJSON: p.Pretty,
	})
	return &Handler{
		Schema: p.Schema,
		render: r,
	}
}
