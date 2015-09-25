package graphql_relay_test

import (
	"fmt"
	"github.com/chris-ramon/graphql-go"
	"github.com/chris-ramon/graphql-go/testutil"
	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-relay-go"
	"reflect"
	"testing"
)

type photo2 struct {
	PhotoId int `json:"photoId"`
	Width   int `json:"width"`
}

var globalIDTestUserData = map[string]*user{
	"1": &user{1, "John Doe"},
	"2": &user{2, "Jane Smith"},
}
var globalIDTestPhotoData = map[string]*photo2{
	"1": &photo2{1, 300},
	"2": &photo2{2, 400},
}

// declare types first, define later in init()
// because they all depend on nodeTestDef
var globalIDTestUserType *types.GraphQLObjectType
var globalIDTestPhotoType *types.GraphQLObjectType

var globalIDTestDef = graphql_relay.NewNodeDefinitions(graphql_relay.NodeDefinitionsConfig{
	IdFetcher: func(globalId string, info types.GraphQLResolveInfo) interface{} {
		resolvedGlobalId := graphql_relay.FromGlobalId(globalId)
		if resolvedGlobalId == nil {
			return nil
		}
		if resolvedGlobalId.Type == "User" {
			return globalIDTestUserData[resolvedGlobalId.Id]
		} else {
			return globalIDTestPhotoData[resolvedGlobalId.Id]
		}
	},
	TypeResolve: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
		switch value.(type) {
		case *user:
			return globalIDTestUserType
		case *photo2:
			return globalIDTestPhotoType
		default:
			panic(fmt.Sprintf("Unknown object type `%v`", value))
		}
	},
})
var globalIDTestQueryType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
	Name: "Query",
	Fields: types.GraphQLFieldConfigMap{
		"node": globalIDTestDef.NodeField,
		"allObjects": &types.GraphQLFieldConfig{
			Type: types.NewGraphQLList(globalIDTestDef.NodeInterface),
			Resolve: func(p types.GQLFRParams) interface{} {
				return []interface{}{
					globalIDTestUserData["1"],
					globalIDTestUserData["2"],
					globalIDTestPhotoData["1"],
					globalIDTestPhotoData["2"],
				}
			},
		},
	},
})

// becareful not to define schema here, since globalIDTestUserType and globalIDTestPhotoType wouldn't be defined till init()
var globalIDTestSchema types.GraphQLSchema

func init() {
	globalIDTestUserType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "User",
		Fields: types.GraphQLFieldConfigMap{
			"id": graphql_relay.GlobalIdField("User", nil),
			"name": &types.GraphQLFieldConfig{
				Type: types.GraphQLString,
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{globalIDTestDef.NodeInterface},
	})
	photoIdFetcher := func(obj interface{}, info types.GraphQLResolveInfo) string {
		switch obj := obj.(type) {
		case *photo2:
			return fmt.Sprintf("%v", obj.PhotoId)
		}
		return ""
	}
	globalIDTestPhotoType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Photo",
		Fields: types.GraphQLFieldConfigMap{
			"id": graphql_relay.GlobalIdField("Photo", photoIdFetcher),
			"width": &types.GraphQLFieldConfig{
				Type: types.GraphQLInt,
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{globalIDTestDef.NodeInterface},
	})

	globalIDTestSchema, _ = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query: globalIDTestQueryType,
	})
}

func TestGlobalIDFields_GivesDifferentIDs(t *testing.T) {
	query := `{
      allObjects {
        id
      }
    }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"allObjects": []interface{}{
				map[string]interface{}{
					"id": "VXNlcjox",
				},
				map[string]interface{}{
					"id": "VXNlcjoy",
				},
				map[string]interface{}{
					"id": "UGhvdG86MQ==",
				},
				map[string]interface{}{
					"id": "UGhvdG86Mg==",
				},
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        globalIDTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}

func TestGlobalIDFields_RefetchesTheIDs(t *testing.T) {
	query := `{
      user: node(id: "VXNlcjox") {
        id
        ... on User {
          name
        }
      },
      photo: node(id: "UGhvdG86MQ==") {
        id
        ... on Photo {
          width
        }
      }
    }`
	expected := &types.GraphQLResult{
		Data: map[string]interface{}{
			"user": map[string]interface{}{
				"id":   "VXNlcjox",
				"name": "John Doe",
			},
			"photo": map[string]interface{}{
				"id":    "UGhvdG86MQ==",
				"width": 300,
			},
		},
	}
	result := graphql(t, gql.GraphqlParams{
		Schema:        globalIDTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
