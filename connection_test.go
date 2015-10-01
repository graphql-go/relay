package gqlrelay_test

import (
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
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
var connectionTestUserType *types.GraphQLObjectType
var connectionTestQueryType *types.GraphQLObjectType
var connectionTestSchema types.GraphQLSchema
var connectionTestConnectionDef *gqlrelay.GraphQLConnectionDefinitions

func init() {
	connectionTestUserType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "User",
		Fields: types.GraphQLFieldConfigMap{
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
			// re-define `friends` field later because `connectionTestUserType` has `connectionTestConnectionDef` has `connectionTestUserType` (cyclic-reference)
			"friends": &types.GraphQLFieldConfig{},
		},
	})

	connectionTestConnectionDef = gqlrelay.ConnectionDefinitions(gqlrelay.ConnectionConfig{
		Name:     "Friend",
		NodeType: connectionTestUserType,
		EdgeFields: types.GraphQLFieldConfigMap{
			"friendshipTime": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
				Resolve: func(p types.GQLFRParams) interface{} {
					return "Yesterday"
				},
			},
		},
		ConnectionFields: types.GraphQLFieldConfigMap{
			"totalCount": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
				Resolve: func(p types.GQLFRParams) interface{} {
					return len(connectionTestAllUsers)
				},
			},
		},
	})

	// define `friends` field here after getting connection definition
	connectionTestUserType.AddFieldConfig("friends", &types.GraphQLFieldConfig{
		Type: connectionTestConnectionDef.ConnectionType,
		Args: gqlrelay.ConnectionArgs,
		Resolve: func(p types.GQLFRParams) interface{} {
			arg := gqlrelay.NewConnectionArguments(p.Args)
			res := gqlrelay.ConnectionFromArray(connectionTestAllUsers, arg)
			return res
		},
	})

	connectionTestQueryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"user": &types.GraphQLFieldConfig{
				Type: connectionTestUserType,
				Resolve: func(p types.GQLFRParams) interface{} {
					return connectionTestAllUsers[0]
				},
			},
		},
	})
	var err error
	connectionTestSchema, err = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
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
	expected := &types.GraphQLResult{
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
	result := graphql(t, gql.GraphqlParams{
		Schema:        connectionTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
