package relay

import (
	"github.com/graphql-go/graphql"
)

type MutationFn func(inputMap map[string]interface{}, info graphql.ResolveInfo) map[string]interface{}

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
	Name                string                            `json:"name"`
	InputFields         graphql.InputObjectConfigFieldMap `json:"inputFields"`
	OutputFields        graphql.FieldConfigMap            `json:"outputFields"`
	MutateAndGetPayload MutationFn                        `json:"mutateAndGetPayload"`
}

/*
Returns a GraphQLFieldConfig for the mutation described by the
provided MutationConfig.
*/

func MutationWithClientMutationID(config MutationConfig) *graphql.FieldConfig {

	augmentedInputFields := config.InputFields
	if augmentedInputFields == nil {
		augmentedInputFields = graphql.InputObjectConfigFieldMap{}
	}
	augmentedInputFields["clientMutationId"] = &graphql.InputObjectFieldConfig{
		Type: graphql.NewNonNull(graphql.String),
	}
	augmentedOutputFields := config.OutputFields
	if augmentedOutputFields == nil {
		augmentedOutputFields = graphql.FieldConfigMap{}
	}
	augmentedOutputFields["clientMutationId"] = &graphql.FieldConfig{
		Type: graphql.NewNonNull(graphql.String),
	}

	inputType := graphql.NewInputObject(graphql.InputObjectConfig{
		Name:   config.Name + "Input",
		Fields: augmentedInputFields,
	})
	outputType := graphql.NewObject(graphql.ObjectConfig{
		Name:   config.Name + "Payload",
		Fields: augmentedOutputFields,
	})
	return &graphql.FieldConfig{
		Type: outputType,
		Args: graphql.FieldConfigArgument{
			"input": &graphql.ArgumentConfig{
				Type: graphql.NewNonNull(inputType),
			},
		},
		Resolve: func(p graphql.GQLFRParams) interface{} {
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
			if clientMutationID, ok := input["clientMutationId"]; ok {
				payload["clientMutationId"] = clientMutationID
			}
			return payload
		},
	}
}
