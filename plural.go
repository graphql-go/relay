package gqlrelay

import (
	"github.com/chris-ramon/graphql-go/types"
)

type ResolveSingleInputFn func(input interface{}) interface{}
type PluralIdentifyingRootFieldConfig struct {
	ArgName            string                  `json:"argName"`
	InputType          types.GraphQLInputType  `json:"inputType"`
	OutputType         types.GraphQLOutputType `json:"outputType"`
	ResolveSingleInput ResolveSingleInputFn    `json:"resolveSingleInput"`
	Description        string                  `json:"description"`
}

func PluralIdentifyingRootField(config PluralIdentifyingRootFieldConfig) *types.GraphQLFieldConfig {
	inputArgs := types.GraphQLFieldConfigArgumentMap{}
	if config.ArgName != "" {
		inputArgs[config.ArgName] = &types.GraphQLArgumentConfig{
			Type: types.NewGraphQLNonNull(types.NewGraphQLList(types.NewGraphQLNonNull(config.InputType))),
		}
	}

	return &types.GraphQLFieldConfig{
		Description: config.Description,
		Type:        types.NewGraphQLList(config.OutputType),
		Args:        inputArgs,
		Resolve: func(p types.GQLFRParams) interface{} {
			inputs, ok := p.Args[config.ArgName]
			if !ok {
				return nil
			}

			if config.ResolveSingleInput == nil {
				return nil
			}
			switch inputs := inputs.(type) {
			case []interface{}:
				res := []interface{}{}
				for _, input := range inputs {
					r := config.ResolveSingleInput(input)
					res = append(res, r)
				}
				return res
			}
			return nil
		},
	}
}
