package starwars_test

import (
	"github.com/chris-ramon/graphql"
	"github.com/chris-ramon/graphql/testutil"
	"github.com/sogko/graphql-relay-go/examples/starwars"
	"reflect"
	"testing"
)

func TestMutation_CorrectlyMutatesTheDataSet(t *testing.T) {
	query := `
      mutation AddBWingQuery($input: IntroduceShipInput!) {
        introduceShip(input: $input) {
          ship {
            id
            name
          }
          faction {
            name
          }
          clientMutationID
        }
      }
    `
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"shipName":         "B-Wing",
			"factionId":        "1",
			"clientMutationID": "abcde",
		},
	}
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"introduceShip": map[string]interface{}{
				"ship": map[string]interface{}{
					"id":   "U2hpcDoxMA==",
					"name": "B-Wing",
				},
				"faction": map[string]interface{}{
					"name": "Alliance to Restore the Republic",
				},
				"clientMutationID": "abcde",
			},
		},
	}
	result := testGraphql(t, graphql.Params{
		Schema:         starwars.Schema,
		RequestString:  query,
		VariableValues: params,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
