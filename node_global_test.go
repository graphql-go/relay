package gqlrelay_test

import (
	"fmt"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql-relay-go"
	"github.com/graphql-go/graphql/testutil"
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
var globalIDTestUserType *graphql.Object
var globalIDTestPhotoType *graphql.Object

var globalIDTestDef = gqlrelay.NewNodeDefinitions(gqlrelay.NodeDefinitionsConfig{
	IDFetcher: func(globalID string, info graphql.ResolveInfo) interface{} {
		resolvedGlobalId := gqlrelay.FromGlobalID(globalID)
		if resolvedGlobalId == nil {
			return nil
		}
		if resolvedGlobalId.Type == "User" {
			return globalIDTestUserData[resolvedGlobalId.ID]
		} else {
			return globalIDTestPhotoData[resolvedGlobalId.ID]
		}
	},
	TypeResolve: func(value interface{}, info graphql.ResolveInfo) *graphql.Object {
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
var globalIDTestQueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Query",
	Fields: graphql.FieldConfigMap{
		"node": globalIDTestDef.NodeField,
		"allObjects": &graphql.FieldConfig{
			Type: graphql.NewList(globalIDTestDef.NodeInterface),
			Resolve: func(p graphql.GQLFRParams) interface{} {
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
var globalIDTestSchema graphql.Schema

func init() {
	globalIDTestUserType = graphql.NewObject(graphql.ObjectConfig{
		Name: "User",
		Fields: graphql.FieldConfigMap{
			"id": gqlrelay.GlobalIDField("User", nil),
			"name": &graphql.FieldConfig{
				Type: graphql.String,
			},
		},
		Interfaces: []*graphql.Interface{globalIDTestDef.NodeInterface},
	})
	photoIDFetcher := func(obj interface{}, info graphql.ResolveInfo) string {
		switch obj := obj.(type) {
		case *photo2:
			return fmt.Sprintf("%v", obj.PhotoId)
		}
		return ""
	}
	globalIDTestPhotoType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Photo",
		Fields: graphql.FieldConfigMap{
			"id": gqlrelay.GlobalIDField("Photo", photoIDFetcher),
			"width": &graphql.FieldConfig{
				Type: graphql.Int,
			},
		},
		Interfaces: []*graphql.Interface{globalIDTestDef.NodeInterface},
	})

	globalIDTestSchema, _ = graphql.NewSchema(graphql.SchemaConfig{
		Query: globalIDTestQueryType,
	})
}

func TestGlobalIDFields_GivesDifferentIDs(t *testing.T) {
	query := `{
      allObjects {
        id
      }
    }`
	expected := &graphql.Result{
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
	result := graphql.Graphql(graphql.Params{
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
	expected := &graphql.Result{
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
	result := graphql.Graphql(graphql.Params{
		Schema:        globalIDTestSchema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
