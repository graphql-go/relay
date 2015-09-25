package graphql_relay_test

import (
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/errors"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-relay-go"
	"reflect"
	"testing"
	"time"
)

func testAsyncDataMutation(resultChan *chan int) {
	// simulate async data mutation
	time.Sleep(time.Second * 1)
	*resultChan <- int(1)
}

var simpleMutationTest = graphql_relay.MutationWithClientMutationId(graphql_relay.MutationConfig{
	Name:        "SimpleMutation",
	InputFields: types.InputObjectConfigFieldMap{},
	OutputFields: types.GraphQLFieldConfigMap{
		"result": &types.GraphQLFieldConfig{
			Type: types.GraphQLInt,
		},
	},
	MutateAndGetPayload: func(object map[string]interface{}, info types.GraphQLResolveInfo) map[string]interface{} {
		return map[string]interface{}{
			"result": 1,
		}
	},
})

// async mutation
var simplePromiseMutationTest = graphql_relay.MutationWithClientMutationId(graphql_relay.MutationConfig{
	Name:        "SimplePromiseMutation",
	InputFields: types.InputObjectConfigFieldMap{},
	OutputFields: types.GraphQLFieldConfigMap{
		"result": &types.GraphQLFieldConfig{
			Type: types.GraphQLInt,
		},
	},
	MutateAndGetPayload: func(object map[string]interface{}, info types.GraphQLResolveInfo) map[string]interface{} {
		c := make(chan int)
		go testAsyncDataMutation(&c)
		result := <-c
		return map[string]interface{}{
			"result": result,
		}
	},
})

var mutationTestType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Mutation",
	Fields: types.GraphQLFieldConfigMap{
		"simpleMutation":        simpleMutationTest,
		"simplePromiseMutation": simplePromiseMutationTest,
	},
})

var mutationTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
	Query:    mutationTestType,
	Mutation: mutationTestType,
})

func TestMutation_WithClientMutationId_BehavesCorrectly_RequiresAnArgument(t *testing.T) {
	t.Skipf("Pending `validator` implementation")
	query := `
        mutation M {
          simpleMutation {
            result
          }
        }
      `
	expected := &types.GraphQLResult{
		Errors: []graphqlerrors.GraphQLFormattedError{
			graphqlerrors.GraphQLFormattedError{
				Message: `Field "simpleMutation" argument "input" of type "SimpleMutationInput!" is required but not provided.`,
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        mutationTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestMutation_WithClientMutationId_BehavesCorrectly_ReturnsTheSameClientMutationId(t *testing.T) {
	query := `
        mutation M {
          simpleMutation(input: {clientMutationId: "abc"}) {
            result
            clientMutationId
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"simpleMutation": map[string]interface{}{
				"result":           1,
				"clientMutationId": "abc",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        mutationTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

// Async mutation using channels
func TestMutation_WithClientMutationId_BehavesCorrectly_SupportsPromiseMutations(t *testing.T) {
	query := `
        mutation M {
          simplePromiseMutation(input: {clientMutationId: "abc"}) {
            result
            clientMutationId
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"simplePromiseMutation": map[string]interface{}{
				"result":           1,
				"clientMutationId": "abc",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        mutationTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestMutation_IntrospectsCorrectly_ContainsCorrectInput(t *testing.T) {
	query := `{
        __type(name: "SimpleMutationInput") {
          name
          kind
          inputFields {
            name
            type {
              name
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
				"name": "SimpleMutationInput",
				"kind": "INPUT_OBJECT",
				"inputFields": []interface{}{
					map[string]interface{}{
						"name": "clientMutationId",
						"type": map[string]interface{}{
							"name": nil,
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
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        mutationTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestMutation_IntrospectsCorrectly_ContainsCorrectPayload(t *testing.T) {
	t.Skipf("Need a util to test for slice equality as an unordered set")
	query := `{
        __type(name: "SimpleMutationPayload") {
          name
          kind
          fields {
            name
            type {
              name
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
				"name": "SimpleMutationPayload",
				"kind": "OBJECT",
				"fields": []interface{}{
					map[string]interface{}{
						"name": "result",
						"type": map[string]interface{}{
							"name":   "Int",
							"kind":   "SCALAR",
							"ofType": nil,
						},
					},
					map[string]interface{}{
						"name": "clientMutationId",
						"type": map[string]interface{}{
							"name": nil,
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
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        mutationTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestMutation_IntrospectsCorrectly_ContainsCorrectField(t *testing.T) {
	t.Skipf("Need a util to test for slice equality as an unordered set")

	query := `{
        __schema {
          mutationType {
            fields {
              name
              args {
                name
                type {
                  name
                  kind
                  ofType {
                    name
                    kind
                  }
                }
              }
              type {
                name
                kind
              }
            }
          }
        }
      }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"__schema": map[string]interface{}{
				"mutationType": map[string]interface{}{
					"fields": []interface{}{
						map[string]interface{}{
							"name": "simpleMutation",
							"args": []interface{}{
								map[string]interface{}{
									"name": "input",
									"type": map[string]interface{}{
										"name": nil,
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "SimpleMutationInput",
											"kind": "INPUT_OBJECT",
										},
									},
								},
							},
							"type": map[string]interface{}{
								"name": "SimpleMutationPayload",
								"kind": "OBJECT",
							},
						},
						map[string]interface{}{
							"name": "simplePromiseMutation",
							"args": []interface{}{
								map[string]interface{}{
									"name": "input",
									"type": map[string]interface{}{
										"name": nil,
										"kind": "NON_NULL",
										"ofType": map[string]interface{}{
											"name": "SimplePromiseMutationInput",
											"kind": "INPUT_OBJECT",
										},
									},
								},
							},
							"type": map[string]interface{}{
								"name": "SimpleMutationPayload",
								"kind": "SimplePromiseMutationPayload",
							},
						},
					},
				},
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        mutationTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
