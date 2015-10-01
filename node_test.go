package gqlrelay_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-relay-go"
	"reflect"
	"testing"
)

type user struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}
type photo struct {
	Id    int `json:"id"`
	Width int `json:"width"`
}

var nodeTestUserData = map[string]*user{
	"1": &user{1, "John Doe"},
	"2": &user{2, "Jane Smith"},
}
var nodeTestPhotoData = map[string]*photo{
	"3": &photo{3, 300},
	"4": &photo{4, 400},
}

// declare types first, define later in init()
// because they all depend on nodeTestDef
var nodeTestUserType *types.GraphQLObjectType
var nodeTestPhotoType *types.GraphQLObjectType

var nodeTestDef = gqlrelay.NewNodeDefinitions(gqlrelay.NodeDefinitionsConfig{
	IdFetcher: func(id string, info types.GraphQLResolveInfo) interface{} {
		if user, ok := nodeTestUserData[id]; ok {
			return user
		}
		if photo, ok := nodeTestPhotoData[id]; ok {
			return photo
		}
		return nil
	},
	TypeResolve: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
		switch value.(type) {
		case *user:
			return nodeTestUserType
		case *photo:
			return nodeTestPhotoType
		default:
			panic(fmt.Sprintf("Unknown object type `%v`", value))
		}
	},
})
var nodeTestQueryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"node": nodeTestDef.NodeField,
	},
})

// becareful not to define schema here, since nodeTestUserType and nodeTestPhotoType wouldn't be defined till init()
var nodeTestSchema types.GraphQLSchema

func init() {
	nodeTestUserType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "User",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLNonNull(types.GraphQLID),
			},
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{nodeTestDef.NodeInterface},
	})
	nodeTestPhotoType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Photo",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type: types.NewGraphQLNonNull(types.GraphQLID),
			},
			"width": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{nodeTestDef.NodeInterface},
	})

	nodeTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: nodeTestQueryType,
	})
}

func graphql(t *testing.T, p gql.GraphqlParams) *types.GraphQLResult {
	resultChannel := make(chan *types.GraphQLResult)
	go gql.Graphql(p, resultChannel)
	result := <-resultChannel
	return result
}
func TestNodeInterfaceAndFields_AllowsRefetching_GetsTheCorrectIDForUsers(t *testing.T) {
	query := `{
        node(id: "1") {
          id
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id": "1",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_GetsTheCorrectIDForPhotos(t *testing.T) {
	query := `{
        node(id: "4") {
          id
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id": "4",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_GetsTheCorrectNameForUsers(t *testing.T) {
	query := `{
        node(id: "1") {
          id
          ... on User {
            name
          }
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":   "1",
				"name": "John Doe",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_GetsTheCorrectWidthForPhotos(t *testing.T) {
	query := `{
        node(id: "4") {
          id
          ... on Photo {
            width
          }
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":    "4",
				"width": 400,
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_GetsTheCorrectTypeNameForUsers(t *testing.T) {
	query := `{
        node(id: "1") {
          id
          __typename
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":         "1",
				"__typename": "User",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_GetsTheCorrectTypeNameForPhotos(t *testing.T) {
	query := `{
        node(id: "4") {
          id
          __typename
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":         "4",
				"__typename": "Photo",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_IgnoresPhotoFragmentsOnUser(t *testing.T) {
	query := `{
        node(id: "1") {
          id
          ... on Photo {
            width
          }
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id": "1",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_AllowsRefetching_ReturnsNullForBadIDs(t *testing.T) {
	query := `{
        node(id: "5") {
          id
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": nil,
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_CorrectlyIntrospects_HasCorrectNodeInterface(t *testing.T) {
	query := `{
        __type(name: "Node") {
          name
          kind
          fields {
            name
            type {
              kind
              ofType {
                name
                kind
              }
            }
          }
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"__type": map[string]interface{}{
				"name": "Node",
				"kind": "INTERFACE",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "id",
						"type": map[string]interface{}{
							"kind": "NON_NULL",
							"ofType": map[string]interface{}{
								"name": "ID",
								"kind": "SCALAR",
							},
						},
					},
				},
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestNodeInterfaceAndFields_CorrectlyIntrospects_HasCorrectNodeRootField(t *testing.T) {
	query := `{
        __schema {
          queryType {
            fields {
              name
              type {
                name
                kind
              }
              args {
                name
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
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"__schema": map[string]interface{}{
				"queryType": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"name": "node",
							"type": map[string]interface{}{
								"name": "Node",
								"kind": "INTERFACE",
							},
							"args": []interface{}{
								map[string]interface{}{
									"name": "id",
									"type": map[string]interface{}{
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "ID",
											"kind": "SCALAR",
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        nodeTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
