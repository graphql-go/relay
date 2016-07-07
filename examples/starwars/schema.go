package starwars

import (
	"errors"

	"github.com/graphql-go/graphql"
	"github.com/graphql-go/relay"
	"golang.org/x/net/context"
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
 *   clientMutationID: string!
 *   shipName: string!
 *   factionId: ID!
 * }
 *
 * input IntroduceShipPayload {
 *   clientMutationID: string!
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

var nodeDefinitions *relay.NodeDefinitions
var shipType *graphql.Object
var factionType *graphql.Object

// exported schema, defined in init()
var Schema graphql.Schema

func init() {

	/**
	 * We get the node interface and field from the relay library.
	 *
	 * The first method is the way we resolve an ID to its object. The second is the
	 * way we resolve an object that implements node to its type.
	 */
	nodeDefinitions = relay.NewNodeDefinitions(relay.NodeDefinitionsConfig{
		IDFetcher: func(id string, info graphql.ResolveInfo, ctx context.Context) (interface{}, error) {
			// resolve id from global id
			resolvedID := relay.FromGlobalID(id)

			// based on id and its type, return the object
			switch resolvedID.Type {
			case "Faction":
				return GetFaction(resolvedID.ID), nil
			case "Ship":
				return GetShip(resolvedID.ID), nil
			default:
				return nil, errors.New("Unknown node type")
			}
		},
		TypeResolve: func(p graphql.ResolveTypeParams) *graphql.Object {
			// based on the type of the value, return GraphQLObjectType
			switch p.Value.(type) {
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
	shipType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Ship",
		Description: "A ship in the Star Wars saga",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Ship", nil),
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the ship.",
			},
		},
		Interfaces: []*graphql.Interface{
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
	shipConnectionDefinition := relay.ConnectionDefinitions(relay.ConnectionConfig{
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
	factionType = graphql.NewObject(graphql.ObjectConfig{
		Name:        "Faction",
		Description: "A faction in the Star Wars saga",
		Fields: graphql.Fields{
			"id": relay.GlobalIDField("Faction", nil),
			"name": &graphql.Field{
				Type:        graphql.String,
				Description: "The name of the faction.",
			},
			"ships": &graphql.Field{
				Type: shipConnectionDefinition.ConnectionType,
				Args: relay.ConnectionArgs,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					// convert args map[string]interface into ConnectionArguments
					args := relay.NewConnectionArguments(p.Args)

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
					return relay.ConnectionFromArray(ships, args), nil
				},
			},
		},
		Interfaces: []*graphql.Interface{
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
	queryType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"rebels": &graphql.Field{
				Type: factionType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetRebels(), nil
				},
			},
			"empire": &graphql.Field{
				Type: factionType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					return GetEmpire(), nil
				},
			},
			"node": nodeDefinitions.NodeField,
		},
	})

	/**
	 * This will return a GraphQLField for our ship
	 * mutation.
	 *
	 * It creates these two types implicitly:
	 *   input IntroduceShipInput {
	 *     clientMutationID: string!
	 *     shipName: string!
	 *     factionId: ID!
	 *   }
	 *
	 *   input IntroduceShipPayload {
	 *     clientMutationID: string!
	 *     ship: Ship
	 *     faction: Faction
	 *   }
	 */
	shipMutation := relay.MutationWithClientMutationID(relay.MutationConfig{
		Name: "IntroduceShip",
		InputFields: graphql.InputObjectConfigFieldMap{
			"shipName": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.String),
			},
			"factionId": &graphql.InputObjectFieldConfig{
				Type: graphql.NewNonNull(graphql.ID),
			},
		},
		OutputFields: graphql.Fields{
			"ship": &graphql.Field{
				Type: shipType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetShip(payload["shipId"].(string)), nil
					}
					return nil, nil
				},
			},
			"faction": &graphql.Field{
				Type: factionType,
				Resolve: func(p graphql.ResolveParams) (interface{}, error) {
					if payload, ok := p.Source.(map[string]interface{}); ok {
						return GetFaction(payload["factionId"].(string)), nil
					}
					return nil, nil
				},
			},
		},
		MutateAndGetPayload: func(inputMap map[string]interface{}, info graphql.ResolveInfo, ctx context.Context) (map[string]interface{}, error) {
			// `inputMap` is a map with keys/fields as specified in `InputFields`
			// Note, that these fields were specified as non-nullables, so we can assume that it exists.
			shipName := inputMap["shipName"].(string)
			factionId := inputMap["factionId"].(string)

			// This mutation involves us creating (introducing) a new ship
			newShip := CreateShip(shipName, factionId)
			// return payload
			return map[string]interface{}{
				"shipId":    newShip.ID,
				"factionId": factionId,
			}, nil
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

	mutationType := graphql.NewObject(graphql.ObjectConfig{
		Name: "Mutation",
		Fields: graphql.Fields{
			"introduceShip": shipMutation,
		},
	})

	/**
	 * Finally, we construct our schema (whose starting query type is the query
	 * type we defined above) and export it.
	 */
	var err error
	Schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query:    queryType,
		Mutation: mutationType,
	})
	if err != nil {
		// panic if there is an error in schema
		panic(err)
	}
}
