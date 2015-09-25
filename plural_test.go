package graphql_relay_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-relay-go"
	"reflect"
	"testing"
)

var pluralTestUserType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "User",
	Fields: types.GraphQLFieldConfigMap{
		"username": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
		"url": &types.GraphQLFieldConfig{
			Type: types.GraphQLString,
		},
	},
})

var pluralTestQueryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"usernames": graphql_relay.PluralIdentifyingRootField(graphql_relay.PluralIdentifyingRootFieldConfig{
			ArgName:     "usernames",
			Description: "Map from a username to the user",
			InputType:   types.GraphQLString,
			OutputType:  pluralTestUserType,
			ResolveSingleInput: func(username interface{}) interface{} {
				return map[string]interface{}{
					"username": fmt.Sprintf("%v", username),
					"url":      fmt.Sprintf("www.facebook.com/%v", username),
				}
			},
		}),
	},
})

var pluralTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query: pluralTestQueryType,
})

func TestPluralIdentifyingRootField_AllowsFetching(t *testing.T) {
	query := `{
      usernames(usernames:["dschafer", "leebyron", "schrockn"]) {
        username
        url
      }
    }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"usernames": []interface{}{
				map[string]interface{}{
					"username": "dschafer",
					"url":      "www.facebook.com/dschafer",
				},
				map[string]interface{}{
					"username": "leebyron",
					"url":      "www.facebook.com/leebyron",
				},
				map[string]interface{}{
					"username": "schrockn",
					"url":      "www.facebook.com/schrockn",
				},
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        pluralTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestPluralIdentifyingRootField_CorrectlyIntrospects(t *testing.T) {
	query := `{
      __schema {
        queryType {
          fields {
            name
            args {
              name
              type {
                kind
                ofType {
                  kind
                  ofType {
                    kind
                    ofType {
                      name
                      kind
                    }
                  }
                }
              }
            }
            type {
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
      }
    }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"__schema": map[string]interface{}{
				"queryType": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"name": "usernames",
							"args": []interface{}{
								map[string]interface{}{
									"name": "usernames",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"kind": "LIST",
											"ofType": map[string]interface{}{
												"kind": "NON_NULL",
												"ofType": map[string]interface{}{
													"name": "String",
													"kind": "SCALAR",
												},
											},
										},
									},
								},
							},
							"type": map[string]interface{}{
								"kind": "LIST",
								"ofType": map[string]interface{}{
									"name": "User",
									"kind": "OBJECT",
								},
							},
						},
					},
				},
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        pluralTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestPluralIdentifyingRootField_Configuration_ResolveSingleInputIsNil(t *testing.T) {

	var pluralTestQueryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"usernames": graphql_relay.PluralIdentifyingRootField(graphql_relay.PluralIdentifyingRootFieldConfig{
				ArgName:     "usernames",
				Description: "Map from a username to the user",
				InputType:   types.GraphQLString,
				OutputType:  pluralTestUserType,
			}),
		},
	})

	var pluralTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: pluralTestQueryType,
	})

	query := `{
      usernames(usernames:["dschafer", "leebyron", "schrockn"]) {
        username
        url
      }
    }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"usernames": nil,
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        pluralTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestPluralIdentifyingRootField_Configuration_ArgNames_WrongArgNameSpecified(t *testing.T) {

	query := `{
      usernames(usernamesMisspelled:["dschafer", "leebyron", "schrockn"]) {
        username
        url
      }
    }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"usernames": nil,
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        pluralTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
