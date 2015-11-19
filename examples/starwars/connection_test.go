package starwars_test

import (
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/testutil"
	"github.com/graphql-go/relay/examples/starwars"
	"reflect"
	"testing"
)

func TestConnection_TestFetching_CorrectlyFetchesTheFirstShipOfTheRebels(t *testing.T) {
	query := `
        query RebelsShipsQuery {
          rebels {
            name,
            ships(first: 1) {
              edges {
                node {
                  name
                }
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"name": "X-Wing",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyFetchesTheFirstTwoShipsOfTheRebelsWithACursor(t *testing.T) {
	query := `
        query MoreRebelShipsQuery {
          rebels {
            name,
            ships(first: 2) {
              edges {
                cursor,
                node {
                  name
                }
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjA=",
							"node": map[string]interface{}{
								"name": "X-Wing",
							},
						},
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjE=",
							"node": map[string]interface{}{
								"name": "Y-Wing",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyFetchesTheNextThreeShipsOfTheRebelsWithACursor(t *testing.T) {
	query := `
        query EndOfRebelShipsQuery {
          rebels {
            name,
            ships(first: 3 after: "YXJyYXljb25uZWN0aW9uOjE=") {
              edges {
                cursor,
                node {
                  name
                }
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjI=",
							"node": map[string]interface{}{
								"name": "A-Wing",
							},
						},
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjM=",
							"node": map[string]interface{}{
								"name": "Millenium Falcon",
							},
						},
						map[string]interface{}{
							"cursor": "YXJyYXljb25uZWN0aW9uOjQ=",
							"node": map[string]interface{}{
								"name": "Home One",
							},
						},
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyFetchesNoShipsOfTheRebelsAtTheEndOfTheConnection(t *testing.T) {
	query := `
        query RebelsQuery {
          rebels {
            name,
            ships(first: 3 after: "YXJyYXljb25uZWN0aW9uOjQ=") {
              edges {
                cursor,
                node {
                  name
                }
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"ships": map[string]interface{}{
					"edges": []interface{}{},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
func TestConnection_TestFetching_CorrectlyIdentifiesTheEndOfTheList(t *testing.T) {
	query := `
        query EndOfRebelShipsQuery {
          rebels {
            name,
            originalShips: ships(first: 2) {
              edges {
                node {
                  name
                }
              }
              pageInfo {
                hasNextPage
              }
            }
            moreShips: ships(first: 3 after: "YXJyYXljb25uZWN0aW9uOjE=") {
              edges {
                node {
                  name
                }
              }
              pageInfo {
                hasNextPage
              }
            }
          }
        }
      `
	expected := &graphql.Result{
		Data: map[string]interface{}{
			"rebels": map[string]interface{}{
				"name": "Alliance to Restore the Republic",
				"originalShips": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"name": "X-Wing",
							},
						},
						map[string]interface{}{
							"node": map[string]interface{}{
								"name": "Y-Wing",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": true,
					},
				},
				"moreShips": map[string]interface{}{
					"edges": []interface{}{
						map[string]interface{}{
							"node": map[string]interface{}{
								"name": "A-Wing",
							},
						},
						map[string]interface{}{
							"node": map[string]interface{}{
								"name": "Millenium Falcon",
							},
						},
						map[string]interface{}{
							"node": map[string]interface{}{
								"name": "Home One",
							},
						},
					},
					"pageInfo": map[string]interface{}{
						"hasNextPage": false,
					},
				},
			},
		},
	}
	result := graphql.Do(graphql.Params{
		Schema:        starwars.Schema,
		RequestString: query,
	})
	if !reflect.DeepEqual(result, expected) {
		t.Fatalf("wrong result, graphql result diff: %v", testutil.Diff(expected, result))
	}
}
