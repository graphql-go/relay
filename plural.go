package gqlrelay

import (
	"github.com/chris-ramon/graphql"
)

type ResolveSingleInputFn func(input interface{}) interface{}
type PluralIdentifyingRootFieldConfig struct {
	ArgName            string               `json:"argName"`
	InputType          graphql.Input        `json:"inputType"`
	OutputType         graphql.Output       `json:"outputType"`
	ResolveSingleInput ResolveSingleInputFn `json:"resolveSingleInput"`
	Description        string               `json:"description"`
}

func PluralIdentifyingRootField(config PluralIdentifyingRootFieldConfig) *graphql.FieldConfig {
	inputArgs := graphql.FieldConfigArgument{}
	if config.ArgName != "" {
		inputArgs[config.ArgName] = &graphql.ArgumentConfig{
			Type: graphql.NewNonNull(graphql.NewList(graphql.NewNonNull(config.InputType))),
		}
	}

	return &graphql.FieldConfig{
		Description: config.Description,
		Type:        graphql.NewList(config.OutputType),
		Args:        inputArgs,
		Resolve: func(p graphql.GQLFRParams) interface{} {
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
