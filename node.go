package graphql_relay

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/chris-ramon/graphql-go/types"
	"strings"
)

type NodeDefinitions struct {
	NodeInterface *types.GraphQLInterfaceType
	NodeField     *types.GraphQLFieldConfig
}

type NodeDefinitionsConfig struct {
	IdFetcher   IdFetcherFn
	TypeResolve types.ResolveTypeFn
}
type IdFetcherFn func(id string, info types.GraphQLResolveInfo) interface{}
type GlobalIdFetcherFn func(obj interface{}, info types.GraphQLResolveInfo) string

/*
 Given a function to map from an ID to an underlying object, and a function
 to map from an underlying object to the concrete GraphQLObjectType it
 corresponds to, constructs a `Node` interface that objects can implement,
 and a field config for a `node` root field.

 If the typeResolver is omitted, object resolution on the interface will be
 handled with the `isTypeOf` method on object types, as with any GraphQL
interface without a provided `resolveType` method.
*/
func NewNodeDefinitions(config NodeDefinitionsConfig) *NodeDefinitions {
	nodeInterface := types.NewGraphQLInterfaceType(types.GraphQLInterfaceTypeConfig{
		Name:        "Node",
		Description: "An object with an ID",
		Fields: types.GraphQLFieldConfigMap{
			"id": &types.GraphQLFieldConfig{
				Type:        types.NewGraphQLNonNull(types.GraphQLID),
				Description: "The id of the object",
			},
		},
		ResolveType: config.TypeResolve,
	})

	nodeField := &types.GraphQLFieldConfig{
		Name:        "Node",
		Description: "Fetches an object given its ID",
		Type:        nodeInterface,
		Args: types.GraphQLFieldConfigArgumentMap{
			"id": &types.GraphQLArgumentConfig{
				Type:        types.NewGraphQLNonNull(types.GraphQLID),
				Description: "The 111111 ID of an object",
			},
		},
		Resolve: func(p types.GQLFRParams) interface{} {
			if config.IdFetcher == nil {
				return nil
			}
			id := ""
			if iid, ok := p.Args["id"]; ok {
				id = fmt.Sprintf("%v", iid)
			}
			fetchedId := config.IdFetcher(id, p.Info)
			return fetchedId
		},
	}
	return &NodeDefinitions{
		NodeInterface: nodeInterface,
		NodeField:     nodeField,
	}
}

type ResolvedGlobalId struct {
	Type string `json:"type"`
	Id   string `json:"id"`
}

/*
Takes a type name and an ID specific to that type name, and returns a
"global ID" that is unique among all types.
*/
func ToGlobalId(ttype string, id string) string {
	str := ttype + ":" + id
	encStr := base64.StdEncoding.EncodeToString([]byte(str))
	return encStr
}

/*
Takes the "global ID" created by toGlobalID, and returns the type name and ID
used to create it.
*/
func FromGlobalId(globalId string) *ResolvedGlobalId {
	strId := ""
	b, err := base64.StdEncoding.DecodeString(globalId)
	if err == nil {
		strId = string(b)
	}
	tokens := strings.Split(strId, ":")
	if len(tokens) < 2 {
		return nil
	}
	return &ResolvedGlobalId{
		Type: tokens[0],
		Id:   tokens[1],
	}
}

/*
Creates the configuration for an id field on a node, using `toGlobalId` to
construct the ID from the provided typename. The type-specific ID is fetcher
by calling idFetcher on the object, or if not provided, by accessing the `id`
property on the object.
*/
func GlobalIdField(typeName string, idFetcher GlobalIdFetcherFn) *types.GraphQLFieldConfig {
	return &types.GraphQLFieldConfig{
		Name:        "id",
		Description: "The ID of an object",
		Type:        types.NewGraphQLNonNull(types.GraphQLID),
		Resolve: func(p types.GQLFRParams) interface{} {
			id := ""
			if idFetcher != nil {
				fetched := idFetcher(p.Source, p.Info)
				id = fmt.Sprintf("%v", fetched)
			} else {
				// try to get from p.Source (data)
				var objMap interface{}
				b, _ := json.Marshal(p.Source)
				_ = json.Unmarshal(b, &objMap)
				switch obj := objMap.(type) {
				case map[string]interface{}:
					if iid, ok := obj["id"]; ok {
						id = fmt.Sprintf("%v", iid)
					}
				}
			}
			globalId := ToGlobalId(typeName, id)
			return globalId
		},
	}
}
