package graphql_relay

import (
	"github.com/chris-ramon/graphql-go/types"
)

type MutationFn func(object map[string]interface{}, info types.GraphQLResolveInfo) map[string]interface{}

/*
A description of a mutation consumable by mutationWithClientMutationId
to create a GraphQLFieldConfig for that mutation.

The inputFields and outputFields should not include `clientMutationId`,
as this will be provided automatically.

An input object will be created containing the input fields, and an
object will be created containing the output fields.

mutateAndGetPayload will receive an Object with a key for each
input field, and it should return an Object with a key for each
output field. It may return synchronously, or return a Promise.
*/
type MutationConfig struct {
	Name                string                          `json:"name"`
	InputFields         types.InputObjectConfigFieldMap `json:"inputFields"`
	OutputFields        types.GraphQLFieldConfigMap     `json:"outputFields"`
	MutateAndGetPayload MutationFn                      `json:"mutateAndGetPayload"`
}

/*
Returns a GraphQLFieldConfig for the mutation described by the
provided MutationConfig.
*/

func MutationWithClientMutationId(config MutationConfig) *types.GraphQLFieldConfig {

	augmentedInputFields := config.InputFields
	augmentedInputFields["clientMutationId"] = &types.InputObjectFieldConfig{
		Type: types.NewGraphQLNonNull(types.GraphQLString),
	}
	augmentedOutputFields := config.OutputFields
	augmentedOutputFields["clientMutationId"] = &types.GraphQLFieldConfig{
		Type: types.NewGraphQLNonNull(types.GraphQLString),
	}

	inputType := types.NewGraphQLInputObjectType(types.InputObjectConfig{
		Name:   config.Name + "Input",
		Fields: augmentedInputFields,
	})
	outputType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:   config.Name + "Payload",
		Fields: augmentedOutputFields,
	})
	return &types.GraphQLFieldConfig{
		Type: outputType,
		Args: types.GraphQLFieldConfigArgumentMap{
			"input": &types.GraphQLArgumentConfig{
				Type: types.NewGraphQLNonNull(inputType),
			},
		},
		Resolve: func(p types.GQLFRParams) interface{} {
			if config.MutateAndGetPayload == nil {
				return nil
			}
			input := map[string]interface{}{}
			if inputVal, ok := p.Args["input"]; ok {
				if inputVal, ok := inputVal.(map[string]interface{}); ok {
					input = inputVal
				}
			}
			payload := config.MutateAndGetPayload(input, p.Info)
			if clientMutationId, ok := input["clientMutationId"]; ok {
				payload["clientMutationId"] = clientMutationId
			}
			return payload
		},
	}
}
