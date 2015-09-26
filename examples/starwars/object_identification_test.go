package starwars_test

import (
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-relay-go/examples/starwars"
	"reflect"
	"testing"
)

func TestObjectIdentification_TestFetching_CorrectlyFetchesTheIDAndTheNameOfTheRebels(t *testing.T) {
	query := `
        query RebelsQuery {
          rebels {
            id
            name
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"id":   "RmFjdGlvbjox",
				"name": "Alliance to Restore the Republic",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestObjectIdentification_TestFetching_CorrectlyRefetchesTheRebels(t *testing.T) {
	query := `
        query RebelsRefetchQuery {
          node(id: "RmFjdGlvbjox") {
            id
            ... on Faction {
              name
            }
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":   "RmFjdGlvbjox",
				"name": "Alliance to Restore the Republic",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestObjectIdentification_TestFetching_CorrectlyFetchesTheIDAndTheNameOfTheEmpire(t *testing.T) {
	query := `
        query EmpireQuery {
          empire {
            id
            name
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"empire": map[string]interface{}{
				"id":   "RmFjdGlvbjoy",
				"name": "Galactic Empire",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestObjectIdentification_TestFetching_CorrectlyRefetchesTheEmpire(t *testing.T) {
	query := `
        query EmpireRefetchQuery {
          node(id: "RmFjdGlvbjoy") {
            id
            ... on Faction {
              name
            }
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":   "RmFjdGlvbjoy",
				"name": "Galactic Empire",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestObjectIdentification_TestFetching_CorrectlyRefetchesTheXWing(t *testing.T) {
	query := `
        query XWingRefetchQuery {
          node(id: "U2hpcDox") {
            id
            ... on Ship {
              name
            }
          }
        }
      `
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"node": map[string]interface{}{
				"id":   "U2hpcDox",
				"name": "X-Wing",
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
