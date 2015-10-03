package gqlrelay

import "github.com/chris-ramon/graphql-go/types"

/*
Returns a GraphQLFieldConfigArgumentMap appropriate to include
on a field whose return type is a connection type.
*/
var ConnectionArgs = types.GraphQLFieldConfigArgumentMap{
	"before": &types.GraphQLArgumentConfig{
		Type: types.GraphQLString,
	},
	"after": &types.GraphQLArgumentConfig{
		Type: types.GraphQLString,
	},
	"first": &types.GraphQLArgumentConfig{
		Type: types.GraphQLInt,
	},
	"last": &types.GraphQLArgumentConfig{
		Type: types.GraphQLInt,
	},
}

func NewConnectionArgs(configMap types.GraphQLFieldConfigArgumentMap) types.GraphQLFieldConfigArgumentMap {
	for fieldName, argConfig := range ConnectionArgs {
		configMap[fieldName] = argConfig
	}
	return configMap
}

type ConnectionConfig struct {
	Name             string                      `json:"name"`
	NodeType         *types.GraphQLObjectType    `json:"nodeType"`
	EdgeFields       types.GraphQLFieldConfigMap `json:"edgeFields"`
	ConnectionFields types.GraphQLFieldConfigMap `json:"connectionFields"`
}

type EdgeType struct {
	Node   interface{}      `json:"node"`
	Cursor ConnectionCursor `json:"cursor"`
}
type GraphQLConnectionDefinitions struct {
	EdgeType       *types.GraphQLObjectType `json:"edgeType"`
	ConnectionType *types.GraphQLObjectType `json:"connectionType"`
}

/*
The common page info type used by all connections.
*/
var pageInfoType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name:        "PageInfo",
	Description: "Information about pagination in a connection.",
	Fields: types.GraphQLFieldConfigMap{
		"hasNextPage": &types.GraphQLFieldConfig{
			Type:        types.NewGraphQLNonNull(types.GraphQLBoolean),
			Description: "When paginating forwards, are there more items?",
		},
		"hasPreviousPage": &types.GraphQLFieldConfig{
			Type:        types.NewGraphQLNonNull(types.GraphQLBoolean),
			Description: "When paginating backwards, are there more items?",
		},
		"startCursor": &types.GraphQLFieldConfig{
			Type:        types.GraphQLString,
			Description: "When paginating backwards, the cursor to continue.",
		},
		"endCursor": &types.GraphQLFieldConfig{
			Type:        types.GraphQLString,
			Description: "When paginating forwards, the cursor to continue.",
		},
	},
})

/*
Returns a GraphQLObjectType for a connection with the given name,
and whose nodes are of the specified type.
*/

func ConnectionDefinitions(config ConnectionConfig) *GraphQLConnectionDefinitions {

	edgeType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:        config.Name + "Edge",
		Description: "An edge in a connection",
		Fields: types.GraphQLFieldConfigMap{
			"node": &types.GraphQLFieldConfig{
				Type:        config.NodeType,
				Description: "The item at the end of the edge",
			},
			"cursor": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLNonNull(types.GraphQLString),
				Description: " cursor for use in pagination",
			},
		},
	})
	for fieldName, fieldConfig := range config.EdgeFields {
		edgeType.AddFieldConfig(fieldName, fieldConfig)
	}

	connectionType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:        config.Name + "Connection",
		Description: "A connection to a list of items.",

		Fields: types.GraphQLFieldConfigMap{
			"pageInfo": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLNonNull(pageInfoType),
				Description: "Information to aid in pagination.",
			},
			"edges": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLList(edgeType),
				Description: "Information to aid in pagination.",
			},
		},
	})
	for fieldName, fieldConfig := range config.ConnectionFields {
		connectionType.AddFieldConfig(fieldName, fieldConfig)
	}

	return &GraphQLConnectionDefinitions{
		EdgeType:       edgeType,
		ConnectionType: connectionType,
	}
}
