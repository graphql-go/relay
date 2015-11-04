package gqlrelay_test

import (
	"github.com/chris-ramon/graphql"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/sogko/graphql-relay-go"
	"reflect"
	"testing"
)

var connectionTestAllUsers = []interface{}{
	&user{Name: "Dan"},
	&user{Name: "Nick"},
	&user{Name: "Lee"},
	&user{Name: "Joe"},
	&user{Name: "Tim"},
}
var connectionTestUserType *graphql.Object
var connectionTestQueryType *graphql.Object
var connectionTestSchema graphql.Schema
var connectionTestConnectionDef *gqlrelay.GraphQLConnectionDefinitions

func init() {
	connectionTestUserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.FieldConfigMap{
			"name": &graphql.FieldConfig{
				Type: graphql.String,
			},
			// re-define `friends` field later because `connectionTestUserType` has `connectionTestConnectionDef` has `connectionTestUserType` (cyclic-reference)
			"friends": &graphql.FieldConfig{},
		},
	})

	connectionTestConnectionDef = gqlrelay.ConnectionDefinitions(gqlrelay.ConnectionConfig{
		Name:     "Friend",
		NodeType: connectionTestUserType,
		EdgeFields: graphql.FieldConfigMap{
			"friendshipTime": &graphql.FieldConfig{
				Type: graphql.String,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					return "Yesterday"
				},
			},
		},
		ConnectionFields: graphql.FieldConfigMap{
			"totalCount": &graphql.FieldConfig{
				Type: graphql.Int,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					return len(connectionTestAllUsers)
				},
			},
		},
	})

	// define `friends` field here after getting connection definition
	connectionTestUserType.AddFieldConfig("friends", &graphql.FieldConfig{
		Type: connectionTestConnectionDef.ConnectionType,
		Args: gqlrelay.ConnectionArgs,
		Resolve: func(p graphql.GQLFRParams) interface{} {
			arg := gqlrelay.NewConnectionArguments(p.Args)
			res := gqlrelay.ConnectionFromArray(connectionTestAllUsers, arg)
			return res
		},
	})

	connectionTestQueryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.FieldConfigMap{
			"user": &graphql.FieldConfig{
				Type: connectionTestUserType,
				Resolve: func(p graphql.GQLFRParams) interface{} {
					return connectionTestAllUsers[0]
				},
			},
		},
	})
	var err error
	connectionTestSchema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: connectionTestQueryType,
	})
	if err != nil {
		panic(err)
	}

}

func TestConnectionDefinition_IncludesConnectionAndEdgeFields(t *testing.T) {
	query := `
      query FriendsQuery {
        user {
          friends(first: 2) {
            totalCount
            edges {
              friendshipTime
              node {
                name
              }
            }
          }
        }
      }
    `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"friends": map[string]interface{}{
					"totalCount": 5,
					"edges": []interface{}{
						map[string]interface{}{
							"friendshipTime": "Yesterday",
							"node": map[string]interface{}{
								"name": "Dan",
							},
						},
						map[string]interface{}{
							"friendshipTime": "Yesterday",
							"node": map[string]interface{}{
								"name": "Nick",
							},
						},
					},
				},
			},
		},
	}
	result := testGraphql(t, graphql.Params{
		Schema:        connectionTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
