package starwars

import (
	"github.com/chris-ramon/graphql-go/types"
	"github.com/sogko/graphql-relay-go"
)

/**
 * This is a basic end-to-end test, designed to demonstrate the various
 * capabilities of a Relay-compliant GraphQL server.
 *
 * It is recommended that readers of this test be familiar with
 * the end-to-end test in GraphQL.js first, as this test skips
 * over the basics covered there in favor of illustrating the
 * key aspects of the Relay spec that this test is designed to illustrate.
 *
 * We will create a GraphQL schema that describes the major
 * factions and ships in the original Star Wars trilogy.
 *
 * NOTE: This may contain spoilers for the original Star
 * Wars trilogy.
 */

/**
 * Using our shorthand to describe type systems, the type system for our
 * example will be the following:
 *
 * interface Node {
 *   id: ID!
 * }
 *
 * type Faction : Node {
 *   id: ID!
 *   name: String
 *   ships: ShipConnection
 * }
 *
 * type Ship : Node {
 *   id: ID!
 *   name: String
 * }
 *
 * type ShipConnection {
 *   edges: [ShipEdge]
 *   pageInfo: PageInfo!
 * }
 *
 * type ShipEdge {
 *   cursor: String!
 *   node: Ship
 * }
 *
 * type PageInfo {
 *   hasNextPage: Boolean!
 *   hasPreviousPage: Boolean!
 *   startCursor: String
 *   endCursor: String
 * }
 *
 * type Query {
 *   rebels: Faction
 *   empire: Faction
 *   node(id: ID!): Node
 * }
 *
 * input IntroduceShipInput {
 *   clientMutationId: string!
 *   shipName: string!
 *   factionId: ID!
 * }
 *
 * input IntroduceShipPayload {
 *   clientMutationId: string!
 *   ship: Ship
 *   faction: Faction
 * }
 *
 * type Mutation {
 *   introduceShip(input IntroduceShipInput!): IntroduceShipPayload
 * }
 */

// declare definitions first, and initialize them in init() to break `initialization loop`
// i.e.:
// - nodeDefinitions refers to
// - shipType refers to
// - nodeDefinitions

var nodeDefinitions *gqlrelay.NodeDefinitions
var shipType *types.GraphQLObjectType
var factionType *types.GraphQLObjectType

// exported schema, defined in init()
var Schema types.GraphQLSchema

