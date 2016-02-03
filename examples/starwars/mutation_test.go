package starwars_test

import (
	"reflect"
	"testing"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/graphql-go/relay/examples/starwars"
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
          clientMutationId
        }
      }
    `
	params := map[string]interface{}{
		"input": map[string]interface{}{
			"shipName":         "B-Wing",
			"factionId":        "1",
			"clientMutationId": "abcde",
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
				"clientMutationId": "abcde",
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:         starwars.Schema,
		RequestString:  query,
		VariableValues: params,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
