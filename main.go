// credit - go-graphql hello world example
package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/graphql-go/graphql"
	"github.com/pborman/getopt"
	"github.com/pborman/uuid"
)

var playerType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Player",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"team": &graphql.Field{
				Type: teamType,
			},
			"statistics": &graphql.Field{
				Type: statsType,
			},
			"age": &graphql.Field{
				Type: graphql.Int,
			},
			"number": &graphql.Field{
				Type: graphql.String,
			},
			"email": &graphql.Field{
				Type: graphql.String,
			},
		},
	},
)

var teamType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Team",
		Fields: graphql.Fields{
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"wins": &graphql.Field{
				Type: graphql.Int,
			},
			"loses": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var seasonType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Season",
		Fields: graphql.Fields{
			"games": &graphql.Field{
				Type: graphql.NewList(gameType),
			},
		},
	},
)

var statsType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Stats",
		Fields: graphql.Fields{
			"goals": &graphql.Field{
				Type: graphql.Int,
			},
			"assists": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var gameType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "Game",
		Fields: graphql.Fields{
			"location": &graphql.Field{
				Type: graphql.String,
			},
			"winner": &graphql.Field{
				Type: teamType,
			},
		},
	},
)

func main() {
	dbUser := getopt.StringLong("db-user", 'u', "", "databse username")
	dbPass := getopt.StringLong("db-pass", 'p', "", "database password")
	optHelp := getopt.BoolLong("help", 0, "help menu")
	getopt.Parse()

	fmt.Println(*dbUser, *dbPass)

	if *dbUser == "" || *dbPass == "" {
		getopt.Usage()
		os.Exit(1)
	}

	if *optHelp {
		getopt.Usage()
		os.Exit(0)
	}

	config := config{
		db: &db{
			user: *dbUser,
			pass: *dbPass,
		},
	}

	db, err := newDatabase(&config)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	id := uuid.NewUUID().String()

	np := player{
		ID:     id,
		Email:  "brian.downs@gmail.com",
		Age:    20,
		Number: "28",
		Stats: &statistics{
			Goals:   1,
			Assists: 10,
		},
	}
	if err := db.AddPlayer(&np); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// Schema
	fields := graphql.Fields{
		"player": &graphql.Field{
			Type: playerType,
			// it's good form to add a description
			// to each field.
			Description: "Get player By ID",
			// We can define arguments that allow us to
			// pick specific tutorials. In this case
			// we want to be able to specify the ID of the
			// tutorial we want to retrieve
			Args: graphql.FieldConfigArgument{
				"id": &graphql.ArgumentConfig{
					Type: graphql.String,
				},
			},
			Resolve: func(p graphql.ResolveParams) (interface{}, error) {
				// take in the ID argument
				id, ok := p.Args["id"].(string)
				if ok {
					return db.PlayerByID(id)
					// Parse our tutorial array for the matching id
					// for _, player := range player {
					// 	if int(player.ID) == id {
					// 		// return our tutorial
					// 		return player, nil
					// 	}
					// }
				}
				return nil, nil
			},
		},
		// this is our `list` endpoint which will return all
		// tutorials available
		"list": &graphql.Field{
			Type:        graphql.NewList(playerType),
			Description: "Get Tutorial List",
			Resolve: func(params graphql.ResolveParams) (interface{}, error) {
				return db.Players()
			},
		},
	}
	rootQuery := graphql.ObjectConfig{
		Name:   "RootQuery",
		Fields: fields,
	}
	schemaConfig := graphql.SchemaConfig{
		Query: graphql.NewObject(rootQuery),
	}
	schema, err := graphql.NewSchema(schemaConfig)
	if err != nil {
		log.Fatalf("failed to create new schema, error: %v", err)
	}

	// Query
	query := fmt.Sprintf(`{player(id: "%s"){id,number}}`, id)
	params := graphql.Params{Schema: schema, RequestString: query}
	r := graphql.Do(params)
	if len(r.Errors) > 0 {
		fmt.Printf("failed to execute graphql operation, errors: %+v\n", r.Errors)
		os.Exit(1)
	}
	b, err := json.Marshal(r)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	fmt.Printf("ID: %s \n", string(b))

	router := mux.NewRouter()

	router.HandleFunc("/lacrosse", func(w http.ResponseWriter, r *http.Request) {
		result := executeQuery(r.URL.Query().Get("query"), schema)
		if err := json.NewEncoder(w).Encode(result); err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	})

	fmt.Println("Test with Get      : curl -g 'http://localhost:8080/lacrosse?query={player(id: <id>){id,number}}'")

	if err := http.ListenAndServe(":8080", router); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func executeQuery(query string, schema graphql.Schema) *graphql.Result {
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: query,
	})
	if len(result.Errors) > 0 {
		fmt.Printf("wrong result, unexpected errors: %v", result.Errors)
		return nil
	}
	fmt.Println("result", result)
	return result
}