func init() {

	/**
	 * We get the node interface and field from the relay library.
	 *
	 * The first method is the way we resolve an ID to its object. The second is the
	 * way we resolve an object that implements node to its type.
	 */
	nodeDefinitions = gqlrelay.NewNodeDefinitions(gqlrelay.NodeDefinitionsConfig{
		IdFetcher: func(id string, info types.GraphQLResolveInfo) interface{} {
			// resolve id from global id
			resolvedId := gqlrelay.FromGlobalId(id)

			// based on id and its type, return the object
			if resolvedId.Type == "Faction" {
				return GetFaction(resolvedId.Id)
			} else {
				return GetShip(resolvedId.Id)
			}
		},
		TypeResolve: func(value interface{}, info types.GraphQLResolveInfo) *types.GraphQLObjectType {
			// based on the type of the value, return GraphQLObjectType
			switch value.(type) {
			case *Faction:
				return factionType
			default:
				return shipType
			}
		},
	})

	/**
	 * We define our basic ship type.
	 *
	 * This implements the following type system shorthand:
	 *   type Ship : Node {
	 *     id: String!
	 *     name: String
	 *   }
	 */
	shipType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:        "Ship",
		Description: "A ship in the Star Wars saga",
		Fields: types.GraphQLFieldConfigMap{
			"id": gqlrelay.GlobalIdField("Ship", nil),
			"name": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The name of the ship.",
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{
			nodeDefinitions.NodeInterface,
		},
	})

	/**
	 * We define a connection between a faction and its ships.
	 *
	 * connectionType implements the following type system shorthand:
	 *   type ShipConnection {
	 *     edges: [ShipEdge]
	 *     pageInfo: PageInfo!
	 *   }
	 *
	 * connectionType has an edges field - a list of edgeTypes that implement the
	 * following type system shorthand:
	 *   type ShipEdge {
	 *     cursor: String!
	 *     node: Ship
	 *   }
	 */
	shipConnectionDefinition := gqlrelay.ConnectionDefinitions(gqlrelay.ConnectionConfig{
		Name:     "Ship",
		NodeType: shipType,
	})

	/**
	 * We define our faction type, which implements the node interface.
	 *
	 * This implements the following type system shorthand:
	 *   type Faction : Node {
	 *     id: String!
	 *     name: String
	 *     ships: ShipConnection
	 *   }
	 */
	factionType = types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name:        "Faction",
		Description: "A faction in the Star Wars saga",
		Fields: types.GraphQLFieldConfigMap{
			"id": gqlrelay.GlobalIdField("Faction", nil),
			"name": &types.GraphQLFieldConfig{
				Type:        types.GraphQLString,
				Description: "The name of the faction.",
			},
			"ships": &types.GraphQLFieldConfig{
				Type: shipConnectionDefinition.ConnectionType,
				Args: gqlrelay.ConnectionArgs,
				Resolve: func(p types.GQLFRParams) interface{} {
					// convert args map[string]interface into ConnectionArguments
					args := gqlrelay.NewConnectionArguments(p.Args)

					// get ship objects from current faction
					ships := []interface{}{}
					if faction, ok := p.Source.(*Faction); ok {
						for _, shipId := range faction.Ships {
							ships = append(ships, GetShip(shipId))
						}
					}
					// let relay library figure out the result, given
					// - the list of ships for this faction
					// - and the filter arguments (i.e. first, last, after, before)
					return gqlrelay.ConnectionFromArray(ships, args)
				},
			},
		},
		Interfaces: []*types.GraphQLInterfaceType{
			nodeDefinitions.NodeInterface,
		},
	})

	/**
	 * This is the type that will be the root of our query, and the
	 * entry point into our schema.
	 *
	 * This implements the following type system shorthand:
	 *   type Query {
	 *     rebels: Faction
	 *     empire: Faction
	 *     node(id: String!): Node
	 *   }
	 */
	queryType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Query",
		Fields: types.GraphQLFieldConfigMap{
			"rebels": &types.GraphQLFieldConfig{
				Type: factionType,
				Resolve: func(p types.GQLFRParams) interface{} {
					return GetRebels()
				},
			},
			"empire": &types.GraphQLFieldConfig{
				Type: factionType,
				Resolve: func(p types.GQLFRParams) interface{} {
					return GetEmpire()
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	/**
	 * This will return a GraphQLFieldConfig for our ship
	 * mutation.
	 *
	 * It creates these two types implicitly:
	 *   input IntroduceShipInput {
	 *     clientMutationId: string!
	 *     shipName: string!
	 *     factionId: ID!
	 *   }
	 *
	 *   input IntroduceShipPayload {
	 *     clientMutationId: string!
	 *     ship: Ship
	 *     faction: Faction
	 *   }
	 */
	shipMutation := gqlrelay.MutationWithClientMutationId(gqlrelay.MutationConfig{
		Name: "IntroduceShip",
		InputFields: types.InputObjectConfigFieldMap{
			"shipName": &types.InputObjectFieldConfig{
				Type: types.NewGraphQLNonNull(types.GraphQLString),
			},
			"factionId": &types.InputObjectFieldConfig{
				Type: types.NewGraphQLNonNull(types.GraphQLID),
			},
		},
		OutputFields: types.GraphQLFieldConfigMap{
			"ship": &types.GraphQLFieldConfig{
				Type: shipType,
				Resolve: func(p types.GQLFRParams) interface{} {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetShip(payload["shipId"].(string))
					}
					return nil
				},
			},
			"faction": &types.GraphQLFieldConfig{
				Type: factionType,
				Resolve: func(p types.GQLFRParams) interface{} {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetFaction(payload["factionId"].(string))
					}
					return nil
				},
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info types.GraphQLResolveInfo) map[string]interface{} {
			// `inputMap` is a map with keys/fields as specified in `InputFields`
			// Note, that these fields were specified as non-nullables, so we can assume that it exists.
			shipName := inputMap["shipName"].(string)
			factionId := inputMap["factionId"].(string)

			// This mutation involves us creating (introducing) a new ship
			newShip := CreateShip(shipName, factionId)
			// return payload
			return map[string]interface{}{
				"shipId":    newShip.Id,
				"factionId": factionId,
			}
		},
	})

	/**
	 * This is the type that will be the root of our mutations, and the
	 * entry point into performing writes in our schema.
	 *
	 * This implements the following type system shorthand:
	 *   type Mutation {
	 *     introduceShip(input IntroduceShipInput!): IntroduceShipPayload
	 *   }
	 */

	mutationType := types.NewGraphQLObjectType(types.GraphQLObjectTypeConfig{
		Name: "Mutation",
		Fields: types.GraphQLFieldConfigMap{
			"introduceShip": shipMutation,
		},
	})

	/**
	 * Finally, we construct our schema (whose starting query type is the query
	 * type we defined above) and export it.
	 */
	var err error
	Schema, err = types.NewGraphQLSchema(types.GraphQLSchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		// panic if there is an error in schema
		panic(err)
	}
}
